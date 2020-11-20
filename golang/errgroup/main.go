package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var urls = []string{
	"http://www.baidu.com",
	"http://www.sougou.com",
	"http://www.fjsdklf23.com",
}

func main() {
	var g errgroup.Group

	for _, url := range urls {
		g.Go(func() error {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}

	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs")
	} else {
		fmt.Println(err)
	}
}

func parallelEx() {
	Google := func(ctx context.Context, query string) (ress []interface{}, err error) {
		g, ctx := errgroup.WithContext(ctx)

		for _, url := range urls {
			url := url
			g.Go(func() error {
				res, err := http.Get(url)
				if err == nil {
					ress = append(ress, res.Body)
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return ress, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

type result struct {
	path string
	sum [md5.Size]byte
}

func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	g, ctx := errgroup.WithContext(ctx)
	paths := make(chan string)

	g.Go(func() error {
		defer close(paths)
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	c := make(chan result)
	const numDigesters = 20
	for i := 0; i < numDigesters; i++ {
		g.Go(func() error {
			for path := range paths {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				select {
				case c <- result{path, md5.Sum(data)}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	go func() {
		g.Wait()
		close(c)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c {
		m[r.path] = r.sum
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return m, nil
}
