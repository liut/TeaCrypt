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

package xxtea

import (
	"encoding/binary"
	"strconv"
	"os"
	//"fmt"
)

var end binary.ByteOrder = binary.BigEndian
//var end binary.ByteOrder = binary.LittleEndian

type xxteaCipher struct {
	key []byte
	keys []uint32
	size uint32
}

type KeySizeError int

func (k KeySizeError) String() string {
	return "xxtea: invalid key size " + strconv.Itoa(int(k))
}

func NewXXTea(key []byte, size uint32) (*xxteaCipher, os.Error) {
	if len(key) != 16 {
		return nil, KeySizeError(len(key))
	}

	cipher := new(xxteaCipher)
	cipher.key = key
	cipher.keys = make([]uint32, 4)
	cipher.size = size

	for i := 0; i < 4; i++ {
		cipher.keys[i] = end.Uint32(key[i * 4:])
	}

	return cipher, nil
}

func (c *xxteaCipher) BlockSize() int {
	return int(c.size)
}

func (c *xxteaCipher) Encrypt(dst, src []byte) {
	var (
		n              uint32 = c.size / 4
		words        []uint32 = make([]uint32, n)
		v0, v1         uint32
		sum            uint32 = 0
		delta          uint32 = 0x9E3779B9
		q              uint32 = 6 + 52 / n
		i, e, p        uint32
	)

	for i = 0; i < n; i++ {
		words[i] = end.Uint32(src[i*4:])
	}
	v0, v1 = words[0], words[n-1]

	/*for i = 0; i < n; i++ {
		fmt.Printf("%08X", words[i])
	}
	fmt.Printf("\n")*/

	for i = 0; i < q; i++ {
		sum += delta;
		e = (sum >> 2) & 3;
		for p = 0; p < n - 1; p++ {
			v0 = words[p + 1]
			words[p] += c.mx(v0, v1, sum, p, e)
			v1 = words[p]
		}
		v0 = words[0]
		words[n - 1] += c.mx(v0, v1, sum, p, e)
		v1 = words[n - 1]
	}

	for i = 0; i < n; i++ {
		end.PutUint32(dst[4*i:], words[i])
	}
}

func (c *xxteaCipher) Decrypt(dst, src []byte) {
	var (
		n              uint32 = c.size / 4
		words        []uint32 = make([]uint32, n)
		v0, v1         uint32
		q              uint32 = 6 + 52 / n
		delta          uint32 = 0x9E3779B9
		sum            uint32 = q * delta
		i, e, p        uint32
	)

	for i = 0; i < n; i++ {
		words[i] = end.Uint32(src[i*4:])
	}
	v0, v1 = words[0], words[n-1]

	for i = 0; i < q; i++ {
		e = (sum >> 2) & 3
		for p = n - 1; p > 0; p-- {
			v1 = words[p - 1]
			words[p] -= c.mx(v0, v1, sum, p, e)
			v0 = words[p]
		}
		v1 = words[n - 1]
		words[0] -= c.mx(v0, v1, sum, p, e)
		v0 = words[0]
		sum -= delta
	}

	for i = 0; i < n; i++ {
		end.PutUint32(dst[4*i:], words[i])
	}
}

func (c *xxteaCipher) mx(v0, v1, sum, p, e uint32) uint32 {
	var r uint32 = ((((v1 >> 5) ^ (v0 << 2)) + ((v0 >> 3) ^ (v1 << 4))) ^ ((sum ^ v0) + (c.keys[(p & 3) ^ e] ^ v1)))
	return r
}
