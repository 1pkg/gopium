package fmtio

import (
	"bytes"
	"fmt"

	"1pkg/gopium/gopium"
	"1pkg/gopium/collections"
)

// SizeAlignMdt defines diff implementation
// which compares two categorized collections
// to formatted markdown table byte slice
func SizeAlignMdt(o gopium.Categorized, r gopium.Categorized) ([]byte, error) {
	// prepare buffer and collections
	var buf bytes.Buffer
	fo, fr := o.Full(), r.Full()
	// write header
	// no error should be
	// checked as it uses
	// buffered writer
	_, _ = buf.WriteString("| Struct Name | Original Size With Pad | Original Align | Current Size With Pad | Current Align |\n")
	_, _ = buf.WriteString("| :---: | :---: | :---: | :---: | :---: |\n")
	for id, sto := range fo {
		// if both collections contains
		// struct, compare them
		if stf, ok := fr[id]; ok {
			// get aligned size and align
			sizeo, aligno := collections.SizeAlign(sto)
			sizer, alignr := collections.SizeAlign(stf)
			// write diff info
			// no error should be
			// checked as it uses
			// buffered writer
			_, _ = buf.WriteString(
				fmt.Sprintf(
					"| %s | %d | %d | %d | %d |\n",
					sto.Name,
					sizeo,
					aligno,
					sizer,
					alignr,
				),
			)
		}
	}
	return buf.Bytes(), nil
}
