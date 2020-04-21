package mocks

import "go/types"

// Type defines mock type
// data transfer object
type Type struct {
	Name  string
	Size  int64
	Align int64
}

// Maven defines mock maven implementation
type Maven struct {
	Types                   map[string]Type
	SysCacheVals            []int64
	SysWordVal, SysAlignVal int64
}

// SysWord mock implementation
func (m Maven) SysWord() int64 {
	return m.SysWordVal
}

// SysAlign mock implementation
func (m Maven) SysAlign() int64 {
	return m.SysAlignVal
}

// SysCache mock implementation
func (m Maven) SysCache(level uint) int64 {
	// decrement level to match index
	l := int(level) - 1
	// check if we have it in vals
	if l >= 0 && l < len(m.SysCacheVals) {
		return m.SysCacheVals[l]
	}
	// otherwise return default val
	return 0
}

// Name mock implementation
func (m Maven) Name(t types.Type) string {
	// check if we have it in vals
	if t, ok := m.Types[t.String()]; ok {
		return t.Name
	}
	// otherwise return default val
	return ""
}

// Size mock implementation
func (m Maven) Size(t types.Type) int64 {
	// check if we have it in vals
	if t, ok := m.Types[t.String()]; ok {
		return t.Size
	}
	// otherwise return default val
	return 0
}

// Align mock implementation
func (m Maven) Align(t types.Type) int64 {
	// check if we have it in vals
	if t, ok := m.Types[t.String()]; ok {
		return t.Align
	}
	// otherwise return default val
	return 0
}
