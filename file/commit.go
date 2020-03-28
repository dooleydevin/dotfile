package file

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// Commit represents a revision for a tracked file.
type Commit struct {
	Hash       string        `json:"hash"`
	Message    string        `json:"message"`
	Timestamp  int64         `json:"timestamp"`
	Compressed *bytes.Buffer `json:"-"`
}

func newCommit(contents []byte, message string) (*Commit, error) {
	compressed := new(bytes.Buffer)
	w := zlib.NewWriter(compressed)
	defer w.Close()

	if _, err := w.Write(contents); err != nil {
		return nil, errors.Wrap(err, "compressing file for commit")
	}

	return &Commit{
		Hash:       fmt.Sprintf("%x", sha1.Sum(contents)),
		Message:    message,
		Timestamp:  time.Now().Unix(),
		Compressed: compressed,
	}, nil
}

func uncompress(compressed []byte) (*bytes.Buffer, error) {
	uncompressed := new(bytes.Buffer)

	r, err := zlib.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, errors.Wrap(err, "uncompressing commit revision")
	}
	defer r.Close()

	if _, err = io.Copy(uncompressed, r); err != nil {
		return nil, errors.Wrap(err, "copying uncompressed data")
	}

	return uncompressed, nil
}