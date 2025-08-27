import type React from "react";
import { createContext, useCallback, useContext, useReducer } from "react";
import { GetTableColumns } from "../../wailsjs/go/handlers/GetTableColumnsHandler";
import { GetTables } from "../../wailsjs/go/handlers/GetTablesHandler";
import { useActiveConnectionStore } from "./ActiveConnectionStore";

interface TableColumn {
	name: string;
	dataType: string;
	isNullable: boolean;
	defaultValue?: string;
}

interface DatabaseState {
	databases: string[];
	selectedDatabase?: string;
	selectedTable?: string;
	databaseTables: Record<string, string[]>;
	tableColumns: Record<string, TableColumn[]>; // key: "database.table"
	loadingTables: Set<string>;
	loadingColumns: Set<string>; // key: "database.table"
	expandedDatabases: Set<string>;
	expandedTables: Set<string>; // key: "database.table"
}

type DatabaseAction =
	| { type: "SET_DATABASES"; payload: string[] }
	| { type: "SELECT_DATABASE"; payload: string }
	| { type: "SELECT_TABLE"; payload: { database: string; table: string } }
	| { type: "TOGGLE_DATABASE"; payload: string }
	| { type: "TOGGLE_TABLE"; payload: { database: string; table: string } }
	| {
			type: "SET_LOADING_TABLES";
			payload: { database: string; loading: boolean };
	  }
	| {
			type: "SET_LOADING_COLUMNS";
			payload: { database: string; table: string; loading: boolean };
	  }
	| {
			type: "SET_DATABASE_TABLES";
			payload: { database: string; tables: string[] };
	  }
	| {
			type: "SET_TABLE_COLUMNS";
			payload: { database: string; table: string; columns: TableColumn[] };
	  }
	| { type: "RESET_STATE" };

const initialState: DatabaseState = {
	databases: [],
	selectedDatabase: undefined,
	selectedTable: undefined,
	databaseTables: {},
	tableColumns: {},
	loadingTables: new Set(),
	loadingColumns: new Set(),
	expandedDatabases: new Set(),
	expandedTables: new Set(),
};

function databaseReducer(
	state: DatabaseState,
	action: DatabaseAction,
): DatabaseState {
	switch (action.type) {
		case "SET_DATABASES":
			return {
				...state,
				databases: action.payload,
			};

		case "SELECT_DATABASE":
			return {
				...state,
				selectedDatabase: action.payload,
				selectedTable: undefined,
			};

		case "SELECT_TABLE":
			return {
				...state,
				selectedDatabase: action.payload.database,
				selectedTable: action.payload.table,
			};

		case "TOGGLE_DATABASE": {
			const newExpanded = new Set(state.expandedDatabases);
			if (newExpanded.has(action.payload)) {
				newExpanded.delete(action.payload);
			} else {
				newExpanded.add(action.payload);
			}
			return {
				...state,
				expandedDatabases: newExpanded,
			};
		}

		case "SET_LOADING_TABLES": {
			const newLoading = new Set(state.loadingTables);
			if (action.payload.loading) {
				newLoading.add(action.payload.database);
			} else {
				newLoading.delete(action.payload.database);
			}
			return {
				...state,
				loadingTables: newLoading,
			};
		}

		case "TOGGLE_TABLE": {
			const tableKey = `${action.payload.database}.${action.payload.table}`;
			const newExpandedTables = new Set(state.expandedTables);
			if (newExpandedTables.has(tableKey)) {
				newExpandedTables.delete(tableKey);
			} else {
				newExpandedTables.add(tableKey);
			}
			return {
				...state,
				expandedTables: newExpandedTables,
			};
		}

		case "SET_LOADING_COLUMNS": {
			const tableKey = `${action.payload.database}.${action.payload.table}`;
			const newLoadingColumns = new Set(state.loadingColumns);
			if (action.payload.loading) {
				newLoadingColumns.add(tableKey);
			} else {
				newLoadingColumns.delete(tableKey);
			}
			return {
				...state,
				loadingColumns: newLoadingColumns,
			};
		}

		case "SET_DATABASE_TABLES":
			return {
				...state,
				databaseTables: {
					...state.databaseTables,
					[action.payload.database]: action.payload.tables,
				},
			};

		case "SET_TABLE_COLUMNS": {
			const tableKey = `${action.payload.database}.${action.payload.table}`;
			return {
				...state,
				tableColumns: {
					...state.tableColumns,
					[tableKey]: action.payload.columns,
				},
			};
		}

		case "RESET_STATE":
			return initialState;

		default:
			return state;
	}
}

interface DatabaseContextValue {
	state: DatabaseState;
	setDatabases: (databases: string[]) => void;
	selectDatabase: (database: string) => void;
	selectTable: (database: string, table: string) => void;
	toggleDatabase: (database: string) => Promise<void>;
	toggleTable: (database: string, table: string) => Promise<void>;
	resetState: () => void;
}

const DatabaseContext = createContext<DatabaseContextValue | undefined>(
	undefined,
);

interface DatabaseProviderProps {
	children: React.ReactNode;
}

export const DatabaseProvider: React.FC<DatabaseProviderProps> = ({
	children,
}) => {
	const [state, dispatch] = useReducer(databaseReducer, initialState);
	const { state: activeConnection } = useActiveConnectionStore();

	const setDatabases = useCallback((databases: string[]) => {
		dispatch({ type: "SET_DATABASES", payload: databases });
	}, []);

	const selectDatabase = useCallback((database: string) => {
		dispatch({ type: "SELECT_DATABASE", payload: database });
	}, []);

	const selectTable = useCallback((database: string, table: string) => {
		dispatch({ type: "SELECT_TABLE", payload: { database, table } });
	}, []);

	const toggleDatabase = useCallback(
		async (database: string) => {
			if (!activeConnection.connectionId) {
				console.error("No active connection ID available for fetching tables");
				return;
			}

			dispatch({ type: "TOGGLE_DATABASE", payload: database });

			// If expanding and tables not loaded, fetch them
			if (
				!state.expandedDatabases.has(database) &&
				!state.databaseTables[database]
			) {
				dispatch({
					type: "SET_LOADING_TABLES",
					payload: { database, loading: true },
				});
				try {
					const result = await GetTables({ 
						id: activeConnection.connectionId,
						database 
					});
					const tables = result?.success && result?.tables ? result.tables : [];
					dispatch({
						type: "SET_DATABASE_TABLES",
						payload: { database, tables },
					});
				} catch (error) {
					console.error(`Failed to fetch tables for ${database}:`, error);
					dispatch({
						type: "SET_DATABASE_TABLES",
						payload: { database, tables: [] },
					});
				} finally {
					dispatch({
						type: "SET_LOADING_TABLES",
						payload: { database, loading: false },
					});
				}
			}
		},
		[state.expandedDatabases, state.databaseTables, activeConnection.connectionId],
	);

	const toggleTable = useCallback(
		async (database: string, table: string) => {
			if (!activeConnection.connectionId) {
				console.error("No active connection ID available for fetching columns");
				return;
			}

			const tableKey = `${database}.${table}`;
			dispatch({ type: "TOGGLE_TABLE", payload: { database, table } });

			// If expanding and columns not loaded, fetch them
			if (
				!state.expandedTables.has(tableKey) &&
				!state.tableColumns[tableKey]
			) {
				dispatch({
					type: "SET_LOADING_COLUMNS",
					payload: { database, table, loading: true },
				});
				try {
					const result = await GetTableColumns({ 
						id: activeConnection.connectionId,
						database, 
						table 
					});
					const columns =
						result?.success && result?.columns ? result.columns : [];
					dispatch({
						type: "SET_TABLE_COLUMNS",
						payload: { database, table, columns },
					});
				} catch (error) {
					console.error(
						`Failed to fetch columns for ${database}.${table}:`,
						error,
					);
					dispatch({
						type: "SET_TABLE_COLUMNS",
						payload: { database, table, columns: [] },
					});
				} finally {
					dispatch({
						type: "SET_LOADING_COLUMNS",
						payload: { database, table, loading: false },
					});
				}
			}
		},
		[state.expandedTables, state.tableColumns, activeConnection.connectionId],
	);

	const resetState = useCallback(() => {
		dispatch({ type: "RESET_STATE" });
	}, []);

	const contextValue: DatabaseContextValue = {
		state,
		setDatabases,
		selectDatabase,
		selectTable,
		toggleDatabase,
		toggleTable,
		resetState,
	};

	return (
		<DatabaseContext.Provider value={contextValue}>
			{children}
		</DatabaseContext.Provider>
	);
};

export const useDatabaseStore = (): DatabaseContextValue => {
	const context = useContext(DatabaseContext);
	if (!context) {
		throw new Error("useDatabaseStore must be used within a DatabaseProvider");
	}
	return context;
};
