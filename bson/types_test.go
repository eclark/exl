// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bson

import (
	"bytes"
	"fmt"
	"github.com/eclark/gomongo/mongo"
	"testing"
)

func TestA(t *testing.T) {
	oid, _ := mongo.NewOID()
	bs, err := mongo.Marshal(map[string]interface{}{"root": "fs", "files_id": oid})
	if err != nil {
		t.Fatal(err)
	}

	mongob := bs.Bytes()
	fmt.Println(mongob)

	buf := bytes.NewBuffer(mongob)

	doc := new(Document)

	_, err = doc.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(doc)

	newb := bytes.NewBuffer(nil)
	doc.WriteTo(newb)

	fmt.Println(newb.Bytes())

	if !bytes.Equal(mongob, newb.Bytes()) {
		t.Fatal("unmatched")
	}
}

func TestB(t *testing.T) {
	doc := new(Document)
	oid, err := NewObjectId()
	if err != nil {
		t.Fatal(err)
	}
	str := String("fs")
	doc.Append("files_id", oid)
	doc.Append("root", &str)

	b := bytes.NewBuffer(nil)
	_, err = doc.WriteTo(b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(b.Bytes())
}

func TestDouble(t *testing.T) {
	doc := new(Document)
	f := Double(1.1243)
	doc.Append("num", &f)

	b := bytes.NewBuffer(nil)
	_, err := doc.WriteTo(b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(b.Bytes())
}
