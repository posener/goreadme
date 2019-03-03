package goreadme

import "io"

// multiNewLineEliminator implements io.Writer interface, and makes
// sure that no more than 2 new lines are written in a row.
type multiNewLineEliminator struct {
	w io.Writer
	// newLines should not be set, it counts the number of new lines
	// that were written in a row.
	newLines int
}

func (e *multiNewLineEliminator) Write(in []byte) (int, error) {
	out := make([]byte, 0, len(in))
	n := 0
	for _, c := range in {
		if c == '\n' {
			e.newLines++
			if e.newLines > 2 {
				continue
			}
		} else {
			e.newLines = 0
		}
		out = append(out, c)
		n++
	}
	return e.w.Write(out[:n])
}
