package converter

type T int

const (
	Null T = iota
	Nested
)

type Converter interface {
	Convert(map[string]interface{}, func(string) string) map[string]interface{}
}
