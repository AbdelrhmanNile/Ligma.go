package main

import (
	"fmt"
	"ligma/repl"
	"os"
)

func main(){
	
	// ccheck if a file was passed as an argument
	if len(os.Args) > 1 {
		repl.RunFile(os.Args[1])
		return
	}

	fmt.Printf("Hello! This is the Ligma programming language!\n")
	fmt.Printf("Let's get ballin'!\n\n")
	repl.Start(os.Stdin, os.Stdout)
}