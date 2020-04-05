package mocks

import "go/types"

// TypeMock defines mock type
// data transfer object
type TypeMock struct {
	Name        string
	Size, Align int64
}

// MavenMock defines mock maven implementation
type MavenMock struct {
	Types                   map[string]TypeMock
	SysCacheVals            []int64
	SysWordVal, SysAlignVal int64
}

// SysWord mock implementation
func (m MavenMock) SysWord() int64 {
	return m.SysWordVal
}

// SysAlign mock implementation
func (m MavenMock) SysAlign() int64 {
	return m.SysAlignVal
}

// SysCache mock implementation
func (m MavenMock) SysCache(level uint) int64 {
	// check if we have it in vals
	if int(level) < len(m.SysCacheVals) {
		return m.SysCacheVals[level]
	}
	// otherwise return default val
	return 0
}

// Name mock implementation
func (m MavenMock) Name(t types.Type) string {
	// check if we have it in vals
	if t, ok := m.Types[t.String()]; ok {
		return t.Name
	}
	// otherwise return default val
	return ""
}

// Size mock implementation
func (m MavenMock) Size(t types.Type) int64 {
	// check if we have it in vals
	if t, ok := m.Types[t.String()]; ok {
		return t.Size
	}
	// otherwise return default val
	return 0
}

// Align mock implementation
func (m MavenMock) Align(t types.Type) int64 {
	// check if we have it in vals
	if t, ok := m.Types[t.String()]; ok {
		return t.Align
	}
	// otherwise return default val
	return 0
}
