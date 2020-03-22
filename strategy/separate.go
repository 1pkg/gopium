package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// separating defines strategy implementation
// that separates structure with
// additional sys/cpu cache padding
// by adding one padding before and one padding after
// structure fields list
type separate struct {
	c   gopium.Curator
	l   uint // cache line num
	sys bool // should sys padding be used
}

// Apply separate implementation
func (stg separate) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// get separator size
	sep := stg.c.SysAlign()
	// if we wanna use
	// non max system separator
	if !stg.sys {
		sep = stg.c.SysCache(stg.l)
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
