package main

import "io"
import "log"

import "bazil.org/fuse"
import "bazil.org/fuse/fs"

type Handle struct {
	*File
	sink io.WriteCloser
}

func NewHandle(file *File, prefix string) *Handle {
	return &Handle{file, NewPrefixingWriter(file.Dir.FS.Writer, prefix)}
}

func (h *Handle) Write(req *fuse.WriteRequest, resp *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	n, err := h.sink.Write(req.Data)
	resp.Size = n
	h.File.Written(n)
	if err == nil {
		return nil
	} else {
		log.Println("ERROR writing chunk:", err)
		return fuse.EIO
	}
}

func (h *Handle) Read(*fuse.ReadRequest, *fuse.ReadResponse, fs.Intr) fuse.Error {
	return fuse.EPERM
}

func (h *Handle) Release(*fuse.ReleaseRequest, fs.Intr) fuse.Error {
	if err := h.sink.Close() ; err != nil {
		log.Println("ERROR closing:", err)
		return fuse.EIO
	}
	return nil
}
