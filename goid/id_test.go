package goid

import (
	"testing"
)

func TestGetGid(t *testing.T) {
	SetOpenGID()
	id := GetGID()
	t.Logf("gid %v", id)
}
