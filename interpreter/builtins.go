package interpreter

import (
	"fmt"
	"math/rand"
)

func registerBuiltins(env *Environment) {
	env.Set("panjang", &BuiltinFunc{
		Name: "panjang",
		Fn: func(args []interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("panjang() expects 1 argument — bagi satu je lah")
			}
			switch v := args[0].(type) {
			case string:
				return float64(len(v)), nil
			case []interface{}:
				return float64(len(v)), nil
			default:
				return nil, fmt.Errorf("panjang() tak boleh guna for type %T", args[0])
			}
		},
	}, false)

	env.Set("jenis", &BuiltinFunc{
		Name: "jenis",
		Fn: func(args []interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("jenis() expects 1 argument")
			}
			switch args[0].(type) {
			case float64:
				return "nombor", nil
			case string:
				return "ayat", nil
			case bool:
				return "boleh_tak_boleh", nil
			case []interface{}:
				return "senarai", nil
			case *Function, *BuiltinFunc:
				return "kerja", nil
			case nil:
				return "kosong", nil
			default:
				return "tak tahu", nil
			}
		},
	}, false)

	env.Set("masuk", &BuiltinFunc{
		Name: "masuk",
		Fn: func(args []interface{}) (interface{}, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("masuk() needs array and value — bagi dua lah")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				return nil, fmt.Errorf("masuk() first arg must be senarai")
			}
			return append(arr, args[1:]...), nil
		},
	}, false)

	env.Set("potong", &BuiltinFunc{
		Name: "potong",
		Fn: func(args []interface{}) (interface{}, error) {
			if len(args) < 2 || len(args) > 3 {
				return nil, fmt.Errorf("potong(senarai, mula, habis) — 2 or 3 args")
			}
			arr, ok := args[0].([]interface{})
			if !ok {
				return nil, fmt.Errorf("potong() first arg must be senarai")
			}
			start := int(toNumber(args[1]))
			end := len(arr)
			if len(args) == 3 {
				end = int(toNumber(args[2]))
			}
			if start < 0 {
				start = 0
			}
			if end > len(arr) {
				end = len(arr)
			}
			result := make([]interface{}, end-start)
			copy(result, arr[start:end])
			return result, nil
		},
	}, false)

	env.Set("nombor", &BuiltinFunc{
		Name: "nombor",
		Fn: func(args []interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("nombor() expects 1 argument")
			}
			switch v := args[0].(type) {
			case float64:
				return v, nil
			case string:
				var n float64
				_, err := fmt.Sscanf(v, "%f", &n)
				if err != nil {
					return nil, fmt.Errorf("tak boleh convert '%s' to nombor lah", v)
				}
				return n, nil
			case bool:
				if v {
					return float64(1), nil
				}
				return float64(0), nil
			default:
				return nil, fmt.Errorf("tak boleh convert to nombor")
			}
		},
	}, false)

	env.Set("cincai", &BuiltinFunc{
		Name: "cincai",
		Fn: func(args []interface{}) (interface{}, error) {
			switch len(args) {
			case 0:
				// No args: return random float between 0 and 1
				return rand.Float64(), nil
			case 1:
				// One arg: return random int from 0 to max (exclusive)
				max := int(toNumber(args[0]))
				if max <= 0 {
					return nil, fmt.Errorf("cincai() max must be > 0 lah")
				}
				return float64(rand.Intn(max)), nil
			case 2:
				// Two args: return random int from min to max (inclusive)
				min := int(toNumber(args[0]))
				max := int(toNumber(args[1]))
				if min >= max {
					return nil, fmt.Errorf("cincai() min must be < max lah")
				}
				return float64(min + rand.Intn(max-min+1)), nil
			default:
				return nil, fmt.Errorf("cincai() takes 0, 1, or 2 arguments — jangan bagi lebih")
			}
		},
	}, false)

	env.Set("ayat", &BuiltinFunc{
		Name: "ayat",
		Fn: func(args []interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("ayat() expects 1 argument")
			}
			return stringify(args[0]), nil
		},
	}, false)
}
