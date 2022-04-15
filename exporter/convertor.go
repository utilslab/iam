package exporter

type TypeConverter func(*BasicType, string) string

var _ TypeConverter = typescriptTypeConverter

func typescriptTypeConverter(bt *BasicType, o string) string {
	switch o {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "string", "decimal.Decimal":
		return "string"
	default:
		if bt != nil && bt.Mapping != nil {
			if v, ok := bt.Mapping["ts"]; ok {
				return v.Type
			}
		}
		return o
	}
}
