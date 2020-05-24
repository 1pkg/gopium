'use strict'

import * as cp from 'child_process'
import * as path from 'path'
import * as vscode from 'vscode'
import * as goinstall from './vscode-go/src/goInstallTools'
import * as gomain from './vscode-go/src/goMain'
import { GO_MODE } from './vscode-go/src/goMode'
import * as gooutline from './vscode-go/src/goOutline'
import * as gostatus from './vscode-go/src/goStatus'
import * as gotools from './vscode-go/src/goTools'
import * as gotelemetry from './vscode-go/src/telemetry'
import * as goutil from './vscode-go/src/util'

// out defines gopium extension out channel
const out = vscode.window.createOutputChannel('gopium')
// gsettings defines global gopium workspace setting
let gsettings: settings.Workspace
// patchgo defines vscode-go functionality patch function
function patchgo() {
	// dummy patch func
	// that simply does nothing
	const patchfunc = (): null => null
	// patch main go module
	// to disable activation
	let pmain = gomain as any
	pmain.activate = patchfunc
	pmain.deactivate = patchfunc
	// patch telemetry go module
	// to disable any event sending
	let ptelementry = gotelemetry as any
	for (const key in ptelementry) {
		if (key.includes('sendTelemetry')) {
			ptelementry[key] = patchfunc
		}
	}
	// patch tools go module
	// to replace tools with gopium tools
	let ptools = gotools as any
	ptools.getConfiguredTools = (): gotools.Tool[] => {
		// we need only gopium and outline
		return [
			{
				name: 'gopium',
				importPath: 'github.com/1pkg/gopim',
				isImportant: true,
				description: '',
			},
			{
				name: 'go-outline',
				importPath: 'github.com/ramya-rao-a/go-outline',
				isImportant: true,
				description: '',
			},
		]
	}
	ptools.getTool = (name: string): gotools.Tool => {
		// we need only gopium and outline
		switch (name) {
			case 'gopium':
				return {
					name: 'gopium',
					importPath: 'github.com/1pkg/gopim',
					isImportant: true,
					description: '',
				}
			case 'go-outline':
				return {
					name: 'go-outline',
					importPath: 'github.com/ramya-rao-a/go-outline',
					isImportant: true,
					description: '',
				}
		}
		return null
	}
	// patch status go module
	// to replace out chan with gopium chan
	let pstatus = gostatus as any
	pstatus.outputChannel = out
}
// activate root extension registration hook
export async function activate(context: vscode.ExtensionContext) {
	// patch vscode-go facilities at start
	patchgo()
	// try to offer to install missing tools
	await tools.Offer(await tools.Missing())
	// set global gopim setting
	gsettings = new settings.Workspace()
	// add gopim actions to vscode
	context.subscriptions.push(
		// on config update
		vscode.workspace.onDidChangeConfiguration((e: vscode.ConfigurationChangeEvent) => {
			// check that update touches gopium
			if (!e.affectsConfiguration('gopium')) {
				return
			}
			// and update global config
			gsettings = new settings.Workspace()
		}),
		// register gopium codelens
		vscode.languages.registerCodeLensProvider(GO_MODE, new Codelens()),
		// register gopium command
		vscode.commands.registerCommand(
			'gopium',
			async (preset, path, pkg, struct) => await gopiumcli.Run(preset, gsettings, path, pkg, struct),
		),
	)
}
// deactivate root extension registration hook
export function deactivate() {}
// Codelens gopium codelens provider implementation
class Codelens implements vscode.CodeLensProvider {
	// provideCodeLenses implements vscode provider interface
	public provideCodeLenses(
		document: vscode.TextDocument,
		token: vscode.CancellationToken,
	): vscode.CodeLens[] | Thenable<vscode.CodeLens[]> {
		// check that go-outline tool exists
		return tools.Getbin('go-outline').then((bin) => {
			if (bin != null) {
				// grab all document symbols from outline
				const dsp = new gooutline.GoDocumentSymbolProvider(true)
				return dsp.provideDocumentSymbols(document, token).then((symbols) => {
					// first symbol should be package
					const pkg = symbols[0]
					if (pkg) {
						// then collect all codelens for package and structs
						const pdir = path.dirname(document.fileName)
						return [...this.package(pdir, pkg), ...this.structs(pdir, pkg)]
					}
					return []
				})
			}
		})
	}
	// package collects all package codelens
	private package(path: string, pkg: vscode.DocumentSymbol): vscode.CodeLens[] {
		let codelens: vscode.CodeLens[] = []
		// for all global settings actions presets
		// add new package codelens
		for (const preset in gsettings.presets) {
			codelens.push(
				new vscode.CodeLens(pkg.range, {
					title: `gopium ${preset}`,
					command: 'gopium',
					arguments: [preset, path, pkg.name],
					tooltip: `gopium package action ${preset}`,
				}),
			)
		}
		return codelens
	}
	// package collects all structs codelens
	private structs(path: string, pkg: vscode.DocumentSymbol): vscode.CodeLens[] {
		let codelens: vscode.CodeLens[] = []
		// go through all package children
		for (const child of pkg.children) {
			// if chils is struct collect codelens
			if (child.kind == vscode.SymbolKind.Struct) {
				// for all global settings actions presets
				// add new package codelens
				for (const preset in gsettings.presets) {
					codelens.push(
						new vscode.CodeLens(child.range, {
							title: `gopium ${preset}`,
							command: 'gopium',
							arguments: [preset, path, pkg.name, child.name],
							tooltip: `gopium struct action ${preset}`,
						}),
					)
				}
			}
		}
		return codelens
	}
}

// tools namespace holds all tools related functionality
namespace tools {
	// missing retus list of missing required tools
	export async function Missing(): Promise<gotools.Tool[]> {
		// grab all required tools
		const gov = await goutil.getGoVersion()
		let all = gotools.getConfiguredTools(gov)
		// filter all not found tools
		return all.filter((t) => !path.isAbsolute(goutil.getBinPath(t.name)))
	}
	// offer tries to offer requested list of tools
	export async function Offer(tools: gotools.Tool[]) {
		if (tools.length > 0) {
			const gov = await goutil.getGoVersion()
			vscode.window
				.showInformationMessage('Required gopium extension tools have not been found.', {
					title: 'Try to install tools',
					command: () => goinstall.installTools(tools, gov), // use go intall facilities
				})
				.then((selection) => {
					// if section has been chosen run its command
					selection && selection.command()
				})
		}
	}
	// getbin tries to get binary tool path
	// if no tool has been found it offers tool
	// and returns null binary path
	export async function Getbin(name: string): Promise<string | null> {
		// use go util helper
		const bin = goutil.getBinPath(name)
		// if no tool has been found
		if (!path.isAbsolute(bin)) {
			// try to offer the toool
			await Offer([gotools.getTool(name)])
			return null
		}
		return bin
	}
}

// settings namespace holds all settings related functionality
namespace settings {
	// Arguments defines list of gopium walker and strategies arguments
	export interface Arguments {
		// gopium walker field
		readonly walker: string
		// gopium strategies list field
		readonly strategies: string[]
	}
	// Flags defines list of gopium flags
	export interface Flags {
		// target platform fields
		readonly c?: string
		readonly a?: string
		readonly l?: number[]
		// package parser fields
		readonly e?: string[]
		readonly f?: string[]
		// gopium walker fields
		readonly d?: boolean
		readonly b?: boolean
		// gopium printer fields
		readonly i?: number
		readonly w?: number
		readonly s?: boolean
		// gopium global fields
		readonly t?: number
	}
	// Workspace defines vscode workspace extension settings implementation
	export class Workspace {
		// presets defines vscode workspace configs actions presets list
		readonly presets: { [key: string]: Arguments }
		// flags defines vscode workspace configs flags list
		readonly flags: Flags
		// constructor acquires workspace configs
		public constructor() {
			// grab root and actions presets workspace configs
			let root = vscode.workspace.getConfiguration('gopium')
			let actions = root.get<any[]>('actions', [])
			// fill presets map from workspace configs
			this.presets = {}
			for (const action of actions) {
				this.presets[action.name] = {
					walker: action.walker,
					strategies: action.strategies,
				}
			}
			// fill flags from workspace configs
			this.flags = {
				// target platform vars
				c: root.get<string>('target_compiler'),
				a: root.get<string>('target_architecture'),
				l: root.get<number[]>('target_cpu_cache_line_sizes'),
				// package parser vars
				e: root.get<string[]>('package_build_envs'),
				f: root.get<string[]>('package_build_flags'),
				// gopium walker vars
				d: root.get<boolean>('walker_deep'),
				b: root.get<boolean>('walker_backref'),
				// gopium printer vars
				i: root.get<number>('printer_indent'),
				w: root.get<number>('printer_tab_width'),
				s: root.get<boolean>('printer_use_space'),
				// gopium global vars
				t: root.get<number>('timeout'),
			}
		}
		// Build simply builds acquired workspace configs
		// and provided configs to gopium args string
		public Build(preset: string, path: string, pkg: string, regex: string): string[] {
			// if such preset doesn't exist build nothing
			if (!(preset in this.presets)) {
				return []
			}
			// get the preset arguments
			let pargs = this.presets[preset]
			// and fill the gopium args
			let args: Array<string> = [pargs.walker, pkg, ...pargs.strategies]
			// fill the gopium flags
			let flags: Array<string> = []
			// push provided parameters first
			if (path != '') {
				flags.push('-p', path)
			}
			if (regex != '') {
				flags.push('-r', regex)
			}
			// then collect all flags vaues
			let fobj = this.flags as { [key: string]: any }
			for (const fkey in fobj) {
				// skip null values
				let fval = fobj[fkey]
				if (fval == null) {
					continue
				}
				// depends on flag type
				// do different serialization
				if (typeof fval === 'boolean' && fval != false) {
					flags.push(`-${fkey}`)
				} else if (Array.isArray(fval)) {
					flags.push(`-${fkey}`, fval.join(' '))
				} else {
					flags.push(`-${fkey}`, String(fval))
				}
			}
			// concat all settings together
			return [...flags, ...args]
		}
	}
}

// gopiumcli namespace holds all gomiump cli related functionality
namespace gopiumcli {
	// run starts gopium cli execution
	// on settings with presets
	// using path, package and struct
	export async function Run(
		preset: string,
		settings: settings.Workspace,
		path: string,
		pkg: string,
		struct?: string,
	): Promise<boolean> {
		// grab gopium cli tool binary
		let gopium = await tools.Getbin('gopium')
		if (gopium == null) {
			return false
		}
		// build and prepare the regex
		let regex = /.*/
		if (struct != null) {
			regex = new RegExp(`^${struct}$`)
		}
		// prepare final cli args
		let args = settings.Build(preset, path, pkg, regex.source)
		// refresh out channel
		out.clear()
		out.show(true)
		out.appendLine(`gopium 🌺: ${gopium} ${args.join(' ')}`)
		// start gopium process
		// and subscribe to all outs
		const proc = cp.spawn(gopium, args)
		proc.stdout.on('data', (chunk) => out.append(chunk.toString()))
		proc.stderr.on('data', (chunk) => out.append(chunk.toString()))
		proc.on('error', (err) => out.appendLine(err.message))
		return true
	}
}
