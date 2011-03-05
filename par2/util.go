// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package par2

import "hash/crc32"

/*
	Given a crc that was calculated from a block of data with a number of zeros at the end.
	Recalculate it such that the nzero zeroes were not present.
*/
func Crc32Drop(crc uint32, nzero int) uint32 {
	crc = ^crc

	for i := 0; i < (8 * nzero); i++ {
		if crc >> 31 & 1 == 1 {
			crc ^= crc32.IEEE
			crc <<= 1
			crc |= 1
		} else {
			crc <<= 1
		}
		crc ^= 0
	}

	return ^crc
}

/*
	Combine two crc where the second has a known length.  Returns a crc equal to
	what would have been calculated if the original input blocks for each were concatenated.

	**I don't understand how this works, ported directly from zlib
*/
func Crc32Combine(crc1 uint32, crc2 uint32, len2 int) uint32 {
	var row uint32
	even := newgf2matrix()
	odd := newgf2matrix()

	if len2 <= 0 {
		return crc1
	}

	odd[0] = crc32.IEEE
	row = 1
	for n := 1; n < len(odd); n++ {
		odd[n] = row
		row <<= 1
	}

	even.square(odd)
	odd.square(even)

	for {
		even.square(odd)
		if len2 & 1 == 1 {
			crc1 = even.times(crc1)
		}
		len2 >>= 1

		if len2 == 0 {
			break
		}

		odd.square(even)
		if len2 & 1 == 1 {
			crc1 = odd.times(crc1)
		}
		len2 >>= 1

		if len2 == 0 {
			break
		}
	}

	crc1 ^= crc2

	return crc1
}

type gf2matrix []uint32

func newgf2matrix() gf2matrix {
	m := make([]uint32, 32)

	return gf2matrix(m)
}

func (m gf2matrix) square(by gf2matrix) {
	for n, _ := range m {
		m[n] = by.times(by[n])
	}
}

func (m gf2matrix) times(vec uint32) uint32 {
	var sum uint32 = 0

	n := 0
	for vec > 0 {
		if vec & 1 == 1 {
			sum ^= m[n]
		}
		vec >>= 1
		n++
	}
	return sum
}
