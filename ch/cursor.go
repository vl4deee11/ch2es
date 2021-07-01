package ch

import "bytes"

type cursorT int

const (
	offsetCursor cursorT = iota
	timeStampCursor
	fileCursor
)

type Cursor interface {
	Next() *bytes.Buffer
}
