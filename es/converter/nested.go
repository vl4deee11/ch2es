package converter

type nested struct {
	field   string
	addNull bool
}

type NestedConverterConf struct {
	Field   string `desc:"nested array field name"`
	AddNull bool   `desc:"add null"`
}

func NewNested(cfg *NestedConverterConf) Converter {
	return &nested{field: cfg.Field, addNull: cfg.AddNull}
}

func (c *nested) Convert(m map[string]interface{}, cleaner func(string) string) map[string]interface{} {
	maxL := -1
	resM := make(map[string]interface{}, len(m))
	for k := range m {
		s, ok := m[k].([]interface{})
		if !ok {
			resM[cleaner(k)] = m[k]
			delete(m, k)
			continue
		}

		if len(s) > maxL {
			maxL = len(s)
		}
	}

	resM[c.field] = make([]interface{}, 0, maxL)
	for i := 0; i < maxL; i++ {
		currM := make(map[string]interface{})
		for k := range m {
			s := m[k].([]interface{})

			if i < len(s) && s[i] != nil {
				currM[cleaner(k)] = s[i]
			} else if c.addNull {
				currM[cleaner(k)] = nil
			}
		}
		resM[c.field] = append(resM[c.field].([]interface{}), currM)
	}
	return resM
}
