package xtermkeyboard

import (
	"bytes"
	"encoding/hex"
)

type KeyboardType string

var (
	KeyboardTypeDelete  KeyboardType = "delete"
	KeyboardTypeEnter   KeyboardType = "enter"
	KeyboardTypeSpace   KeyboardType = "space"
	KeyboardTypeTab     KeyboardType = "tab"
	KeyboardTypeEsc     KeyboardType = "esc"
	KeyboardTypeUp      KeyboardType = "up"
	KeyboardTypeDown    KeyboardType = "down"
	KeyboardTypeLeft    KeyboardType = "left"
	KeyboardTypeRight   KeyboardType = "right"
	KeyboardTypeUnknown KeyboardType = "unknown"
)

type KeyboardCode string

var (
	keyboardCodeDelete = "7f"
	keyboardCodeEnter  = "0d"
	keyboardCodeSpace  = "20"
	keyboardCodeTab    = "09"
	keyboardCodeEsc    = "1b"
	keyboardCodeUp     = "1b5b41"
	keyboardCodeDown   = "1b5b42"
	keyboardCodeRight  = "1b5b43"
	keyboardCodeLeft   = "1b5b44"
)

func ToCodeString(bs []byte) string {
	return hex.EncodeToString(bs)
}

func ToCode(bs []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(bs)))
	hex.Encode(dst, bs)
	return dst
}

func GetCode(bs []byte) KeyboardType {
	kbcode := ToCodeString(bs)
	switch kbcode {
	case keyboardCodeDelete:
		return KeyboardTypeDelete
	case keyboardCodeEnter:
		return KeyboardTypeEnter
	case keyboardCodeSpace:
		return KeyboardTypeSpace
	case keyboardCodeTab:
		return KeyboardTypeTab
	case keyboardCodeEsc:
		return KeyboardTypeEsc
	case keyboardCodeUp:
		return KeyboardTypeUp
	case keyboardCodeDown:
		return KeyboardTypeDown
	case keyboardCodeRight:
		return KeyboardTypeRight
	case keyboardCodeLeft:
		return KeyboardTypeLeft
	default:
		return KeyboardTypeUnknown
	}
}

func IsKeyboardDelete(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeDelete)) {
		return true
	}
	return false
}

func IsKeyboardEnter(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeEnter)) {
		return true
	}
	return false
}

func IsKeyboardSpace(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeSpace)) {
		return true
	}
	return false
}

func IsKeyboardTab(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeTab)) {
		return true
	}
	return false
}

func IsKeyboardEsc(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeEsc)) {
		return true
	}
	return false
}

func IsKeyboardUp(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeUp)) {
		return true
	}
	return false
}

func IsKeyboardDown(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeDown)) {
		return true
	}
	return false
}

func IsKeyboardRight(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeRight)) {
		return true
	}
	return false
}

func IsKeyboardLeft(bs []byte) bool {
	if bytes.Equal(ToCode(bs), []byte(keyboardCodeLeft)) {
		return true
	}
	return false
}
