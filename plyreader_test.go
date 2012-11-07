package plyreader

import (
	"testing"
	"fmt"
)

func Test(t *testing.T){
	p := new(PLY)
	e := p.Load("test.ply")
	msg := fmt.Sprintf("\nNumber of elements: %v\n", len(p.Elements))
	for _, vElem := range p.Elements {
		msg += fmt.Sprintf("\telement: %v\n", vElem.Name)
		for _, vProp := range vElem.Properties {
			msg += fmt.Sprintf("\t\tproperty: %v ", vProp.Name)
			if vProp.IsList {
				msg += "[" + vProp.ListSizeType +"]list "
			}
			msg += vProp.Type + "\n"
		}
	}
	msg += fmt.Sprintf("Number of obj_infos: %v\n", len(p.ObjInfoItems))
	msg += fmt.Sprintf("filename: %v\n", p.filename)
	msg += fmt.Sprintf("FileType: %v\n", p.FileType)
	t.Log(msg)
	if e != nil{
		t.Error(e)
	}
	t.Error("oops")
}
