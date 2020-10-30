package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/studio-b12/gowebdav"
)

func download(ctx context.Context, c *gowebdav.Client, path string) error {
	_, err := c.Stat(path)
	if err != nil {
		return err
	}

	src, err := c.ReadStream(path)
	if err != nil {
		return err
	}
	defer src.Close()

	output := filepath.Join(*out, strings.TrimLeft(path, *id))
	if err := os.MkdirAll(strings.TrimRight(output, filepath.Base(output)), 0777); err != nil {
		return err
	}
	dst, err := os.Create(output)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
