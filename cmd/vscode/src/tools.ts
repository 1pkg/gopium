// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as path from 'path'
import * as vscode from 'vscode'
import * as install from './vscode-go/src/goInstallTools'
import * as tools from './vscode-go/src/goTools'
import * as util from './vscode-go/src/util'

export async function offerTools() {
	const goVersion = await util.getGoVersion()
	let all = tools.getConfiguredTools(goVersion)
	let missing = all.filter((t) => isMissing(t))
	if (missing.length > 0) {
		vscode.window
			.showInformationMessage('Required gopium extension tools have been failed to found.', {
				title: 'Try to install tools',
				command() {
					install.installTools(missing, goVersion)
				},
			})
			.then((selection) => {
				if (selection) {
					selection.command()
				}
			})
	} else {
		vscode.window
			.showInformationMessage('Required gopium extension tools have been found successfully.', {
				title: 'Try to update tools',
				command() {
					install.installTools(all, goVersion)
				},
			})
			.then((selection) => {
				if (selection) {
					selection.command()
				}
			})
	}
}

export async function getInstallTool(tname: string) {
	let tool = tools.getTool(tname)
	if (isMissing(tool)) {
		await install.installTools([tool], await util.getGoVersion())
	}
	return util.getBinPath(tname)
}

export function isMissing(tool: tools.Tool): boolean {
	return !path.isAbsolute(util.getBinPath(tool.name))
}
