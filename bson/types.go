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

type BinSubtype byte

const (
	GenericType BinSubtype = iota
	FunctionType
	OldBinaryType
	UUIDType
	MD5Type  BinSubtype = 5
	UserType BinSubtype = 0x80
)

func read(r io.Reader, data ...interface{}) (n int64, err os.Error) {
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

func write(w io.Writer, data ...interface{}) (n int64, err os.Error) {
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
	write(buf, byte(0))
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

	bufd := make([]byte, 0, tBytelen-4)
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

type ArrayDocument struct {
	Document
}

func (ad *ArrayDocument) Append(value Element) {
	ad.append(typeof(value), strconv.Itoa(len(ad.Document)), value)
}

type Double float64

func (d *Double) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, float64(*d))
}

func (d *Double) ReadFrom(r io.Reader) (n int64, err os.Error) {
	return read(r, d)
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
	b := make([]byte, l)
	m, err = read(r, b)
	n += m
	if err != nil {
		return
	}

	*s = String(b[:len(b)-1])
	return
}

type Binary struct {
	Subtype byte
	Data    []byte
}

func (b *Binary) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, len(b.Data), b.Subtype, b.Data)
}

func (b *Binary) ReadFrom(r io.Reader) (n int64, err os.Error) {
	//TODO handle subtype 0x02
	var l int32

	m, err := read(r, &l, &b.Subtype)
	n += m
	if err != nil {
		return
	}

	b.Data = make([]byte, l)

	m, err = read(r, b.Data)
	n += m
	return
}

type ObjectId []byte

func (o *ObjectId) String() string {
	return fmt.Sprintf("%x", *o)
}

func (o *ObjectId) WriteTo(w io.Writer) (n int64, err os.Error) {
	if *o == nil {
		*o = make([]byte, 12)
	}
	m, err := w.Write([]byte(*o))
	n = int64(m)
	return
}

func (o *ObjectId) ReadFrom(r io.Reader) (n int64, err os.Error) {
	if *o == nil {
		*o = make([]byte, 12)
	}
	m, err := io.ReadFull(r, []byte(*o))
	n = int64(m)
	return
}

type Boolean bool

func (b *Boolean) WriteTo(w io.Writer) (n int64, err os.Error) {
	var v byte
	if bool(*b) == true {
		v = 1
	}
	return write(w, v)
}

func (b *Boolean) ReadFrom(r io.Reader) (n int64, err os.Error) {
	var v byte
	n, err = read(r, &v)
	if err != nil {
		return
	}

	switch v {
	case 0:
		*b = false
	case 1:
		*b = true
	default:
		err = os.NewError("bad boolean code")
	}

	return
}

type Time int64

func (t *Time) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, *t)
}

func (t *Time) ReadFrom(r io.Reader) (n int64, err os.Error) {
	return read(r, t)
}

type Null struct{}

func (x *Null) WriteTo(w io.Writer) (n int64, err os.Error) {
	return
}

func (x *Null) ReadFrom(r io.Reader) (n int64, err os.Error) {
	return
}

type Regex struct {
	Pattern string
	Options string
}

func (re *Regex) WriteTo(w io.Writer) (n int64, err os.Error) {
	return
}

func (re *Regex) ReadFrom(r io.Reader) (n int64, err os.Error) {
	return
}

type Code struct {
	String
}

type Symbol struct {
	String
}

func newElement(typ fieldType) (e Element) {
	switch typ {
	case DoubleType:
		e = new(Double)
	case StringType:
		e = new(String)
	case DocumentType:
		e = new(Document)
	case ArrayDocumentType:
		e = new(ArrayDocument)
	case BinaryType:
		e = new(Binary)
	case ObjectIdType:
		e = new(ObjectId)
	case BooleanType:
		e = new(Boolean)
	case NullType:
		e = new(Null)
	case RegexType:
		e = new(Regex)
	case dbPointerType:
		panic("db pointers are not supported")
	case CodeType:
		e = new(Code)
	case SymbolType:
		e = new(Symbol)
	default:
		panic("unknown type" + strconv.Itoa(int(typ)))
	}
	return
}

func typeof(e Element) (f fieldType) {
	switch e.(type) {
	case *Double:
		f = DoubleType
	case *String:
		f = StringType
	case *Document:
		f = DocumentType
	case *ArrayDocument:
		f = ArrayDocumentType
	case *Binary:
		f = BinaryType
	case *ObjectId:
		f = ObjectIdType
	case *Boolean:
		f = BooleanType
	case *Null:
		f = NullType
	case *Regex:
		f = RegexType
	//TODO something with dbPointer
	case *Code:
		f = CodeType
	case *Symbol:
		f = SymbolType
	default:
		panic("unknown element")
	}
	return
}
