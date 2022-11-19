package comm

// a revised copy of github.com/allan-simon/go-singleinstance v0.0.0-20210120080615-d0997106ab37

import (
	"strconv"

	"github.com/spf13/afero"
)

// If filename is a lock file, returns the PID of the process locking it
func GetLockFilePid(fs afero.Fs, filename string) (pid int, err error) {
	contents, err := afero.ReadFile(fs, filename)
	if err != nil {
		return
	}

	pid, err = strconv.Atoi(string(contents))
	return
}
