package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var f1, f2 string
	flag.StringVar(&f1, "f1", "", "FilePath of old file")
	flag.StringVar(&f2, "f2", "", "FilePath of the old file")
	if len(os.Args) < 5 {
		fmt.Println("Specify the filepath: -f1 filePathOld -f2 filePathNew")
		os.Exit(1)
	}
	flag.Parse()

	if len(f1) == 0 || len(f2) == 0 {
		fmt.Println("Specify the filepath: -f1 and -f2 filePaths")
		os.Exit(1)
	}
	filePath := os.Args[2]
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	s := string(file)
	arr := strings.Split(s, "\n")

	fileMap := make(map[string]bool)

	for _, line := range arr {
		if len(line) != 0 {
			fileMap[line] = false
		}
	}

	filePath = os.Args[4]
	newFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(newFile)
	for scanner.Scan() {
		line := scanner.Text()
		_, isContains := fileMap[line]
		if !isContains {
			fmt.Printf("ADDED %s\n", line)
		} else {
			fileMap[line] = true
		}
	}
	newFile.Close()
	for line, isContains := range fileMap {
		if !isContains {
			fmt.Printf("REMOVED %s\n", line)
		}
	}
}
