package gopium

// Field defines single structure field
// data transfer object abstraction
type Field struct {
	Name     string   `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Type     string   `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Size     int64    `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Align    int64    `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Ptr      int64    `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Tag      string   `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Exported bool     `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Embedded bool     `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Doc      []string `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Comment  []string `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
} // struct size: 122 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; struct ptr scan size: 106 bytes; - 🌺 gopium @1pkg

// Struct defines single structure
// data transfer object abstraction
type Struct struct {
	Name    string   `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Doc     []string `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Comment []string `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
	Fields  []Field  `gopium:"filter_pads,struct_annotate_comment,add_tag_group_force"`
} // struct size: 88 bytes; struct align: 8 bytes; struct aligned size: 88 bytes; struct ptr scan size: 72 bytes; - 🌺 gopium @1pkg
