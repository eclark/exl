// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

type fieldType byte

const (
	DoubleType fieldType = iota + 1
	StringType
	DocumentType
	ArrayDocumentType
	BinaryType
	undefType
	ObjectIdType
	BooleanType
	TimeType
	NullType
	RegexType
	dbPointerType
	CodeType
	SymbolType
	ScopedCodeType
	Int32Type
	TimestampType
	Int64Type
	MinType
	MaxType
)

func read(r io.Reader, data ... interface{}) (n int64, err os.Error) {
	for _, d := range data {
		switch dt := reflect.NewValue(d).(type) {
		case *reflect.PtrValue:
			n += int64(binary.TotalSize(dt.Elem()))
		case *reflect.SliceValue:
			n += int64(binary.TotalSize(dt))
		default:
			panic("decode type error")
		}
		err = binary.Read(r, binary.LittleEndian, d)
		if err != nil {
			return
		}
	}
	return
}

func write(w io.Writer, data ... interface{}) (n int64, err os.Error) {
	for _, d := range data {
		n += int64(binary.TotalSize(reflect.NewValue(d)))
		err = binary.Write(w, binary.LittleEndian, d)
		if err != nil {
			return
		}
	}
	return
}

type Element interface {
	WriteTo(io.Writer) (int64, os.Error)
	ReadFrom(io.Reader) (int64, os.Error)
}

type dElement struct {
	Element
	typ fieldType
	key string
}

type Document []dElement

func (d *Document) Append(key string, value Element) {
	d.append(typeof(value), key, value)
}

func (d *Document) append(typ fieldType, key string, e Element) {
	*d = append([]dElement(*d), dElement{e, typ, key})
}

func (d *Document) WriteTo(w io.Writer) (n int64, err os.Error) {
	buf := bytes.NewBuffer(nil)
	for _, de := range []dElement(*d) {
		write(buf, de.typ, []byte(de.key), byte(0))
		m, err := de.WriteTo(buf)
		n += m
		if err != nil {
			return
		}
	}
	write(buf,byte(0))
	write(w, int32(buf.Len()+4))
	io.Copy(w, buf)
	return
}

func (d *Document) ReadFrom(r io.Reader) (n int64, err os.Error) {
	var tBytelen int32
	m, err := read(r, &tBytelen)
	n += m
	if err != nil {
		return
	}

	bufd := make([]byte, 0, tBytelen - 4)
	buf := bytes.NewBuffer(bufd)

	m, err = buf.ReadFrom(r)

	for {
		if buf.Len() == 1 {
			break
		}
		var typ fieldType
		m, err = read(buf, &typ)
		n += m
		if err != nil {
			if err == os.EOF {
				break
			} else {
				return
			}
		}

		elem := newElement(typ)

		nameb, err := buf.ReadBytes(0)
		if err != nil {
			return
		}
		n += int64(len(nameb))

		key := string(nameb[:len(nameb)-1])

		m, err = elem.ReadFrom(buf)
		n += m
		if err != nil {
			return
		}

		d.append(typ, key, elem)
	}

	return
}

type String string

func (s *String) String() string {
	return string(*s)
}

func (s *String) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, int32(len(string(*s))+1), []byte(string(*s)), byte(0))
}

func (s *String) ReadFrom(r io.Reader) (n int64, err os.Error) {
	var l int32
	m, err := read(r, &l)
	n += m
	if err != nil {
		return
	}
	b := make([]byte,l)
	m, err = read(r, b)
	n += m
	if err != nil {
		return
	}

	*s = String(b[:len(b)-1])
	return
}

type ObjectId []byte

func (o *ObjectId) String() string {
	return fmt.Sprintf("%x", *o)
}

func (o *ObjectId) WriteTo(w io.Writer) (n int64, err os.Error) {
	if *o == nil {
		*o = make([]byte,12)
	}
	m, err := w.Write([]byte(*o))
	n = int64(m)
	return
}

func (o *ObjectId) ReadFrom(r io.Reader) (n int64, err os.Error) {
	if *o == nil {
		*o = make([]byte,12)
	}
	m, err := io.ReadFull(r, []byte(*o))
	n = int64(m)
	return
}

func newElement(typ fieldType) (e Element) {
	switch typ {
	case StringType:
		e = new(String)
	case ObjectIdType:
		e = new(ObjectId)
	default:
		panic("unknown type" + strconv.Itoa(int(typ)))
	}
	return
}

func typeof(e Element) (f fieldType) {
	switch e.(type) {
		case *String:
			f = StringType
		case *ObjectId:
			f = ObjectIdType
		default:
			panic("unknown element")
	}
	return
}
