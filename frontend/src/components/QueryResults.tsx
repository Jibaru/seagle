import { CheckCircle, Clock, Database, XCircle } from "lucide-react";
import type React from "react";
import { useCallback, useRef, useState } from "react";

interface QueryResult {
	columns: string[];
	//biome-ignore lint/suspicious/noExplicitAny: database values can be of any type
	rows: any[][];
	rowsAffected: number;
	duration: number;
}

interface QueryResultsProps {
	result?: QueryResult;
	error?: string;
	isLoading?: boolean;
	query?: string;
}

export const QueryResults: React.FC<QueryResultsProps> = ({
	result,
	error,
	isLoading,
	query,
}) => {
	// Column resizing state
	const [columnWidths, setColumnWidths] = useState<Record<string, number>>({});
	const [isResizing, setIsResizing] = useState(false);
	const [resizingColumn, setResizingColumn] = useState<string | null>(null);
	const resizeStartX = useRef<number>(0);
	const resizeStartWidth = useRef<number>(0);

	// Cell editing state
	const [editingCell, setEditingCell] = useState<{
		row: number;
		col: number;
	} | null>(null);
	const [editValue, setEditValue] = useState<string>("");
	const editInputRef = useRef<HTMLTextAreaElement>(null);

	const getColumnWidth = useCallback(
		(columnName: string, defaultWidth = 100) => {
			return columnWidths[columnName] || defaultWidth;
		},
		[columnWidths],
	);

	const handleResizeStart = useCallback(
		(e: React.MouseEvent, columnName: string) => {
			e.preventDefault();
			e.stopPropagation();
			
			console.log('Resize start for column:', columnName); // Debug log
			
			setIsResizing(true);
			setResizingColumn(columnName);
			resizeStartX.current = e.clientX;
			resizeStartWidth.current = getColumnWidth(columnName);

			const handleMouseMove = (moveEvent: MouseEvent) => {
				moveEvent.preventDefault();
				const diff = moveEvent.clientX - resizeStartX.current;
				const newWidth = Math.max(60, resizeStartWidth.current + diff);
				console.log('Resizing to:', newWidth); // Debug log
				setColumnWidths((prev) => ({ ...prev, [columnName]: newWidth }));
			};

			const handleMouseUp = () => {
				console.log('Resize end'); // Debug log
				setIsResizing(false);
				setResizingColumn(null);
				document.removeEventListener("mousemove", handleMouseMove);
				document.removeEventListener("mouseup", handleMouseUp);
			};

			document.addEventListener("mousemove", handleMouseMove);
			document.addEventListener("mouseup", handleMouseUp);
		},
		[getColumnWidth],
	);

	const handleCellDoubleClick = (
		rowIndex: number,
		cellIndex: number,
		//biome-ignore lint/suspicious/noExplicitAny: cell values can be of any database type
		cellValue: any,
	) => {
		setEditingCell({ row: rowIndex, col: cellIndex });
		setEditValue(String(cellValue ?? ""));

		// Focus the input after state update
		setTimeout(() => {
			if (editInputRef.current) {
				editInputRef.current.focus();
				editInputRef.current.select();
			}
		}, 0);
	};

	const handleEditComplete = () => {
		setEditingCell(null);
		setEditValue("");
	};

	const handleEditKeyDown = (e: React.KeyboardEvent) => {
		if (e.key === "Enter" && !e.shiftKey) {
			e.preventDefault();
			handleEditComplete();
		} else if (e.key === "Escape") {
			handleEditComplete();
		}
	};
	if (isLoading) {
		return (
			<div className="flex h-64 flex-col items-center justify-center text-gray-500">
				<div className="mb-4 h-8 w-8 animate-spin rounded-full border-blue-600 border-b-2" />
				<div>Executing query...</div>
				{query && <div className="mt-2 max-w-md truncate text-sm">{query}</div>}
			</div>
		);
	}

	if (error) {
		return (
			<div className="rounded-md border border-red-200 bg-red-50 p-4">
				<div className="mb-2 flex items-center text-red-700">
					<XCircle className="mr-2 h-5 w-5" />
					<span className="font-medium">Query Error</span>
				</div>
				<pre className="whitespace-pre-wrap font-mono text-red-800 text-sm">
					{error}
				</pre>
			</div>
		);
	}

	if (!result) {
		return (
			<div className="flex h-64 flex-col items-center justify-center text-gray-500">
				<Database className="mb-4 h-12 w-12" />
				<div className="mb-2 text-lg">No query executed</div>
				<div className="text-sm">Execute a SQL query to see results here</div>
			</div>
		);
	}

	const hasRows = result.rows && result.rows.length > 0;
	const hasColumns = result.columns && result.columns.length > 0;

	return (
		<div className="flex h-full flex-col overflow-hidden">
			{/* Result Status Bar */}
			<div className="flex flex-shrink-0 items-center justify-between rounded-t-md border border-green-200 bg-green-50 p-3">
				<div className="flex items-center text-green-700">
					<CheckCircle className="mr-2 h-5 w-5" />
					<span className="font-medium">Query executed successfully</span>
				</div>
				<div className="flex items-center space-x-4 text-green-600 text-sm">
					<div className="flex items-center">
						<Clock className="mr-1 h-4 w-4" />
						{result.duration}ms
					</div>
					<div>
						{hasRows
							? `${result.rows.length} rows`
							: `${result.rowsAffected} rows affected`}
					</div>
				</div>
			</div>

			{/* Results Content */}
			<div className="flex-1 overflow-hidden bg-white" style={{ minHeight: 0 }}>
				{hasColumns && hasRows ? (
					/* Table Results */
					<div
						className={`h-full border border-gray-300 ${isResizing ? "select-none" : ""}`}
						style={{
							overflow: "auto",
							maxWidth: "100%",
							maxHeight: "100%",
							cursor: isResizing ? "col-resize" : "default",
						}}
					>
						<table className="border-collapse" style={{ tableLayout: "auto", width: "100%" }}>
							<thead className="sticky top-0 z-10 bg-gray-50">
								<tr>
									<th
										className="sticky left-0 z-20 border border-gray-300 bg-gray-100 px-2 py-1 text-left font-medium text-gray-900 text-xs"
										style={{ width: "50px" }}
									>
										#
									</th>
									{result.columns.map((column) => (
										<th
											key={`col-${column}`}
											className="relative border border-gray-300 bg-gray-100 px-2 py-1 text-left font-medium text-gray-900 text-sm"
											style={{
												width: getColumnWidth(column),
												minWidth: "60px",
												maxWidth: getColumnWidth(column),
											}}
										>
											<div className="truncate" title={column}>
												{column}
											</div>
											{/* Resize handle */}
											<div
												className="absolute top-0 right-0 h-full w-1 cursor-col-resize bg-transparent hover:bg-blue-400 z-30"
												onMouseDown={(e) => handleResizeStart(e, column)}
												title="Drag to resize column"
												style={{ 
													right: '-1px',
													userSelect: 'none',
													borderRight: isResizing && resizingColumn === column ? '2px solid #3b82f6' : 'none'
												}}
											/>
										</th>
									))}
								</tr>
							</thead>
							<tbody>
								{result.rows.map((row, rowIndex) => (
									<tr
										key={`row-${rowIndex}-${row[0]}`}
										className={rowIndex % 2 === 0 ? "bg-white" : "bg-gray-50"}
									>
										<td
											className="sticky left-0 z-10 border border-gray-300 bg-gray-100 px-2 py-1 text-gray-500 text-xs"
											style={{ width: "50px" }}
										>
											{rowIndex + 1}
										</td>
										{row.map((cell, cellIndex) => {
											const isEditing =
												editingCell?.row === rowIndex &&
												editingCell?.col === cellIndex;
											const columnName = result.columns[cellIndex];

											return (
												<td
													key={`cell-${rowIndex}-${cellIndex}-${String(cell).slice(0, 10)}`}
													className="relative border border-gray-300 px-2 py-1 font-mono text-gray-900 text-sm"
													style={{
														width: getColumnWidth(columnName),
														minWidth: "60px",
														maxWidth: getColumnWidth(columnName),
													}}
													onDoubleClick={() =>
														handleCellDoubleClick(rowIndex, cellIndex, cell)
													}
												>
													{isEditing ? (
														<textarea
															ref={editInputRef}
															value={editValue}
															onChange={(e) => setEditValue(e.target.value)}
															onBlur={handleEditComplete}
															onKeyDown={handleEditKeyDown}
															className="h-full w-full resize-none border-2 border-blue-400 border-none bg-yellow-50 px-2 py-1 font-mono text-sm outline-none"
															style={{
																minHeight: "32px",
																maxHeight: "200px",
																overflow: "auto",
															}}
															placeholder="Edit value..."
														/>
													) : (
														<div
															className="flex h-full cursor-pointer items-center px-2 py-1 hover:bg-blue-50"
															title={`Double-click to edit. Value: ${String(cell ?? "")}`}
														>
															<div className="w-full truncate">
																{cell === null ? (
																	<span className="text-gray-400 italic">
																		NULL
																	</span>
																) : cell === "" ? (
																	<span className="text-gray-400 italic">
																		(empty)
																	</span>
																) : (
																	<span className="text-gray-900">
																		{String(cell)}
																	</span>
																)}
															</div>
														</div>
													)}
												</td>
											);
										})}
									</tr>
								))}
							</tbody>
						</table>
					</div>
				) : (
					/* Non-SELECT Results */
					<div className="flex h-32 items-center justify-center text-gray-600">
						<div className="text-center">
							<div className="font-medium text-lg">
								{result.rowsAffected} rows affected
							</div>
							<div className="mt-1 text-gray-500 text-sm">
								Query completed in {result.duration}ms
							</div>
						</div>
					</div>
				)}
			</div>

			{/* Footer Info */}
			{hasRows && (
				<div className="flex-shrink-0 border-gray-200 border-t bg-gray-50 px-3 py-2 text-gray-600 text-xs">
					<div className="flex justify-between">
						<span>
							Showing {result.rows.length} of {result.rows.length} rows
						</span>
						<span>{result.columns.length} columns</span>
					</div>
				</div>
			)}
		</div>
	);
};
