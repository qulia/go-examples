package webcrawler

import (
	"fmt"
	"strings"
)

type compositeerror struct {
	errors []error
}

func (ce *compositeerror) Error() string {
	sb := &strings.Builder{}
	for _, e := range ce.errors {
		fmt.Fprintf(sb, "%s\n", e.Error())
	}

	return sb.String()
}

func (ce *compositeerror) Add(e error) {
	ce.errors = append(ce.errors, e)
}
