package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/studio-b12/gowebdav"
)

func download(ctx context.Context, c *gowebdav.Client, path string) error {
	fileInfo, err := c.Stat(path)
	if err != nil {
		return err
	}
	if time.Since(fileInfo.ModTime()) > *lastUpdate {
		return nil
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
