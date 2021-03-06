// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package key

// DO NOT EDIT.
// FILE GENERATED BY BSONGEN.

import (
	"bytes"

	"github.com/youtube/vitess/go/bson"
	"github.com/youtube/vitess/go/bytes2"
)

// MarshalBson bson-encodes KeyRange.
func (keyRange *KeyRange) MarshalBson(buf *bytes2.ChunkedWriter, key string) {
	bson.EncodeOptionalPrefix(buf, bson.Object, key)
	lenWriter := bson.NewLenWriter(buf)

	keyRange.Start.MarshalBson(buf, "Start")
	keyRange.End.MarshalBson(buf, "End")

	lenWriter.Close()
}

// UnmarshalBson bson-decodes into KeyRange.
func (keyRange *KeyRange) UnmarshalBson(buf *bytes.Buffer, kind byte) {
	switch kind {
	case bson.EOO, bson.Object:
		// valid
	case bson.Null:
		return
	default:
		panic(bson.NewBsonError("unexpected kind %v for KeyRange", kind))
	}
	bson.Next(buf, 4)

	for kind := bson.NextByte(buf); kind != bson.EOO; kind = bson.NextByte(buf) {
		switch bson.ReadCString(buf) {
		case "Start":
			keyRange.Start.UnmarshalBson(buf, kind)
		case "End":
			keyRange.End.UnmarshalBson(buf, kind)
		default:
			bson.Skip(buf, kind)
		}
	}
}
