package resp3

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type testCase struct {
	Description       string
	Input             interface{}
	ExpectedEncOutput interface{}
	ExpectedDecOutput interface{}
	IsEncodeTC        bool
	IsDecodeTC        bool
}

var testCases = []testCase{
	{
		Description:       "Simple String",
		Input:             "hello",
		ExpectedEncOutput: "$5\r\nhello\r\n",
		ExpectedDecOutput: "hello",
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Bulk String",
		Input:             "bulk string example",
		ExpectedEncOutput: "$19\r\nbulk string example\r\n",
		ExpectedDecOutput: "bulk string example",
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Integer",
		Input:             42,
		ExpectedEncOutput: ":42\r\n",
		ExpectedDecOutput: 42,
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Float",
		Input:             42.04,
		ExpectedEncOutput: ",42.040000\r\n",
		ExpectedDecOutput: 42.04,
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Boolean (true)",
		Input:             true,
		ExpectedEncOutput: "#t\r\n",
		ExpectedDecOutput: true,
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Boolean (false)",
		Input:             false,
		ExpectedEncOutput: "#f\r\n",
		ExpectedDecOutput: false,
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Nil",
		Input:             nil,
		ExpectedEncOutput: "_\r\n",
		ExpectedDecOutput: nil,
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Array",
		Input:             []interface{}{"msg", 123},
		ExpectedEncOutput: "*2\r\n$3\r\nmsg\r\n:123\r\n",
		ExpectedDecOutput: []interface{}{"msg", 123},
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Error",
		Input:             errors.New("some error"),
		ExpectedEncOutput: "-some error\r\n",
		ExpectedDecOutput: errors.New("some error"),
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description:       "Array of Interfaces",
		Input:             []interface{}{true, nil, "test"},
		ExpectedEncOutput: "*3\r\n#t\r\n_\r\n$4\r\ntest\r\n",
		ExpectedDecOutput: []interface{}{true, nil, "test"},
		IsEncodeTC:        true,
		IsDecodeTC:        true,
	},
	{
		Description: "Map map[string]interface{}",
		Input: map[string]interface{}{
			"age":       30,
			"isStudent": false,
			"grades": map[string]interface{}{
				"math":    95,
				"science": 90,
			},
		},
		ExpectedEncOutput: "%6\r\n+age\r\n:30\r\n+isStudent\r\n#f\r\n+grades\r\n%4\r\n+math\r\n:95\r\n+science\r\n:90\r\n",
		ExpectedDecOutput: map[string]interface{}{
			"age":       30,
			"isStudent": false,
			"grades": map[string]interface{}{
				"math":    95,
				"science": 90,
			},
		},
		IsEncodeTC: false,
		IsDecodeTC: true,
	},
	{
		Description: "ScalarRecord",
		Input: &ScalarRecord{
			Value:  30,
			Type:   8,
			LAT:    2434,
			Expiry: 3443,
		},
		ExpectedEncOutput: "%8\r\n+Value\r\n:30\r\n+Type\r\n:8\r\n+LAT\r\n:2434\r\n+Expiry\r\n:3443\r\n",
		ExpectedDecOutput: map[string]interface{}{
			"Value":  30,
			"Type":   8,
			"LAT":    2434,
			"Expiry": 3443,
		},
		IsEncodeTC: false,
		IsDecodeTC: true,
	},
	{
		Description: "RecordResponse",
		Input: &RecordResponse{
			Value: 30,
			Code:  8,
		},
		ExpectedEncOutput: "%4\r\n+Value\r\n:30\r\n+Code\r\n:8\r\n",
		ExpectedDecOutput: map[string]interface{}{
			"Value": 30,
			"Code":  8,
		},
		IsEncodeTC: false,
		IsDecodeTC: true,
	},
}

func TestEncode(t *testing.T) {

	for _, tc := range testCases {
		actualEncOutput, err := Encode(tc.Input)

		if err != nil {
			t.Errorf("%s: error while encoding %v", tc.Description, tc.Input)
		}

		if tc.IsEncodeTC == true {
			if actualEncOutput != tc.ExpectedEncOutput {
				t.Errorf("%s: Encoding assertion failed [%v != %v", tc.Description, actualEncOutput, tc.ExpectedEncOutput)
			}
		}
	}
}

func TestDecode(t *testing.T) {

	for _, tc := range testCases {
		actualDecOutput, err := Decode(bufio.NewReader(strings.NewReader(tc.ExpectedEncOutput.(string))))

		if err != nil {
			t.Errorf("%s: error while decoding %v", tc.Description, tc.ExpectedEncOutput)
		}

		fActual := fmt.Sprintf("%#v", actualDecOutput)
		fExpected := fmt.Sprintf("%#v", tc.ExpectedDecOutput)

		if tc.IsDecodeTC == true {
			if fActual != fExpected {
				t.Errorf("%s: Decoding assertion failed for [ %v != %v ]", tc.Description, fActual, fExpected)
			}
		}
	}
}
