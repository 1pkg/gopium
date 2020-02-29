package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of fshare presets
var (
	sepsys = sep{sys: true}
	sepl1  = sep{line: 1}
	sepl2  = sep{line: 2}
	sepl3  = sep{line: 3}
)

// sep defines strategy implementation
// that separates structure with
// additional sys/cpu cache padding
// by adding one padding before and one padding after
// structure fields list
type sep struct {
	curator gopium.Curator
	line    uint
	sys     bool
}

// Curator erich sep strategy with curator instance
func (stg sep) Curator(curator gopium.Curator) sep {
	stg.curator = curator
	return stg
}

// Apply sep implementation
func (stg sep) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// get separator size
	sep := stg.curator.SysAlign()
	// if we wanna use
	// non max system separator
	if !stg.sys {
		sep = stg.curator.SysCache(stg.line)
	}
	// add field before and after
	r.Fields = append([]gopium.Field{
		gopium.Field{
			Name:  "_",
			Type:  fmt.Sprintf("[%d]byte", sep),
			Size:  sep,
			Align: 1,
		},
	}, r.Fields...)
	r.Fields = append(r.Fields, gopium.Field{
		Name:  "_",
		Type:  fmt.Sprintf("[%d]byte", sep),
		Size:  sep,
		Align: 1,
	})
	return
}
