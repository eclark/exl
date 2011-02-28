// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bson

import (
	"crypto/md5"
	"encoding/binary"
	"os"
	"sync"
)

var lock sync.Locker = &sync.Mutex{}
var accum int
var hosthash []byte

func init() {
	hostname, _ := os.Hostname()
	h := md5.New()

	h.Write([]byte(hostname))
	hosthash = h.Sum()
}

func NewObjectId() (Element, os.Error) {
	t64, _, err := os.Time()
	if err != nil {
		return nil, err
	}

	oid := make([]byte, 12)
	tmp := make([]byte, 4)

	binary.BigEndian.PutUint32(oid[0:4], uint32(t64))
	copy(oid[4:7], hosthash)

	binary.BigEndian.PutUint32(tmp, uint32(os.Getpid()))
	copy(oid[7:9], tmp[2:4])

	lock.Lock()
	accum++
	inc := accum
	lock.Unlock()

	binary.BigEndian.PutUint32(tmp, uint32(inc))
	copy(oid[9:12], tmp[1:4])

	x := ObjectId(oid)
	return &x, nil
}
