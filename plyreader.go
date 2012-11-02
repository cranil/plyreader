// cc: Anil C Raghuramu <anil.c.raghuramu@gmail.com>

package plyreader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sort"
)

type Property struct {
	data [][]byte
}

type Element struct {
	Properties map[string]Property
}

const (
	BinaryBigEndian = 0
	BinaryLittleEndian = 1
	Ascii = 2
)

type obj_info

type PLY struct {
	Elements map[string]Element
	FileType int8
}

func NewPLY(filename string) ( p PLY, err error){
	file, file_io_error := os.Open(filename)
	if file_io_error != nil {
		return nil, file_io_error
	}
	
	return PLY{}, nil
}

var types []string = []string{"invalid", "int8", "int16", "int32", "uint8", "uint16", "uint32", "float32", "float64"}
var old_types []string = []string{"invalid", "char", "short", "int", "uchar", "ushort", "uint", "float", "double"}

func read_type( type_name, data string ) []byte {
	return 0
}

func parse_header(file *os.File, p *PLY) {
	r := bufio.NewReader(file)
}

