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

package tea

import (
	"encoding/binary"
	"strconv"
	"os"
)

type teaCipher struct {
	key []byte
}

type KeySizeError int

func (k KeySizeError) String() string {
	return "tea: invalid key size " + strconv.Itoa(int(k))
}

func NewTea(key []byte) (*teaCipher, os.Error) {
	if len(key) != 16 {
		return nil, KeySizeError(len(key))
	}

	cipher := new(teaCipher)
	cipher.key = key

	return cipher, nil
}

func (c *teaCipher) BlockSize() int {
	return 8
}

func (c *teaCipher) Encrypt(dst, src []byte) {
	var (
		end                   = binary.BigEndian
		v0, v1         uint32 = end.Uint32(src), end.Uint32(src[4:])
		sum            uint32 = 0
		delta          uint32 = 0x9E3779B9
		k0, k1, k2, k3 uint32 = end.Uint32(c.key[0:]),
			end.Uint32(c.key[4:]),
			end.Uint32(c.key[8:]),
			end.Uint32(c.key[12:])
	)

	for i := 0; i < 32; i++ {
		sum += delta
		v0 += ((v1 << 4) + k0) ^ (v1 + sum) ^ ((v1 >> 5) + k1)
		v1 += ((v0 << 4) + k2) ^ (v0 + sum) ^ ((v0 >> 5) + k3)
	}

	end.PutUint32(dst, v0)
	end.PutUint32(dst[4:], v1)
}

func (c *teaCipher) Decrypt(dst, src []byte) {
	var (
		end                   = binary.BigEndian
		v0, v1         uint32 = end.Uint32(src[0:4]), end.Uint32(src[4:8])
		delta          uint32 = 0x9E3779B9
		sum            uint32 = delta << 5
		k0, k1, k2, k3 uint32 = end.Uint32(c.key[0:]),
			end.Uint32(c.key[4:]),
			end.Uint32(c.key[8:]),
			end.Uint32(c.key[12:])
	)

	for i := 0; i < 32; i++ {
		v1 -= ((v0 << 4) + k2) ^ (v0 + sum) ^ ((v0 >> 5) + k3)
		v0 -= ((v1 << 4) + k0) ^ (v1 + sum) ^ ((v1 >> 5) + k1)
		sum -= delta
	}

	end.PutUint32(dst, v0)
	end.PutUint32(dst[4:], v1)
}
