package main

import (
	"context"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

func walk(ctx context.Context, c *gowebdav.Client, path string) (<-chan string, <-chan error) {
	pathChan := make(chan string, 1)
	errChan := make(chan error, 1)

	var dfs func(string) error
	dfs = func(path string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		files, err := c.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			if err := ctx.Err(); err != nil {
				return err
			}

			cur := filepath.Join(path, file.Name())
			if file.IsDir() {
				if err := dfs(cur); err != nil {
					return err
				}
				continue
			}

			pathChan <- cur
		}

		return nil
	}

	go func() {
		errChan <- dfs(path)
		close(pathChan)
	}()

	return pathChan, errChan
}
