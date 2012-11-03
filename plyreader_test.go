package plyreader

import (
	"testing"
)

func Test(t *testing.T) {
	p := PLY{}
	e := p.Load("test.ply")
	if e != nil {
		t.Error(e)
	}
}