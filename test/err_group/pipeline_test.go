package err_group

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io/fs"
	"path/filepath"

	"github/justWindy/go_demo/utils"
	"golang.org/x/sync/errgroup"
)

type result struct {
	path string
	sum  [md5.Size]byte
}

func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	g, ctx := errgroup.WithContext(ctx)

	paths := make(chan string)

	g.Go(func() error {
		defer close(paths)
		return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			paths <- path + info.Name()
			return err
		})
	})

	c := make(chan result)
	const numDigesters = 20
	for i := 0; i < numDigesters; i++ {
		g.Go(func() error {
			for path := range paths {
				c <- result{
					path: path,
					sum:  calculateMD5(path),
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

func calculateMD5(input string) [md5.Size]byte {
	data := utils.Str2byte(input)
	hash := md5.Sum(data)
	dst := make([]byte, hex.EncodedLen(len(hash[:])))
	hex.Encode(dst, hash[:])
	return ([md5.Size]byte)(dst)
}
