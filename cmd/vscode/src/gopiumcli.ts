import * as cp from 'child_process';
import * as vscode from 'vscode';
import * as extension from './extension';

export default class Gopiumcli implements extension.Runner {
	async run(preset: string, settings: extension.Settings, path: string, pkg: string, struct?: string): Promise<boolean> {
		return new Promise<boolean>(async (resolve, reject) => {
			// prepare out channel
			const out = vscode.window.createOutputChannel('gopium');
			out.clear();
			out.show(true);
			// build the regex
			let regex = /.*/;
			if (struct != null) {
				regex = new RegExp(`^${struct}$`);
			}
			// prepare final args and cmd
			let args = settings.build(preset, path, pkg, regex.source)
			let cmd = 'gopium ' + args.join(' ')
			// start gopium process
			out.appendLine(`executing gopium cli: ${cmd}`);
			const proc = cp.spawn("gopium", args);
			proc.stdout.on('data', (chunk) => out.append(chunk.toString()));
			proc.stderr.on('data', (chunk) => out.append(chunk.toString()));
			proc.on('close', (code, signal) => {
				out.appendLine("gopium cli finished execution");
				// out.hide();
				resolve(code === 0);
			});
		});
	}
}