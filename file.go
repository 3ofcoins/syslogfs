package main

import "fmt"
import "syscall"

import "bazil.org/fuse"
import "bazil.org/fuse/fs"

const ENOTDIR = fuse.Errno(syscall.ENOTDIR)

type File struct {
	*Dir
	Name string
	BytesWritten uint64
	writtenByteChan chan int 
}

func (f *File) doUpdateBytes() {
	for {
		f.BytesWritten += uint64(<- f.writtenByteChan)
	}
}

func NewFile (dir *Dir, filename string) *File {
	f := &File{dir, filename, 0, make(chan int)}
	go f.doUpdateBytes()
	return f
}

func (f *File) Written(bytes int) {
	f.writtenByteChan <- bytes
}

func (f *File) Handle(hdr fuse.Header) *Handle {
	return NewHandle(f, fmt.Sprintf("%s[%d]: %s:", f.Dir.Name, hdr.Pid, f.Name))
}

func (f *File) Attr() fuse.Attr {
	return fuse.Attr{Size: f.BytesWritten}
}

func (f *File) Open (req *fuse.OpenRequest, resp *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	switch {
	case req.Dir:
		return nil,                  ENOTDIR
	case (req.Flags|syscall.O_WRONLY) == 0:
		return nil,                  fuse.EPERM
	default:
		return f.Handle(req.Header), nil
	}
}
