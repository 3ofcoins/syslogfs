package main

import "encoding/json"
import "io"
import "io/ioutil"
import "os"

import "bazil.org/fuse/fs"

type FS struct {
	fs.Tree
	io.Writer
}

func NewFS(writer io.Writer) *FS {
	fs := new(FS)
	fs.Writer = writer
	return fs
}

func (fs *FS) NewDir(name string) *Dir {
	dir := &Dir{fs, name, make(map[string]*File)}
	fs.Add(name, dir)
	return dir
}

func (fs *FS) LoadConfig(path string) error {
	var config map[string]interface{}

	cff, err := os.Open(path)
	if err != nil {
		return err
	}

	cfb, err := ioutil.ReadAll(cff)
	if err != nil {
		return err
	}

	json.Unmarshal(cfb, &config)
	for dir, _ := range config {
		fs.NewDir(dir)
	}

	return nil
}
