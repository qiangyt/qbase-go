package comm

import (
	"net/url"
	"time"

	"github.com/goodsru/go-universal-network-adapter/models"
	"github.com/goodsru/go-universal-network-adapter/services"
	"github.com/pkg/errors"
)

var _uniNwAdapter *services.UniversalNetworkAdapter

type RemoteFileT struct {
	backend *models.RemoteFile
}

type RemoteFile = *RemoteFileT

func init() {
	_uniNwAdapter = services.NewUniversalNetworkAdapter()
}

func NewRemoteFileP(url string, credentials Credentials, timeout time.Duration) RemoteFile {
	r, err := NewRemoteFile(url, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func NewRemoteFile(url string, credentials Credentials, timeout time.Duration) (RemoteFile, error) {
	remoteFile, err := models.NewRemoteFile(models.NewDestination(url, credentials, &timeout))
	if err != nil {
		return nil, errors.Wrapf(err, "new remote file object")
	}

	return &RemoteFileT{remoteFile}, nil
}

func (me RemoteFile) Name() string {
	return me.backend.Name
}

// filepath
func (me RemoteFile) Dir() string {
	r := me.backend.Path

	lastPos := len(r) - 1
	if r[lastPos] == '/' {
		return r[:lastPos]
	}

	return r
}

// remote file/dir url
func (me RemoteFile) Url() string {
	return me.backend.ParsedDestination.Url
}

// remote protocol. May be used to explicitly tell what protocol to use (i.e. "http", "ftp", "etc").
func (me RemoteFile) Protocol() string {
	return me.backend.ParsedDestination.ParsedUrl.Scheme
}

func (me RemoteFile) URL() *url.URL {
	return me.backend.ParsedDestination.ParsedUrl
}

func (me RemoteFile) Credentials() Credentials {
	return &me.backend.ParsedDestination.Credentials
}

func (me RemoteFile) Timeout() time.Duration {
	return me.backend.ParsedDestination.Timeout
}

func (me RemoteFile) DownloadP() Content {
	r, err := me.Download()
	if err != nil {
		panic(err)
	}
	return r
}

func (me RemoteFile) Download() (Content, error) {
	r, err := _uniNwAdapter.Download(me.backend)
	if err != nil {
		return nil, errors.Wrapf(err, "download %s", me.Url())
	}

	return r, nil
}
