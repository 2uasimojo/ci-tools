package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	hivev1 "github.com/openshift/hive/apis/hive/v1"

	"github.com/openshift/ci-tools/pkg/testhelper"
)

func TestEnsureLabels(t *testing.T) {
	testCases := []struct {
		name             string
		given            hivev1.ClusterPool
		expected         hivev1.ClusterPool
		expectedModified bool
	}{
		{
			name: "basic case",
			given: hivev1.ClusterPool{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "dpp",
					},
				},
			},
			expected: hivev1.ClusterPool{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "dpp",
					},
				},
				Spec: hivev1.ClusterPoolSpec{
					Labels: map[string]string{"tp.openshift.io/owner": "dpp"},
				},
			},
			expectedModified: true,
		},
		{
			name: "not modified",
			given: hivev1.ClusterPool{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "dpp",
					},
				},
				Spec: hivev1.ClusterPoolSpec{
					Labels: map[string]string{"tp.openshift.io/owner": "dpp"},
				},
			},
			expected: hivev1.ClusterPool{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "dpp",
					},
				},
				Spec: hivev1.ClusterPoolSpec{
					Labels: map[string]string{"tp.openshift.io/owner": "dpp"},
				},
			},
		},
		{
			name: "modified",
			given: hivev1.ClusterPool{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "dpp",
					},
				},
				Spec: hivev1.ClusterPoolSpec{
					Labels: map[string]string{"tp.openshift.io/owner": "not-dpp"},
				},
			},
			expected: hivev1.ClusterPool{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "dpp",
					},
				},
				Spec: hivev1.ClusterPoolSpec{
					Labels: map[string]string{"tp.openshift.io/owner": "dpp"},
				},
			},
			expectedModified: true,
		},
		{
			name: "given has no labels",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, actualModified := ensureLabels(tc.given)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("%s differs from expected:\n%s", tc.name, diff)
			}
			if diff := cmp.Diff(tc.expectedModified, actualModified); diff != "" {
				t.Errorf("%s differs from expected:\n%s", tc.name, diff)
			}
		})
	}
}

func TestEnsureLabelsOnClusterPools(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestEnsureLabelsOnClusterPools")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	testCases := []struct {
		name            string
		input           string
		output          string
		expected        error
		expectedContent string
	}{
		{
			name:   "basic case",
			input:  filepath.Join("testdata", "pools", "cvp-ocp-4-9-amd64-aws-us-west-2_clusterpool.yaml"),
			output: filepath.Join(dir, "cvp-ocp-4-9-amd64-aws-us-west-2_clusterpool.yaml"),
			expectedContent: `apiVersion: hive.openshift.io/v1
kind: ClusterPool
metadata:
  creationTimestamp: null
  labels:
    architecture: amd64
    cloud: aws
    owner: cvp
    product: ocp
    region: us-west-2
    version: "4.9"
    version_lower: 4.9.0-0
    version_upper: 4.10.0-0
  name: cvp-ocp-4-9-amd64-aws-us-west-2
  namespace: cvp-cluster-pool
spec:
  baseDomain: cpaas-ci.devcluster.openshift.com
  hibernationConfig:
    resumeTimeout: 15m0s
  imageSetRef:
    name: ocp-release-4.9.57-x86-64-for-4.9.0-0-to-4.10.0-0
  installAttemptsLimit: 1
  installConfigSecretTemplateRef:
    name: install-config-aws-us-west-2
  labels:
    tp.openshift.io/owner: cvp
  maxSize: 10
  platform:
    aws:
      credentialsSecretRef:
        name: cvp-aws-credentials
      region: us-west-2
  pullSecretRef:
    name: pull-secret
  runningCount: 1
  size: 5
  skipMachinePools: true
status:
  ready: 0
  size: 0
  standby: 0
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := addSchemes(); err != nil {
				t.Fatal("Failed to set up scheme")
			}
			s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme,
				scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: false, Strict: false})
			actual := ensureLabelsOnClusterPool(s, tc.input, tc.output)
			if diff := cmp.Diff(tc.expected, actual, testhelper.EquateErrorMessage); diff != "" {
				t.Errorf("%s differs from expected:\n%s", tc.name, diff)
			}
			if actual == nil {
				raw, err := os.ReadFile(tc.output)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
				}
				if diff := cmp.Diff(tc.expectedContent, string(raw)); diff != "" {
					t.Errorf("%s differs from expected:\n%s", tc.name, diff)
				}
			}
		})
	}
}
