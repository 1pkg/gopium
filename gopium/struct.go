package gopium

// Field defines single structure field
// data transfer object abstraction
type Field struct {
	Name     string   `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Type     string   `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Size     int64    `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Align    int64    `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Tag      string   `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Exported bool     `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Embedded bool     `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Doc      []string `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Comment  []string `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
} // struct size: 114 bytes; struct align: 8 bytes; struct aligned size: 120 bytes; - 🌺 gopium @1pkg

// Struct defines single structure
// data transfer object abstraction
type Struct struct {
	Name    string   `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Doc     []string `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Comment []string `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
	Fields  []Field  `gopium:"filter_pads,comment_struct_annotate,add_tag_group_force"`
} // struct size: 88 bytes; struct align: 8 bytes; struct aligned size: 88 bytes; - 🌺 gopium @1pkg
