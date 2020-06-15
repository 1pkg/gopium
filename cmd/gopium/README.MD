# Gopium 🌺 CLI

Gopium CLI is the cobra cli tool for [Gopium](../../README.MD) that was designed to automate and simplify some common performance transformations for structs, such as:

- cpu cache alignment
- memory packing
- false sharing guarding
- auto annotation
- generic fields management, etc.

In order to use gopium cli you need to provide at least package name (full package name is expected), list of strategies which is applied one by one and single walker.
Outcome of execution is fully defined by list of strategies and walker combination.
List of strategies modifies structs inside the package, walker facilitates and insures, that outcome is formatted written to one of provided destinations.

Gopium CLI uses next parameters schema:

```bash
gopium -flag_0 -flag_n walker package strategy_1 strategy_2 strategy_3 ...
```

where:

- walker defines destination for execution outcome, it should contain one value from [full walkers list](#walkers-and-formatters).
- package defines target package name for execution, note that full package name is expected.
- strategies [1..n] define transformations list that should be applied to package, they should contain at least one values from [full transformations list](#strategies-and-transformations).
- flags [0..n] define modificators for transfromations and walker, see [full flags list](#options-and-flags).

Real Gopium CLI commands examples that have been used to create [examples](../examples):

```bash
###
# creates new folder `transaction_gopium` inside /Users/1pkg/proj/src/1pkg/gopium/examples/
#
# with all inside structures transformed with:
# - rearrange structure fields in order to get optimal memory layout
# - add structure comment size annotation
# - add gopium tags to structure to save this list of transformation
# (it saves just list of transformation inside the tags and doesn't affect previous transformation at all, this way this list of strategies could be reused later)
# (force means override gopium tags even if previous exist)
#
# -p defines full path to package `1pkg/gopium/examples/transaction`
# it could be omitted if full package name would be used instead `gopium ast_go_tree 1pkg/gopium/examples/transaction memory_pack ...`
###
gopium ast_go_tree transaction memory_pack comment_struct_annotate add_tag_group_force -p /Users/1pkg/proj/src/1pkg/gopium/examples/transaction
```

```bash
###
# creates new folder `transaction_gopium` inside /Users/1pkg/proj/src/1pkg/gopium/examples/
#
# with all structures transformed with:
# - filter explicit fields paddings
# - add cache line #1 false sharing paddings after each field
# - add structure comment size annotation
# - add gopium tags to structure to save this list of transformation
# (it saves just list of transformation inside the tags and doesn't affect previous transformation at all, this way this list of strategies could be reused later)
# (force means override gopium tags even if previous exist)
#
# -p defines full path to package `1pkg/gopium/examples/transaction`
# it could be omitted if full package name would be used instead `gopium ast_go_tree 1pkg/gopium/examples/transaction memory_pack ...`
###
gopium ast_go_tree transaction filter_pads false_sharing_cpu_l1 add_tag_group_force -p /Users/1pkg/proj/src/1pkg/gopium/examples/transaction
```

```bash
###
# creates new folder `transaction_gopium` inside /Users/1pkg/proj/src/1pkg/gopium/examples/
#
# with all structures transformed to:
# - filter explicit fields paddings
# - add explicit system alignment padding after each field
# - add explicit padding in the end of structure to round its size to cache line #1
# - add structure comment size annotation
# - add gopium tags to structure to save this list of transformation
# (it saves just list of transformation inside the tags and doesn't affect previous transformation at all, this way this list of strategies could be reused later)
# (force means override gopium tags even if previous exist)
#
# -p defines full path to package `1pkg/gopium/examples/transaction`
# it could be omitted if full package name would be used instead `gopium ast_go_tree 1pkg/gopium/examples/transaction memory_pack ...`
#
# -l -l -l defines cache lines #1 #2 #3 sizes that are used for transformations
###
gopium -l 64 -l 64 -l 64 ast_go_tree transaction filter_pads explicit_padings_system_alignment cache_rounding_cpu_l1 comment_struct_annotate add_tag_group_force -p /Users/1pkg/proj/src/1pkg/gopium/examples/transaction
```

Gopium also has rich vscode extension to provide better experience for usage and simplify interactions with cli tool, [see more](../vscode/README.MD).

## Walkers and Formatters

Gopium provides next walkers:

- ast_go (directly syncs result as go code to orinal file)
- ast_go_tree (directly syncs result as go code to copy package)
- ast_std (prints result as go code to stdout)
- ast_gopium (directly syncs result as go code to copy gopium files)
- json_file (prints json encoded results to single file inside package directory)
- xml_file (prints xml encoded results to single file inside package directory)
- csv_file (prints csv encoded results to single file inside package directory)
- md_table_file (prints markdown table encoded results to single file inside package directory)
- size_align_md_table_file (prints markdown encoded table of sizes and aligns difference for results to single file inside package directory)
- fields_html_table_file (prints html encoded table of fields difference for results to single file inside package directory)

## Strategies and Transformations

Gopium provides next strategies:

- process_tag_group (uses gopium fields tags annotation in order to process different set of strategies on different groups and then combine results in single struct result)
- memory_pack (rearranges structure fields to obtain optimal memory utilization)
- memory_unpack (rearranges structure field list to obtain inflated memory utilization)
- cache_rounding_cpu_l1 (fits structure into cpu cache line #1 by adding bottom partial rounding cpu cache padding)
- cache_rounding_cpu_l2 (fits structure into cpu cache line #2 by adding bottom partial rounding cpu cache padding)
- cache_rounding_cpu_l3 (fits structure into cpu cache line #3 by adding bottom partial rounding cpu cache padding)
- full_cache_rounding_cpu_l1 (fits structure into full cpu cache line #1 by adding bottom rounding cpu cache padding)
- full_cache_rounding_cpu_l2 (fits structure into full cpu cache line #2 by adding bottom rounding cpu cache padding)
- full_cache_rounding_cpu_l3 (fits structure into full cpu cache line #3 by adding bottom rounding cpu cache padding)
- false_sharing_cpu_l1 (guards structure from false sharing by adding extra cpu cache line #1 paddings for each structure field)
- false_sharing_cpu_l2 (guards structure from false sharing by adding extra cpu cache line #1 paddings for each structure field)
- false_sharing_cpu_l3 (guards structure from false sharing by adding extra cpu cache line #1 paddings for each structure field)
- separate_padding_system_alignment_top (separates structure with extra system alignment padding by adding the padding at the top)
- separate_padding_cpu_l1_top (separates structure with extra cpu cache line #1 padding by adding the padding at the top)
- separate_padding_cpu_l2_top (separates structure with extra cpu cache line #2 padding by adding the padding at the top)
- separate_padding_cpu_l3_top (separates structure with extra cpu cache line #3 padding by adding the padding at the top)
- separate_padding_system_alignment_bottom (separates structure with extra system alignment padding by adding the padding at the bottom)
- separate_padding_cpu_l1_bottom (separates structure with extra cpu cache line #1 padding by adding the padding at the bottom)
- separate_padding_cpu_l2_bottom (separates structure with extra cpu cache line #2 padding by adding the padding at the bottom)
- separate_padding_cpu_l3_bottom (separates structure with extra cpu cache line #3 padding by adding the padding at the bottom)
- explicit_padings_system_alignment (explicitly aligns each structure field to system alignment padding by adding missing paddings for each field)
- explicit_padings_type_natural (explicitly aligns each structure field to max type alignment padding by adding missing paddings for each field)
- add_tag_group_soft (adds gopium fields tags annotation if no previous annotation found)
- add_tag_group_force (adds gopium fields tags annotation if previous annotation found overwrites it)
- add_tag_group_discrete (discretely adds gopium fields tags annotation if no previous annotation found)
- add_tag_group_force_discrete (discretely adds gopium fields tags annotation if previous annotation found overwrites it)
- remove_tag_group (removes gopium fields tags annotation)
- doc_fields_annotate (adds align and size doc annotation for each structure field)
- comment_fields_annotate adds align and size comment annotation for each structure field)
- doc_struct_annotate (adds aggregated align and size doc annotation for structure)
- comment_struct_annotate (adds aggregated align and size comment annotation for structure)
- name_lexicographical_ascending (sorts fields accordingly to their names in ascending order)
- name_lexicographical_descending (sorts fields accordingly to their names descending order)
- type_lexicographical_ascending (sorts fields accordingly to their types in ascending order)
- type_lexicographical_descending (sorts fields accordingly to their types in descending order)
- filter_pads (filters out all structure padding fields)
- ignore (does nothing by returning original structure)

## Gopium and Tags

Gopium CLI usues structure fields tags strategies for two purposes:

- save transformations list for later reuse
- to have abillity to apply different strategies groups inside one structure

For example by applying `gopium -l 64 -l 128 ast_stg transaction process_tag_group` to:

```go
// transaction defines business transaction
type transaction struct {
	amount   float64 `gopium:"group:critical;filter_pads,false_sharing_cpu_l1,separate_padding_cpu_l2_bottom"`
	serial   uint64  `gopium:"group:other;memory_pack"`
	discount float64 `gopium:"group:critical;filter_pads,false_sharing_cpu_l1,separate_padding_cpu_l2_bottom"`
	void     bool    `gopium:"group:other;memory_pack"`
	skip     bool    `gopium:"group:other;memory_pack"`
}
```

we can expect next Gopium CLI groups result:

```go
// transaction defines business transaction
type transaction struct {
        amount   float64 `gopium:"group:critical;filter_pads,false_sharing_cpu_l1,separate_padding_cpu_l2_bottom"`
        _        [56]byte
        discount float64 `gopium:"group:critical;filter_pads,false_sharing_cpu_l1,separate_padding_cpu_l2_bottom"`
        _        [56]byte
        _        [64]byte
        serial   uint64 `gopium:"group:other;memory_pack"`
        void     bool   `gopium:"group:other;memory_pack"`
        skip     bool   `gopium:"group:other;memory_pack"`
}
```

In this example fields `amount` and `discount` (group `critical`) were processed independently from fields `serial`, `void`, `skip` (group `other`) and two different sets of transformations were applied to each of them.

## Additional Notes

- it might be useful to use filter_pads in pipes with other strategies to clean paddings first.
- process_tag_group currently supports only next fields tags annotation formats:
  - gopium:"stg,stg,stg" processed as default group
  - gopium:"group:def;stg,stg,stg" processed as named group
- by specifying tag_type you can automatically generate fields tags annotation suitable for process_tag_group.
- `add_tag_*` strategies just add list of applied transformations to structure fields tags and NOT change results of other strategies, you can execute `process_tag_group` strategy afterwards to reuse saved strategies list.

## Options and Flags

|             Full             | Short |   Type   |     Default     | Description                                                                                                                                                                                                                                        |
| :--------------------------: | :---: | :------: | :-------------: | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|      --target_compiler       |  -c   |  string  |       gc        | Gopium target platform compiler, possible values are: gc or gccgo.                                                                                                                                                                                 |
|    --target_architecture     |  -a   |  string  |      amd64      | Gopium target platform architecture, possible values are: 386, arm, arm64, amd64, mips, etc.                                                                                                                                                       |
| target_cpu_cache_lines_sizes |  -l   |  []int   |  [64, 64, 64]   | Gopium target platform CPU cache line sizes in bytes, cache line size is set one by one l1,l2,l3,... For now only 3 lines of cache are supported by strategies.                                                                                    |
|        --package_path        |  -p   |  string  | src/{{package}} | Gopium go package path, either relative or absolute path to root of the package is expected. To obtain full path from relative, package path is concatenated with current GOPATH env var. Template {{package}} part is replaced with package name. |
|     --package_build_envs     |  -e   | []string |       [ ]       | Gopium go package build envs, additional list of building envs is expected.                                                                                                                                                                        |
|    --package_build_flags     |  -f   | []string |       [ ]       | Gopium go package build flags, additional list of building flags is expected.                                                                                                                                                                      |
|       --walker_regexp        |  -r   |  string  |       .\*       | Gopium walker regexp, regexp that defines which structures are subjects for visiting. Visiting is done only if structure name matches the regexp.                                                                                                  |
|        --walker_deep         |  -d   |   bool   |      true       | Gopium walker deep flag, flag that defines type of nested scopes visiting. By default it visits all nested scopes.                                                                                                                                 |
|       --walker_backref       |  -b   |   bool   |      true       | Gopium walker backref flag, flag that defines type of names referencing. By default any previous visited types have affect on future relevant visits.                                                                                              |
|       --printer_indent       |  -i   |   int    |        0        | Gopium printer width of tab, defines the least code indent.                                                                                                                                                                                        |
|     --printer_tab_width      |  -w   |   int    |        8        | Gopium printer width of tab, defines width of tab in spaces for printer.                                                                                                                                                                           |
|     --printer_use_space      |  -s   |   bool   |      false      | Gopium printer use space flag, flag that defines if all formatting should be done by spaces.                                                                                                                                                       |
|     --printer_use_gofmt      |  -g   |   bool   |      true       | Gopium printer use gofmt flag, flag that defines if canonical gofmt tool should be used for formatting. By default it is used and overrides other printer formatting parameters.                                                                   |
|           timeout            |  -t   |   int    |        0        | Gopium global timeout of cli command in seconds, considered only if value greater than 0.                                                                                                                                                          |
