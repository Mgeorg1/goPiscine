package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	fFlag := flag.Bool("f", false, "find files")
	dFlag := flag.Bool("d", false, "find directories")
	slFlag := flag.Bool("sl", false, "find symlink")
	extFlag := flag.String("ext", "", "filter by extension")
	flag.Parse()

	if *extFlag != "" && !*fFlag {
		fmt.Println("Flag error. Use -f for using -ext")
		os.Exit(1)
	}
	if !*fFlag && !*dFlag && !*slFlag {
		*fFlag, *dFlag, *slFlag = true, true, true
	}

	f := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if (info.Mode()&fs.ModeSymlink != 0) && *slFlag {
			link, err := filepath.EvalSymlinks(path)
			if err != nil {
				link = "[broken]"
			}
			fmt.Printf("%s -> %s\n", path, link)
		} else if info.IsDir() && *dFlag {
			fmt.Println(path)
		} else if *fFlag && info.Mode()>>9 == 0 {
			if *extFlag != "" {
				ext := "." + *extFlag
				if filepath.Ext(path) == ext {
					fmt.Println(path)
				}
			} else if *extFlag == "" {
				fmt.Println(path)
			}
		}
		return nil
	}
	err := filepath.Walk(flag.Arg(0), f)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
