// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

func init() {
	register(creatorType,newCreator)
}

var creatorType = []byte{'P','A','R',' ','2','.','0',0,'C','r','e','a','t','o','r',0}
type Creator struct {
	Message string
}

func newCreator(data []byte) interface{} {
	c := new(Creator)

	c.Message = string(data)

	return c
}


