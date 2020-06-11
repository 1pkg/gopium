# Gopium 🌺 VSCode Extension

Gopium vscode extension for [Gopium](../../README.MD) provides better experience for usage and simplify interactions with [Gopium CLI](../gopium/README.MD).
![](../../vscode.gif)

## Requirements and Installation

To install gopium vscode extension go to [vscode marketplace](https://marketplace.visualstudio.com/items?itemName=golang.Gopium).

In order to work, extension requires latest version of [Gopium CLI](../gopium/README.MD) and [Goutline](https://github.com/1pkg/goutline) tools.  
Extension automatically tries to install required tools if they have been not detected in your system.

To install and update gopium cli use:

```bash
go get -u github.com/1pkg/gopium/cmd/gopium
```

To install and update goutline use:

```bash
go get -u github.com/1pkg/goutline
```

## Origin

Extension is based on and heavily using [vscode-go](https://github.com/microsoft/vscode-go) extension functionality and keeps it as [submodule](src/vscode-go). Such decision to vendor extension was made in favor to simplify origin Gopium VSCode Extension development by reusing common components like: goutline codelens, tools installation, etc.

Nevertheless it's impossible to just seamlessly reuse another extension to build your own, so a lot of functionality has been patched and changed, most noticeable: vscode-go telemetry is disabled, vscode-go extension entrypoits are removed, [go-outline](https://github.com/ramya-rao-a/go-outline) is replaced with [go-outline](https://github.com/1pkg/goutline), etc.
