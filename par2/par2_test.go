// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

import (
	"os"
	"bytes"
	"testing"
)

func toasciichar(rune int) int {
	if rune > 0x1f && rune < 0x7f {
		return rune
	}
	return 0x20
}

func toascii(in []byte) string {
	return string(bytes.Map(toasciichar, in))
}

func TestOpen(t *testing.T) {
	filename := "test.par2"
	if len(os.Args) > 2 {
		filename = os.Args[2]
	}

	sets, err := OpenSet(filename)
	if err != nil {
		t.Fatal(err)
	}

	for _, set := range sets {
		var slicesize = set.Main.SliceSize
		t.Logf("%s %x %d", filename, set.Id, slicesize)

		for _, fileid := range set.Main.RecoveryIds {
			filedesc := set.FileDesc[fileid]
			t.Logf(" %s %x %x", filedesc.Name, filedesc.Id, filedesc.Hash)

			good, bad, err := set.Verify(fileid)
			if err != nil {
				t.Logf("%s %s", filedesc.Name, err)
			} else {
				t.Logf("%s %d %d", filedesc.Name, good, bad)
			}
		}
	}
}

func TestRead(t *testing.T) {
	filename := "test.par2"
	if len(os.Args) > 2 {
		filename = os.Args[2]
	}

	file, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for {
		_, p, err := readPacket(file)
		if err != nil {
			if err == os.EOF {
				break
			}
			if v, ok := err.(MD5Error); ok {

				t.Log(err)
				t.Log(v.FileSum)
				t.Fatal(v.CalcSum)
			}

			t.Fatal(err)
		}
		continue

		switch v := p.(type) {
		case *FileDesc:
			t.Logf("FileDesc %x %s %x", v.Id, v.Name, v.Hash)
		case *Ifsc:
			t.Logf("Ifsc     %x %d", v.Id, len(v.Slices))
		case *Main:
			t.Logf("Main     C:%d SS:%d", v.FileCount, v.SliceSize)
		case *Creator:
			t.Logf("Creator  %s", v.Message)
		default:
			t.Logf("Unknown")
		}
	}

}
