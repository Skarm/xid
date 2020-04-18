package xid

import (
	"bytes"
	"sort"
)

// Compare returns an integer comparing two IDs. It behaves just like `bytes.Compare`.
// The result will be 0 if two IDs are identical, -1 if current id is less than the other one,
// and 1 if current id is greater than the other.
func (id ID) Compare(other ID) int {
	return bytes.Compare(id[:], other[:])
}

type sorter []ID

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Less(i, j int) bool {
	return s[i].Compare(s[j]) < 0
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort sorts an array of IDs inplace.
// It works by wrapping `[]ID` and use `sort.Sort`.
func Sort(ids []ID) {
	sort.Sort(sorter(ids))
}
