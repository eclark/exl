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
	bson_append_start_object( buf, "d" );
	bson_append_finish_object( buf );
	bson_from_buffer( b, buf );
	hexprint(b);
	bson_destroy( b );

	// ArrayDocument
	bson_buffer_init( buf );
	bson_append_start_array( buf, "a" );
	bson_append_finish_object( buf );
	bson_from_buffer( b, buf );
	hexprint(b);
	bson_destroy( b );

	// Binary
	char s[] = "1234567890abcdefghijklmnop";
	int s_len = 26;
	bson_buffer_init( buf );
	bson_append_binary( buf, "b", 0, s, s_len );
	bson_append_binary( buf, "b2", 2, s, s_len );
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

	// Code
	bson_buffer_init( buf );
	bson_append_code( buf, "c", "function(a, b) { return a + b }" );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Symbol
	bson_buffer_init( buf );
	bson_append_symbol( buf, "s", "sex" );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// ScopedCode
	bson inner[1];
	bson_buffer inner_buf[1];
	bson_buffer_init( inner_buf );
	bson_append_double( inner_buf, "a", 6 );
	bson_append_double( inner_buf, "b", 4 );
	bson_from_buffer( inner, inner_buf );
	bson_buffer_init( buf );
	bson_append_code_w_scope( buf, "sc", "a+b", inner );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );
	bson_destroy( inner );
	
	// Int32
	bson_buffer_init( buf );
	bson_append_int( buf, "i", 31337 );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Timestamp
	bson_timestamp_t ts[1];
	bson_buffer_init( buf );
	bson_append_timestamp( buf, "t", ts );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );

	// Int64
	bson_buffer_init( buf );
	bson_append_long( buf, "i", 31337L );
	bson_from_buffer( b, buf );
	hexprint( b );
	bson_destroy( b );
}
