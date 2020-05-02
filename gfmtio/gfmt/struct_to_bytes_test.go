package gfmt

import (
	"1pkg/gopium"
	"reflect"
	"strings"
	"testing"
)

func TestStructToBytes(t *testing.T) {
	// prepare
	table := map[string]struct {
		fmt StructToBytes
		st  gopium.Struct
		r   []byte
		err error
	}{
		"pretty json should return valid result for empty struct": {
			fmt: PrettyJson,
			st:  gopium.Struct{},
			r: []byte(`
{
	"Name": "",
	"Doc": null,
	"Comment": null,
	"Fields": null
}
`),
		},
		"pretty json should return valid result for non empty struct": {
			fmt: PrettyJson,
			st: gopium.Struct{
				Name:    "Test",
				Doc:     []string{"doctest"},
				Comment: []string{"comtest"},
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "string",
						Size:     16,
						Align:    8,
						Tag:      "test-tag",
						Exported: true,
						Embedded: true,
						Doc:      []string{"fdoctest"},
						Comment:  []string{"fcomtest"},
					},
					{
						Name:  "test-2",
						Type:  "test_type",
						Size:  12,
						Align: 4,
					},
				},
			},
			r: []byte(`
{
	"Name": "Test",
	"Doc": [
		"doctest"
	],
	"Comment": [
		"comtest"
	],
	"Fields": [
		{
			"Name": "test-1",
			"Type": "string",
			"Size": 16,
			"Align": 8,
			"Tag": "test-tag",
			"Exported": true,
			"Embedded": true,
			"Doc": [
				"fdoctest"
			],
			"Comment": [
				"fcomtest"
			]
		},
		{
			"Name": "test-2",
			"Type": "test_type",
			"Size": 12,
			"Align": 4,
			"Tag": "",
			"Exported": false,
			"Embedded": false,
			"Doc": null,
			"Comment": null
		}
	]
}
`),
		},
		"pretty xml should return valid result for empty struct": {
			fmt: PrettyXml,
			st:  gopium.Struct{},
			r: []byte(`
<Struct>
	<Name></Name>
</Struct>
`),
		},
		"pretty xml should return valid result for non empty struct": {
			fmt: PrettyXml,
			st: gopium.Struct{
				Name:    "Test",
				Doc:     []string{"doctest"},
				Comment: []string{"comtest"},
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "string",
						Size:     16,
						Align:    8,
						Tag:      "test-tag",
						Exported: true,
						Embedded: true,
						Doc:      []string{"fdoctest"},
						Comment:  []string{"fcomtest"},
					},
					{
						Name:  "test-2",
						Type:  "test_type",
						Size:  12,
						Align: 4,
					},
				},
			},
			r: []byte(`
<Struct>
	<Name>Test</Name>
	<Doc>doctest</Doc>
	<Comment>comtest</Comment>
	<Fields>
		<Name>test-1</Name>
		<Type>string</Type>
		<Size>16</Size>
		<Align>8</Align>
		<Tag>test-tag</Tag>
		<Exported>true</Exported>
		<Embedded>true</Embedded>
		<Doc>fdoctest</Doc>
		<Comment>fcomtest</Comment>
	</Fields>
	<Fields>
		<Name>test-2</Name>
		<Type>test_type</Type>
		<Size>12</Size>
		<Align>4</Align>
		<Tag></Tag>
		<Exported>false</Exported>
		<Embedded>false</Embedded>
	</Fields>
</Struct>
`),
		},
		"pretty csv should return valid result for empty struct": {
			fmt: PrettyCsv,
			st:  gopium.Struct{},
			r: []byte(`
Struct Name,Struct Doc,Struct Comment,Field Name,Field Type,Field Size,Field Align,Field Tag,Field Exported,Field Embedded,Field Doc,Field Comment
`),
		},
		"pretty csv should return valid result for non empty struct": {
			fmt: PrettyCsv,
			st: gopium.Struct{
				Name:    "Test",
				Doc:     []string{"doctest"},
				Comment: []string{"comtest"},
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "string",
						Size:     16,
						Align:    8,
						Tag:      "test-tag",
						Exported: true,
						Embedded: true,
						Doc:      []string{"fdoctest"},
						Comment:  []string{"fcomtest"},
					},
					{
						Name:  "test-2",
						Type:  "test_type",
						Size:  12,
						Align: 4,
					},
				},
			},
			r: []byte(`
Struct Name,Struct Doc,Struct Comment,Field Name,Field Type,Field Size,Field Align,Field Tag,Field Exported,Field Embedded,Field Doc,Field Comment
Test,doctest,comtest,test-1,string,16,8,test-tag,true,true,fdoctest,fcomtest
Test,doctest,comtest,test-2,test_type,12,4,,false,false,,
`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.fmt(tcase.st)
			// check
			// format actual and expected identically
			stract, strexp := strings.Trim(string(r), "\n"), strings.Trim(string(tcase.r), "\n")
			if stract != strexp {
				t.Errorf("actual %v doesn't equal to expected %v", stract, strexp)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
