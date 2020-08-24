package main

import "crypto/sha256"

var characterSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")


func RandomNumber(seed uint64) uint64 {
	seed ^= seed << 21
	seed ^= seed >> 35
	seed ^= seed << 4
	return seed
}

func RandomString(str []byte, offset int, seed uint64) uint64 {
	for i := offset; i < len(str); i++ {
		seed = RandomNumber(seed)
		str[i] = characterSet[seed%62]
	}
	return seed
}

func Hash(data []byte, bits int) bool {
	bs := sha256.Sum256(data)
	nbytes := bits / 8
	nbits := bits % 8
	idx := 0
	for ; idx < nbytes; idx++ {
		if bs[idx] > 0 {
			return false
		}
	}
	return (bs[idx] >> (8 - nbits)) == 0
}
