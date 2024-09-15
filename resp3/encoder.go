package resp3

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Encode converts a Go data type into its corresponding RESP3 encoded string format.
// It supports a wide range of Go types, including basic types, composite types,
// and user-defined types, ensuring flexible and efficient serialization.
//
// Supported Types:
//
//   - **String**: Encodes Go strings as RESP3 bulk strings.
//     Example: "hello" -> "$5\r\nhello\r\n"
//
//   - **Integers**: Supports all Go integer types (int, int8, int16, int32, int64, uint, uint8, etc.)
//     and encodes them as RESP3 integers.
//     Example: 123 -> ":123\r\n"
//
//   - **Floats**: Encodes float32 and float64 types as RESP3 floating-point numbers.
//     Example: 3.14 -> ",3.140000\r\n"
//
//   - **Booleans**: Encodes booleans (true/false) as RESP3 boolean values.
//     Example: true -> "#t\r\n", false -> "#f\r\n"
//
//   - **Nil**: Encodes nil as RESP3 null.
//     Example: nil -> "_\r\n"
//
//   - **Errors**: Encodes Go error types as RESP3 errors.
//     Example: errors.New("error message") -> "-error message\r\n"
//
//   - **Slices**: Supports slices of any type (e.g., []string, []int, []float64, etc.) and encodes them as RESP3 arrays.
//     Each element of the slice is recursively encoded using the same rules.
//     Example: []int{1, 2, 3} -> "*3\r\n:1\r\n:2\r\n:3\r\n"
//
//   - **Maps**: Supports maps with either string keys or interface{} keys (e.g., map[string]interface{}, map[interface{}]interface{}).
//     The key-value pairs are encoded as RESP3 maps. The keys and values are recursively encoded.
//     Example: map[string]interface{}{"a": 1, "b": 2} -> "%4\r\n+a\r\n:1\r\n+b\r\n:2\r\n"
//
//   - **Structs**: Encodes Go structs by treating field names as map keys and field values as map values.
//     Each field is recursively encoded.
//     Example: struct{ Name string; Age int } -> "%4\r\n+Name\r\n$5\r\nAlice\r\n+Age\r\n:25\r\n"
//
//   - **time.Time**: Encodes time.Time values as Unix timestamps in milliseconds.
//     Example: time.Now() -> ":1620832335000\r\n"
//
//   - **Custom Types**: Custom types (like ScalarRecord or RecordResponse) are handled by converting them to maps and encoding them recursively.
//
// Parameters:
//   - value: The Go value to be encoded. This value can be of any supported type, including
//     basic types (like int, string, float), composite types (like slices, maps, structs), or custom types.
//
// Returns:
//
//   - string: The RESP3-encoded string representation of the input value.
//
//   - error: An error is returned if the value type is not supported for encoding, or if
//     any issue arises during the encoding process.
//
// Example Usage:
//
//	encoded, err := Encode("hello")
//	encoded, err := Encode(42)
//	encoded, err := Encode(map[string]interface{}{"a": 1, "b": true})
//	encoded, err := Encode(time.Now())
//
// The Encode function is designed to be recursive, meaning composite types (e.g., slices, maps, structs)
// will have their elements or fields encoded individually according to their respective types.
// This ensures that nested data structures can be efficiently serialized into RESP3 format.
func Encode(value interface{}) (string, error) {
	switch v := value.(type) {

	// Strings
	case string:
		// If the string is short enough (less than 16 chars), use Simple String
		if len(v) <= (1 << 4) {
			return "+" + v + "\r\n", nil // Simple String
		}
		// Otherwise, treat it as a Bulk String
		return "$" + strconv.Itoa(len(v)) + "\r\n" + v + "\r\n", nil

	// Integers and their variations
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return ":" + fmt.Sprintf("%d", v) + "\r\n", nil

	// Floats
	case float32, float64:
		return "," + fmt.Sprintf("%f", v) + "\r\n", nil

	// Boolean
	case bool:
		if v {
			return "#t\r\n", nil
		}
		return "#f\r\n", nil

	// Nil
	case nil:
		return "_\r\n", nil

	// Error
	case error:
		return "-" + v.Error() + "\r\n", nil

		// Arrays of interface{}
	case []interface{}:
		resp := "*" + strconv.Itoa(len(v)) + "\r\n"
		for _, elem := range v {
			// Handle strings separately to use Simple Strings for short text
			if str, ok := elem.(string); ok && len(str) <= 12 {
				resp += "+" + str + "\r\n" // Use Simple String for short strings
			} else {
				encodedElem, err := Encode(elem)
				if err != nil {
					return "", err
				}
				resp += encodedElem
			}
		}
		return resp, nil

		// Arrays of strings
	case []string:
		resp := "*" + strconv.Itoa(len(v)) + "\r\n"
		for _, elem := range v {
			encodedElem := "+" + elem + "\r\n" // Change to Simple String
			resp += encodedElem
		}
		return resp, nil

	// Arrays of integers (all int types)
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
		val := reflect.ValueOf(v)
		resp := "*" + strconv.Itoa(val.Len()) + "\r\n"
		for i := 0; i < val.Len(); i++ {
			encodedElem, err := Encode(val.Index(i).Interface())
			if err != nil {
				return "", err
			}
			resp += encodedElem
		}
		return resp, nil

	// Arrays of bools
	case []bool:
		resp := "*" + strconv.Itoa(len(v)) + "\r\n"
		for _, elem := range v {
			encodedElem, err := Encode(elem)
			if err != nil {
				return "", err
			}
			resp += encodedElem
		}
		return resp, nil

	// Arrays of float32 and float64
	case []float32, []float64:
		val := reflect.ValueOf(v)
		resp := "*" + strconv.Itoa(val.Len()) + "\r\n"
		for i := 0; i < val.Len(); i++ {
			encodedElem, err := Encode(val.Index(i).Interface())
			if err != nil {
				return "", err
			}
			resp += encodedElem
		}
		return resp, nil

	// Map with string keys and interface values
	case map[string]interface{}:
		resp := "%" + strconv.Itoa(len(v)*2) + "\r\n"
		for kx, vx := range v {
			resp += "+" + kx + "\r\n"
			valueStr, err := Encode(vx)
			if err != nil {
				return "", err
			}
			resp += valueStr
		}
		return resp, nil

		// Map with interface{} keys and values (map[interface{}]interface{})
	case map[interface{}]interface{}:
		resp := "%" + strconv.Itoa(len(v)*2) + "\r\n"
		for kx, vx := range v {
			var keyStr string
			var err error

			// Check the type of the key and encode accordingly
			switch key := kx.(type) {
			case string:
				keyStr = "+" + key + "\r\n" // Simple string

			default:
				keyStr, err = Encode(key) // Other types
			}

			if err != nil {
				return "", err
			}

			valueStr, err := Encode(vx)
			if err != nil {
				return "", err
			}

			resp += keyStr + valueStr
		}
		return resp, nil

	// time.Time encoded as Unix timestamp in milliseconds
	case time.Time:
		return ":" + strconv.FormatInt(v.UnixMilli(), 10) + "\r\n", nil

	// Handle structs
	case struct{}:
		return encodeStruct(v)

	default:
		// Handle structs through reflection if no direct case matches
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Struct {
			return encodeStruct(rv.Interface())
		}

		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

func encodeStruct(s interface{}) (string, error) {
	val := reflect.ValueOf(s)
	typ := val.Type()

	// Create the response map based on the number of exported fields
	resp := "%" + strconv.Itoa(val.NumField()*2) + "\r\n"

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		fieldValue := val.Field(i).Interface()

		resp += "+" + fieldName + "\r\n"

		encodedValue, err := Encode(fieldValue)
		if err != nil {
			return "", err
		}

		resp += encodedValue
	}

	return resp, nil
}
