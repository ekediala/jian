package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ekediala/jian/lexer"
	"github.com/ekediala/jian/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		if ok := scanner.Scan(); !ok {
			break
		}
		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}

	}
}
