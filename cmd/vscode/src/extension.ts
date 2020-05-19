// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import Gopiumcli from "./gopiumcli";
import Workspace from "./workspace";

export interface Settings {
	readonly presets: { [key: string]: Arguments };
	readonly flags: Flags;
	build(preset: string, path: string, pkg: string, regex: string): string[];
}

export interface Arguments {
	readonly walker: string;
	readonly strategies: string[];
}

export interface Flags {
	// target platform vars
	readonly c?: string;
	readonly a?: string;
	readonly l?: number[];
	// package parser vars
	readonly e?: string[];
	readonly f?: string[];
	// gopium walker vars
	readonly d?: boolean;
	readonly b?: boolean;
	// gopium printer vars
	readonly i?: number;
	readonly w?: number;
	readonly s?: boolean;
	// gopium global vars
	readonly t?: number;
}

export interface Runner {
	run(preset: string, configs: Settings, path: string, pkg: string, struct?: string): Promise<boolean>;
}

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
		let cli = new Gopiumcli();
		await cli.run("test", new Workspace(), "", "");
	});

	context.subscriptions.push(disposable);
}

// this method is called when your extension is deactivated
export function deactivate() { }
