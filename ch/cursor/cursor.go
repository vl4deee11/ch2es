package cursor

import "bytes"

type T int

const (
	Offset T = iota
	Timestamp
	JSONFile
	Stdin
)

type Cursor interface {
	Next() *bytes.Buffer
}
