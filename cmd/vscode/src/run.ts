import * as cp from 'child_process';
import * as vscode from 'vscode';
import Confs from './confs';

export default async function run(confs: Confs): Promise<boolean> {
	return await new Promise<boolean>(async (resolve, reject) => {
		// prepare out channel
		const out = vscode.window.createOutputChannel('gopium');
		out.clear();
		out.show(true);
		// fill the gopium args
		let args: Array<string> = [
			confs.args.walker,
			confs.args.package,
			...confs.args.startegies,
		];
		// fill the gopium flags
		let flags: Array<string> = [];
		let fobj = confs.flags as { [key: string]: any }
		for (const key in fobj) {
			let val = fobj[key];
			if (val != null) {
				flags.push(`-${key}`, String(val));
			}
		}
		// prepare final args and cmd
		args = [...flags, ...args]
		let cmd = 'gopium ' + args.join(' ')
		// start gopium process
		out.appendLine(`executing gopium cli: ${cmd}`);
		const proc = cp.spawn("gopium", args);
		proc.stdout.on('data', (chunk) => out.append(chunk.toString()));
		proc.stderr.on('data', (chunk) => out.append(chunk.toString()));
		proc.on('close', (code, signal) => {
			out.appendLine("gopium cli finished execution");
			resolve(code === 0);
		});
	});
}