// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

import (
	"encoding/binary"
)

func init() {
	register(ifscType,newIfsc)
}

var ifscType = []byte{'P','A','R',' ','2','.','0',0,'I','F','S','C',0,0,0,0}
type Ifsc struct {
	Id string
	Slices []Slice
}

type Slice struct {
	Hash []byte
	Crc uint32
}

func newIfsc(data []byte) interface{} {
	i := new(Ifsc)
	i.Id = string(data[:16])

	nslices := (len(data) - 16) / 20
	i.Slices = make([]Slice, nslices)

	window := data[16:]
	for j := 0; j < nslices; j++ {
		i.Slices[j].Hash = make([]byte,16)
		copy(i.Slices[j].Hash,window[:16])
		i.Slices[j].Crc = binary.LittleEndian.Uint32(window[16:20])
		window = window[20:]
	}

	return i
}


