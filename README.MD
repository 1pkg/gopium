# Gopium 🌺: Smart Go Structures Optimizer and Manager

[![lint](https://github.com/1pkg/gopium/workflows/lint/badge.svg)](https://github.com/1pkg/gopium/actions?query=workflow%3Alint+branch%3Amaster+)
[![build](https://github.com/1pkg/gopium/workflows/build/badge.svg)](https://github.com/1pkg/gopium/actions?query=workflow%3Abuild+branch%3Amaster+)
[![test](https://github.com/1pkg/gopium/workflows/test/badge.svg)](https://github.com/1pkg/gopium/actions?query=workflow%3Atest+branch%3Amaster+)

[![gopium](gopher.png)](https://github.com/1pkg/gopium)

![](vscode.gif)

## Introduction

Gopium is the tool that was designed to automate and simplify some common performance transformations for structs, as:

- cpu cache alignment
- memory packing
- false sharing guarding
- auto annotation
- generic fields management, etc.

To install gopium cli:

```bash
go get -u github.com/1pkg/gopium/cmd/gopium
```

Gopium usage example:

```bash
gopium -E -r ^A json_std 1pkg/gopium filter_pads memory_pack separate_padding_cpu_l1_top separate_padding_cpu_l1_bottom
```

## Features

- Gopium supports numbers of common transformations for stucts, which could be combined and piped together see full list of them TODO
- Results of all actions could be routed to multiple places among which are gocode itself TODO
- Gopium has vscode extension to provide better user expirience and simplify interactions

## Contribution

Do you have an idea to improve Gopium? -> [Create an issue](https://github.com/1pkg/gopium/issues/new/choose)

Do you have already coded something for Gopium? -> Create a pull request.

Did you discover a bug? -> [Create an issue](https://github.com/1pkg/gopium/issues/new/choose)

## Licence

Gopium is licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.
