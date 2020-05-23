// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode'
import Codelens from './codelens'
import Gopiumcli from './gopiumcli'
import { patch } from './patch'
import * as tools from './tools'
import { GO_MODE } from './vscode-go/src/goMode'
import Workspace from './workspace'

export interface Settings {
	readonly presets: { [key: string]: Arguments }
	readonly flags: Flags
	build(preset: string, path: string, pkg: string, regex: string): string[]
}

export interface Arguments {
	readonly walker: string
	readonly strategies: string[]
}

export interface Flags {
	// target platform vars
	readonly c?: string
	readonly a?: string
	readonly l?: number[]
	// package parser vars
	readonly e?: string[]
	readonly f?: string[]
	// gopium walker vars
	readonly d?: boolean
	readonly b?: boolean
	// gopium printer vars
	readonly i?: number
	readonly w?: number
	readonly s?: boolean
	// gopium global vars
	readonly t?: number
}

export interface Runner {
	run(preset: string, configs: Settings, path: string, pkg: string, struct?: string): Promise<boolean>
}

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export async function activate(context: vscode.ExtensionContext) {
	patch()

	await tools.offerTools()

	context.subscriptions.push(
		vscode.workspace.onDidChangeConfiguration((e: vscode.ConfigurationChangeEvent) => {
			if (!e.affectsConfiguration('gopium')) {
				return
			}
			vscode.commands.executeCommand('workbench.action.reloadWindow')
		}),

		vscode.languages.registerCodeLensProvider(GO_MODE, new Codelens(new Workspace())),

		vscode.commands.registerCommand('gopium', async (preset, path, pkg, struct) => {
			let gopiumcli = new Gopiumcli()
			await gopiumcli.run(preset, new Workspace(), path, pkg, struct)
		}),
	)
}

// this method is called when your extension is deactivated
export function deactivate() {}
