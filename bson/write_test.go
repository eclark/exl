// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bson

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"testing"
)

const testFile = "testdata/hexdata"

var testData [][]byte

func readHexdata(filename string) os.Error {
	var err os.Error
	file, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	brdr := bufio.NewReader(file)
	eof := false
	var line []byte
	for !eof {
		line, err = brdr.ReadBytes('\n')
		if err == nil {
			line = line[:len(line)-1] // chomp off line ending
		} else if err == os.EOF {
			eof = true
		} else {
			return err
		}

		dec := make([]byte, len(line)/2)

		_, err = hex.Decode(dec, line)
		if err != nil {
			return err
		}

		testData = append(testData, dec)
	}

	return nil
}

var buf *bytes.Buffer

func init() {
	err := readHexdata(testFile)
	if err != nil {
		panic(err)
	}

	buf = bytes.NewBuffer(nil)
}

func TestWriteEmpty(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	if n, err := doc.WriteTo(buf); err != nil || n != 5 {
		t.Errorf("doc.WriteTo(%v) = (%d, %v), want (%d, %v)", buf, n, err, 5, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[0], gen) {
		t.Errorf("Document(%v) does not match Empty from gen.c\nhave: %x\nwant: %x", doc, gen, testData[0])
	}
}

func TestWriteDouble(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	f := Double(22.0 / 7.0)
	doc.Append("d", &f)

	if n, err := doc.WriteTo(buf); err != nil || n != 16 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 16, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[1], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Double from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[1])
	}
}

func TestWriteString(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	s := String("bcdefg")
	doc.Append("s", &s)

	if n, err := doc.WriteTo(buf); err != nil || n != 19 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 19, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[2], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != String from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[2])
	}
}

func TestWriteDocument(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	doc2 := new(Document)
	doc.Append("d", doc2)

	if n, err := doc.WriteTo(buf); err != nil || n != 13 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 13, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[3], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Document from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[3])
	}
}

func TestWriteArrayDocument(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	doc2 := new(ArrayDocument)
	doc.Append("a", doc2)

	if n, err := doc.WriteTo(buf); err != nil || n != 13 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 13, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[4], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != ArrayDocument from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[4])
	}
}

func TestWriteBinary(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	bin := &Binary{Subtype: 0, Data: []byte("1234567890abcdefghijklmnop")}
	doc.Append("b", bin)
	bin = &Binary{Subtype: 2, Data: []byte("1234567890abcdefghijklmnop")}
	doc.Append("b2", bin)

	if n, err := doc.WriteTo(buf); err != nil || n != 78 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 78, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[5], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Binary from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[5])
	}
}

func TestWriteObjectId(t *testing.T) {
	buf.Reset()
	doc := new(Document)

	oid := new(ObjectId)
	oid.FromString("4d6d4cee9433e95b30cd38ec")
	doc.Append("o", oid)

	if n, err := doc.WriteTo(buf); err != nil || n != 20 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 20, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[6], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != ObjectId from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[6])
	}
}

func TestWriteBoolean(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	doc.Append("b", False())
	doc.Append("c", True())

	if n, err := doc.WriteTo(buf); err != nil || n != 13 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 13, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[7], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Boolean from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[7])
	}
}

func TestWriteTime(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	time := Time(20000)
	doc.Append("t", &time)

	if n, err := doc.WriteTo(buf); err != nil || n != 16 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 16, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[8], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Time from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[8])
	}
}

func TestWriteNull(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	doc.Append("n", new(Null))

	if n, err := doc.WriteTo(buf); err != nil || n != 8 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 8, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[9], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Null from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[9])
	}
}

func TestWriteRegex(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	re := &Regex{Pattern: "[a-z]+", Options: "i"}
	doc.Append("r", re)

	if n, err := doc.WriteTo(buf); err != nil || n != 17 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 17, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[10], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Regex from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[10])
	}
}

func TestWriteCode(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	doc.Append("c", &Code{String("function(a, b) { return a + b }")})

	if n, err := doc.WriteTo(buf); err != nil || n != 44 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 44, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[11], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Code from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[11])
	}
}

func TestWriteSymbol(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	doc.Append("s", &Symbol{String("sex")})

	if n, err := doc.WriteTo(buf); err != nil || n != 16 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 16, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[12], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Symbol from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[12])
	}
}

func TestWriteScopedCode(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	inner := new(Document)
	a := Double(6)
	b := Double(4)
	inner.Append("a", &a)
	inner.Append("b", &b)
	doc.Append("sc", &ScopedCode{Code: &Code{String("a+b")}, Scope: inner})

	if n, err := doc.WriteTo(buf); err != nil || n != 48 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 48, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[13], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != ScopedCode from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[13])
	}
}

func TestWriteInt32(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	i := Int32(31337)
	doc.Append("i", &i)

	if n, err := doc.WriteTo(buf); err != nil || n != 12 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 12, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[14], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Int32 from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[14])
	}
}

func TestWriteTimestamp(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	doc.Append("t", &Timestamp{Int64(0)})

	if n, err := doc.WriteTo(buf); err != nil || n != 16 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 16, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[15], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Timestamp from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[15])
	}
}

func TestWriteInt64(t *testing.T) {
	buf.Reset()
	doc := new(Document)
	i := Int64(31337)
	doc.Append("i", &i)

	if n, err := doc.WriteTo(buf); err != nil || n != 16 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 16, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[16], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Int64 from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[16])
	}
}
