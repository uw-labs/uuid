package uuid

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"
)

type UUID [16]byte

// Parse parses a string to a UUID object. The input must be the standard 36
// char format with '-' chars separators in the appropriate places.
func Parse(s string) (UUID, error) {
	var uuid UUID
	if len(s) != 36 {
		return uuid, fmt.Errorf("Expected length 36, got %d", len(s))
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return uuid, errors.New("invalid UUID format")
	}
	s = strings.ReplaceAll(s, "-", "")
	dec, err := hex.DecodeString(s)
	if err != nil {
		return uuid, err
	}
	copy(uuid[:], dec)
	return uuid, nil
}

// MustParse parses a string to a UUID object, and will panic if it fails.
func MustParse(s string) UUID {
	uuid, err := Parse(s)
	if err != nil {
		panic(`uuid: Parse(` + s + `): ` + err.Error())
	}
	return uuid
}

// NewMD5Gen returns a new UUID generator for type 3 uuids using the given
// namespace.
func NewMD5Gen(namespace UUID) *Generator {
	return &Generator{md5.New(), namespace[:], make([]byte, 16), 3}
}

// NewSHA1Gen returns a new UUID generator for type 5 uuids using the given
// namespace.
func NewSHA1Gen(namespace UUID) *Generator {
	return &Generator{sha1.New(), namespace[:], make([]byte, 20), 5}
}

// Generator is a generator of hash based (type 3 + 5) UUIDs.
type Generator struct {
	h         hash.Hash
	namespace []byte
	buf       []byte
	version   int
}

// Generate generates a UUID for the given input, writing the result in the
// given target object
func (hg *Generator) Generate(target *UUID, input []byte) {
	hg.h.Reset()
	hg.h.Write(hg.namespace[:])
	hg.h.Write(input)
	hg.buf = hg.h.Sum(hg.buf[0:0])
	copy((*target)[:], hg.buf)
	target[6] = (target[6] & 0x0f) | uint8((hg.version&0xf)<<4)
	target[8] = (target[8] & 0x3f) | 0x80 // RFC 4122 variant
}

// String returns the UUID in standard UUID string format.
func (uuid UUID) String() string {
	formatted := [36]byte{
		hextable[uuid[0]>>4], hextable[uuid[0]&0x0f],
		hextable[uuid[1]>>4], hextable[uuid[1]&0x0f],
		hextable[uuid[2]>>4], hextable[uuid[2]&0x0f],
		hextable[uuid[3]>>4], hextable[uuid[3]&0x0f],
		'-',
		hextable[uuid[4]>>4], hextable[uuid[4]&0x0f],
		hextable[uuid[5]>>4], hextable[uuid[5]&0x0f],
		'-',
		hextable[uuid[6]>>4], hextable[uuid[6]&0x0f],
		hextable[uuid[7]>>4], hextable[uuid[7]&0x0f],
		'-',
		hextable[uuid[8]>>4], hextable[uuid[8]&0x0f],
		hextable[uuid[9]>>4], hextable[uuid[9]&0x0f],
		'-',
		hextable[uuid[10]>>4], hextable[uuid[10]&0x0f],
		hextable[uuid[11]>>4], hextable[uuid[11]&0x0f],
		hextable[uuid[12]>>4], hextable[uuid[12]&0x0f],
		hextable[uuid[13]>>4], hextable[uuid[13]&0x0f],
		hextable[uuid[14]>>4], hextable[uuid[14]&0x0f],
		hextable[uuid[15]>>4], hextable[uuid[15]&0x0f],
	}

	return string(formatted[:])
}

// AppendFormatted appends the standard hex string representation of this UUID
// to a byte slice and return the result.
func (uuid UUID) AppendFormatted(buf []byte) []byte {
	formatted := [36]byte{
		hextable[uuid[0]>>4], hextable[uuid[0]&0x0f],
		hextable[uuid[1]>>4], hextable[uuid[1]&0x0f],
		hextable[uuid[2]>>4], hextable[uuid[2]&0x0f],
		hextable[uuid[3]>>4], hextable[uuid[3]&0x0f],
		'-',
		hextable[uuid[4]>>4], hextable[uuid[4]&0x0f],
		hextable[uuid[5]>>4], hextable[uuid[5]&0x0f],
		'-',
		hextable[uuid[6]>>4], hextable[uuid[6]&0x0f],
		hextable[uuid[7]>>4], hextable[uuid[7]&0x0f],
		'-',
		hextable[uuid[8]>>4], hextable[uuid[8]&0x0f],
		hextable[uuid[9]>>4], hextable[uuid[9]&0x0f],
		'-',
		hextable[uuid[10]>>4], hextable[uuid[10]&0x0f],
		hextable[uuid[11]>>4], hextable[uuid[11]&0x0f],
		hextable[uuid[12]>>4], hextable[uuid[12]&0x0f],
		hextable[uuid[13]>>4], hextable[uuid[13]&0x0f],
		hextable[uuid[14]>>4], hextable[uuid[14]&0x0f],
		hextable[uuid[15]>>4], hextable[uuid[15]&0x0f],
	}
	return append(buf, formatted[:]...)
}

const hextable = "0123456789abcdef"
