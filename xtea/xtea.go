/*
Copyright (c) 2012, Logan J. Drews

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

package xtea

import (
	"encoding/binary"
	"strconv"
)

var end binary.ByteOrder = binary.BigEndian

type xteaCipher struct {
	key  []byte
	keys []uint32
}

type KeySizeError int

func (k KeySizeError) Error() string {
	return "xtea: invalid key size " + strconv.Itoa(int(k))
}

func NewCipher(key []byte) (*xteaCipher, error) {
	if len(key) != 16 {
		return nil, KeySizeError(len(key))
	}

	cipher := new(xteaCipher)
	cipher.key = key
	cipher.keys = make([]uint32, 4)

	for i := 0; i < 4; i++ {
		cipher.keys[i] = end.Uint32(key[i*4:])
	}

	return cipher, nil
}

func (c *xteaCipher) BlockSize() int {
	return 8
}

func (c *xteaCipher) Encrypt(dst, src []byte) {
	var (
		v0, v1 uint32 = end.Uint32(src), end.Uint32(src[4:])
		sum    uint32 = 0
		delta  uint32 = 0x9E3779B9
	)

	for i := 0; i < 32; i++ {
		v0 += (((v1 << 4) ^ (v1 >> 5)) + v1) ^ (sum + c.keys[sum&3])
		sum += delta
		v1 += (((v0 << 4) ^ (v0 >> 5)) + v0) ^ (sum + c.keys[(sum>>11)&3])
	}

	end.PutUint32(dst, v0)
	end.PutUint32(dst[4:], v1)
}

func (c *xteaCipher) Clear() {
	for i := 0; i < len(c.key); i++ {
		c.key[i] = 0
	}

	for i := 0; i < len(c.keys); i++ {
		c.key[i] = 0
	}
}

func (c *xteaCipher) Decrypt(dst, src []byte) {
	var (
		v0, v1 uint32 = end.Uint32(src[0:4]), end.Uint32(src[4:8])
		delta  uint32 = 0x9E3779B9
		sum    uint32 = delta << 5
	)

	for i := 0; i < 32; i++ {
		v1 -= ((v0<<4 ^ v0>>5) + v0) ^ (sum + c.keys[(sum>>11)&3])
		sum -= delta
		v0 -= ((v1<<4 ^ v1>>5) + v1) ^ (sum + c.keys[sum&3])
	}

	end.PutUint32(dst, v0)
	end.PutUint32(dst[4:], v1)
}
