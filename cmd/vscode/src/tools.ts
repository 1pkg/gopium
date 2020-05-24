// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as path from 'path'
import * as vscode from 'vscode'
import * as install from './vscode-go/src/goInstallTools'
import * as tools from './vscode-go/src/goTools'
import * as util from './vscode-go/src/util'

export function missing(): tools.Tool[] {
	let all = tools.getConfiguredTools(new util.GoVersion('any'))
	return all.filter((t) => !path.isAbsolute(util.getBinPath(t.name)))
}

export async function getb(name: string): Promise<string | null> {
	const bpath = util.getBinPath(name)
	if (!path.isAbsolute(util.getBinPath(name))) {
		await offer([tools.getTool(name)])
		return null
	}
	return bpath
}

export async function offer(tools: tools.Tool[]) {
	if (tools.length == 0) {
		return
	}
	const gov = await util.getGoVersion()
	vscode.window
		.showInformationMessage('Required gopium extension tools have been not found.', {
			title: 'Try to install tools',
			command() {
				install.installTools(tools, gov)
			},
		})
		.then((selection) => {
			if (selection) {
				selection.command()
			}
		})
}
