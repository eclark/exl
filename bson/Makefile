# Copyright 2011 Eric Clark. All rights reserved.
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# These targets are only useful for regenerating test data
# Check the first line of gen.c for the gcc options

testdata/gen: testdata/gen.c
	$(shell cd testdata; $(shell head -1 testdata/gen.c | cut -b 4-); cd ..)

testdata/hexdata: testdata/gen
	testdata/gen >testdata/hexdata

testdata: testdata/hexdata

include $(GOROOT)/src/Make.inc

TARG=github.com/eclark/exl/bson
GOFILES=\
	types.go\
	objectid.go\

include $(GOROOT)/src/Make.pkg
