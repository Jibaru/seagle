import { Loader2, Play, Square, Sparkles } from "lucide-react";
import type React from "react";
import { useEffect, useRef, useState } from "react";
import Editor, { type OnMount } from "@monaco-editor/react";
import type { editor } from "monaco-editor";
import { Button } from "./ui/button";
import { GenerateQuery } from "../../wailsjs/go/handlers/GenQueryHandler";
import { useActiveConnectionStore } from "../store/ActiveConnectionStore";
import { useTheme } from "../contexts/ThemeContext";

interface SqlEditorProps {
	value: string;
	onChange: (value: string) => void;
	onExecute: (query: string) => void;
	onStop?: () => void;
	isExecuting?: boolean;
	database?: string;
}

export const SqlEditor: React.FC<SqlEditorProps> = ({
	value,
	onChange,
	onExecute,
	onStop,
	isExecuting = false,
	database,
}) => {
	const { state: activeConnection } = useActiveConnectionStore();
	const { actualTheme } = useTheme();
	const [selectedText, setSelectedText] = useState("");
	const [isGenerating, setIsGenerating] = useState(false);
	const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
	
	
	const getSelectedText = (editor: editor.IStandaloneCodeEditor) => {
		const selection = editor.getSelection();
		if (selection && !selection.isEmpty()) {
			return (editor.getModel()?.getValueInRange(selection) || "").trim();
		}
		return "";
	};

	const handleSelectionChange = () => {
		if (editorRef.current) {
			const editor = editorRef.current;
			const selected = getSelectedText(editor);
			setSelectedText(selected);
		}
	};

	const handleExecute = async () => {
		const queryToExecute = selectedText || value;
		if (queryToExecute.trim()) {
			await executeQuery(queryToExecute);
		}
	};

	const executeQuery = async (queryOrPrompt: string) => {
		if (queryOrPrompt.trim().toLowerCase().startsWith("gen:")) {
			await handleGenerateQuery(queryOrPrompt.trim());
		} else {
			onExecute(queryOrPrompt.trim());
		}
	};

	const handleGenerateQuery = async (genCommand: string) => {
		if (!database) {
			alert("No database selected");
			return;
		}

		// Extract the prompt from "gen: prompt"
		const prompt = genCommand.substring(4).trim();
		if (!prompt) {
			alert("Please provide a prompt after 'gen:'");
			return;
		}

		if (!activeConnection.connectionId) {
			alert("No active connection available");
			return;
		}

		setIsGenerating(true);
		try {
			const response = await GenerateQuery({
				id: activeConnection.connectionId,
				database: database,
				prompt: prompt,
			});

			if (response.success && response.result) {
				// Find the gen: line and replace it with the generated query
				const lines = value.split('\n');
				const genLineIndex = lines.findIndex(line => 
					line.trim().toLowerCase().startsWith('gen:')
				);

				if (genLineIndex !== -1) {
					// Replace the gen: line with the generated query
					lines[genLineIndex] = response.result.generatedQuery;
					const newValue = lines.join('\n');
					onChange(newValue);
				} else {
					// If gen: line not found, append the generated query
					const newValue = value + response.result.generatedQuery;
					onChange(newValue);
				}
			} else {
				alert(`Query generation failed: ${response.message}`);
			}
		} catch (error) {
			console.error("Error generating query:", error);
			alert("Failed to generate query. Please try again.");
		} finally {
			setIsGenerating(false);
		}
	};

	const handleEditorDidMount: OnMount = (editor, monaco) => {
		editorRef.current = editor;

		// Add Ctrl+Enter keybinding using addCommand instead of addAction
		editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, async () => {
			const selectedText = getSelectedText(editor);
			if (selectedText) {
				await executeQuery(selectedText);
			}
		});

		// Listen to selection changes
		editor.onDidChangeCursorSelection(handleSelectionChange);
	};

	//biome-ignore lint/correctness/useExhaustiveDependencies: handleSelectionChange doesn't need to be in deps
	useEffect(() => {
		handleSelectionChange();
	}, []);

	return (
		<div className="flex h-full flex-col">
			<div className="flex items-center justify-between border-gray-200 border-b bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-700">
				<div className="flex items-center space-x-3">
					<Button
						onClick={handleExecute}
						disabled={isExecuting || isGenerating || !value.trim()}
						size="sm"
						className={`text-white ${
							(selectedText || value).trim().toLowerCase().startsWith("gen:")
								? "bg-purple-600 hover:bg-purple-700 dark:bg-purple-700 dark:hover:bg-purple-800"
								: "bg-green-600 hover:bg-green-700 dark:bg-green-700 dark:hover:bg-green-800"
						}`}
					>
						{isGenerating ? (
							<>
								<Loader2 className="mr-2 h-4 w-4 animate-spin" />
								Generating
							</>
						) : isExecuting ? (
							<>
								<Loader2 className="mr-2 h-4 w-4 animate-spin" />
								Running
							</>
						) : (selectedText || value).trim().toLowerCase().startsWith("gen:") ? (
							<>
								<Sparkles className="mr-2 h-4 w-4" />
								Generate
							</>
						) : (
							<>
								<Play className="mr-2 h-4 w-4" />
								Execute
							</>
						)}
					</Button>

					{isExecuting && onStop && (
						<Button onClick={onStop} size="sm" variant="destructive">
							<Square className="mr-2 h-4 w-4" />
							Stop
						</Button>
					)}

					<div className="text-gray-600 text-sm dark:text-gray-300">
						{database && `Database: ${database}`}
					</div>
				</div>

				<div className="text-gray-500 text-xs dark:text-gray-400">
					{selectedText
						? `Selected: ${selectedText.length} chars`
						: `Total: ${value.length} chars`}
					<span className="ml-2">
						{(selectedText || value).trim().toLowerCase().startsWith("gen:")
							? "Ctrl+Enter to generate with AI"
							: "Ctrl+Enter to execute"}
					</span>
				</div>
			</div>

			<div className="relative flex-1">
				<Editor
					height="100%"
					value={value}
					onChange={(newValue) => onChange(newValue || "")}
					language="sql"
					theme={actualTheme === 'dark' ? 'vs-dark' : 'vs'}
					onMount={handleEditorDidMount}
					options={{
						minimap: { enabled: false },
						lineNumbers: 'on',
						roundedSelection: false,
						scrollBeyondLastLine: false,
						automaticLayout: true,
						fontFamily: 'ui-monospace, SFMono-Regular, "SF Mono", Consolas, "Liberation Mono", Menlo, monospace',
						fontSize: 14,
						lineHeight: 22,
						tabSize: 2,
						insertSpaces: true,
						wordWrap: 'on',
						bracketPairColorization: { enabled: true },
						suggestOnTriggerCharacters: true,
						acceptSuggestionOnCommitCharacter: true,
						acceptSuggestionOnEnter: 'on',
						quickSuggestions: true,
						parameterHints: { enabled: true },
						colorDecorators: true,
						codeLens: false,
						folding: true,
						foldingHighlight: true,
						showFoldingControls: 'mouseover',
						matchBrackets: 'always',
						selectionHighlight: true,
						occurrencesHighlight: 'singleFile'
					}}
				/>
			</div>

			<div className="border-gray-200 border-t bg-gray-50 px-4 py-2 text-gray-600 text-xs dark:border-gray-600 dark:bg-gray-700 dark:text-gray-400">
				<div className="flex justify-between">
					<span>
						{(selectedText || value).trim().toLowerCase().startsWith("gen:") ? (
							<span className="flex items-center">
								<Sparkles className="mr-1 h-3 w-3" />
								{selectedText
									? `Will generate from selected prompt (${selectedText.split("\n").length} lines)`
									: `Will generate from prompt (${value.split("\n").length} lines)`}
							</span>
						) : (
							selectedText
								? `Will execute selected text (${selectedText.split("\n").length} lines)`
								: `Will execute entire query (${value.split("\n").length} lines)`
						)}
					</span>
					<span>Lines: {value.split("\n").length}</span>
				</div>
			</div>
		</div>
	);
};