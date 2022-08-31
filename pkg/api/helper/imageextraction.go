package helper

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/ci-tools/pkg/api"
)

// ImageStreamTagMap is a map [types.NamespacedName.String()]types.NamespacedName of
// ImagestreamTags
type ImageStreamTagMap map[string]types.NamespacedName

func (istm ImageStreamTagMap) String() string {
	var ret []string

	for fullTag := range istm {
		ret = append(ret, fullTag)
	}

	return strings.Join(ret, ", ")
}

// MergeImageStramTagMaps merges multiple ImageStreamTagMaps together
func MergeImageStreamTagMaps(target ImageStreamTagMap, toMerge ...ImageStreamTagMap) {
	for _, toMerge := range toMerge {
		for k, v := range toMerge {
			target[k] = v
		}
	}
}

func TestInputImageStreamsFromResolvedConfig(cfg api.ReleaseBuildConfiguration) []types.NamespacedName {
	s := map[types.NamespacedName]struct{}{}
	add := func(ns, name string) {
		s[types.NamespacedName{Namespace: ns, Name: name}] = struct{}{}
	}
	if c := cfg.ReleaseTagConfiguration; c != nil {
		add(c.Namespace, c.Name)
	}
	for _, r := range cfg.Releases {
		if i := r.Integration; i != nil {
			add(i.Namespace, i.Name)
		}
	}
	var ret []types.NamespacedName
	for k := range s {
		ret = append(ret, k)
	}
	return ret
}

// TestInputImageStreamTagsFromResolvedConfig returns all ImageStreamTags referenced anywhere in the config as input.
// It only returns their namespace and name and drops the cluster field, as we plan to remove that.
// The key is in namespace/name format.
// It assumes that the config is already resolved, i.E. that MultiStageTestConfiguration is always nil
// and MultiStageTestConfigurationLiteral gets set instead
func TestInputImageStreamTagsFromResolvedConfig(cfg api.ReleaseBuildConfiguration) (ImageStreamTagMap, error) {
	result := map[string]types.NamespacedName{}

	imageStreamTagReferenceMapIntoMap(cfg.BaseImages, result)
	imageStreamTagReferenceMapIntoMap(cfg.BaseRPMImages, result)
	if cfg.BuildRootImage != nil {
		if cfg.BuildRootImage.ImageStreamTagReference != nil {
			insert(*cfg.BuildRootImage.ImageStreamTagReference, result)
		}
		if cfg.BuildRootImage.UseBuildCache {
			insert(api.BuildCacheFor(cfg.Metadata).ImageStreamTagReference, result)
		}
	}

	var errs []error
	for _, testStep := range cfg.Tests {
		if testStep.MultiStageTestConfigurationLiteral != nil {
			insertTagReferencesFromSteps(*testStep.MultiStageTestConfigurationLiteral, result)
		}
		if testStep.MultiStageTestConfiguration != nil && testStep.MultiStageTestConfigurationLiteral == nil {
			errs = append(errs, errors.New("got unresolved config"))
		}
	}

	for _, rawStep := range cfg.RawSteps {
		if rawStep.InputImageTagStepConfiguration != nil {
			insert(rawStep.InputImageTagStepConfiguration.BaseImage.ImageStreamTagReference, result)
		}
		if rawStep.SourceStepConfiguration != nil {
			insert(rawStep.SourceStepConfiguration.ClonerefsImage, result)
		}
		if rawStep.TestStepConfiguration != nil {
			if rawStep.TestStepConfiguration.MultiStageTestConfigurationLiteral != nil {
				insertTagReferencesFromSteps(*rawStep.TestStepConfiguration.MultiStageTestConfigurationLiteral, result)
			}
			if rawStep.TestStepConfiguration.MultiStageTestConfiguration != nil && rawStep.TestStepConfiguration.MultiStageTestConfigurationLiteral == nil {
				errs = append(errs, errors.New("got unresolved config"))
			}
		}
	}

	return ImageStreamTagMap(result), utilerrors.NewAggregate(errs)
}

func imageStreamTagReferenceMapIntoMap(i map[string]api.ImageStreamTagReference, m map[string]types.NamespacedName) {
	for _, item := range i {
		insert(item, m)
	}
}

func imageStreamTagReferenceToString(istr api.ImageStreamTagReference) string {
	return fmt.Sprintf("%s/%s:%s", istr.Namespace, istr.Name, istr.Tag)
}

func insertTagReferencesFromSteps(config api.MultiStageTestConfigurationLiteral, m map[string]types.NamespacedName) {
	for _, subStep := range append(append(config.Pre, config.Test...), config.Post...) {
		if subStep.FromImage != nil {
			insert(*subStep.FromImage, m)
		}
	}
	for _, observer := range config.Observers {
		insert(*observer.FromImage, m)
	}
}

func insert(item api.ImageStreamTagReference, m map[string]types.NamespacedName) {
	if _, ok := m[imageStreamTagReferenceToString(item)]; ok {
		return
	}
	m[imageStreamTagReferenceToString(item)] = types.NamespacedName{
		Namespace: item.Namespace,
		Name:      fmt.Sprintf("%s:%s", item.Name, item.Tag),
	}
}
