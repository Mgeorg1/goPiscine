package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	stat, _ := os.Stdin.Stat()
	var stdin []byte
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	argStr := string(stdin)
	args := strings.Split(argStr, " ")

	var inputArgs []string
	if len(os.Args) != 0 {
		inputArgs = os.Args[1:]
	}
	inputArgs = append(inputArgs, args...)
	name := inputArgs[0]
	newArgs := inputArgs[1:]

	com := exec.Command(name, newArgs...)
	com.Stdout = os.Stdout
	err := com.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
