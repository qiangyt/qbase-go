package comm

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/goodsru/go-universal-network-adapter/models"
	"github.com/goodsru/go-universal-network-adapter/services"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	HTTP  = services.HTTP + "://"
	HTTPS = services.HTTPS + "://"
	FTP   = services.FTP + "://"
	FTPS  = services.FTPS + "://"
	SFTP  = services.SFTP + "://"
	S3    = services.S3 + "://"
	FILE  = "file://"
)

type File interface {
	Name() string
	Dir() string
	Url() string
	Protocol() string
	URL() *url.URL
	Credentials() Credentials
	Timeout() time.Duration
	DownloadP() Content
	Download() (Content, error)
}

type (
	CredentialsT = models.Credentials
	Credentials  = *CredentialsT
)

type (
	ContentT = models.RemoteFileContent
	Content  = *ContentT
)

func NewFileP(afs afero.Fs, url string, credentials Credentials, timeout time.Duration) File {
	r, err := NewFile(afs, url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func NewFile(afs afero.Fs, url string, credentials Credentials, timeout time.Duration) (File, error) {
	if IsRemote(url) {
		return NewRemoteFile(url, credentials, timeout)
	}
	return NewAferoFile(afs, url, credentials, timeout)
}

func IsFileProtocol(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), FILE)
}

func IsRemote(url string) bool {
	lc := strings.ToLower(url)

	if strings.HasPrefix(lc, HTTP) ||
		strings.HasPrefix(lc, HTTPS) ||
		strings.HasPrefix(lc, FTP) ||
		strings.HasPrefix(lc, FTPS) ||
		strings.HasPrefix(lc, SFTP) ||
		strings.HasPrefix(lc, S3) {
		return true
	}
	return false
}

func WorkDir(url string, defaultDir string) string {
	if IsRemote(url) {
		return defaultDir
	}
	if IsFileProtocol(url) {
		url = url[len(FILE):]
	}

	r := filepath.Dir(url)
	if r == "." {
		return defaultDir
	}
	if !filepath.IsAbs(url) {
		r = filepath.Join(defaultDir, r)
	}
	return r
}

func ShortDescription(url string) string {
	r := url

	protocol := ""
	posOfProtocolSep := strings.Index(strings.ToLower(r), "://")
	if posOfProtocolSep >= 0 {
		protocol = r[:posOfProtocolSep+3]
		r = r[posOfProtocolSep+len("://"):]
	}

	lem := len(r)
	if lem <= 3+8+1+5 /* ...12345678.hosts */ {
		return protocol + r
	}

	return protocol + r[:3] + "..." + r[lem-(8+1+5): /* 12345678.hosts */]
}

func DownloadBytesP(logger Logger, fallbackDir string, fs afero.Fs, url string, credentials Credentials, timeout time.Duration) []byte {
	r, err := DownloadBytes(logger, fallbackDir, fs, url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func DownloadBytes(logger Logger, fallbackDir string, fs afero.Fs, url string, credentials Credentials, timeout time.Duration) (result []byte, err error) {
	result, err = downloadBytes(fs, url, credentials, timeout)

	if len(fallbackDir) > 0 {
		if err == nil {
			if logger != nil {
				logger.Info().Str("fallbackDir", fallbackDir).Str("url", url).Msg("save download files to fallback dir")
			}
			fallbackFilePath, fallbackErr := WriteFallbackFile(fallbackDir, fs, url, result)
			if fallbackErr != nil {
				if logger != nil {
					logger.Warn().Err(fallbackErr).Str("url", url).Str("fallbackFilePath", fallbackFilePath).Msg("save fallback file failed")
				}
			}
			return
		} else {
			if logger != nil {
				logger.Warn().Err(err).Str("url", url).Msg("fallbacking due to failed to download the file")
			}
			fallbackFilePath, r, fallbackErr := ReadFallbackFile(fallbackDir, fs, url)
			if fallbackErr != nil {
				if logger != nil {
					logger.Warn().Err(fallbackErr).Str("url", url).Str("fallbackFilePath", fallbackFilePath).Msg("get fallbacked file failed too")
				}
			} else {
				result = r
				err = nil
			}
		}
	}

	return
}

func downloadBytes(fs afero.Fs, url string, credentials Credentials, timeout time.Duration) ([]byte, error) {
	f, err := NewFile(fs, url, credentials, timeout)
	if err != nil {
		return nil, err
	}

	c, err := f.Download()
	if err != nil {
		return nil, err
	}

	blob := c.Blob
	defer blob.Close()

	return ReadBytes(blob)
}

func DownloadTextP(logger Logger, fallbackDir string, fs afero.Fs, url string, credentials Credentials, timeout time.Duration) string {
	r, err := DownloadText(logger, fallbackDir, fs, url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func DownloadText(logger Logger, fallbackDir string, fs afero.Fs, url string, credentials Credentials, timeout time.Duration) (string, error) {
	bytes, err := DownloadBytes(logger, fallbackDir, fs, url, credentials, timeout)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MapFromYamlFileP(fs afero.Fs, path string, envsubt bool) map[string]any {
	r, err := MapFromYamlFile(fs, path, envsubt)
	if err != nil {
		panic(err)
	}
	return r
}

func MapFromYamlFile(fs afero.Fs, path string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromYamlFile(fs, path, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func MapFromYamlP(yamlText string, envsubt bool) map[string]any {
	r, err := MapFromYaml(yamlText, envsubt)
	if err != nil {
		panic(err)
	}
	return r
}

func MapFromYaml(yamlText string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromYaml(yamlText, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func MapFromJsonFileP(fs afero.Fs, path string, envsubt bool) map[string]any {
	r, err := MapFromJsonFile(fs, path, envsubt)
	if err != nil {
		panic(err)
	}
	return r
}

func MapFromJsonFile(fs afero.Fs, path string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromJsonFile(fs, path, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func MapFromJsonP(yamlText string, envsubt bool) map[string]any {
	r, err := MapFromJson(yamlText, envsubt)
	if err != nil {
		panic(err)
	}
	return r
}

func MapFromJson(yamlText string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromJson(yamlText, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func FromYamlFileP(fs afero.Fs, path string, envsubt bool, result any) {
	if err := FromYamlFile(fs, path, envsubt, result); err != nil {
		panic(err)
	}
}

func FromYamlFile(fs afero.Fs, path string, envsubt bool, result any) error {
	yamlText, err := ReadFileText(fs, path)
	if err != nil {
		return err
	}

	if err := FromYaml(yamlText, envsubt, result); err != nil {
		return errors.Wrapf(err, "parse yaml file: %s", path)
	}
	return nil
}

func FromYamlP(yamlText string, envsubt bool, result any) {
	if err := FromYaml(yamlText, envsubt, result); err != nil {
		panic(err)
	}
}

func FromYaml(yamlText string, envsubt bool, result any) (err error) {
	if envsubt {
		yamlText, err = EnvSubst(yamlText, nil)
		if err != nil {
			return
		}
	}

	if err = yaml.Unmarshal([]byte(yamlText), result); err != nil {
		return errors.Wrapf(err, "parse yaml: \n\n%s", yamlText)
	}
	return nil
}

func FromJsonFileP(fs afero.Fs, path string, envsubt bool, result any) {
	if err := FromJsonFile(fs, path, envsubt, result); err != nil {
		panic(err)
	}
}

func FromJsonFile(fs afero.Fs, path string, envsubt bool, result any) error {
	yamlText, err := ReadFileText(fs, path)
	if err != nil {
		return err
	}

	if err := FromJson(yamlText, envsubt, result); err != nil {
		return errors.Wrapf(err, "parse json file: %s", path)
	}
	return nil
}

func FromJsonP(jsonText string, envsubt bool, result any) {
	if err := FromJson(jsonText, envsubt, result); err != nil {
		panic(err)
	}
}

func FromJson(jsonText string, envsubt bool, result any) (err error) {
	if envsubt {
		jsonText, err = EnvSubst(jsonText, nil)
		if err != nil {
			return
		}
	}

	if err = json.Unmarshal([]byte(jsonText), result); err != nil {
		return errors.Wrapf(err, "parse json: \n\n%s", jsonText)
	}
	return nil
}
