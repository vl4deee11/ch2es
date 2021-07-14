package converter

type null struct{}

func NewNull() Converter {
	return new(null)
}

func (c *null) Convert(d map[string]interface{}, cleaner func(string) string) map[string]interface{} {
	m := d
	for k := range m {
		kk := cleaner(k)
		if kk != k {
			m[kk] = m[k]
			delete(m, k)
		}
	}
	return m
}
