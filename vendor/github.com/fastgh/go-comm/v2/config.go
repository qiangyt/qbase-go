package comm

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Derived from mapstructure.DecodeConfig
type ConfigConfig struct {
	// If ErrorUnused is true, then it is an error for there to exist
	// keys in the original map that were unused in the decoding process
	// (extra keys).
	ErrorUnused bool

	// If ErrorUnset is true, then it is an error for there to exist
	// fields in the result that were not set in the decoding process
	// (extra fields). This only applies to decoding to a struct. This
	// will affect all nested structs as well.
	ErrorUnset bool

	// ZeroFields, if set to true, will zero fields before writing them.
	// For example, a map will be emptied before decoded values are put in
	// it. If this is false, a map will be merged.
	ZeroFields bool

	// If WeaklyTypedInput is true, the decoder will make the following
	// "weak" conversions:
	//
	//   - bools to string (true = "1", false = "0")
	//   - numbers to string (base 10)
	//   - bools to int/uint (true = 1, false = 0)
	//   - strings to int/uint (base implied by prefix)
	//   - int to bool (true if value != 0)
	//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
	//     FALSE, false, False. Anything else is an error)
	//   - empty array = empty map and vice versa
	//   - negative numbers to overflowed uint values (base 10)
	//   - slice of maps to a merged map
	//   - single values are converted to slices if required. Each
	//     element is weakly decoded. For example: "4" can become []int{4}
	//     if the target type is an int slice.
	//
	WeaklyTypedInput bool

	// Squash will squash embedded structs.  A squash tag may also be
	// added to an individual struct field using a tag.  For example:
	//
	//  type Parent struct {
	//      Child `mapstructure:",squash"`
	//  }
	Squash bool

	// IgnoreUntaggedFields ignores all struct fields without explicit
	// TagName, comparable to `mapstructure:"-"` as default behaviour.
	IgnoreUntaggedFields bool
}

func StrictConfigConfig() *ConfigConfig {
	return &ConfigConfig{
		ErrorUnused:          true,
		ErrorUnset:           true,
		ZeroFields:           true,
		WeaklyTypedInput:     false,
		Squash:               false,
		IgnoreUntaggedFields: true,
	}
}

func DynamicConfigConfig() *ConfigConfig {
	return &ConfigConfig{
		ErrorUnused:          false,
		ErrorUnset:           false,
		ZeroFields:           false,
		WeaklyTypedInput:     true,
		Squash:               true,
		IgnoreUntaggedFields: true,
	}
}

func (me *ConfigConfig) ToMapstruct() *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		DecodeHook:           nil,
		ErrorUnused:          me.ErrorUnused,
		ErrorUnset:           me.ErrorUnset,
		ZeroFields:           me.ZeroFields,
		WeaklyTypedInput:     me.WeaklyTypedInput,
		Squash:               me.Squash,
		Metadata:             nil,
		Result:               nil,
		TagName:              "",
		IgnoreUntaggedFields: me.IgnoreUntaggedFields,
		MatchName:            nil,
	}
}

func DecodeWithYamlP[T any](yamlText string, cfgcfg *ConfigConfig, result *T, devault map[string]any) *T {
	r, err := DecodeWithYaml(yamlText, cfgcfg, result, devault)
	if err != nil {
		panic(err)
	}
	return r
}

func DecodeWithYaml[T any](yamlText string, cfgcfg *ConfigConfig, result *T, devault map[string]any) (*T, error) {
	input, err := MapFromYaml(yamlText, false)
	if err != nil {
		return nil, err
	}

	return DecodeWithMap(input, cfgcfg, result, devault)
}

func DecodeWithMapP[T any](input map[string]any, cfgcfg *ConfigConfig, result *T, devault map[string]any) *T {
	r, err := DecodeWithMap(input, cfgcfg, result, devault)
	if err != nil {
		panic(err)
	}
	return r
}

func DecodeWithMap[T any](input map[string]any, cfgcfg *ConfigConfig, result *T, devault map[string]any) (*T, error) {
	backend := MergeMap(devault, input)

	ms := cfgcfg.ToMapstruct()
	ms.Result = result

	decoder, err := mapstructure.NewDecoder(ms)
	// TODO: for better user-friendly error message, use DecoderConfig{Metadata} to find Unset
	if err != nil {
		return nil, errors.Wrap(err, "create mapstructure decoder")
	}

	if err = decoder.Decode(backend); err != nil {
		return nil, errors.Wrapf(err, "decode map: %v", backend)
	}

	return result, nil
}
