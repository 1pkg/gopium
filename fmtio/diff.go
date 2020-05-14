package fmtio

import "1pkg/gopium/collections"

// Diff defines abstraction for formatting
// gopium collections difference to byte slice
type Diff func(collections.Hierarchic, collections.Hierarchic) ([]byte, error)
