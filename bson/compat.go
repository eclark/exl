// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bson

import (
	"bytes"
	"github.com/eclark/gomongo/mongo"
	"time"
)

type wrapped struct {
	e Element
}

func Wrap(e Element) mongo.BSON {
	return &wrapped{e}
}

func (w *wrapped) Kind() int {
	return int(typeof(w.e))
}

func (w *wrapped) Number() float64 {
	if d, ok := w.e.(*Double); ok {
		return float64(*d)
	}
	return 0
}

func (w *wrapped) String() string {
	if s, ok := w.e.(*String); ok {
		return string(*s)
	}
	return ""
}

func (w *wrapped) OID() []byte {
	if o, ok := w.e.(*ObjectId); ok {
		return []byte(*o)
	}
	return make([]byte, 12)
}

func (w *wrapped) Bool() bool {
	if b, ok := w.e.(*Boolean); ok {
		return bool(*b)
	}
	return false
}

func (w *wrapped) Date() *time.Time {
	if t, ok := w.e.(*Time); ok {
		return time.SecondsToUTC(int64(*t))
	}
	return nil
}

func (w *wrapped) Regex() (string, string) {
	return "", ""
}

func (w *wrapped) Int() int32 {
	return 0
}

func (w *wrapped) Long() int64 {
	return 0
}

func (w *wrapped) Get(s string) mongo.BSON {
	return mongo.Null
}

func (w *wrapped) Elem(i int) mongo.BSON {
	return mongo.Null
}

func (w *wrapped) Len() int {
	return 0
}

func (w *wrapped) Binary() []byte {
	return nil
}

func (w *wrapped) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	_, err := w.e.WriteTo(buf)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
