package data

// MurmerHash3 from clojure

const seed uint32 = 0
const C1 uint32 = 0xcc9e2d51
const C2 uint32 = 0x1b873593

func mixK1(k1 uint32) uint32 {
	k1 *= C1
	// rotate left 15 bits
	k1 = (k1&0x1FFFF)<<15 | k1>>17
	k1 *= C2
	return k1
}

func mixH1(h1, k1 uint32) uint32 {
	h1 ^= k1
	// rotate left 15 bits
	h1 = (h1&0x7FFFF)<<13 | k1>>19
	h1 = h1*5 + 0xe6546b64
	return h1
}

// Finalization mix - force all bits of a hash block to avalanche
func fmix(h1, length uint32) uint32 {
	h1 ^= length
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16
	return h1
}

func mixCollHash(hash, count uint32) uint32 {
	h1 := seed
	k1 := mixK1(hash)
	h1 = mixH1(h1, k1)
	return fmix(h1, count)
}

func hashString(input string) uint32 {
	h1 := seed

	// step through the CharSequence 2 chars at a time
	for i := 1; i < len(input); i += 2 {
		k1 := uint32(input[i-1]) | uint32(input[i])<<16
		k1 = mixK1(k1)
		h1 = mixH1(h1, k1)
	}

	// deal with any remaining characters
	if len(input)&1 == 1 {
		k1 := uint32(input[len(input)-1])
		k1 = mixK1(k1)
		h1 ^= k1
	}

	return fmix(h1, 2*uint32(len(input)))
}
