package repl

import (
	"bufio"
	"fmt"
	"io"
	"ligma/lexer"
	"ligma/parser"
	"ligma/runtime"
	"os"
)

const PROMPT = ">> "

// Start starts the REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	//env := runtime.NewEnvironment()
	i := runtime.NewInterpreter()
	r := runtime.NewResolver(i)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)


		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			os.Exit(1)
		}

		r.Resolve(program.Statements)

		evaluated := i.Interpret(program)

		if evaluated != nil {
			instance := evaluated.(*runtime.LigmaInstance)
			repr_func, _ :=instance.Get("__repr__")
			repr_fun := repr_func.(runtime.LigmaCallable)
			io.WriteString(out, repr_fun.Call(nil, nil).Inspect())
			//io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t" + msg + "\n")
	}

	os.Exit(69)
}

// run a script file
func RunFile(path string) {
	data, err := os.ReadFile(path)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    
	l := lexer.New(string(data))
	p := parser.New(l)
	i := runtime.NewInterpreter()
	r := runtime.NewResolver(i)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	r.Resolve(program.Statements)

	evaluated := i.Interpret(program)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}