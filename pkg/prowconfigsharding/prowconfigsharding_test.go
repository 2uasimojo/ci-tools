package prowconfigsharding

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"

	pluginsflagutil "k8s.io/test-infra/prow/flagutil/plugins"
	"k8s.io/test-infra/prow/plugins"
	"sigs.k8s.io/yaml"
)

func TestShardPluginConfig(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		in   *plugins.Configuration

		expectedConfig     *plugins.Configuration
		expectedShardFiles map[string]string
	}{
		{
			name: "Plugin config gets sharded",
			in: &plugins.Configuration{
				Plugins: plugins.Plugins{
					"openshift":         plugins.OrgPlugins{Plugins: []string{"foo"}},
					"openshift/release": plugins.OrgPlugins{Plugins: []string{"bar"}},
				},
				Cat: plugins.Cat{KeyPath: "/etc/raw"},
			},

			expectedConfig: &plugins.Configuration{
				Plugins: plugins.Plugins{},
				Cat:     plugins.Cat{KeyPath: "/etc/raw"},
			},
			expectedShardFiles: map[string]string{
				"openshift/_pluginconfig.yaml": strings.Join([]string{
					"plugins:",
					"  openshift:",
					"    plugins:",
					"    - foo",
					"",
				}, "\n"),
				"openshift/release/_pluginconfig.yaml": strings.Join([]string{
					"plugins:",
					"  openshift/release:",
					"    plugins:",
					"    - bar",
					"",
				}, "\n"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serializedInitialConfig, err := yaml.Marshal(tc.in)
			if err != nil {
				t.Fatalf("failed to serialize initial config: %v", err)
			}

			afs := afero.NewMemMapFs()

			updated, err := WriteShardedPluginConfig(tc.in, afs)
			if err != nil {
				t.Fatalf("failed to shard plugin config: %v", err)
			}
			if diff := cmp.Diff(tc.expectedConfig, updated); diff != "" {
				t.Errorf("updated plugin config differs from expected: %s", diff)
			}

			shardedConfigFiles := map[string]string{}
			if err := afero.Walk(afs, "", func(path string, info fs.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				if filepath.Base(path) != "_pluginconfig.yaml" {
					t.Errorf("found file %s which doesn't have the expected _prowconfig.yaml name", path)
				}
				data, err := afero.ReadFile(afs, path)
				if err != nil {
					t.Errorf("failed to read file %s: %v", path, err)
				}
				shardedConfigFiles[path] = string(data)
				return nil
			}); err != nil {
				t.Errorf("waking the fs failed: %v", err)
			}

			if diff := cmp.Diff(tc.expectedShardFiles, shardedConfigFiles); diff != "" {
				t.Fatalf("actual sharded config differs from expected:\n%s", diff)
			}

			// Test that when we load the sharded config its identical to the config with which we started
			tempDir := t.TempDir()

			// We need to write and load the initial config to put it through defaulting
			if err := ioutil.WriteFile(filepath.Join(tempDir, "_original_config.yaml"), serializedInitialConfig, 0644); err != nil {
				t.Fatalf("failed to write out serialized initial config: %v", err)
			}
			// Defaulting is unexported and only happens inside plugins.ConfigAgent.Load()
			initialConfigAgent := plugins.ConfigAgent{}
			if err := initialConfigAgent.Start(filepath.Join(tempDir, "_original_config.yaml"), nil, "", false); err != nil {
				t.Fatalf("failed to start old plugin config agent: %v", err)
			}

			serializedNewConfig, err := yaml.Marshal(updated)
			if err != nil {
				t.Fatalf("failed to marshal the new config: %v", err)
			}
			if err := ioutil.WriteFile(filepath.Join(tempDir, "_plugins.yaml"), serializedNewConfig, 0644); err != nil {
				t.Fatalf("failed to write new config: %v", err)
			}

			for name, content := range shardedConfigFiles {
				path := filepath.Join(tempDir, name)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("failed to create directories for %s: %v", path, err)
				}
				if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write file %s: %v", name, err)
				}
			}

			fs := &flag.FlagSet{}
			opts := pluginsflagutil.PluginOptions{}
			opts.AddFlags(fs)
			if err := fs.Parse([]string{
				"--plugin-config=" + filepath.Join(tempDir, "_plugins.yaml"),
				"--supplemental-plugin-config-dir=" + tempDir,
			}); err != nil {
				t.Fatalf("faield to parse flags")
			}

			pluginAgent, err := opts.PluginAgent()
			if err != nil {
				t.Fatalf("failed to construct plugin agent: %v", err)
			}
			if diff := cmp.Diff(initialConfigAgent.Config(), pluginAgent.Config(), cmp.Exporter(func(_ reflect.Type) bool { return true })); diff != "" {
				t.Errorf("initial config differs from what we got when sharding: %s", diff)
			}
		})
	}
}
