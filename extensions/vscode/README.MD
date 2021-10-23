# Gopium 🌺 VSCode Extension

Gopium vscode extension for [Gopium](https://github.com/1pkg/gopium) provides better experience for usage and simplify interactions with Gopium CLI.
![](vscode.gif)

## Requirements and Installation

To install gopium vscode extension go to [vscode marketplace](https://marketplace.visualstudio.com/items?itemName=1pkg.gopium).

In order to work, extension requires latest version of [Gopium CLI](https://github.com/1pkg/gopium) and [Goutline](https://github.com/1pkg/goutline) tools.  
Extension automatically tries to install required tools if they have been not detected in your system.

To install and update gopium cli use:

```bash
go install github.com/1pkg/gopium@latest
```

To install and update goutline use:

```bash
go install github.com/1pkg/goutline@latest
```

## Transformations Customization

You can easily adjust list of Gopium transformations to your needs by updating `gopium.actions` settings in vscode, see also Gopium available [transformations](https://github.com/1pkg/gopium#strategies-and-transformations) and [walkers](https://github.com/1pkg/gopium#walkers-and-formatters).
By default extension uses next transformations presets:

`gopium pack | gopium cache | gopium pack cache | gopium false sharing | gopium pack verbose | gopium cache verbose`

which are defined by next settings in vscode:

```json
{
	"gopium.actions": [
		{
			"name": "pack",
			"walker": "ast_go",
			"strategies": ["memory_pack"]
		},
		{
			"name": "cache",
			"walker": "ast_go",
			"strategies": ["cache_rounding_cpu_l1_discrete"]
		},
		{
			"name": "pack cache",
			"walker": "ast_go",
			"strategies": ["filter_pads", "memory_pack", "explicit_paddings_type_natural", "cache_rounding_cpu_l1_discrete"]
		},
		{
			"name": "false sharing",
			"walker": "ast_go",
			"strategies": ["filter_pads", "false_sharing_cpu_l1"]
		},
		{
			"name": "pack verbose",
			"walker": "ast_go",
			"strategies": [
				"filter_pads",
				"memory_pack",
				"explicit_paddings_type_natural",
				"fields_annotate_comment",
				"struct_annotate_comment",
				"add_tag_group_soft"
			]
		},
		{
			"name": "cache verbose",
			"walker": "ast_go",
			"strategies": [
				"filter_pads",
				"explicit_paddings_type_natural",
				"cache_rounding_cpu_l1_discrete",
				"fields_annotate_comment",
				"struct_annotate_comment",
				"add_tag_group_soft"
			]
		}
	]
}
```

## Origin and Decisions

Extension is based on and heavily using [vscode-go](https://github.com/microsoft/vscode-go) extension functionality and keeps it as [submodule](https://github.com/1pkg/vscode-go). Such decision to vendor extension was made in favor to simplify origin Gopium VSCode Extension development by reusing common components like: goutline codelens, tools installation, etc.

Nevertheless it's impossible to just seamlessly reuse another extension to build your own, so a lot of functionality has been patched and changed, most noticeable: vscode-go telemetry is disabled, vscode-go extension entrypoits are removed, [go-outline](https://github.com/ramya-rao-a/go-outline) is replaced with [go-outline](https://github.com/1pkg/goutline), etc.

## Licence

Gopium is licensed under the MIT License.  
See [LICENSE](LICENSE) for the full license text.
