package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext     string
	size    int64
	list    bool
	del     bool
	wLog    io.Writer
	archive string
}

func main() {
	root := flag.String("root", ".", "the root directory to start")
	logFile := flag.String("log", "deleted_files.log", "log deletes to this file")
	ext := flag.String("ext", ".txt", "the extension of files to filter out")
	list := flag.Bool("list", false, "list only files")
	size := flag.Int64("size", 0, "minimum file size to filter out")
	del := flag.Bool("del", true, "delete files")
	arch := flag.String("arch", "", "archive directory path")
	flag.Parse()

	var (
		f   = os.Stdout
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext:     *ext,
		list:    *list,
		size:    *size,
		del:     *del,
		wLog:    f,
		archive: *arch,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)

	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}
			// If list was explicitly set, don't do anything else
			if cfg.list {
				return listFile(path, out)
			}
			// Archive files
			if cfg.archive != "" {
				if err := archiveFiles(cfg.archive, root, path); err != nil {
					return err
				}
			}
			// Delete files
			if cfg.del {
				return delFile(path, delLogger)
			}

			// List is the default option if nothing else was set
			return listFile(path, out)
		})
}
