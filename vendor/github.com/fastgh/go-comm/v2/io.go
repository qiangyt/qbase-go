package comm

import (
	"bufio"
	"io"
	"os"

	"github.com/pkg/errors"
)

func DefaultOutput() io.Writer {
	if IsTerminal() {
		return os.Stdout
	} else {
		return io.Discard
	}
}

// ReadBytesP ...
func ReadBytesP(reader io.Reader) []byte {
	r, err := ReadBytes(reader)
	if err != nil {
		panic(err)
	}
	return r
}

// ReadBytes ...
func ReadBytes(reader io.Reader) ([]byte, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "read from Reader")
	}
	return r, nil
}

// ReadText ...
func ReadTextP(reader io.Reader) string {
	r, err := ReadText(reader)
	if err != nil {
		panic(err)
	}
	return r
}

// ReadText ...
func ReadText(reader io.Reader) (string, error) {
	byts, err := io.ReadAll(reader)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return string(byts), nil
}

func ReadLines(reader io.Reader) []string {
	r := make([]string, 0, 32)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		r = append(r, line)
	}

	return r
}
