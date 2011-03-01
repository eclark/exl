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

	if n, err := doc.WriteTo(buf); err != nil || n != 34 {
		t.Errorf("doc.WriteTo(buf) = (%d, %v), want (%d, %v)", n, err, 34, nil)
	}

	gen := buf.Bytes()
	if !bytes.Equal(testData[5], gen) {
		t.Errorf("Document(%v).WriteTo(buf) != Binary from %s\nhave: %x\nwant: %x", doc, testFile, gen, testData[5])
	}
}
