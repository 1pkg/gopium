package fmtio

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/tests"
)

func TestBytes(t *testing.T) {
	// prepare
	table := map[string]struct {
		fmt gopium.Xbytes
		f   collections.Flat
		r   []byte
		err error
	}{
		"json should return expected result for empty collection": {
			fmt: Jsonb,
			f:   collections.Flat{},
			r:   []byte(`[]`),
		},
		"json should return expected result for empty struct in collection": {
			fmt: Jsonb,
			f:   collections.Flat{"test": gopium.Struct{}},
			r: []byte(`
[
	{
		"Name": "",
		"Doc": null,
		"Comment": null,
		"Fields": null
	}
]
`),
		},
		"json should return expected result for non collection": {
			fmt: Jsonb,
			f: collections.Flat{
				"test-2": gopium.Struct{
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
				"test-1": gopium.Struct{
					Name: "Test-1",
					Fields: []gopium.Field{
						{
							Name:  "test-3",
							Type:  "test",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
			r: []byte(`
[
	{
		"Name": "Test-1",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "test-3",
				"Type": "test",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	},
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
]
`),
		},
		"xml should return expected result for empty collection": {
			fmt: Xmlb,
			f:   collections.Flat{},
			r:   []byte(``),
		},
		"xml should return expected result for empty struct in collection": {
			fmt: Xmlb,
			f:   collections.Flat{"test": gopium.Struct{}},
			r: []byte(`
<Struct>
	<Name></Name>
</Struct>
`),
		},
		"xml should return valid expected for non collection": {
			fmt: Xmlb,
			f: collections.Flat{
				"test-2": gopium.Struct{
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
				"test-1": gopium.Struct{
					Name: "Test-1",
					Fields: []gopium.Field{
						{
							Name:  "test-3",
							Type:  "test",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
			r: []byte(`
<Struct>
	<Name>Test-1</Name>
	<Fields>
		<Name>test-3</Name>
		<Type>test</Type>
		<Size>1</Size>
		<Align>1</Align>
		<Tag></Tag>
		<Exported>false</Exported>
		<Embedded>false</Embedded>
	</Fields>
</Struct>
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
		"csv should return expected result for empty collection": {
			fmt: Csvb(Buffer()),
			f:   collections.Flat{},
			r:   []byte(``),
		},
		"csv should return expected result for empty struct in collection": {
			fmt: Csvb(Buffer()),
			f:   collections.Flat{"test": gopium.Struct{}},
			r: []byte(`
Struct Name,Struct Doc,Struct Comment,Field Name,Field Type,Field Size,Field Align,Field Tag,Field Exported,Field Embedded,Field Doc,Field Comment
`),
		},
		"csv should return error on writer error": {
			fmt: Csvb(&tests.RWC{Werr: errors.New("test")}),
			f: collections.Flat{
				"test-2": gopium.Struct{
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
				"test-1": gopium.Struct{
					Name: "Test-1",
					Fields: []gopium.Field{
						{
							Name:  "test-3",
							Type:  "test",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
			err: errors.New("test"),
		},
		"csv should return expected result for non empty collection": {
			fmt: Csvb(Buffer()),
			f: collections.Flat{
				"test-2": gopium.Struct{
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
				"test-1": gopium.Struct{
					Name: "Test-1",
					Fields: []gopium.Field{
						{
							Name:  "test-3",
							Type:  "test",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
			r: []byte(`
Struct Name,Struct Doc,Struct Comment,Field Name,Field Type,Field Size,Field Align,Field Tag,Field Exported,Field Embedded,Field Doc,Field Comment
Test-1,,,test-3,test,1,1,,false,false,,
Test,doctest,comtest,test-1,string,16,8,test-tag,true,true,fdoctest,fcomtest
Test,doctest,comtest,test-2,test_type,12,4,,false,false,,
`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.fmt(tcase.f.Sorted())
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// format actual and expected identically
			actual := strings.Trim(string(r), "\n")
			expected := strings.Trim(string(tcase.r), "\n")
			if err == nil && !reflect.DeepEqual(actual, expected) {
				t.Errorf("actual %v doesn't equal to expected %v", actual, expected)
			}
		})
	}
}
