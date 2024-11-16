package static_analysis

import (
	"fmt"
	"strings"
	"text/scanner"
)

func lexicalAnalysis(phpCode string) {
	var s scanner.Scanner
	s.Init(strings.NewReader(phpCode))
	s.Filename = "file"

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("Token: %s\tValue: %s\tPosition: %s\n",
			s.TokenText(), scanner.TokenString(tok), s.Pos(),
		)
	}
}
