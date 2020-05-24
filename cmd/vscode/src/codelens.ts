'use strict'

import * as path from 'path'
import * as vscode from 'vscode'
import * as extension from './extension'
import * as tools from './tools'
import { GoDocumentSymbolProvider } from './vscode-go/src/goOutline'
import { GoRunTestCodeLensProvider } from './vscode-go/src/goRunTestCodelens'

export default class Codelens extends GoRunTestCodeLensProvider {
	private settings: extension.Settings

	public constructor(settings: extension.Settings) {
		super()
		this.settings = settings
	}

	public provideCodeLenses(
		document: vscode.TextDocument,
		token: vscode.CancellationToken,
	): vscode.CodeLens[] | Thenable<vscode.CodeLens[]> {
		return tools.getb('go-outline').then((bpath) => {
			if (bpath != null) {
				const dsp = new GoDocumentSymbolProvider(true)
				return dsp.provideDocumentSymbols(document, token).then((symbols) => {
					const pkg = symbols[0]
					if (!pkg) {
						return []
					}
					const pdir = path.dirname(document.fileName)
					return [...this.package(pdir, pkg), ...this.structs(pdir, pkg)]
				})
			}
		})
	}

	private package(path: string, pkg: vscode.DocumentSymbol): vscode.CodeLens[] {
		let codelens: vscode.CodeLens[] = []
		for (const preset in this.settings.presets) {
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

	private structs(path: string, pkg: vscode.DocumentSymbol): vscode.CodeLens[] {
		let codelens: vscode.CodeLens[] = []
		for (const struct of pkg.children) {
			if (struct.kind == vscode.SymbolKind.Struct) {
				for (const preset in this.settings.presets) {
					codelens.push(
						new vscode.CodeLens(struct.range, {
							title: `gopium ${preset}`,
							command: 'gopium',
							arguments: [preset, path, pkg.name, struct.name],
							tooltip: `gopium struct action ${preset}`,
						}),
					)
				}
			}
		}
		return codelens
	}
}
