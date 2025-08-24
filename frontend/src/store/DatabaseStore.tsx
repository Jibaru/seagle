import type React from "react";
import { createContext, useCallback, useContext, useReducer } from "react";
import { GetTables } from "../../wailsjs/go/handlers/GetTablesHandler";

interface DatabaseState {
	databases: string[];
	selectedDatabase?: string;
	selectedTable?: string;
	databaseTables: Record<string, string[]>;
	loadingTables: Set<string>;
	expandedDatabases: Set<string>;
}

type DatabaseAction =
	| { type: "SET_DATABASES"; payload: string[] }
	| { type: "SELECT_DATABASE"; payload: string }
	| { type: "SELECT_TABLE"; payload: { database: string; table: string } }
	| { type: "TOGGLE_DATABASE"; payload: string }
	| {
			type: "SET_LOADING_TABLES";
			payload: { database: string; loading: boolean };
	  }
	| {
			type: "SET_DATABASE_TABLES";
			payload: { database: string; tables: string[] };
	  }
	| { type: "RESET_STATE" };

const initialState: DatabaseState = {
	databases: [],
	selectedDatabase: undefined,
	selectedTable: undefined,
	databaseTables: {},
	loadingTables: new Set(),
	expandedDatabases: new Set(),
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

		case "SET_DATABASE_TABLES":
			return {
				...state,
				databaseTables: {
					...state.databaseTables,
					[action.payload.database]: action.payload.tables,
				},
			};

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
					const result = await GetTables({ database });
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
		[state.expandedDatabases, state.databaseTables],
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
