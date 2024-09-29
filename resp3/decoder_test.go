package resp3

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func newReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

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

	if !errors.Is(err, ErrUnsupportedRespDataType) {
		t.Fatalf("expected ErrUnsupportedRespDataType, got %v", err)
	}
}

func TestDecodeIncompleteBulkString(t *testing.T) {
	input := "$6\r\nfoo"

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestDecodeIncompleteArray(t *testing.T) {
	input := "*2\r\n$3\r\nfoo\r\n$3\r\n"

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestDecodeIncompleteMap(t *testing.T) {
	// Incomplete map: Only one key-value pair is fully provided, second pair is incomplete
	input := "%4\r\n+key1\r\n$6\r\nvalue1\r\n+key2\r\n"

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestDecodeIncompleteBlobError(t *testing.T) {
	input := "!20\r\nThis is a " // Incomplete blob error: length 20, but only 10 bytes of data provided

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestDecodePartialInteger(t *testing.T) {
	input := ":\r\n" // Incomplete integer: Colon is provided, but no number

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestDecodeCompleteAfterIncomplete(t *testing.T) {
	// First part is an incomplete array, second part completes it
	incompleteInput := "*2\r\n$3\r\nfoo\r\n" // Missing second element
	completeInput := "$3\r\nbar\r\n"

	combinedInput := incompleteInput + completeInput
	reader := bufio.NewReader(strings.NewReader(combinedInput))

	result, err := Decode(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []interface{}{"foo", "bar"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDecodeArrayWithMissingElement(t *testing.T) {
	input := "*2\r\n$3\r\nfoo\r\n" // Array size 2 but only one element provided

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestDecodeMapWithMissingValue(t *testing.T) {
	input := "%4\r\n+key1\r\n" // Map size specifies 4 elements but no values provided

	reader := newReader(input)
	_, err := Decode(reader)

	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}
