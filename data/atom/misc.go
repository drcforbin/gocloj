package atom

/*
// -1 if x <  y
//  0 if x == y
// +1 if x >  y
func Compare(atom1 atom.Atom, atom2 atom.Atom) (cmp int) {
// NOTE: T > everything else
// NOTE: nil < everything else
}
*/

func SeqEquals(seqa Seq, seqb Seq) bool {
	ita, itb := seqa.Iterator(), seqb.Iterator()

	// walk a's
	for ita.Next() {
		// are we out of b's?
		if !itb.Next() {
			return false
		}

		// TODO: need package-level Equals?
		if !ita.Value().Equals(itb.Value()) {
			return false
		}
	}

	// do we still have more b's?
	if itb.Next() {
		return false
	}

	return true
}

func SeqHash(seq Seq) uint32 {
	hash := uint32(0)
	count := uint32(0)

	it := seq.Iterator()
	for it.Next() {
		hash += it.Value().Hash()
		count++
	}

	return mixCollHash(hash, count)
}
