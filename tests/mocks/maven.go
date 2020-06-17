package mocks

import "go/types"

// Type defines mock type
// data transfer object
type Type struct {
	Name  string `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Size  int64  `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Align int64  `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Maven defines mock maven implementation
type Maven struct {
	SCache []int64         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Types  map[string]Type `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	SWord  int64           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	SAlign int64           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_      [16]byte        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// SysWord mock implementation
func (m Maven) SysWord() int64 {
	return m.SWord
}

// SysAlign mock implementation
func (m Maven) SysAlign() int64 {
	return m.SAlign
}

// SysCache mock implementation
func (m Maven) SysCache(level uint) int64 {
	// decrement level to match index
	l := int(level) - 1
	// check if we have it in vals
	if l >= 0 && l < len(m.SCache) {
		return m.SCache[l]
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
