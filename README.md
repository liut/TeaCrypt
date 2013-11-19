# TeaCrypt
This is a library that implements the several of the TEA family of encryption algorithms in the Go language. Each algorithm in this library conforms to the crypto/cipher/Block interface.

## TEA
TEA is implemented in the teacrypt/tea package.

Ciphers are initialized such:
`cipher := tea.NewCipher(key)`

## XTEA
XTEA is implemented in the teacrypt/xtea package.

Ciphers are initialized such:
`cipher := xtea.NewCipher(key)`

## XXTEA
XXTEA is implemented in the teacrypt/xxtea package. Unlike TEA and XTEA, which are fixed block width ciphers, XXTEA can operate on block sizes of variables length (minimum 64-bits). In addition to taking a key, initializing an XXTEA cipher takes a block size.

Ciphers are initialized such:
`cipher := xxtea.NewCipher(key, size)`

# Install
## prepare
* `mkdir ~/gocode/src && cd ~/gocode/src`
* `git clone https://github.com/liut/TeaCrypt`

To test, run `GOPATH=~/gocode go test TeaCrypt/xxtea` in the root directory.

To install, run `GOPATH=~/gocode go install TeaCrypt/xxtea` in the root directory.

# Compatibility
Due to the lack of official test vectors for these algorithms, I cannot guarantee these algorithms are fully compatible with all other imlementations.

In the documents that introduced these algorithms, no official byte order is given. However, these implementations use the "encoding/binary" package to interpret the byte arrays. By default, they use binary.BigEndian, but can be changed in source to binary.LittleEndian or any other system that conforms to binary.ByteOrder interface.

# License
The license for each file can be found at the top of the source code. Currently, all the files are licensed under the permissive ISC License and is freely usable for all purposes.

