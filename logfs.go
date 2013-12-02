package main

import "flag"
import "fmt"
import "log"
import "os"
import "strings"

import "bazil.org/fuse"
import "bazil.org/fuse/fs"

var usage = func () {
	fmt.Fprintf(os.Stderr, "Usage: %s TEMPLATE.JSON MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

type logWriter struct{}
func (logWriter) Write(data []byte) (int, error) {
	log.Print(string(data))
	return len(data), nil
}

func main() {
	flag.Usage = usage
	raw_options := flag.String("o", "", "mount options (ignored)")
	debug := flag.Bool("debug", false, "show debugging info")
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	options := make(map[string]string)
	for _, opt := range strings.Split(*raw_options, ",") {
		kv := strings.SplitN(opt, "=", 2)
		if len(kv) == 1 {
			options[kv[0]] = ""
		} else {
			options[kv[0]] = kv[1]
		}
	}
	log.Println(options)

	root := NewFS(logWriter{})
	if err := root.LoadConfig(flag.Arg(0)) ; err != nil {
		log.Fatal(err)
	}

	if *debug {
		fuse.Debugf = log.Printf
	}

	if c, err := fuse.Mount(flag.Arg(1)) ; err != nil {
		log.Fatal(err)
	} else {
		fs.Serve(c, root)
	}
}
