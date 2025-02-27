package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/ekediala/jian/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args
	if len(args) == 2 {
		fileName := args[1]
		f, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		repl.Start(f, os.Stdout)
		return
	}

	fmt.Printf("Hello %s! This is the Jian programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
