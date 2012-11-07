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
}

type Element struct {
	Name       string
	Properties []*Property
	Size       int
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
	if file_io_error != nil {
		return file_io_error
	}
	r := bufio.NewReader(file)
	e := parse_header(p)
	if e != nil {
		return e
	}
	return nil
}

func (p *PLY) Read(elemName, propName string) (data [][]byte, e error) {
	switch p.FileType {
	case BinaryBigEndian:
		data, e = parse_binary_big_endian(p, elemName, propName)
	case BinaryLittleEndian:
		data, e = parse_binary_little_endian(p, elemName, propName)
	case Ascii:
		data, e = parse_ascii(p, elemName, propName)
	default:
		data = make([][]byte,0)
		e = errors.New("File type error")
	}
	return data, e
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

func to_type(type_name, data string) (b []byte, e error) {
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
		n, e = strconv.ParseInt(data, 0, 64)
		t := int64(n)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 8)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[5] || type_name == OldTypes[5]:
		u, e = strconv.ParseUint(data, 0, 8)
		t := uint8(u)
		if e != nil {
			return nil, e
		}
		b = make([]byte, 1)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[6] || type_name == OldTypes[6]:
		u, e = strconv.ParseUint(data, 0, 16)
		t := uint16(u)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 2)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[7] || type_name == OldTypes[7]:
		u, e = strconv.ParseUint(data, 0, 32)
		t := uint32(u)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 4)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[8] || type_name == OldTypes[8]:
		u, e = strconv.ParseUint(data, 0, 64)
		t := uint64(u)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 8)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[9] || type_name == OldTypes[9]:
		f, e = strconv.ParseFloat(data, 32)
		t := float32(f)
		if e != nil {
			return nil, e
		}
		b := make([]byte, 4)
		buf := bytes.NewBuffer(b)
		binary.Write(buf, binary.LittleEndian, &t)
		return b, nil
	case type_name == Types[10] || type_name == OldTypes[10]:
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
		} else if words[0][0] == "property" {
			cnt := 1
			currWord := words[cnt][0]
			prop := new(Property)
			if currWord == "list" {
				prop.IsList = true
				cnt++
				currWord = words[cnt][0]
				prop.ListSizeType = currWord
				cnt++
				currWord = words[cnt][0]
				fmt.Println(words, cnt)
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

func parse_binary_big_endian(p *PLY, elemName, propName string) ([][]byte, error) {
	data := make([][]byte,0)
	e := errors.New("Not yet implemented")
	return data, e
}

func parse_binary_little_endian(p *PLY, elemName, propName string) ([][]byte, error) {
	data := make([][]byte,0)
	e := errors.New("Not yet implemented")
	return data, e
}

func parse_ascii(p *PLY, elemName, propName string) ([][]byte, error) {
	r := p.reader
	data := make([][]byte,0)
	numMatcher, e := regexp.Compile("[\\+\\-]*([0-9]*)+\\.*[0-9]+")
	if e != nil {
		return data, e
	}
	var elem *Element
	var propName
	elem = nil
	for _, element := range p.Elements {
		if elemName == element.Name {
			elem = element
		}
	}
	if elem == nil {
		return data, errors.New("Property does not exist")
	}
	for _, property := range elem.Properties {
		if propName == property.Name {
			prop = property
		}
	}
	return data, nil
}
