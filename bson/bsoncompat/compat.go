// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bsoncompat

import (
	"bytes"
	"github.com/eclark/exl/bson"
	"github.com/eclark/gomongo/mongo"
	"time"
)

type wrapped struct {
	e bson.Element
}

func Wrap(e bson.Element) mongo.BSON {
	return &wrapped{e}
}

func (w *wrapped) Kind() int {
	return int(bson.Typeof(w.e))
}

func (w *wrapped) Number() float64 {
	if d, ok := w.e.(*bson.Double); ok {
		return float64(*d)
	}
	return 0
}

func (w *wrapped) String() string {
	if s, ok := w.e.(*bson.String); ok {
		return string(*s)
	}
	return ""
}

func (w *wrapped) OID() []byte {
	if o, ok := w.e.(*bson.ObjectId); ok {
		return []byte(*o)
	}
	return make([]byte, 12)
}

func (w *wrapped) Bool() bool {
	if b, ok := w.e.(*bson.Boolean); ok {
		return bool(*b)
	}
	return false
}

func (w *wrapped) Date() *time.Time {
	if t, ok := w.e.(*bson.Time); ok {
		return time.SecondsToUTC(int64(*t))
	}
	return nil
}

func (w *wrapped) Regex() (string, string) {
	if r, ok := w.e.(*bson.Regex); ok {
		return r.Pattern, r.Options
	}
	return "", ""
}

func (w *wrapped) Int() int32 {
	if i, ok := w.e.(*bson.Int32); ok {
		return int32(*i)
	}
	return 0
}

func (w *wrapped) Long() int64 {
	if i, ok := w.e.(*bson.Int64); ok {
		return int64(*i)
	}
	return 0
}

func (w *wrapped) Get(s string) mongo.BSON {
	if o, ok := w.e.(*bson.Document); ok {
		for i, _ := range *o {
			if (*o)[i].Key == s {
				return Wrap((*o)[i].Element)
			}
		}
	}
	return mongo.Null
}

func (w *wrapped) Elem(i int) mongo.BSON {
	if a, ok := w.e.(*bson.ArrayDocument); ok {
		return Wrap(a.Document[i].Element)
	}
	return mongo.Null
}

func (w *wrapped) Len() int {
	switch v := w.e.(type) {
	case *bson.Document:
		return len(*v)
	case *bson.ArrayDocument:
		return len(v.Document)
	}
	return 0
}

func (w *wrapped) Binary() []byte {
	if d, ok := w.e.(*bson.Binary); ok {
		return d.Data
	}
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
