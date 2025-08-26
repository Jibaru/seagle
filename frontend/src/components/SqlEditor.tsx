import { Loader2, Play, Square, Sparkles } from "lucide-react";
import type React from "react";
import { useEffect, useRef, useState } from "react";
import { Button } from "./ui/button";
import { GenerateQuery } from "../../wailsjs/go/handlers/GenQueryHandler";

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
	const [selectedText, setSelectedText] = useState("");
	const [isGenerating, setIsGenerating] = useState(false);
	const textareaRef = useRef<HTMLTextAreaElement>(null);

	const handleSelectionChange = () => {
		if (textareaRef.current) {
			const start = textareaRef.current.selectionStart;
			const end = textareaRef.current.selectionEnd;
			const selected = value.substring(start, end);
			setSelectedText(selected.trim());
		}
	};

	const handleExecute = async () => {
		const queryToExecute = selectedText || value;
		if (queryToExecute.trim()) {
			// Check if this is a generation request
			if (queryToExecute.trim().toLowerCase().startsWith("gen:")) {
				await handleGenerateQuery(queryToExecute.trim());
			} else {
				onExecute(queryToExecute.trim());
			}
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

		setIsGenerating(true);
		try {
			const response = await GenerateQuery({
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
					const newValue = value + '\n\n' + response.result.generatedQuery;
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

	const handleKeyDown = (e: React.KeyboardEvent) => {
		// Ctrl+Enter or Cmd+Enter to execute query
		if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
			e.preventDefault();
			handleExecute();
		}

		// Handle tab insertion
		if (e.key === "Tab") {
			e.preventDefault();
			const target = e.currentTarget as HTMLTextAreaElement;
			const start = target.selectionStart;
			const end = target.selectionEnd;
			const newValue = `${value.substring(0, start)}  ${value.substring(end)}`;
			onChange(newValue);

			// Set cursor position after tab
			setTimeout(() => {
				if (textareaRef.current) {
					textareaRef.current.selectionStart =
						textareaRef.current.selectionEnd = start + 2;
				}
			}, 0);
		}
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
				<textarea
					ref={textareaRef}
					value={value}
					onChange={(e) => onChange(e.target.value)}
					onKeyDown={handleKeyDown}
					onSelect={handleSelectionChange}
					onMouseUp={handleSelectionChange}
					placeholder={
						"-- Enter your SQL query here\n-- Use Ctrl+Enter to execute\n-- Select text to execute only the selection\n-- Type 'gen: your question' for AI-generated queries\n\nSELECT * FROM your_table LIMIT 10;\n\n-- Examples for AI generation:\n-- gen: find all users who registered last month\n-- gen: get the top 10 products by sales"
					}
					className="h-full w-full resize-none border-none bg-white p-4 font-mono text-gray-900 text-sm placeholder-gray-400 outline-none dark:bg-gray-800 dark:text-gray-200 dark:placeholder-gray-500"
					spellCheck={false}
					style={{
						tabSize: 2,
						lineHeight: "1.5",
					}}
				/>

				{/* Simple syntax highlighting overlay would go here */}
				{/* For now, we'll use CSS-based highlighting through the textarea */}
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
