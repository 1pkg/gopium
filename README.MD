<p align="center">
    <img src="https://raw.githubusercontent.com/1pkg/gopim/master/gopher.png" alt="gopium"/>
</p>

# Gopium 🌺: Smart Go Structures Optimizer and Manager

[![lint](https://github.com/1pkg/gopium/workflows/lint/badge.svg)](https://github.com/1pkg/gopium/actions?query=workflow%3Alint+branch%3Amaster+)
[![build](https://github.com/1pkg/gopium/workflows/build/badge.svg)](https://github.com/1pkg/gopium/actions?query=workflow%3Abuild+branch%3Amaster+)
[![test](https://github.com/1pkg/gopium/workflows/test/badge.svg)](https://github.com/1pkg/gopium/actions?query=workflow%3Atest+branch%3Amaster+)

## Introduction

Gopium is the tool that was designed to automate and simplify common structs performance and optimization transformations and fields management.

## Features

- Gopium provides numbers of common transformations for structs, such as: cpu cache alignment, memory packing, false sharing guarding, [and more](cmd/gopium/README.MD#strategies-and-transformations).
- All Gopium transformations could be combined together to build optimal actions for your specific code and needs.

```bash
gopium walker package strategy_1 strategy_2 strategy_3 ...
```

- All Gopium transformations results and artifacts could be differently formatted and routed to multiple outputs, among which are: direct overwriting or copying to go (code) ast, multiple options for code formatting and printing, comparison tables for results, common textual data serialization formats, [and more](cmd/gopium/README.MD#walkers-and-formatters).
- Gopium has rich vscode extension to provide better experience for usage and simplify interactions with tool, [see more](cmd/vscode/README.MD).
  ![](vscode.gif)
- Gopium supports target configs variety, in case if you need optimizations for specific platform or hardware, such as: target compiler, target architecture, target cpu cache lines sizes, [and more](cmd/gopium/README.MD#options-and-flags).
- Gopium has optional go structs tag literals to keep and reproduce gopium transformations inside your code base, also tags are intended to be used as structs fields grouping mechanism if you need to apply different transformation to different fields groups inside one structure, [see more](cmd/gopium/README.MD#strategies-and-transformations).

## Installation and Usage

To install Gopium VSCode Extension use [vscode marketplace](https://marketplace.visualstudio.com/items?itemName=golang.Gopium).

To install and update Gopium CLI use:

```bash
go get -u github.com/1pkg/gopium/cmd/gopium
```

Gopium CLI usage examples:

```bash
# exec transformations against all structs starting A inside 1pkg/gopium/gopium package
# and print results formatted as json array to standart output
gopium -r ^A json_std 1pkg/gopium/gopium filter_pads memory_pack separate_padding_cpu_l1_top separate_padding_cpu_l1_bottom
```

```bash
# exec transformations against all structs inside gopium package
# at /alternative/go/path/src/ path and write results back directly to go files
gopium ast_go gopium memory_pack -p /alternative/go/path/src/{{package}}
```

```bash
# exec transformations against all structs inside 1pkg/gopium/gopium package
# with target arm architecture and preset cache lines 8,8,8
# and write results to go package 1pkg/gopium/gopium copy
gopium -a arm -l 8 8 8 ast_go_tree 1pkg/gopium/gopium false_sharing_cpu_l1
```

## Examples Benchmarks and Docs

```
    goversion: go1.13.10
    goos: darwin
    goarch: amd64

    === simple business transaction aggregator example ===
    pkg: 1pkg/gopium/examples/transaction
    BenchmarkCompress-8   	      16	  98752864 ns/op	     230 B/op	       3 allocs/op

    === `gopium ast_go_tree 1pkg/gopium/examples/transaction memory_pack` applied ===
    pkg: 1pkg/gopium/examples/memory_pack/transaction
    BenchmarkCompress-8   	      19	  83438394 ns/op	     254 B/op	       3 allocs/op

    === `gopium ast_go_tree 1pkg/gopium/examples/transaction filter_pads false_sharing_cpu_l1` applied ===
    pkg: 1pkg/gopium/examples/false_sharing_cpu_l1/transaction
    BenchmarkCompress-8   	      20	  84066847 ns/op	     336 B/op	       3 allocs/op

    === `gopium ast_go_tree 1pkg/gopium/examples/transaction filter_pads explicit_padings_system_alignment cache_rounding_cpu_l1` applied ===
    pkg: 1pkg/gopium/examples/cache_rounding_cpu_l1/transaction
    BenchmarkCompress-8   	      13	 130450557 ns/op	     209 B/op	       3 allocs/op
```

See [examples](examples) for more examples and benchmarks details.  
See [docs](https://pkg.go.dev/github.com/1pkg/gopium?tab=doc) for godoc.

## Licence

Gopium is licensed under the MIT License.  
See [LICENSE](LICENSE) for the full license text.
