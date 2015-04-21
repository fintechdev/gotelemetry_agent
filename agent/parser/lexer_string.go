package parser

import (
	"fmt"
)

func (t token) String() string {
	return fmt.Sprintf("%s %q [%d:%d]", t.terminal, t.source, t.line, t.start)
}
