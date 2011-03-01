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
	Key string
}

type Document []dElement

func (d *Document) Append(key string, value Element) {
	d.append(Typeof(value), key, value)
}

func (d *Document) append(typ fieldType, key string, e Element) {
	*d = append([]dElement(*d), dElement{e, typ, key})
}

func (d *Document) WriteTo(w io.Writer) (n int64, err os.Error) {
	buf := bytes.NewBuffer(nil)
	for _, de := range []dElement(*d) {
		write(buf, de.typ, []byte(de.Key), byte(0))
		_, err := de.WriteTo(buf)
		if err != nil {
			return
		}
	}
	_, err = write(buf, byte(0))
	if err != nil {
		return
	}
	n, err = write(w, int32(buf.Len()+4))
	if err != nil {
		return
	}
	m, err := io.Copy(w, buf)
	n += m
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
	ad.append(Typeof(value), strconv.Itoa(len(ad.Document)), value)
}

type Double float64

func (d *Double) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, *d)
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
	return write(w, re.Pattern, 0, re.Options, 0)
}

func (re *Regex) ReadFrom(r io.Reader) (n int64, err os.Error) {
	var v byte
	b := make([]byte, 0, 16)

	var pe bool
	for {
		m, err := read(r, &v)
		n += m
		if err != nil {
			return
		}

		if v == 0 {
			if !pe {
				re.Pattern = string(b)
				b = b[:0]
				pe = true
			} else {
				re.Options = string(b)
				b = b[:0]
				return
			}
		} else {
			b = append(b, v)
		}
	}

	return
}

type Code struct {
	String
}

type Symbol struct {
	String
}

type ScopedCode struct {
	Code  *Code
	Scope *Document
}

func (sc *ScopedCode) WriteTo(w io.Writer) (n int64, err os.Error) {
	buf := bytes.NewBuffer(nil)
	_, err = sc.Code.WriteTo(buf)
	if err != nil {
		return
	}
	_, err = sc.Scope.WriteTo(buf)
	if err != nil {
		return
	}

	var l int32 = int32(buf.Len()) + 4

	return write(w, l, buf.Bytes())
}

func (sc *ScopedCode) ReadFrom(r io.Reader) (n int64, err os.Error) {
	var l int32
	m, err := read(r, &l)
	n += m
	if err != nil {
		return
	}
	lr := io.LimitReader(r, int64(l-4))

	sc.Code = new(Code)
	m, err = sc.Code.ReadFrom(lr)
	n += m
	if err != nil {
		return
	}
	sc.Scope = new(Document)
	m, err = sc.Scope.ReadFrom(lr)
	n += m
	return
}

type Int32 int32

func (i *Int32) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, *i)
}

func (i *Int32) ReadFrom(r io.Reader) (n int64, err os.Error) {
	return read(r, i)
}

type Timestamp struct {
	Int64
}
type Int64 int64

func (i *Int64) WriteTo(w io.Writer) (n int64, err os.Error) {
	return write(w, *i)
}

func (i *Int64) ReadFrom(r io.Reader) (n int64, err os.Error) {
	return read(r, i)
}

type Min struct {
	Null
}
type Max struct {
	Null
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
	case ScopedCodeType:
		e = new(ScopedCode)
	case Int32Type:
		e = new(Int32)
	case TimestampType:
		e = new(Timestamp)
	case Int64Type:
		e = new(Int64)
	case MinType:
		e = new(Min)
	case MaxType:
		e = new(Max)
	default:
		panic("unknown type" + strconv.Itoa(int(typ)))
	}
	return
}

func Typeof(e Element) (f fieldType) {
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
	case *ScopedCode:
		f = ScopedCodeType
	case *Int32:
		f = Int32Type
	case *Timestamp:
		f = TimestampType
	case *Int64:
		f = Int64Type
	case *Min:
		f = MinType
	case *Max:
		f = MaxType
	default:
		panic("unknown element")
	}
	return
}
