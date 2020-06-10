package fmtio

import (
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
)

func TestDiff(t *testing.T) {
	// prepare
	oh := collections.NewHierarchic("")
	rh := collections.NewHierarchic("")
	rhb := collections.NewHierarchic("")
	oh.Push("test", "test", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test1",
				Size:  3,
				Align: 1,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test3",
				Size:  3,
				Align: 1,
			},
		},
	})
	rh.Push("test", "test", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test1",
				Size:  3,
				Align: 1,
			},
			{
				Name:  "test3",
				Size:  3,
				Align: 1,
			},
		},
	})
	rhb.Push("test", "test", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
		},
	})
	rhb.Push("test1", "test1", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
		},
	})
	table := map[string]struct {
		fmt gopium.Diff
		o   gopium.Categorized
		r   gopium.Categorized
		b   []byte
		err error
	}{
		"size align md table should return expected result for empty collections": {
			fmt: SizeAlignMdt,
			o:   collections.NewHierarchic(""),
			r:   collections.NewHierarchic(""),
			b: []byte(`
| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |
| :---: | :---: | :---: | :---: | :---: |
`),
		},
		"size align md table should return expected result for non empty collections": {
			fmt: SizeAlignMdt,
			o:   oh,
			r:   rh,
			b: []byte(`
| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |
| :---: | :---: | :---: | :---: | :---: |
| test | 24 bytes | 16 bytes | -8 bytes | -33.33% |
| Total | 24 bytes | 16 bytes | -8 bytes | -33.33% |
`),
		},
		"size align md table should return expected result for non empty overlapping collections": {
			fmt: SizeAlignMdt,
			o:   oh,
			r:   rhb,
			b: []byte(`
| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |
| :---: | :---: | :---: | :---: | :---: |
| test | 24 bytes | 32 bytes | +8 bytes | +33.33% |
| Total | 24 bytes | 32 bytes | +8 bytes | +33.33% |
`),
		},
		"fields html table should return expected result for empty collections": {
			fmt: FieldsHtmlt,
			o:   collections.NewHierarchic(""),
			r:   collections.NewHierarchic(""),
			b: []byte(`
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
		</div>
	<body>
</html>
`),
		},
		"fields html table should return expected result for non empty collections": {
			fmt: FieldsHtmlt,
			o:   oh,
			r:   rh,
			b: []byte(`
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
		<div class="card">
			<div class="card-header" id="headingtest">
				<h2 class="mb-0">
					<button
						class="btn btn-link btn-block"
						type="button"
						data-toggle="collapse"
						data-target="#collapsetest"
						aria-expanded="true"
						aria-controls="collapsetest"
					>
						test
					</button>
				</h2>
			</div>
			<div id="collapsetest" class="collapse" aria-labelledby="headingtest" data-parent="#structs">
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
						<tr class="diff">
							<th scope="row">1</th>
							<td>test1 -> test2</td>
							<td> -> float64</td>
							<td>3 -> 8</td>
							<td>1 -> 8</td>
							<td>"" -> ""</td>
							<td>false -> false</td>
							<td>false -> false</td>
							<td>"" -> ""</td>
							<td>"" -> ""</td>
							</tr>
						<tr class="diff">
							<th scope="row">2</th>
							<td>test2 -> test1</td>
							<td>float64 -> </td>
							<td>8 -> 3</td>
							<td>8 -> 1</td>
							<td>"" -> ""</td>
							<td>false -> false</td>
							<td>false -> false</td>
							<td>"" -> ""</td>
							<td>"" -> ""</td>
							</tr>
						<tr class="diff">
							<th scope="row">3</th>
							<td>test3 -> test3</td>
							<td> -> </td>
							<td>3 -> 3</td>
							<td>1 -> 1</td>
							<td>"" -> ""</td>
							<td>false -> false</td>
							<td>false -> false</td>
							<td>"" -> ""</td>
							<td>"" -> ""</td>
							</tr>
					</tbody>
				</table>
			</div>
		</div>
		</div>
	<body>
</html>
`),
		},
		"fields html table return expected result for non empty overlapping collections": {
			fmt: FieldsHtmlt,
			o:   oh,
			r:   rhb,
			b: []byte(`
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
		<div class="card">
			<div class="card-header" id="headingtest">
				<h2 class="mb-0">
					<button
						class="btn btn-link btn-block"
						type="button"
						data-toggle="collapse"
						data-target="#collapsetest"
						aria-expanded="true"
						aria-controls="collapsetest"
					>
						test
					</button>
				</h2>
			</div>
			<div id="collapsetest" class="collapse" aria-labelledby="headingtest" data-parent="#structs">
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
						<tr class="diff">
							<th scope="row">1</th>
							<td>test1 -> test2</td>
							<td> -> float64</td>
							<td>3 -> 8</td>
							<td>1 -> 8</td>
							<td>"" -> ""</td>
							<td>false -> false</td>
							<td>false -> false</td>
							<td>"" -> ""</td>
							<td>"" -> ""</td>
							</tr>
						<tr class="diff">
							<th scope="row">2</th>
							<td>test2 -> test2</td>
							<td>float64 -> float64</td>
							<td>8 -> 8</td>
							<td>8 -> 8</td>
							<td>"" -> ""</td>
							<td>false -> false</td>
							<td>false -> false</td>
							<td>"" -> ""</td>
							<td>"" -> ""</td>
							</tr>
						<tr class="diff">
							<th scope="row">3</th>
							<td>test3 -> test2</td>
							<td> -> float64</td>
							<td>3 -> 8</td>
							<td>1 -> 8</td>
							<td>"" -> ""</td>
							<td>false -> false</td>
							<td>false -> false</td>
							<td>"" -> ""</td>
							<td>"" -> ""</td>
							</tr>
						<tr class="add">
							<th scope="row">4</th>
							<td>test2</td>
							<td>float64</td>
							<td>8</td>
							<td>8</td>
							<td>""</td>
							<td>false</td>
							<td>false</td>
							<td>""</td>
							<td>""</td>
							</tr>
					</tbody>
				</table>
			</div>
		</div>
		</div>
	<body>
</html>
`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			b, err := tcase.fmt(tcase.o, tcase.r)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// format actual and expected identically
			actual := strings.Trim(string(b), "\n")
			expected := strings.Trim(string(tcase.b), "\n")
			if err == nil && !reflect.DeepEqual(actual, expected) {
				t.Errorf("actual %v doesn't equal to expected %v", actual, expected)
			}
		})
	}
}
