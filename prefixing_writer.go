package main

import "bufio"
import "fmt"
import "io"
import "log"

type PrefixingWriter struct {
	io.Writer
	done <-chan error
}

func NewPrefixingWriter(sink io.Writer, prefix string) *PrefixingWriter {
	rd, wr := io.Pipe()
	done := make(chan error, 1)

	go func() {
		// At the end, close the sink if closeable and signal we're done
		defer func() {
			if closer, ok := sink.(io.Closer) ; ok {
				done <- closer.Close()
			} else {
				done <- nil
			}
		}()

		scanner := bufio.NewScanner(rd)
		defer func() {
			if err := scanner.Err() ; err != nil {
				log.Println("ERROR reading line", err)
			}
		}()

		for scanner.Scan() {
			fmt.Fprintln(sink, prefix, scanner.Text())
		}
	}()

	return &PrefixingWriter{wr, done}
}

func (pw *PrefixingWriter) Close() error {
	pw.Writer.(io.Closer).Close()
	return <- pw.done
}
