package comm

import (
	"net/url"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type AferoFileT struct {
	afs         afero.Fs
	rawPath     string
	url         *url.URL
	name        string
	credentials Credentials
	timeout     time.Duration
}

type AferoFile = *AferoFileT

func NewAferoFileP(afs afero.Fs, apath string, credentials Credentials, timeout time.Duration) AferoFile {
	r, err := NewAferoFile(afs, apath, credentials, timeout)
	if err != nil {
		panic(err)
	}
	return r
}

func NewAferoFile(afs afero.Fs, apath string, credentials Credentials, timeout time.Duration) (AferoFile, error) {
	var rawPath, rawUrl string
	if IsFileProtocol(apath) {
		rawPath = apath[len(FILE):]
		rawUrl = apath
	} else {
		rawPath = apath
		rawUrl = FILE + apath
	}

	_url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "parse url: %s", rawUrl)
	}

	return &AferoFileT{
		afs:         afs,
		url:         _url,
		name:        filepath.Base(rawPath),
		rawPath:     rawPath,
		credentials: credentials,
		timeout:     timeout,
	}, nil
}

func (me AferoFile) Fs() afero.Fs {
	return me.afs
}

func (me AferoFile) Name() string {
	return me.name
}

// filepath
func (me AferoFile) Dir() string {
	return filepath.Dir(me.rawPath)
}

// remote file/dir url
func (me AferoFile) Url() string {
	return me.url.String()
}

// remote protocol. May be used to explicitly tell what protocol to use (i.e. "http", "ftp", "etc").
func (me AferoFile) Protocol() string {
	return me.url.Scheme
}

func (me AferoFile) URL() *url.URL {
	return me.url
}

func (me AferoFile) Credentials() Credentials {
	return me.credentials
}

func (me AferoFile) Timeout() time.Duration {
	return me.timeout
}

func (me AferoFile) DownloadP() Content {
	return &ContentT{
		Name: me.Name(),
		Path: me.rawPath,
		Blob: NewAferoBlob(me.afs, me.rawPath),
	}
}

func (me AferoFile) Download() (Content, error) {
	return me.DownloadP(), nil
}

type AferoBlobT struct {
	path string
	afs  afero.Fs
	file afero.File
}

type AferoBlob = *AferoBlobT

func NewAferoBlob(afs afero.Fs, path string) AferoBlob {
	return &AferoBlobT{afs: afs, path: path}
}

func (me AferoBlob) Path() string {
	return me.path
}

func (me AferoBlob) Fs() afero.Fs {
	return me.afs
}

func (me AferoBlob) Read(p []byte) (n int, err error) {
	if me.file == nil {
		f, err := me.afs.Open(me.path)
		me.file = f
		if err != nil {
			return 0, err
		}
	}

	return me.file.Read(p)
}

func (me AferoBlob) Close() error {
	if me.file != nil {
		err := me.file.Close()
		if err != nil {
			return err
		}
		err = me.afs.Remove(me.path)
		if err != nil {
			return err
		}
		me.file = nil
	}
	return nil
}
