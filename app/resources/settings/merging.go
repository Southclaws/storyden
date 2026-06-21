package settings

import (
	"reflect"
	"strings"

	"dario.cat/mergo"
)

// NOTE: somewhat of a hack to get mergo to merge opt.Optional[T] values.

type settingsMergeTransformer struct{}

func (settingsMergeTransformer) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	if !isOptionalType(t) {
		return nil
	}

	return mergeOptionalValue
}

func mergeOptionalValue(dst, src reflect.Value) error {
	if src.IsNil() || src.Len() == 0 {
		return nil
	}

	if dst.IsNil() || dst.Len() == 0 {
		dst.Set(src)
		return nil
	}

	dstValue := dst.Index(0)
	srcValue := src.Index(0)

	if shouldMergeValue(dstValue.Type()) {
		return mergePresentValue(dstValue, srcValue)
	}

	dstValue.Set(srcValue)

	return nil
}

func mergePresentValue(dst, src reflect.Value) error {
	if dst.Kind() == reflect.Map {
		return mergeMapValue(dst, src)
	}

	return mergo.Merge(
		dst.Addr().Interface(),
		src.Interface(),
		mergo.WithOverride,
		mergo.WithTransformers(settingsMergeTransformer{}),
	)
}

func mergeMapValue(dst, src reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	if src.Len() == 0 {
		dst.Set(src)
		return nil
	}

	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	for _, key := range src.MapKeys() {
		srcValue := src.MapIndex(key)
		dstValue := dst.MapIndex(key)
		if !dstValue.IsValid() || !shouldMergeValue(srcValue.Type()) {
			dst.SetMapIndex(key, srcValue)
			continue
		}

		mergedValue := reflect.New(dstValue.Type()).Elem()
		mergedValue.Set(dstValue)
		if err := mergePresentValue(mergedValue, srcValue); err != nil {
			return err
		}

		dst.SetMapIndex(key, mergedValue)
	}

	return nil
}

func isOptionalType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice &&
		t.PkgPath() == "github.com/Southclaws/opt" &&
		strings.HasPrefix(t.Name(), "Optional[")
}

func shouldMergeValue(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map:
		return true
	case reflect.Struct:
		return t.PkgPath() == "github.com/Southclaws/storyden/app/resources/settings"
	default:
		return false
	}
}
