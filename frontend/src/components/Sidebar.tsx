import {
	ChevronDown,
	ChevronRight,
	Columns,
	Database,
	Loader2,
	Table,
} from "lucide-react";
import type React from "react";
import { useDatabaseStore } from "../store/DatabaseStore";

export const Sidebar: React.FC = () => {
	const { state, selectDatabase, selectTable, toggleDatabase, toggleTable } =
		useDatabaseStore();

	return (
		<div className="flex h-full w-64 flex-col border-gray-200 border-r bg-white shadow-sm">
			<div className="flex-shrink-0 border-gray-200 border-b p-4">
				<h2 className="font-semibold text-gray-800 text-lg">Databases</h2>
			</div>
			<div className="flex-1 overflow-y-auto bg-white p-2">
				{state.databases.length === 0 ? (
					<div className="p-4 text-center text-gray-500">
						No databases available
					</div>
				) : (
					<div className="space-y-1">
						{state.databases.map((database) => (
							<div key={database}>
								<div className="flex items-center">
									<button
										type="button"
										onClick={() => toggleDatabase(database)}
										className="flex h-8 w-8 items-center justify-center rounded hover:bg-gray-100"
									>
										{state.expandedDatabases.has(database) ? (
											<ChevronDown className="h-4 w-4 text-gray-500" />
										) : (
											<ChevronRight className="h-4 w-4 text-gray-500" />
										)}
									</button>
									<button
										type="button"
										onClick={() => selectDatabase(database)}
										className={`flex flex-1 items-center rounded-md px-2 py-2 text-left text-sm transition-colors ${
											state.selectedDatabase === database
												? "bg-blue-100 text-blue-700"
												: "text-gray-700 hover:bg-gray-100"
										}`}
									>
										<Database className="mr-2 h-4 w-4 flex-shrink-0" />
										<span className="truncate">{database}</span>
									</button>
								</div>

								{state.expandedDatabases.has(database) && (
									<div className="ml-6 border-gray-200 border-l pl-2">
										{state.loadingTables.has(database) ? (
											<div className="flex items-center px-3 py-2 text-gray-500 text-sm">
												<Loader2 className="mr-2 h-3 w-3 animate-spin" />
												Loading tables...
											</div>
										) : state.databaseTables[database]?.length === 0 ? (
											<div className="px-3 py-2 text-gray-500 text-sm">
												No tables found
											</div>
										) : (
											state.databaseTables[database]?.map((table) => {
												const tableKey = `${database}.${table}`;
												return (
													<div key={table} className="space-y-1">
														<div className="flex items-center">
															<button
																type="button"
																onClick={() => toggleTable(database, table)}
																className="flex h-6 w-6 items-center justify-center rounded hover:bg-gray-100"
															>
																{state.expandedTables.has(tableKey) ? (
																	<ChevronDown className="h-3 w-3 text-gray-500" />
																) : (
																	<ChevronRight className="h-3 w-3 text-gray-500" />
																)}
															</button>
															<button
																type="button"
																onClick={() => selectTable(database, table)}
																className={`flex flex-1 items-center rounded-md px-2 py-1.5 text-left text-sm transition-colors ${
																	state.selectedDatabase === database &&
																	state.selectedTable === table
																		? "bg-green-100 text-green-700"
																		: "text-gray-600 hover:bg-gray-100"
																}`}
															>
																<Table className="mr-2 h-3 w-3 flex-shrink-0" />
																<span className="truncate">{table}</span>
															</button>
														</div>

														{state.expandedTables.has(tableKey) && (
															<div className="ml-4 border-gray-300 border-l pl-2">
																{state.loadingColumns.has(tableKey) ? (
																	<div className="flex items-center px-3 py-1 text-gray-500 text-xs">
																		<Loader2 className="mr-2 h-2 w-2 animate-spin" />
																		Loading columns...
																	</div>
																) : state.tableColumns[tableKey]?.length ===
																	0 ? (
																	<div className="px-3 py-1 text-gray-500 text-xs">
																		No columns found
																	</div>
																) : (
																	state.tableColumns[tableKey]?.map(
																		(column) => (
																			<div
																				key={column.name}
																				className="flex items-center px-3 py-1 text-gray-600 text-xs"
																			>
																				<Columns className="mr-2 h-2 w-2 flex-shrink-0" />
																				<span className="mr-1 truncate">
																					{column.name}
																				</span>
																				<span className="text-gray-400 text-xs">
																					({column.dataType})
																				</span>
																			</div>
																		),
																	)
																)}
															</div>
														)}
													</div>
												);
											})
										)}
									</div>
								)}
							</div>
						))}
					</div>
				)}
			</div>
		</div>
	);
};
