package xid

import "unsafe"

// FromBytes convert the byte array representation of `ID` back to `ID`
func FromBytes(b []byte) (ID, error) {
	var id ID

	if len(b) != rawLen {
		return id, ErrInvalidID
	}

	copy(id[:], b)

	return id, nil
}

// FromString reads an ID from its string representation
func FromString(id string) (ID, error) {
	i := &ID{}
	err := i.UnmarshalText([]byte(id))

	return *i, err
}

// Bytes returns the byte array representation of `ID`
func (id ID) Bytes() []byte {
	return id[:]
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id ID) String() string {
	text := make([]byte, encodedLen)
	encode(text, id[:])

	return *(*string)(unsafe.Pointer(&text))
}

// MarshalText implements encoding/text TextMarshaler interface
func (id ID) MarshalText() ([]byte, error) {
	text := make([]byte, encodedLen)
	encode(text, id[:])

	return text, nil
}

// MarshalJSON implements encoding/json Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	if id.IsNil() {
		return []byte("null"), nil
	}

	text := make([]byte, encodedLen+2)
	encode(text[1:encodedLen+1], id[:])
	text[0], text[encodedLen+1] = '"', '"'

	return text, nil
}

// encode by unrolling the stdlib base32 algorithm + removing all safe checks
func encode(dst, id []byte) {
	_ = dst[19]
	_ = id[11]

	dst[19] = encoding[(id[11]<<4)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[17] = encoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[15] = encoding[id[9]&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[11] = encoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[9] = encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[7] = encoding[id[4]&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[3] = encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[1] = encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[0] = encoding[id[0]>>3]
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *ID) UnmarshalText(text []byte) error {
	if len(text) != encodedLen {
		return ErrInvalidID
	}

	for _, c := range text {
		if dec[c] == 0xFF {
			return ErrInvalidID
		}
	}

	decode(id, text)

	return nil
}

// UnmarshalJSON implements encoding/json Unmarshaler interface
func (id *ID) UnmarshalJSON(b []byte) error {
	if s := string(b); s == "null" {
		*id = nilID

		return nil
	}

	return id.UnmarshalText(b[1 : len(b)-1])
}

// decode by unrolling the stdlib base32 algorithm + removing all safe checks
func decode(id *ID, src []byte) {
	_ = src[19]
	_ = id[11]

	id[11] = dec[src[17]]<<6 | dec[src[18]]<<1 | dec[src[19]]>>4
	id[10] = dec[src[16]]<<3 | dec[src[17]]>>2
	id[9] = dec[src[14]]<<5 | dec[src[15]]
	id[8] = dec[src[12]]<<7 | dec[src[13]]<<2 | dec[src[14]]>>3
	id[7] = dec[src[11]]<<4 | dec[src[12]]>>1
	id[6] = dec[src[9]]<<6 | dec[src[10]]<<1 | dec[src[11]]>>4
	id[5] = dec[src[8]]<<3 | dec[src[9]]>>2
	id[4] = dec[src[6]]<<5 | dec[src[7]]
	id[3] = dec[src[4]]<<7 | dec[src[5]]<<2 | dec[src[6]]>>3
	id[2] = dec[src[3]]<<4 | dec[src[4]]>>1
	id[1] = dec[src[1]]<<6 | dec[src[2]]<<1 | dec[src[3]]>>4
	id[0] = dec[src[0]]<<3 | dec[src[1]]>>2
}
