package md5_cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Cache interface {
	Set(string, []byte) (string, error)
	Get(string) ([]byte, error)
	GetName(string) string
	Exists(string) bool
}

type cache struct {
	cacheDir string
}

func convert(key string) string {
	sum := md5.Sum([]byte(key))
	fname := hex.EncodeToString(sum[:])
	first_subdir := fname[:2]
	second_subdir := fname[2:4]
	return fmt.Sprintf("%s/%s/%s", first_subdir, second_subdir, fname)
}

func (c cache) path(name string) string {
	return path.Join(c.cacheDir, convert(name))
}

func New(cacheDir string) (Cache, error) {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return &cache{cacheDir: cacheDir}, nil
}

func (c cache) Exists(key string) bool {
	fpath := c.path(key)
	fmt.Println(fpath)
	if _, err := os.Stat(fpath); err == nil {
		return true
	}
	return false
}

func (c cache) Set(key string, data []byte) (string, error) {
	fpath := c.path(key)
	if _, err := os.Stat(path.Dir(fpath)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}
	}
	return fpath, ioutil.WriteFile(fpath, data, 0644)
}

func (c cache) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(c.path(key))
}

func (c cache) GetName(key string) string {
	return c.path(key)
}
