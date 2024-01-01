package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func rotate(dir *string, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	srcFile, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer srcFile.Close()

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Println(err)
		return
	}

	name := fileInfo.Name()
	time := fileInfo.ModTime()

	if *dir != "" {
		if (*dir)[len(*dir)] == '/' {
			*dir += "/"
		}
		name = filepath.Base(name)
		idx := strings.LastIndex(name, ".")
		if idx > 0 {
			name = name[:idx]
		}
		name = *dir + name + "_" + strconv.FormatInt(time.Unix(), 10) + "tag.gz"
	} else {
		idx := strings.LastIndex(name, ".")
		if idx > 0 {
			name = name[:idx]
		}
		name = name + "_" + strconv.FormatInt(time.Unix(), 10) + "tag.gz"
	}
	log.Println(name)
	dst, err := os.Create(name)
	if err != nil {
		log.Println(err)
		return
	}
	defer dst.Close()

	gzWriter := gzip.NewWriter(dst)
	if gzWriter == nil {
		log.Println("Error due creating writer")
		return
	}
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var wg sync.WaitGroup
	dir := flag.String("a", "", "specify an output dir")
	flag.Parse()
	fileNames := flag.Args()
	for _, file := range fileNames {
		wg.Add(1)
		go rotate(dir, file, &wg)
	}
	wg.Wait()
}
