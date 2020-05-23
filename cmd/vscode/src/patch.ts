import * as vscode from 'vscode'
import * as main from './vscode-go/src/goMain'
import * as status from './vscode-go/src/goStatus'
import * as tools from './vscode-go/src/goTools'
import * as telemetry from './vscode-go/src/telemetry'

const pfunc = (): null => null

export function patch() {
	let pmain = main as any
	pmain.activate = pfunc
	pmain.deactivate = pfunc
	let ptelementry = telemetry as any
	for (const key in ptelementry) {
		if (key.includes('sendTelemetry')) {
			ptelementry[key] = pfunc
		}
	}
	let ptools = tools as any
	ptools.getConfiguredTools = (): tools.Tool[] => {
		return [
			{
				name: 'go-outline',
				importPath: 'github.com/ramya-rao-a/go-outline',
				isImportant: true,
				description: 'Go to symbol in file',
			},
			{
				name: 'gopium',
				importPath: 'github.com/1pkg/gopim',
				isImportant: true,
				description: 'Gopium struct manager package',
			},
		]
	}
	ptools.getTool = (name: string): tools.Tool => {
		switch (name) {
			case 'go-outline':
				return {
					name: 'go-outline',
					importPath: 'github.com/ramya-rao-a/go-outline',
					isImportant: true,
					description: 'Go to symbol in file',
				}
			case 'gopium':
				return {
					name: 'gopium',
					importPath: 'github.com/1pkg/gopim',
					isImportant: true,
					description: 'Gopium struct manager package',
				}
		}
		return null
	}
	let pstatus = status as any
	pstatus.outputChannel = vscode.window.createOutputChannel('gopium')
}
