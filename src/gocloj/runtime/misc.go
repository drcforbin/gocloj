package runtime

import (
	// "errors"
	// "fmt"
	"gocloj/data"
	// "strings"
)

// TODO: make this right!
func Truthy(val data.Atom) bool {
	// TODO: empty seq should be false
	// TODO: handle true / false (T/F)

	if val != nil {
		return true
	}

	return false
}

/*
func DumpTree(atom data.Atom, dent int, b *strings.Builder, eol bool) error {
	for i := 0; i < dent; i++ {
		_, err := b.WriteString(" ")
		if err != nil {
			return errors.New(fmt.Sprintf("unable to indent tree item %s: %s",
				atom, err))
		}
	}

	var err error

	switch v := atom.(type) {
	case *data.Nil:
		_, err = b.WriteString("nil")
	case *data.Lst:
		_, err = b.WriteString("List:")
		if len(v.Vals) > 0 {
			for _, item := range v.Vals {
				_, err = b.WriteString("\n")
				if err != nil {
					break
				}
				err = DumpTree(item, dent+1, b, false)
				if err != nil {
					break
				}
			}
		} else {
			_, err = b.WriteString(" (empty)")
		}
	case *data.Pair:
		_, err = b.WriteString(fmt.Sprintf("Pair: %s", v.GetName()))
		// TODO: Val and Props
	case *data.Num:
		_, err = b.WriteString(fmt.Sprintf("Num: %s", v.Val))
	default:
		err = errors.New(fmt.Sprintf("unexpected atom type: %T", atom))
	}

	if err == nil {
		if eol {
			_, err = b.WriteString("\n")
			if err != nil {
				return errors.New(fmt.Sprintf("unable write eol %s: %s",
					atom, err))
			}
		}
	} else {
		return errors.New(fmt.Sprintf("unable to write atom %s: %s",
			atom, err))
	}

	return nil
}
*/
