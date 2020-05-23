import * as cp from 'child_process'
import { isAbsolute } from 'path'
import * as vscode from 'vscode'
import * as extension from './extension'
import { promptForMissingTool } from './vscode-go/src/goInstallTools'
import { getBinPath } from './vscode-go/src/util'

export default class Gopiumcli implements extension.Runner {
	async run(
		preset: string,
		settings: extension.Settings,
		path: string,
		pkg: string,
		struct?: string,
	): Promise<boolean> {
		return new Promise<boolean>(async (resolve, reject) => {
			// prepare out channel
			const out = vscode.window.createOutputChannel('gopium')
			out.clear()
			out.show(true)
			// build the regex
			let regex = /.*/
			if (struct != null) {
				regex = new RegExp(`^${struct}$`)
			}
			// prepare final args and cli
			let gopium = getBinPath('gopium')
			if (!isAbsolute(gopium)) {
				console.log('AAAAA', gopium)
				await promptForMissingTool('gopium')
				gopium = getBinPath('gopium')
			}
			let args = settings.build(preset, path, pkg, regex.source)
			// start gopium process
			out.appendLine(`executing gopium cli: ${gopium} ${args.join(' ')}`)
			const proc = cp.spawn(gopium, args)
			proc.stdout.on('data', (chunk) => out.append(chunk.toString()))
			proc.stderr.on('data', (chunk) => out.append(chunk.toString()))
			proc.on('close', (code, signal) => {
				out.appendLine('gopium cli finished execution')
				// out.hide();
				resolve(code === 0)
			})
		})
	}
}
