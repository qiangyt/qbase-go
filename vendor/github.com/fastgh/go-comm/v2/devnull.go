// copy of github.com/go-task/task/v3/internal/execext/devnull.go
package comm

import (
	"io"
)

var _ io.ReadWriteCloser = devNull{}

type devNull struct{}

func (devNull) Read(p []byte) (int, error)  { return 0, io.EOF }
func (devNull) Write(p []byte) (int, error) { return len(p), nil }
func (devNull) Close() error                { return nil }
