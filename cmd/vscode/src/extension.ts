// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import * as conf from './confs';
import run from "./run";

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	// Use the console to output diagnostic information (console.log) and errors (console.error)
	// This line of code will only be executed once when your extension is activated
	console.log('Congratulations, your extension "gopium" is now active!');

	// The command has been defined in the package.json file
	// Now provide the implementation of the command with registerCommand
	// The commandId parameter must match the command field in package.json
	let disposable = vscode.commands.registerCommand('gopium.helloWorld', async () => {
		const confs: conf.default = {
			args: {
				package: "1pkg/gopium",
				walker: "json_file",
				startegies: ["explicit_padings_type_natural", "add_tag_group_soft"],
			} as conf.Args,
			flags: {
				w: 4,
			} as conf.Flags,
		}
		await run(confs);
	});

	context.subscriptions.push(disposable);
}

// this method is called when your extension is deactivated
export function deactivate() { }
