import type React from "react";
import { useState } from "react";
import { ExecuteQuery } from "../../wailsjs/go/handlers/ExecuteQueryHandler";
import { QueryResults } from "./QueryResults";
import { SqlEditor } from "./SqlEditor";

interface QueryResult {
	columns: string[];
	//biome-ignore lint/suspicious/noExplicitAny: elements can be of any type
	rows: any[][];
	rowsAffected: number;
	duration: number;
}

interface QueryInterfaceProps {
	database: string;
}

export const QueryInterface: React.FC<QueryInterfaceProps> = ({ database }) => {
	const [query, setQuery] = useState("");
	const [result, setResult] = useState<QueryResult>();
	const [error, setError] = useState<string>();
	const [isExecuting, setIsExecuting] = useState(false);
	const [lastExecutedQuery, setLastExecutedQuery] = useState<string>();

	const handleExecuteQuery = async (queryToExecute: string) => {
		if (!queryToExecute.trim()) return;

		setIsExecuting(true);
		setError(undefined);
		setResult(undefined);
		setLastExecutedQuery(queryToExecute);

		try {
			const response = await ExecuteQuery({
				database,
				query: queryToExecute,
			});

			if (response?.success && response?.result) {
				setResult(response.result);
			} else {
				setError(response?.message || "Query execution failed");
			}
		} catch (err) {
			setError(
				err instanceof Error ? err.message : "An unexpected error occurred",
			);
		} finally {
			setIsExecuting(false);
		}
	};

	const handleStopQuery = () => {
		// For now, we'll just set executing to false
		// In a real implementation, you'd cancel the request
		setIsExecuting(false);
	};

	return (
		<div className="flex h-full flex-col overflow-hidden">
			{/* SQL Editor - takes up 40% of height */}
			<div className="h-2/5 flex-shrink-0 border-gray-200 border-b">
				<SqlEditor
					value={query}
					onChange={setQuery}
					onExecute={handleExecuteQuery}
					onStop={handleStopQuery}
					isExecuting={isExecuting}
					database={database}
				/>
			</div>

			{/* Query Results - takes up 60% of height */}
			<div className="flex-1 overflow-hidden">
				<QueryResults
					result={result}
					error={error}
					isLoading={isExecuting}
					query={lastExecutedQuery}
				/>
			</div>
		</div>
	);
};
