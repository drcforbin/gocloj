package atom

import (
	"math/bits"
)

// MurmerHash3 from clojure

const defaultSeed uint32 = 0
const C1 uint32 = 0xcc9e2d51
const C2 uint32 = 0x1b873593

func mixK1(k1 uint32) uint32 {
	k1 *= C1
	k1 = bits.RotateLeft32(k1, 15)
	k1 *= C2
	return k1
}

func mixH1(h1, k1 uint32) uint32 {
	h1 ^= k1
	h1 = bits.RotateLeft32(h1, 13)
	h1 = h1*5 + 0xe6546b64
	return h1
}

func fmix32(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func mixCollHash(hash, count uint32) uint32 {
	h1 := defaultSeed
	k1 := mixK1(hash)
	h1 = mixH1(h1, k1)

	h1 ^= count
	return fmix32(h1)
}

func hashString(input string) uint32 {
	return MurmurHash3([]byte(input), defaultSeed)
}

func MurmurHash3(key []byte, seed uint32) uint32 {
	h1 := seed

	//----------
	// body

	blocks := len(key) / 4
	for i := 0; i < blocks*4; i += 4 {
		// assumes little endian
		k1 := uint32(key[i+0]) |
			uint32(key[i+1])<<8 |
			uint32(key[i+2])<<16 |
			uint32(key[i+3])<<24

		k1 = mixK1(k1)
		h1 = mixH1(h1, k1)
	}

	//----------
	// tail

	tailBytes := len(key) & 3
	if tailBytes > 0 {
		tail := key[blocks*4:]

		k1 := uint32(0)
		switch tailBytes {
		case 3:
			k1 ^= uint32(tail[2]) << 16
			fallthrough
		case 2:
			k1 ^= uint32(tail[1]) << 8
			fallthrough
		case 1:
			k1 ^= uint32(tail[0])
			k1 = mixK1(k1)
			h1 ^= k1
		}
	}

	//----------
	// finalization

	h1 ^= uint32(len(key))
	return fmix32(h1)
}
