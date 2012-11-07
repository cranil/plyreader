// 
// Author: Anil Raghuramu <anil.c.raghuramu@gmail.com>
// Copyright: Copyright (c) 2012, Anil Raghuramu
// 
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
// 
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
// 02110-1301, USA.
// 

package plyreader

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Types []string = []string{"invalid", "int8", "int16", "int32", "uint8", "uint16", "uint32", "float32", "float64"}

var OldTypes []string = []string{"invalid", "char", "short", "int", "uchar", "ushort", "uint", "float", "double"}

var SizeOfType = map[string]int{
	"invalid": 0,
	"int8":    1,
	"int16":   2,
	"int32":   4,
	"uint8":   1,
	"uint16":  2,
	"uint32":  4,
	"float32": 4,
	"float64": 8,
	"char":    1,
	"short":   2,
	"int":     4,
	"uchar":   1,
	"ushort":  2,
	"uint":    4,
	"float":   4,
	"double":  8}

type Property struct {
	Name         string
	IsList       bool
	Data         [][]byte
	Type         string
	ListSizeType string
	pos          int
}

type Element struct {
	Name       string
	Properties []*Property
	Size       int
}

func (p *Property)print() {
	if p.IsList {
		fmt.Printf("\tproperty list %s %s %s\n", p.ListSizeType, p.Type, p.Name)
	} else {
		fmt.Printf("\tproperty %s %s\n", p.Type, p.Name)
	}
}

func (e *Element) print() {
	fmt.Printf("element %s\n", e.Name)
}

const (
	BinaryBigEndian    = 0
	BinaryLittleEndian = 1
	Ascii              = 2
)

type PLY struct {
	Elements     []*Element
	FileType     int8
	ObjInfoItems map[string]string
	currentLine  int
	filename     string
	reader       *bufio.Reader
}

func (p *PLY) Save(filename string) error {
	return nil
}

func (p *PLY) Load(filename string) error {
	p.filename = filename
	file, file_io_error := os.Open(filename)
	p.reader = bufio.NewReader(file)
	if file_io_error != nil {
		return file_io_error
	}
	e := parse_header(p)
	if e != nil {
		return e
	}
	switch p.FileType {
	case BinaryBigEndian:
		e = parse_binary_big_endian(p)
	case BinaryLittleEndian:
		e = parse_binary_little_endian(p)
	case Ascii:
		e = parse_ascii(p)
	default:
		e = errors.New("File type error")
	}
	return e
}

func (p *PLY) ReadVertices() ([][]float32) {
	flag := false
	count := 0
	fmt.Println(count)
	for _, elem := range p.Elements {
		if elem.Name == "vertex" {
			flag = true
			break
		}
		count++
	}
	if flag {
		data := make([][]float32, 3)
		for j:=0; j<3; j++ {
			fmt.Println(count)
			elem := p.Elements[count]
			b := elem.Properties[j].Data
			i := 0
			data[j] = make([]float32, elem.Size)
			for _, v := range b {
				buf := bytes.NewBuffer(v)
				binary.Write(buf, binary.LittleEndian, data[j][i])
				i++
			}
		}
		return data
	}
	return nil
}

func strip(s string) string {
	return strings.TrimSpace(s)
}

func readLine(r *bufio.Reader) (line string, e error) {
	line, e = r.ReadString('\n')
	if e != nil {
		return line, e
	}
	return strip(line), nil
}

func to_type(data, type_name string) (b []byte, e error) {
	var n int64
	var u uint64
	var f float64
	switch {
	case type_name == Types[1] || type_name == OldTypes[1]:
		n, e = strconv.ParseInt(data, 0, 8)
		t := int8(n)
		if e != nil {
			return nil, e
		}
		b = make([]byte, 1)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[2] || type_name == OldTypes[2]:
		n, e = strconv.ParseInt(data, 0, 16)
		t := int16(n)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 2)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[3] || type_name == OldTypes[3]:
		n, e = strconv.ParseInt(data, 0, 32)
		t := int32(n)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 4)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[4] || type_name == OldTypes[4]:
		u, e = strconv.ParseUint(data, 0, 8)
		t := uint8(u)
		if e != nil {
			return nil, e
		}
		b = make([]byte, 1)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[5] || type_name == OldTypes[5]:
		u, e = strconv.ParseUint(data, 0, 16)
		t := uint16(u)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 2)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[6] || type_name == OldTypes[6]:
		u, e = strconv.ParseUint(data, 0, 32)
		t := uint32(u)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 4)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[7] || type_name == OldTypes[7]:
		f, e = strconv.ParseFloat(data, 32)
		t := float32(f)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 4)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[8] || type_name == OldTypes[8]:
		f, e = strconv.ParseFloat(data, 64)
		t := float64(f)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 8)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	}
	return nil, nil
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func parse_header(p *PLY) error {
	r := p.reader
	line, e := readLine(r)
	if e != nil {
		return e
	}
	if line != "ply" {
		return errors.New(p.filename + " is not a ply file.")
	}
	p.currentLine++

	line, e = readLine(r)
	if e != nil {
		return e
	}
	p.currentLine++
	wordMatcher, e := regexp.Compile("((\\d+)*\\.*\\d+)|\\w+")
	words := wordMatcher.FindAllStringSubmatch(line, -1)
	if len(words) != 3 {
		return errors.New("Incorrect format in " +
			p.filename + " at line " + itoa(p.currentLine))
	}
	if words[0][0] != "format" {
		return errors.New("Incorrect format in " +
			p.filename + " at line " + itoa(p.currentLine))
	}
	switch {
	case words[1][0] == "ascii":
		p.FileType = Ascii
	case words[1][0] == "binary_big_endian":
		p.FileType = BinaryBigEndian
	case words[1][0] == "binary_little_endian":
		p.FileType = BinaryLittleEndian
	default:
		return errors.New("Incorrect format in " +
			p.filename + " at line " + itoa(p.currentLine))
	}
	currentElem := -1
	propPos     := 0
	for {
		line, e = readLine(r)
		if e != nil {
			return e
		}
		p.currentLine++
		words = wordMatcher.FindAllStringSubmatch(line, -1)
		if words[0][0] == "comment" {
			// skip
		} else if words[0][0] == "element" {
			elemName := words[1][0]
			elem := new(Element)
			num, e := strconv.ParseInt(words[2][0], 0, 32)
			if e != nil {
				return errors.New(e.Error() + p.filename +
					" at line " + itoa(p.currentLine))
			}
			elem.Size = int(num)
			elem.Name = elemName
			p.Elements = append(p.Elements, elem)
			currentElem++
			propPos = 0
		} else if words[0][0] == "property" {
			cnt := 1
			currWord := words[cnt][0]
			prop := new(Property)
			prop.pos = propPos
			if currWord == "list" {
				prop.IsList = true
				cnt++
				currWord = words[cnt][0]
				prop.ListSizeType = currWord
				cnt++
				currWord = words[cnt][0]
			} else {
				prop.IsList = false
			}
			prop.Type = currWord
			cnt++
			currWord = words[cnt][0]
			propName := currWord
			prop.Name = propName
			p.Elements[currentElem].Properties =
				append(p.Elements[currentElem].Properties, prop)
			propPos++
		} else if words[0][0] == "obj_info" {
			if p.ObjInfoItems == nil {
				p.ObjInfoItems = make(map[string]string)
			}
			p.ObjInfoItems[words[1][0]] = words[2][0]
		} else if words[0][0] == "end_header" {
			break
		}
	}
	return nil
}

func appendBytes(slice, data[]byte) []byte {
    l := len(slice);
    if l + len(data) > cap(slice) {	// reallocate
    	// Allocate double what's needed, for future growth.
    	newSlice := make([]byte, (l+len(data))*2);
    	// Copy data (could use bytes.Copy()).
    	for i, c := range slice {
    		newSlice[i] = c
    	}
    	slice = newSlice;
    }
    slice = slice[0:l+len(data)];
    for i, c := range data {
    	slice[l+i] = c
    }
    return slice;
}

func parse_binary_big_endian(p *PLY) error {
	e := errors.New("Not yet implemented")
	return e
}

func parse_binary_little_endian(p *PLY) error {
	e := errors.New("Not yet implemented")
	return e
}

func parse_ascii(p *PLY) error {
	r := p.reader
	numMatcher, e := regexp.Compile("[\\+\\-]*([0-9]*)+\\.*[0-9]+")
	if e != nil {
		return e
	}
	for _, elem := range p.Elements {
		elem.print()
		for _, prop := range elem.Properties {
			prop.print()
			prop.Data = make([][]byte, elem.Size)
		}
		for i:=0; i<elem.Size; i++ {
			line, e := readLine(r)
			if e != nil {
				return e
			}
			words := numMatcher.FindAllStringSubmatch(line, -1)
			currWord := 0
			if words == nil {
				// skip empty lines
			} else {
				for _, prop := range elem.Properties {
					if prop.IsList {
						num, e := strconv.ParseInt(words[currWord][0], 10, 32)
						if e != nil {
							return e
						}
						numSize := int(num)
						currWord++
						l := make([]byte, numSize*SizeOfType[prop.Type])
						for j:=0; j<numSize; j++ {
							b, e := to_type(words[currWord][0], prop.Type)
							if e!=nil {
								return e
							}
							l = appendBytes(l, b)
							currWord++
						}
						prop.Data[i] = l
					} else {
						b, e := to_type(words[currWord][0], prop.Type)
						if e!=nil {
							fmt.Println("")
							prop.print()
							return e
						}
						prop.Data[i] = b
					}
				}
			}
		}
	}
	return nil
}
