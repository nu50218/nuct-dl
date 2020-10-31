package main

import (
	"context"
	"flag"
	"log"
	"math"
	"runtime"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/studio-b12/gowebdav"
)

var (
	uri        = flag.String("uri", "https://ct.nagoya-u.ac.jp/dav/", "base uri")
	user       = flag.String("user", "", "username")
	pass       = flag.String("pass", "", "password")
	id         = flag.String("id", "", "siteid")
	out        = flag.String("out", "", "output directory")
	lastUpdate = flag.Duration("last_update", math.MaxInt64, "last_update")
)

func init() {
	flag.Parse()
	if *id == "" {
		log.Fatal("id is empty")
	}
	if *user == "" || *pass == "" {
		log.Fatal("user or pass is empty")
	}
	if *out == "" {
		out = id
	}
}

func main() {
	ctx := context.Background()

	client := gowebdav.NewClient(*uri, *user, *pass)
	if err := client.Connect(); err != nil {
		log.Fatalf("connection error: %v", err)
	}

	pathChan, errChan := walk(ctx, client, *id)
	limitChan := make(chan struct{}, runtime.NumCPU())
	bar := pb.StartNew(0).SetMaxWidth(80).SetTemplateString(`{{counters . }} {{ bar . "[" "=" ">" "_" "]"}} {{percent . }}`)
	defer bar.Finish()
	var (
		wg  sync.WaitGroup
		cnt int64
	)

	for path := range pathChan {
		path := path

		select {
		case <-ctx.Done():
			log.Fatal(ctx.Err())

		default:
			cnt++
			bar.SetTotal(cnt)

			wg.Add(1)
			go func() {
				defer wg.Done()

				limitChan <- struct{}{}
				if err := download(ctx, client, path); err != nil {
					log.Fatal(err)
				}
				bar.Increment()
				<-limitChan
			}()
		}
	}

	wg.Wait()

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
