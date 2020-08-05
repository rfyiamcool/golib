package strconv2

import "strconv"

// StrToInt `strconv.Atoi(s)`.
func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// MustStrToInt
func MustStrToInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

// IntToStr `strconv.Itoa(i)`.
func IntToStr(i int) string {
	return strconv.Itoa(i)
}

// IntToFloat64 int convert to float64
func IntToFloat64(i int) (float64, error) {
	return StrToFloat64(IntToStr(i))
}

// MustIntToFloat64
func MustIntToFloat64(i int) float64 {
	v, err := StrToFloat64(IntToStr(i))
	checkError(err)
	return v
}

// Int64ToFloat64 int64 convert to float64
func Int64ToFloat64(i int64) (float64, error) {
	return StrToFloat64(Int64ToStr(i))
}

// MustInt64ToFloat64
func MustInt64ToFloat64(i int64) float64 {
	v, err := StrToFloat64(Int64ToStr(i))
	checkError(err)
	return v
}

// StrToFloat64 str convert to float64
func StrToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// MustStrToFloat64
func MustStrToFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	checkError(err)
	return v
}

// StrToInt64
func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// MustStrToInt64
func MustStrToInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	checkError(err)
	return v
}

// Int64ToStr
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// ByteToInt64
func ByteToInt64(b []byte) (int64, error) {
	return StrToInt64(string(b))
}

// MustByteToInt64
func MustByteToInt64(b []byte) int64 {
	v, err := StrToInt64(string(b))
	checkError(err)
	return v
}

// ByteToInt
func ByteToInt(b []byte) (int, error) {
	return StrToInt(string(b))
}

// MustByteToInt
func MustByteToInt(b []byte) int {
	v, err := StrToInt(string(b))
	checkError(err)
	return v
}

// ByteToFloat64
func ByteToFloat64(b []byte) (float64, error) {
	return StrToFloat64(string(b))
}

// MustByteToFloat64
func MustByteToFloat64(b []byte) (float64, error) {
	return StrToFloat64(string(b))
}

// checkError if err not nil, direct panic
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
