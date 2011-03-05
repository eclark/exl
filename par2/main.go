// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

import (
	"encoding/binary"
)

func init() {
	register(mainType,newMain)
}

var mainType = []byte{'P','A','R',' ','2','.','0',0,'M','a','i','n',0,0,0,0}
type Main struct {
	SliceSize uint64
	FileCount uint32
	RecoveryIds []string
	NonRecoveryIds []string
}

func newMain(data []byte) interface{} {
	m := new(Main)

	m.SliceSize = binary.LittleEndian.Uint64(data[:8])
	m.FileCount = binary.LittleEndian.Uint32(data[8:12])

	m.RecoveryIds = make([]string, m.FileCount)

	var start int
	for i, _ := range m.RecoveryIds {
		start = 12+(i*16)
		m.RecoveryIds[i] = string(data[start:start+16])
	}

	nrs := 12+(16*m.FileCount)
	nrc := (uint32(len(data)) - nrs) / 16;

  _ = nrc

	//implement non-recovery id

	return m
}


