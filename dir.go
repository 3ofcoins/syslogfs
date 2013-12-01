package main

import "os"
import "syscall"

import "bazil.org/fuse"
import "bazil.org/fuse/fs"

const EEXIST = fuse.Errno(syscall.EEXIST)

var Root = new(fs.Tree)

type Dir struct {
	Name string
	files map[string]*File
}

func NewDir(name string) *Dir {
	dir := &Dir{name, make(map[string]*File)}
	Root.Add(name, dir)
	return dir
}


func (dir *Dir) Attr() fuse.Attr {
	return fuse.Attr{Mode: os.ModeDir}
}

func (dir *Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if file, exists := dir.files[name] ; exists {
		return file, nil
	} else {
		return nil, fuse.ENOENT
	}
}

func (dir *Dir) ReadDir(fs.Intr)  ([]fuse.Dirent, fuse.Error) {
	rv := make([]fuse.Dirent, 0, len(dir.files))
	for name, _ := range(dir.files) {
		rv = append(rv, fuse.Dirent{Type: fuse.DT_File, Name: name})
	}
	return rv, nil
}

func (dir *Dir) Create(req *fuse.CreateRequest, res *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
	if _, exists := dir.files[req.Name] ; exists {
		return nil, nil, EEXIST
	}
	
	file := NewFile(dir.Name, req.Name)
	dir.files[req.Name] = file
	return file, file.Handle(req.Header), nil
}
