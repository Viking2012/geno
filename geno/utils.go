package geno

import (
	"fmt"
	"sort"
)

func interfaceToFloat(v any) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return float64(t)
	default:
		return 0
	}
}

func templatizeProps(props map[string]any, assignor, paramPrefix string) []string {
	var (
		keys   []string = make([]string, 0, len(props))
		params []string = make([]string, len(props))
	)

	for key := range props {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		params[i] = fmt.Sprintf("%s%s$%s", key, assignor, paramPrefix+key)
	}
	return params
}
