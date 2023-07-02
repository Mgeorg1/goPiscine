package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"unicode/utf8"
)

var wg sync.WaitGroup

func countLines(filePath string, m *sync.Mutex) {
	defer wg.Done()
	var count uint64
	count = 0

	f, err := os.Open(filePath)

	if err != nil {
		m.Lock()
		fmt.Println(err.Error())
		m.Unlock()
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	m.Lock()
	fmt.Printf("%s %d\n", filePath, count)
	m.Unlock()
}

func countWords(filePath string, m *sync.Mutex) {
	defer wg.Done()
	f, err := os.Open(filePath)
	if err != nil {
		m.Lock()
		fmt.Println(err.Error())
		m.Unlock()
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	if err := scanner.Err(); err != nil {
		m.Lock()
		fmt.Println(err.Error())
		m.Unlock()
		return
	}

	var count int64
	count = 0

	for scanner.Scan() {
		count++
	}
	m.Lock()
	fmt.Printf("%s %d\n", filePath, count)
	m.Unlock()
}

func countCharacters(filePath string, m *sync.Mutex) {
	defer wg.Done()
	text, err := os.ReadFile(filePath)

	if err != nil {
		m.Lock()
		fmt.Println(err.Error())
		m.Unlock()
		return
	}
	m.Lock()
	fmt.Printf("%s %d\n", filePath, utf8.RuneCount(text))
	m.Unlock()
}

func main() {
	lFlag := flag.Bool("l", false, "count lines")
	mFlag := flag.Bool("m", false, "count characters")
	wFlag := flag.Bool("w", false, "count words")
	flag.Parse()

	if *lFlag && *mFlag || *lFlag && *wFlag || *mFlag && *wFlag {
		fmt.Println("Only one flag may be specified!")
		os.Exit(1)
	}

	if !*lFlag && !*mFlag && !*wFlag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No input files")
		os.Exit(1)
	}
	m := &sync.Mutex{}
	for _, path := range args {
		if path != "" {
			wg.Add(1)
			if *lFlag {
				go countLines(path, m)
			} else if *wFlag {
				go countWords(path, m)
			} else {
				go countCharacters(path, m)
			}
		}
	}
	wg.Wait()
}
