// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

import (
	"encoding/binary"
	"strings"
)

func init() {
	register(fileDescType,newFileDesc)
}

var fileDescType = []byte{'P','A','R',' ','2','.','0',0,'F','i','l','e','D','e','s','c'}
type FileDesc struct {
	Id string
	Hash []byte
	StartHash []byte
	Length uint64
	Name string
}

func newFileDesc(data []byte) interface{} {
	f := new(FileDesc)
	f.Hash = make([]byte,16)
	f.StartHash = make([]byte,16)

	f.Id = string(data[:16])
	copy(f.Hash,data[16:32])
	copy(f.StartHash,data[32:48])
	f.Length = binary.LittleEndian.Uint64(data[48:56])

	f.Name = strings.TrimRightFunc(string(data[56:]), func(r int) bool {
		return r == 0
	})

	return f
}

