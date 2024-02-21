package wpsev

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

func SendFile(w http.ResponseWriter, filePath string) error {
	dir, name := filepath.Split(filePath)

	file, err := http.Dir(dir).Open(name)
	if err != nil {
		return err
	}

	defer file.Close()

	d, err := file.Stat()
	if err != nil {
		return err
	}

	if d.Size() < 0 {
		return errors.New("negative content size computed")
	}

	var sendContent io.Reader = file

	if w.Header().Get("Content-Encoding") == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(d.Size(), 10))
	}
	w.Header().Set("Accept-Ranges", "bytes")

	_, err = io.CopyN(w, sendContent, d.Size())
	return err
}
