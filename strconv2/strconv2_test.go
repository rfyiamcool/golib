package strconv2

import (
	"fmt"
	"testing"
)

const (
	expectedIntType     = "int"
	expectedInt64Type   = "int64"
	expectedStringType  = "string"
	expectedFloat64Type = "float64"

	testStr = "123"
	testInt = 123
)

func TestStrToInt(t *testing.T) {

	val, err := StrToInt(testStr)

	if err != nil {
		t.Fatal(err)
	}

	if typeof(val) != expectedIntType && val != testInt {
		t.Errorf("returned unexpected value : got %v want %v", val, testInt)
	}
}

func TestIntToStr(t *testing.T) {
	val := IntToStr(testInt)

	if typeof(val) != expectedStringType && val != testStr {
		t.Errorf("returned unexpected value : got %v want %v", val, testInt)
	}
}

func TestInt64ToFloat64(t *testing.T) {
	var expected int64
	expected = 123
	val, err := Int64ToFloat64(expected)
	if err != nil {
		t.Fatal(err)
	}

	if typeof(val) != expectedFloat64Type {
		t.Errorf("returned unexpected type : got %v want %v", typeof(val), expectedFloat64Type)
	}
}

func TestInt64ToStr(t *testing.T) {
	var expected int64
	expected = 123
	val := Int64ToStr(expected)
	if typeof(val) != expectedStringType {
		t.Errorf("returned unexpected type : got %v want %v", typeof(val), expectedStringType)
	}
}

func TestStrToInt64(t *testing.T) {
	val, err := StrToInt64(testStr)

	if err != nil {
		t.Fatal(err)
	}

	if typeof(val) != expectedInt64Type {
		t.Errorf("returned unexpected type : got %v want %v", typeof(val), expectedInt64Type)
	}
}

func TestIntToFloat64(t *testing.T) {
	val, err := IntToFloat64(testInt)
	if err != nil {
		t.Fatal(err)
	}
	if typeof(val) != expectedFloat64Type {
		t.Errorf("returned unexpected type : got %v want %v", typeof(val), expectedFloat64Type)
	}
}

func TestStrToFloat64(t *testing.T) {
	val, err := StrToFloat64(testStr)
	if err != nil {
		t.Fatal(err)
	}
	if typeof(val) != expectedFloat64Type {
		t.Errorf("returned unexpected type : got %v want %v", typeof(val), expectedFloat64Type)
	}
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
