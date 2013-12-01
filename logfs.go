package main

import "encoding/json"

import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "os"

import "bazil.org/fuse"
import "bazil.org/fuse/fs"


func loadConfig(path string) {
	var config map[string]interface{}

	cff, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open config", path, ":", err)
	}

	cfb, err := ioutil.ReadAll(cff)
	if err != nil {
		log.Fatal("Cannot read config", path, ":", err)
	}

	json.Unmarshal(cfb, &config)
	for dir, _ := range config {
		NewDir(dir)
	}
}

var usage = func () {
	fmt.Fprintf(os.Stderr, "Usage: %s TEMPLATE.JSON MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	_ = flag.String("o", "", "mount options (ignored)")
	debug := flag.Bool("debug", false, "show debugging info")
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	loadConfig(flag.Arg(0))

	if *debug {
		fuse.Debugf = log.Printf
	}

	if c, err := fuse.Mount(flag.Arg(1)) ; err != nil {
		log.Fatal(err)
	} else {
		fs.Serve(c, Root)
	}
}
