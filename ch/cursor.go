package ch

import "bytes"

type cursorT int

const (
	offsetCursor cursorT = iota
	timeStampCursor
	jsonFileCursor
	stdinCursor
)

type Cursor interface {
	Next() *bytes.Buffer
}
