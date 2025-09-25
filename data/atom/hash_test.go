package atom

import (
	"testing"
)

func teshHashBytes(t *testing.T, desc string, key []byte, seed uint32, expected uint32) {
	actual := MurmurHash3(key, seed)
	if actual != expected {
		t.Errorf("%s: expected 0x%X; got 0x%X", desc, actual, expected)
	}
}

func teshHashString(t *testing.T, desc string, key string, seed uint32, expected uint32) {
	var actual uint32
	if seed == 0 {
		actual = hashString(key)
	} else {
		actual = MurmurHash3([]byte(key), seed)
	}

	if actual != expected {
		t.Errorf("%s: expected 0x%X; got %0xX", desc, actual, expected)
	}
}

func TestHash(t *testing.T) {
	/*
		| (no bytes)   | 0          | 0          | with zero data and zero seed, everything becomes zero
		| (no bytes)   | 1          | 0x514E28B7 | ignores nearly all the math
		| (no bytes)   | 0xffffffff | 0x81F16F39 | make sure your seed uses unsigned 32-bit math
		| FF FF FF FF  | 0          | 0x76293B50 | make sure 4-byte chunks use unsigned math
		| 21 43 65 87  | 0          | 0xF55B516B | Endian order. UInt32 should end up as 0x87654321
		| 21 43 65 87  | 0x5082EDEE | 0x2362F9DE | Special seed value eliminates initial key with xor
		| 21 43 65     | 0          | 0x7E4A8634 | Only three bytes. Should end up as 0x654321
		| 21 43        | 0          | 0xA0F7B07A | Only two bytes. Should end up as 0x4321
		| 21           | 0          | 0x72661CF4 | Only one byte. Should end up as 0x21
		| 00 00 00 00  | 0          | 0x2362F9DE | Make sure compiler doesn't see zero and convert to null
		| 00 00 00     | 0          | 0x85F0B427 |
		| 00 00        | 0          | 0x30F4C306 |
		| 00           | 0          | 0x514E28B7 |
	*/

	// empty string with zero seed should give zero
	teshHashBytes(t, "empty bytes, zero seed", []byte{}, 0, 0)
	teshHashString(t, "empty string, zero seed", "", 0, 0)
	teshHashBytes(t, "empty bytes, one seed", []byte{}, 1, 0x514E28B7)
	teshHashString(t, "empty string, one seed", "", 1, 0x514E28B7)

	// make sure seed value is handled unsigned
	teshHashBytes(t, "empty bytes unsigned", []byte{}, 0xffffffff, 0x81F16F39)
	teshHashString(t, "empty string unsigned", "", 0xffffffff, 0x81F16F39)

	// make sure we handle embedded nulls
	teshHashBytes(t, "null bytes", []byte{0, 0, 0, 0}, 0, 0x2362F9DE)
	teshHashString(t, "string with nulls", "\000\000\000\000", 0, 0x2362F9DE)

	teshHashString(t, "four char string", "aaaa", 0x9747b28c, 0x5A97808A)
	teshHashString(t, "three char string", "aaa", 0x9747b28c, 0x283E0130)
	teshHashString(t, "two char string", "aa", 0x9747b28c, 0x5D211726)
	teshHashString(t, "one char string", "a", 0x9747b28c, 0x7FA09EA6)

	// endian order within the chunks
	teshHashString(t, "endian check, four char", "abcd", 0x9747b28c, 0xF0478627) //one full chunk
	teshHashString(t, "endian check, three char", "abc", 0x9747b28c, 0xC84A62DD)
	teshHashString(t, "endian check, two char", "ab", 0x9747b28c, 0x74875592)
	teshHashString(t, "endian check, one char", "a", 0x9747b28c, 0x7FA09EA6)

	teshHashString(t, "hello world check", "Hello, world!", 0x9747b28c, 0x24884CBA)

	//Make sure you handle UTF-8 high characters. A bcrypt implementation messed this up
	teshHashString(t, "utf-8 chars", "ππππππππ", 0x9747b28c, 0xD58063C1) //U+03C0: Greek Small Letter Pi

	// string of 256 characters.
	// TODO: teshHashString(t, "256 chars",  StringOfChar("a", 256), 0x9747b28c, 0x37405BDC)

	teshHashString(t, "abc check", "abc", 0, 0xB3DD93FA)
	teshHashString(t, "long abc check", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", 0, 0xEE925B90)
}

func BenchmarkHash(b *testing.B) {
	benchString := func(str string, expected uint32) {
		actual := hashString(str)
		if actual != expected {
			b.Errorf("expected 0x%X; got %0xX", actual, expected)
		}
	}

	for i := 0; i < b.N; i++ {
		benchString("a", 0x3C2569B2)
		benchString("ab", 0x9BBFD75F)
		benchString("abc", 0xB3DD93FA)
		benchString("four score and Seven7 years ago", 0xC212FF8B)
		benchString("lorem ipsum", 0x2B0916F8)
		benchString("sit dolor", 0xFDE0E55A)

		benchString(
			`When in the Course of human events, it becomes necessary for one
people to dissolve the political bands which have connected them
with another, and to assume among the powers of the earth, the
separate and equal station to which the Laws of Nature and of
Nature's God entitle them, a decent respect to the opinions of
mankind requires that they should declare the causes which impel
them to the separation.

We hold these truths to be self-evident, that all men are created
equal, that they are endowed by their Creator with certain
unalienable Rights, that among these are Life, Liberty and the
pursuit of Happiness.--That to secure these rights, Governments
are instituted among Men, deriving their just powers from the
consent of the governed, --That whenever any Form of Government
becomes destructive of these ends, it is the Right of the People
to alter or to abolish it, and to institute new Government,
laying its foundation on such principles and organizing its
powers in such form, as to them shall seem most likely to effect
their Safety and Happiness. Prudence, indeed, will dictate that
Governments long established should not be changed for light and
transient causes; and accordingly all experience hath shewn, that
mankind are more disposed to suffer, while evils are sufferable,
than to right themselves by abolishing the forms to which they
are accustomed. But when a long train of abuses and usurpations,
pursuing invariably the same Object evinces a design to reduce
them under absolute Despotism, it is their right, it is their
duty, to throw off such Government, and to provide new Guards for
their future security.`, 0x8776A21B)
	}
}
