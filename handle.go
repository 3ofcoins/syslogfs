package main

import "bufio"
import "io"
import "log"
import "strings"

import "bazil.org/fuse"
import "bazil.org/fuse/fs"

type Handle struct { 
	*File	
	Prefix string
	*io.PipeWriter
}

func (h *Handle) doTheLogging(brd *bufio.Reader) {
	for {
		ln, err := brd.ReadString('\n')

		if ln = strings.TrimSpace(ln) ; ln != "" {
			log.Println(h.Prefix, ln)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(h.Prefix, "ERROR reading line:", err)
			break
		}
	}
}

func NewHandle(f *File, prefix string) *Handle {
	rd, wr := io.Pipe()
	h := &Handle{f, prefix, wr}
	go h.doTheLogging(bufio.NewReader(rd))
	return h
}

func (h *Handle) Write(req *fuse.WriteRequest, resp *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	n, err := h.PipeWriter.Write(req.Data)
	resp.Size = n
	h.File.Written(n)
	if err == nil {
		return nil
	} else {
		log.Println(h.Prefix, "ERROR writing chunk:", err)
		return fuse.EIO
	}
}

func (h *Handle) Read(*fuse.ReadRequest, *fuse.ReadResponse, fs.Intr) fuse.Error {
	return fuse.EPERM
}

func (h *Handle) Release(*fuse.ReleaseRequest, fs.Intr) fuse.Error {
	if err := h.PipeWriter.Close() ; err != nil {
		log.Println(h.Prefix, "ERROR closing:", err)
		return fuse.EIO
	}
	return nil
}
