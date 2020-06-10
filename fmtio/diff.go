package fmtio

import (
	"bytes"
	"fmt"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
)

// SizeAlignMdt defines diff implementation
// which compares two categorized collections
// to formatted markdown table byte slice
func SizeAlignMdt(o gopium.Categorized, r gopium.Categorized) ([]byte, error) {
	// prepare buffer and collections
	var buf bytes.Buffer
	var tsizeo, tsizer int64
	fo, fr := o.Full(), r.Full()
	// write header
	// no error should be
	// checked as it uses
	// buffered writer
	_, _ = buf.WriteString("| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |\n")
	_, _ = buf.WriteString("| :---: | :---: | :---: | :---: | :---: |\n")
	for id, sto := range fo {
		// if both collections contains
		// struct, compare them
		if stf, ok := fr[id]; ok {
			// get aligned size and align
			sizeo, _ := collections.SizeAlign(sto)
			sizer, _ := collections.SizeAlign(stf)
			// write diff info
			// no error should be
			// checked as it uses
			// buffered writer
			_, _ = buf.WriteString(
				fmt.Sprintf(
					"| %s | %d bytes | %d bytes | %+d bytes | %+.2f%% |\n",
					sto.Name,
					sizeo,
					sizer,
					sizer-sizeo,
					float64(sizer-sizeo)/float64(sizeo)*100.0,
				),
			)
			// increment total sizes
			tsizeo += sizeo
			tsizer += sizer
		}
	}
	// zero divide guard
	if tsizeo > 0 {
		// write diff info
		// no error should be
		// checked as it uses
		// buffered writer
		_, _ = buf.WriteString(
			fmt.Sprintf(
				"| %s | %d bytes | %d bytes | %+d bytes | %+.2f%% |\n",
				"Total",
				tsizeo,
				tsizer,
				tsizer-tsizeo,
				float64(tsizer-tsizeo)/float64(tsizeo)*100.0,
			),
		)
	}
	return buf.Bytes(), nil
}
