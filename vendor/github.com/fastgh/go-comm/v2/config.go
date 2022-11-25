package comm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ConfigMetadata = mapstructure.Metadata

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

	Metadata ConfigMetadata

	DoValidate bool
}

func StrictConfigConfig() *ConfigConfig {
	return &ConfigConfig{
		ErrorUnused:          true,
		ErrorUnset:           true,
		ZeroFields:           true,
		WeaklyTypedInput:     false,
		Squash:               false,
		IgnoreUntaggedFields: true,
		Metadata:             ConfigMetadata{},
		DoValidate:           false,
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
		Metadata:             ConfigMetadata{},
		DoValidate:           false,
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
		Metadata:             &me.Metadata,
		Result:               nil,
		TagName:              "",
		IgnoreUntaggedFields: me.IgnoreUntaggedFields,
		MatchName:            nil,
	}
}

func DecodeWithYamlP[T any](yamlText string, cfgcfg *ConfigConfig, result *T, devault map[string]any) (*T, *ConfigMetadata) {
	r, m, err := DecodeWithYaml(yamlText, cfgcfg, result, devault)
	if err != nil {
		panic(err)
	}
	return r, m
}

func DecodeWithYaml[T any](yamlText string, cfgcfg *ConfigConfig, result *T, devault map[string]any) (*T, *ConfigMetadata, error) {
	input, err := MapFromYaml(yamlText, false)
	if err != nil {
		return nil, nil, err
	}

	return DecodeWithMap(input, cfgcfg, result, devault)
}

func DecodeWithMapP[T any](input map[string]any, cfgcfg *ConfigConfig, result *T, devault map[string]any) (*T, *ConfigMetadata) {
	r, m, err := DecodeWithMap(input, cfgcfg, result, devault)
	if err != nil {
		panic(err)
	}
	return r, m
}

func DecodeWithMap[T any](input map[string]any, cfgcfg *ConfigConfig, result *T, devault map[string]any) (*T, *ConfigMetadata, error) {
	backend := MergeMap(devault, input)

	ms := cfgcfg.ToMapstruct()
	ms.Result = result

	decoder, err := mapstructure.NewDecoder(ms)
	// TODO: for better user-friendly error message, use DecoderConfig{Metadata} to find Unset
	if err != nil {
		return nil, &cfgcfg.Metadata, errors.Wrap(err, "create mapstructure decoder")
	}

	if err = decoder.Decode(backend); err != nil {
		return nil, &cfgcfg.Metadata, errors.Wrapf(err, "decode map: %v", backend)
	}

	if cfgcfg.DoValidate {
		if err = validate.Struct(result); err != nil {
			return nil, &cfgcfg.Metadata, errors.Wrapf(err, "validation: %v", backend)
		}
	}

	return result, &cfgcfg.Metadata, nil
}

func GetMapValue[T any](m map[string]any, key string, devault func() T) T {
	if i, has := m[key]; has {
		return i.(T)
	}

	r := devault()
	m[key] = r

	return r
}

func LoadEnvScripts(fs afero.Fs, vars map[string]string, filenames ...string) (map[string]string, error) {
	errs := NewErrorGroup(false)

	if len(filenames) == 0 {
		filenames = SysEnvFileNames(fs, "")
	}

	for _, filename := range filenames {
		var err error
		vars, err = LoadEnvScript(fs, vars, filename)
		errs.Add(err)
	}

	return vars, errs.MayError()
}

func LoadEnvScript(fs afero.Fs, vars map[string]string, filename string) (map[string]string, error) {
	if filename == "/etc/paths" {
		paths, err := ReadFileLines(fs, filename)
		if err == nil {
			if len(vars["PATH"]) >= 0 {
				paths = append([]string{vars["PATH"]}, paths...)
			}
			vars["PATH"] = strings.Join(paths, ":")
		}
		return vars, err
	}

	output, err := RunGoshCommand(vars, "", filename, nil)
	if err != nil {
		return vars, err
	}
	return output.Vars, nil
}

func SysEnvFileNames(fs afero.Fs, shell string) []string {
	r := []string{}

	if len(shell) == 0 {
		shell = os.Getenv("SHELL")
	}

	home, _ := ExpandHomePath("~")
	hasHome := (len(home) > 0)

	pth := filepath.Join("/etc/profile")
	if exists, _ := FileExists(fs, pth); exists {
		r = append(r, pth)
	}

	pth = filepath.Join("/etc/paths")
	if exists, _ := FileExists(fs, pth); exists {
		r = append(r, pth)
	}

	if !strings.Contains(shell, "zsh") {
		if exists, _ := FileExists(fs, "/etc/bashrc"); exists {
			r = append(r, pth)
		}

		if hasHome {
			pth = filepath.Join(home, ".bashrc")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			}

			pth = filepath.Join(home, ".bash_profile")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			} else {
				pth = filepath.Join(home, ".bash_login")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
				pth = filepath.Join(home, ".profile")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
			}
		}
	} else {
		if exists, _ := FileExists(fs, "/etc/zshrc"); exists {
			r = append(r, pth)
		}

		if hasHome {
			pth = filepath.Join(home, ".zshrc")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			}

			pth = filepath.Join(home, ".zshenv")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			}

			pth = filepath.Join(home, ".zprofile")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			} else {
				pth = filepath.Join(home, ".zsh_login")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
				pth = filepath.Join(home, ".profile")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
			}
		}
	}

	pth = ".env"
	if exists, _ := FileExists(fs, pth); exists {
		r = append(r, pth)
	}

	return r
}
