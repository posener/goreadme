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
			// Skip not first new line.
			if e.newLines > 1 {
				continue
			}
		} else {
			if e.newLines > 1 {
				// Add second new line if originally there were more
				// than 1 new line, only if there is another character
				// to write after it.
				// This eliminates multiple new lines in the end of
				// the document.
				out = append(out, '\n')
				n++
			}
			e.newLines = 0
		}
		out = append(out, c)
		n++
	}
	return e.w.Write(out[:n])
}
