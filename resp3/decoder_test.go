package resp3

import (
	"bufio"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// Helper function to create a reader from a string
func newReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

// Test for Simple String
func TestDecodeSimpleString(t *testing.T) {
	input := "+OK\r\n"
	expected := "OK"

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Error
func TestDecodeError(t *testing.T) {
	input := "-Error message\r\n"
	expected := errors.New("Error message")

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if errResult, ok := result.(error); ok {
		if errResult.Error() != expected.Error() {
			t.Errorf("expected %v, got %v", expected.Error(), errResult.Error())
		}
	} else {
		t.Errorf("expected error, got non-error result: %v", result)
	}
}

// Test for Integer
func TestDecodeInteger(t *testing.T) {
	input := ":42\r\n"
	expected := int64(42)

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Negative Integer
func TestDecodeNegativeInteger(t *testing.T) {
	input := ":-42\r\n"
	expected := int64(-42)

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Float
func TestDecodeFloat(t *testing.T) {
	input := ",3.14159\r\n"
	expected := 3.14159

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Bulk String
func TestDecodeBulkString(t *testing.T) {
	input := "$6\r\nfoobar\r\n"
	expected := "foobar"

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Null Bulk String
func TestDecodeNullBulkString(t *testing.T) {
	input := "$-1\r\n"
	var expected interface{} = nil

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Empty Bulk String
func TestDecodeEmptyBulkString(t *testing.T) {
	input := "$0\r\n\r\n"
	expected := ""

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Verbatim String
func TestDecodeVerbatimString(t *testing.T) {
	input := "=13\r\nsome verbatim\r\n"
	expected := "some verbatim"

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Empty Verbatim String
func TestDecodeEmptyVerbatimString(t *testing.T) {
	input := "=0\r\n\r\n"
	expected := ""

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Array
func TestDecodeArray(t *testing.T) {
	input := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	expected := []interface{}{"foo", "bar"}

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Empty Array
func TestDecodeEmptyArray(t *testing.T) {
	input := "*0\r\n"
	expected := []interface{}{}

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Null Array
func TestDecodeNullArray(t *testing.T) {
	input := "*-1\r\n"
	var expected interface{} = nil

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Boolean True
func TestDecodeBooleanTrue(t *testing.T) {
	input := "#t\r\n"
	expected := true

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Boolean False
func TestDecodeBooleanFalse(t *testing.T) {
	input := "#f\r\n"
	expected := false

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Map with String Keys
func TestDecodeStringKeyMap(t *testing.T) {
	input := "%4\r\n+key1\r\n$6\r\nvalue1\r\n+key2\r\n$6\r\nvalue2\r\n"
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Map with Int64 Keys
func TestDecodeInt64KeyMap(t *testing.T) {
	input := "%4\r\n:1\r\n$6\r\nvalue1\r\n:2\r\n$6\r\nvalue2\r\n"
	expected := map[int64]interface{}{
		1: "value1",
		2: "value2",
	}

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Map with Mixed Keys
func TestDecodeMixedKeyMap(t *testing.T) {
	input := "%6\r\n+key1\r\n$6\r\nvalue1\r\n:2\r\n$6\r\nvalue2\r\n+key3\r\n$6\r\nvalue3\r\n"
	expected := map[interface{}]interface{}{
		"key1":   "value1",
		int64(2): "value2",
		"key3":   "value3",
	}

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test for Blob Error
func TestDecodeBlobError(t *testing.T) {
	input := "!20\r\nThis is a blob error\r\n"
	expected := errors.New("This is a blob error")

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if errResult, ok := result.(error); ok {
		if errResult.Error() != expected.Error() {
			t.Errorf("expected %v, got %v", expected.Error(), errResult.Error())
		}
	} else {
		t.Errorf("expected error, got non-error result: %v", result)
	}
}

// Test for Null
func TestDecodeNull(t *testing.T) {
	input := "_\r\n"
	var expected interface{} = nil

	reader := newReader(input)
	result, err := Decode(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDecodeUnsupportedType(t *testing.T) {
	input := "&\r\n"

	reader := newReader(input)
	result, err := Decode(reader)

	if err == nil {
		t.Fatalf("expected error, got %v", result)
	}
}
