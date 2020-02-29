package strategy

import (
	"context"
	"reflect"
	"strings"

	"1pkg/gopium"
)

// list of tag presets
var (
	taglexasc   = tag{tag: string(LexAsc)}
	taglexdesc  = tag{tag: string(LexDesc)}
	taglenasc   = tag{tag: string(LenAsc)}
	taglendesc  = tag{tag: string(LenDesc)}
	tagpack     = tag{tag: string(Pack)}
	tagunpack   = tag{tag: string(Unpack)}
	tagpadsys   = tag{tag: string(PadSys)}
	tagpadtnat  = tag{tag: string(PadTnat)}
	tagfsahrel1 = tag{tag: string(FShareL1)}
	tagfsahrel2 = tag{tag: string(FShareL2)}
	tagfsahrel3 = tag{tag: string(FShareL3)}
	tagcachel1  = tag{tag: string(CacheL1)}
	tagcachel2  = tag{tag: string(CacheL2)}
	tagcachel3  = tag{tag: string(CacheL3)}
	tagsepsys   = tag{tag: string(SepSys)}
	tagsepl1    = tag{tag: string(SepL1)}
	tagsepl2    = tag{tag: string(SepL2)}
	tagsepl3    = tag{tag: string(SepL3)}
)

// tag defines strategy implementation
// that adds or updates tag annotation
// that could be parsed by group strategy
type tag struct {
	tag   string
	force bool
}

// Apply tag implementation
func (stg tag) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// iterate through all fields
	for i := range r.Fields {
		f := &r.Fields[i]
		// grab the field tag
		tag, ok := reflect.StructTag(f.Tag).Lookup(gopium.TAG)
		// in case gopium tag already exists
		// and force is set - replace tag
		// in case tag is not empty and
		// gopium tag doesn't exist - append tag
		// in case tag is empty - set tag
		if ok && stg.force {
			f.Tag = strings.Replace(f.Tag, tag, stg.tag, 1)
		} else if len(f.Tag) != 0 {
			f.Tag += " " + stg.tag
		} else {
			f.Tag = stg.tag
		}
	}
	return
}
