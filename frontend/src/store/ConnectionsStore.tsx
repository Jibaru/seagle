import type React from "react";
import { createContext, useCallback, useContext, useReducer } from "react";
import { ListConnections } from "../../wailsjs/go/handlers/ListConnectionsHandler";
import { DeleteConnection } from "../../wailsjs/go/handlers/DeleteConnectionHandler";
import { types, handlers } from "../../wailsjs/go/models";

interface ConnectionsState {
	connections: types.ConnectionSummary[];
	loading: boolean;
	error: string | null;
	lastLoaded: number | null;
	connectingId: string | null;
	deletingId: string | null;
}

type ConnectionsAction =
	| { type: "SET_LOADING"; payload: boolean }
	| { type: "SET_ERROR"; payload: string | null }
	| { type: "SET_CONNECTIONS"; payload: types.ConnectionSummary[] }
	| { type: "SET_CONNECTING_ID"; payload: string | null }
	| { type: "SET_DELETING_ID"; payload: string | null }
	| { type: "REFRESH_CONNECTIONS" }
	| { type: "RESET_STATE" };

const initialState: ConnectionsState = {
	connections: [],
	loading: true,
	error: null,
	lastLoaded: null,
	connectingId: null,
	deletingId: null,
};

function connectionsReducer(
	state: ConnectionsState,
	action: ConnectionsAction,
): ConnectionsState {
	switch (action.type) {
		case "SET_LOADING":
			return {
				...state,
				loading: action.payload,
			};

		case "SET_ERROR":
			return {
				...state,
				error: action.payload,
				loading: false,
			};

		case "SET_CONNECTIONS":
			return {
				...state,
				connections: action.payload,
				loading: false,
				error: null,
				lastLoaded: Date.now(),
			};

		case "SET_CONNECTING_ID":
			return {
				...state,
				connectingId: action.payload,
			};

		case "SET_DELETING_ID":
			return {
				...state,
				deletingId: action.payload,
			};

		case "REFRESH_CONNECTIONS":
			return {
				...state,
				loading: true,
				error: null,
			};

		case "RESET_STATE":
			return initialState;

		default:
			return state;
	}
}

interface ConnectionsContextValue {
	state: ConnectionsState;
	loadConnections: (force?: boolean) => Promise<void>;
	setConnectingId: (id: string | null) => void;
	deleteConnection: (id: string) => Promise<void>;
	refreshConnections: () => Promise<void>;
	resetState: () => void;
}

const ConnectionsContext = createContext<ConnectionsContextValue | undefined>(
	undefined,
);

interface ConnectionsProviderProps {
	children: React.ReactNode;
}

export const ConnectionsProvider: React.FC<ConnectionsProviderProps> = ({
	children,
}) => {
	const [state, dispatch] = useReducer(connectionsReducer, initialState);

	const loadConnections = useCallback(
		async (force = false) => {
			// Skip loading if we have recent data (less than 30 seconds old) unless forced
			const thirtySecondsAgo = Date.now() - 30 * 1000;
			if (
				!force &&
				state.lastLoaded &&
				state.lastLoaded > thirtySecondsAgo &&
				!state.loading &&
				!state.error
			) {
				return;
			}

			try {
				dispatch({ type: "SET_LOADING", payload: true });
				const result = await ListConnections();
				
				if (result.success) {
					dispatch({ 
						type: "SET_CONNECTIONS", 
						payload: result.connections || [] 
					});
				} else {
					dispatch({ 
						type: "SET_ERROR", 
						payload: result.message || "Failed to load saved connections" 
					});
				}
			} catch (err) {
				dispatch({ 
					type: "SET_ERROR", 
					payload: "Failed to load saved connections" 
				});
				console.error("Error loading connections:", err);
			}
		},
		[state.lastLoaded, state.loading, state.error],
	);

	const setConnectingId = useCallback((id: string | null) => {
		dispatch({ type: "SET_CONNECTING_ID", payload: id });
	}, []);

	const deleteConnection = useCallback(async (id: string) => {
		try {
			dispatch({ type: "SET_DELETING_ID", payload: id });
			
			const input = new handlers.DeleteConnectionInput({ id });
			await DeleteConnection(input);
			
			// Refresh connections after successful deletion
			await loadConnections(true);
		} catch (err) {
			dispatch({ 
				type: "SET_ERROR", 
				payload: "Failed to delete connection" 
			});
			console.error("Error deleting connection:", err);
		} finally {
			dispatch({ type: "SET_DELETING_ID", payload: null });
		}
	}, [loadConnections]);

	const refreshConnections = useCallback(async () => {
		dispatch({ type: "REFRESH_CONNECTIONS" });
		await loadConnections(true);
	}, [loadConnections]);

	const resetState = useCallback(() => {
		dispatch({ type: "RESET_STATE" });
	}, []);

	const contextValue: ConnectionsContextValue = {
		state,
		loadConnections,
		setConnectingId,
		deleteConnection,
		refreshConnections,
		resetState,
	};

	return (
		<ConnectionsContext.Provider value={contextValue}>
			{children}
		</ConnectionsContext.Provider>
	);
};

export const useConnectionsStore = (): ConnectionsContextValue => {
	const context = useContext(ConnectionsContext);
	if (!context) {
		throw new Error("useConnectionsStore must be used within a ConnectionsProvider");
	}
	return context;
};