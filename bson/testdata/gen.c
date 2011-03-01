// gcc -g --std=c99 -I ~/git/mongo-c-driver/src gen.c ~/git/mongo-c-driver/libbson.a -o gen 

// Copyright 2011 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdio.h>
#include "bson.h"

// Output hex so its human readable
void hexprint(bson* b) {
	unsigned char* s = b->data;
	unsigned char* e = s + bson_size(b);
	for (;s < e; s++) {
		printf("%.2x", *s);
	}
	printf("\n");
}

int main() {
	bson b[1];
	bson_buffer buf[1];

	// Empty
	bson_buffer_init( buf );
	bson_from_buffer( b, buf );
	hexprint(b);
	bson_destroy( b );

	// Double
	bson_buffer_init( buf );
	bson_append_double( buf, "d", 22.0/7.0 );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// String
	bson_buffer_init( buf );
	bson_append_string( buf, "s", "bcdefg" );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Document
	bson_buffer_init( buf );
	bson_from_buffer( b, buf );
	//TODO
	hexprint(b);
	bson_destroy( b );

	// ArrayDocument
	bson_buffer_init( buf );
	bson_from_buffer( b, buf );
	//TODO
	hexprint(b);
	bson_destroy( b );

	// Binary
	char s[] = "1234567890abcdefghijklmnop";
	int s_len = 26;
	bson_buffer_init( buf );
	bson_append_binary( buf, "b", 0, s, s_len );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// ObjectId
	char oid_str[] = "4d6d4cee9433e95b30cd38ec";
	bson_oid_t oid[1];
	bson_oid_from_string( oid, oid_str );
	bson_buffer_init( buf );
	bson_append_oid( buf, "o", oid );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Boolean
	bson_buffer_init( buf );
	bson_append_bool( buf, "b", 0 );
	bson_append_bool( buf, "c", 1 );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Time
	bson_buffer_init( buf );
	bson_append_time_t( buf, "t", (time_t)20 );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Null
	bson_buffer_init( buf );
	bson_append_null( buf, "n" );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Regex
	bson_buffer_init( buf );
	bson_append_regex( buf, "r", "[a-z]+", "i" );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );
}
