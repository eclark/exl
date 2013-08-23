// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

import (
	"io"
	"bufio"
	"os"
	"bytes"
	"encoding/binary"
	"crypto/md5"
	"path"
	"fmt"
	"hash"
	"errors"
	//"log"
)

type Set struct {
	Id string
	Path string
	Main *Main
	Creator *Creator
	FileDesc map[string]*FileDesc
	Ifsc map[string]*Ifsc
	Recovery []interface{}
}

type File struct {
	FileDesc *FileDesc
	Ifsc *Ifsc
}

type MD5Error struct {
	FileSum []byte
	CalcSum []byte
}

func (m MD5Error) Error() string {
	return fmt.Sprintf("MD5 mismatch: %x != %x", m.FileSum, m.CalcSum)
}

var magicSeq = []byte{'P','A','R','2',0,'P','K','T'}

func OpenSet(file string) (set []*Set, err error) {
	rawRd, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer rawRd.Close()

	set, err = ReadFull(rawRd)
	if err != nil {
		return
	}

	for _, v := range set {
		v.Path, _ = path.Split(file)
	}

	return
}


func ReadFull(rawRd io.Reader)(set []*Set, err error) {
	var rdr io.Reader
	rdr = bufio.NewReaderSize(rawRd, 1024*1024)

	sets := make(map[string]*Set)

	for {
		id, p, err := readPacket(rdr)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		set, ok := sets[id]
		if !ok {
			set = new(Set)
			set.Id = id
			sets[id] = set
			set.FileDesc = make(map[string]*FileDesc)
			set.Ifsc = make(map[string]*Ifsc)
		}

		switch b := p.(type) {
			case *Main:
				set.Main = b
			case *Creator:
				set.Creator = b
			case *FileDesc:
				set.FileDesc[b.Id] = b
			case *Ifsc:
				set.Ifsc[b.Id] = b
		}
	}

	set = make([]*Set, 0, len(sets))
	for _,v := range sets {
		set = append(set, v)
	}

	return
}

func readPacket(in io.Reader) (setid string, p interface{}, err error) {
	rawHeader := make([]byte,64)
	_, err = io.ReadFull(in, rawHeader)
	if err != nil {
		return
	}

	magicstr := rawHeader[0:8]
	length := binary.LittleEndian.Uint64(rawHeader[8:16])
	md5sum := rawHeader[16:32]
	setid = string(rawHeader[32:48])
	typ := rawHeader[48:64]

	if !bytes.Equal(magicstr, magicSeq) {
		return "", nil, errors.New("Magic number does not match, stream is not PAR2")
	}

	body := make([]byte,length - 64)
	_, err = io.ReadFull(in, body)
	if err != nil {
		return
	}

	h := md5.New()
	h.Write(rawHeader[32:])
	h.Write(body)
	calcsum := h.Sum([]byte{})
	if !bytes.Equal(md5sum, calcsum) {
		err = MD5Error{md5sum,calcsum}
		return
	}

	return setid, newbody(typ,body), nil
}

var loaders map[string]bodyloader
type bodyloader func (data []byte) interface{}

func register(typ []byte, loader bodyloader) {
	if loaders == nil {
		loaders = make(map[string]bodyloader)
	}

	loaders[string(typ)] = loader
}

func newbody(typ []byte, data []byte) interface{} {
	if f, ok := loaders[string(typ)]; ok {
		return f(data)
	}
	return nil
}

type Verifier struct {
	Good int
	Bad int
	FileId string
	slicesize int
	ifsc *Ifsc
	sn int
	md5 hash.Hash
	w int
}

func newVerifier(set *Set, fileid string) *Verifier {
	return &Verifier{0,0,fileid,int(set.Main.SliceSize),set.Ifsc[fileid],0,md5.New(),0}
}

func (v *Verifier) Write(p []byte) (n int, err error) {
	n = len(p)
	for len(p) > 0 {
		amt := len(p)
		if amt + v.w > v.slicesize {
			amt = v.slicesize - v.w
		}
		v.md5.Write(p[:amt])
		v.w = v.w + amt
		p = p[amt:]

		if v.w == v.slicesize {
			if bytes.Equal(v.ifsc.Slices[v.sn].Hash, v.md5.Sum([]byte{})) {
				v.Good++
			} else {
				v.Bad++
			}
			v.md5.Reset()
			v.w = 0
			v.sn++
		}

	}
	return
}

func (v *Verifier) Close() error {
	zeros := make([]byte, v.slicesize - v.w)
	v.Write(zeros)

	return nil
}


func (s *Set) Verify(fileid string) (good int, bad int, err error) {
	file, err := os.OpenFile(path.Join(s.Path,s.FileDesc[fileid].Name), os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer file.Close()

	return s.VerifyReader(fileid, file)
}

func (s *Set) VerifyReader(fileid string, in io.Reader) (good int, bad int, err error) {
	v := newVerifier(s, fileid)

	_, err = io.Copy(v, in)
	if err != nil {
		return
	}

	v.Close()

	good = v.Good
	bad = v.Bad

	return
}

