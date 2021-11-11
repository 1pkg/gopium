package fmtio

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/template"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

const (
	cnone     = "none"
	cadd      = "add"
	cdel      = "del"
	cdiff     = "diff"
	fhtmltmpl = `
<html>
	<head>
		<link
			href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css"
			rel="stylesheet"
			integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk"
			crossorigin="anonymous"
		>
		<script
			src="http://code.jquery.com/jquery-3.5.1.min.js"
			integrity="sha256-9/aliU8dGd2tb6OSsuzixeV4y/faTqgFtohetphbbj0="
			crossorigin="anonymous"
		></script>
		<script
			src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.bundle.min.js"
			integrity="sha384-1CmrxMRARb6aLqgBO7yyAxTOQE2AKb9GfXnEo760AUcUmFx3ibVJJAzGytlQcNXd"
			crossorigin="anonymous"
		></script>
		<style>
			.add{ background: green; }
			.del{ background: red; }
			.diff{ background: white; }
		</style>
	</head>
	<body>
		<div class="accordion" id="structs">
		{{- range $struct, $fields := . }}
		<div class="card">
			<div class="card-header" id="heading{{$struct}}">
				<h2 class="mb-0">
					<button
						class="btn btn-link btn-block"
						type="button"
						data-toggle="collapse"
						data-target="#collapse{{$struct}}"
						aria-expanded="true"
						aria-controls="collapse{{$struct}}"
					>
						{{$struct}}
					</button>
				</h2>
			</div>
			<div id="collapse{{$struct}}" class="collapse" aria-labelledby="heading{{$struct}}" data-parent="#structs">
				<table class="table">
					<thead>
						<tr>
							<th scope="col">#</th>
							<th scope="col">Name</th>
							<th scope="col">Type</th>
							<th scope="col">Size</th>
							<th scope="col">Align</th>
							<th scope="col">Tag</th>
							<th scope="col">Exported</th>
							<th scope="col">Embedded</th>
							<th scope="col">Doc</th>
							<th scope="col">Comment</th>
						</tr>
					</thead>
					<tbody>
						{{- range $fields }}
						<tr class="{{.Class}}">
							<th scope="row">{{.Index}}</th>
							{{- if eq .Class "diff" }}
							<td>{{.OldName}} -> {{.NewName}}</td>
							<td>{{.OldType}} -> {{.NewType}}</td>
							<td>{{.OldSize}} -> {{.NewSize}}</td>
							<td>{{.OldAlign}} -> {{.NewAlign}}</td>
							<td>{{.OldTag}} -> {{.NewTag}}</td>
							<td>{{.OldExported}} -> {{.NewExported}}</td>
							<td>{{.OldEmbedded}} -> {{.NewEmbedded}}</td>
							<td>{{.OldDoc}} -> {{.NewDoc}}</td>
							<td>{{.OldComment}} -> {{.NewComment}}</td>
							{{else if eq .Class "add"}}
							<td>{{.NewName}}</td>
							<td>{{.NewType}}</td>
							<td>{{.NewSize}}</td>
							<td>{{.NewAlign}}</td>
							<td>{{.NewTag}}</td>
							<td>{{.NewExported}}</td>
							<td>{{.NewEmbedded}}</td>
							<td>{{.NewDoc}}</td>
							<td>{{.NewComment}}</td>
							{{else if eq .Class "del"}}
							<td>{{.OldName}}</td>
							<td>{{.OldType}}</td>
							<td>{{.OldSize}}</td>
							<td>{{.OldAlign}}</td>
							<td>{{.OldTag}}</td>
							<td>{{.OldExported}}</td>
							<td>{{.OldEmbedded}}</td>
							<td>{{.OldDoc}}</td>
							<td>{{.OldComment}}</td>
							{{ end -}}
						</tr>
						{{- end }}
					</tbody>
				</table>
			</div>
		</div>
		{{- end }}
		</div>
	<body>
</html>
`
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
	_, _ = buf.WriteString("| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Size Difference | Relative Size Difference | Original Ptr Size with Pad | Current Ptr Size with Pad | Absolute Ptr Size Difference | Relative Ptr Size Difference |\n")
	_, _ = buf.WriteString("| :---: | :---: | :---: | :---: | :---: |\n")
	for id, sto := range fo {
		// if both collections contains
		// struct, compare them
		if stf, ok := fr[id]; ok {
			// get aligned size and align
			sizeo, _, ptro := collections.SizeAlignPtr(sto)
			sizer, _, ptrr := collections.SizeAlignPtr(stf)
			// write diff info
			// no error should be
			// checked as it uses
			// buffered writer
			_, _ = buf.WriteString(
				fmt.Sprintf(
					"| %s | %d bytes | %d bytes | %+d bytes | %+.2f%% | %d bytes | %d bytes | %+d bytes | %+.2f%% |\n",
					sto.Name,
					sizeo,
					sizer,
					sizer-sizeo,
					float64(sizer-sizeo)/float64(sizeo)*100.0,
					ptro,
					ptrr,
					ptrr-ptro,
					float64(ptrr-ptro)/float64(ptro)*100.0,
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

// FieldsHtmlt defines diff implementation
// which compares two categorized collections
// to formatted struct fields html table byte slice
func FieldsHtmlt(o gopium.Categorized, r gopium.Categorized) ([]byte, error) {
	// prepare buffer and collections
	var buf bytes.Buffer
	fo, fr := o.Full(), r.Full()
	// prepare data set for template
	data := make(map[string][]interface{}, len(fo))
	// go through original collection
	for id, sto := range fo {
		// if both collections contains
		// struct, compare them
		if stf, ok := fr[id]; ok {
			// precalculate fields sizes for both structs
			stol, stfl := len(sto.Fields), len(stf.Fields)
			// find bigger size and use it as max index
			index := int(math.Max(float64(stol), float64(stfl)))
			// create resulted fields set for template
			fields := make([]interface{}, 0, index)
			for i := 0; i < index; i++ {
				// set field class
				// base on index
				// also grab original
				// and resulted field
				class := cnone
				fo := gopium.Field{}
				if i < stol {
					class = cdel
					fo = sto.Fields[i]
				}
				fr := gopium.Field{}
				if i < stfl {
					class = cadd
					fr = stf.Fields[i]
				}
				if i < stol && i < stfl {
					class = cdiff
				}
				// push fields comparison with meta
				// to struct template data bucket
				fields = append(fields, struct {
					Index       int
					Class       string
					OldName     string
					OldType     string
					OldSize     int64
					OldAlign    int64
					OldTag      string
					OldExported bool
					OldEmbedded bool
					OldDoc      string
					OldComment  string
					NewName     string
					NewType     string
					NewSize     int64
					NewAlign    int64
					NewTag      string
					NewExported bool
					NewEmbedded bool
					NewDoc      string
					NewComment  string
				}{
					Index:       i + 1,
					Class:       class,
					OldName:     fo.Name,
					OldType:     fo.Type,
					OldSize:     fo.Size,
					OldAlign:    fo.Align,
					OldTag:      fmt.Sprintf("%q", fo.Tag),
					OldExported: fo.Exported,
					OldEmbedded: fo.Exported,
					OldDoc:      fmt.Sprintf("%q", strings.Join(fo.Doc, " ")),
					OldComment:  fmt.Sprintf("%q", strings.Join(fo.Comment, " ")),
					NewName:     fr.Name,
					NewType:     fr.Type,
					NewSize:     fr.Size,
					NewAlign:    fr.Align,
					NewTag:      fmt.Sprintf("%q", fr.Tag),
					NewExported: fr.Exported,
					NewEmbedded: fr.Exported,
					NewDoc:      fmt.Sprintf("%q", strings.Join(fr.Doc, " ")),
					NewComment:  fmt.Sprintf("%q", strings.Join(fr.Comment, " ")),
				})
			}
			// set struct template data bucket
			data[sto.Name] = fields
		}
	}
	// parse and execute template
	tmpl := template.Must(template.New("tmpl").Parse(fhtmltmpl))
	err := tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}
