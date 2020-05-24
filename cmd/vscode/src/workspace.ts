import * as vscode from 'vscode'
import * as extension from './extension'

// Workspace defines vscode workspace extension settings implementation
export default class Workspace implements extension.Settings {
	// presets defines vscode workspace configs actions presets list
	readonly presets: { [key: string]: extension.Arguments }
	// flags defines vscode workspace configs flags list
	readonly flags: extension.Flags

	// constructor acquires workspace configs
	constructor() {
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

	// build simply builds acquired workspace configs
	// and provided configs to gopium args string
	build(preset: string, path: string, pkg: string, regex: string): string[] {
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
