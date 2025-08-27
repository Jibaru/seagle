import type React from "react";
import { createContext, useCallback, useContext, useReducer } from "react";

interface ActiveConnectionState {
	connectionId: string | null;
	isConnected: boolean;
}

type ActiveConnectionAction =
	| { type: "SET_CONNECTION"; payload: { connectionId: string } }
	| { type: "CLEAR_CONNECTION" };

const initialState: ActiveConnectionState = {
	connectionId: null,
	isConnected: false,
};

function activeConnectionReducer(
	state: ActiveConnectionState,
	action: ActiveConnectionAction,
): ActiveConnectionState {
	switch (action.type) {
		case "SET_CONNECTION":
			return {
				connectionId: action.payload.connectionId,
				isConnected: true,
			};

		case "CLEAR_CONNECTION":
			return {
				connectionId: null,
				isConnected: false,
			};

		default:
			return state;
	}
}

interface ActiveConnectionContextValue {
	state: ActiveConnectionState;
	setConnection: (connectionId: string) => void;
	clearConnection: () => void;
}

const ActiveConnectionContext = createContext<ActiveConnectionContextValue | undefined>(
	undefined,
);

interface ActiveConnectionProviderProps {
	children: React.ReactNode;
}

export const ActiveConnectionProvider: React.FC<ActiveConnectionProviderProps> = ({
	children,
}) => {
	const [state, dispatch] = useReducer(activeConnectionReducer, initialState);

	const setConnection = useCallback((connectionId: string) => {
		dispatch({ type: "SET_CONNECTION", payload: { connectionId } });
	}, []);

	const clearConnection = useCallback(() => {
		dispatch({ type: "CLEAR_CONNECTION" });
	}, []);

	const contextValue: ActiveConnectionContextValue = {
		state,
		setConnection,
		clearConnection,
	};

	return (
		<ActiveConnectionContext.Provider value={contextValue}>
			{children}
		</ActiveConnectionContext.Provider>
	);
};

export const useActiveConnectionStore = (): ActiveConnectionContextValue => {
	const context = useContext(ActiveConnectionContext);
	if (!context) {
		throw new Error("useActiveConnectionStore must be used within an ActiveConnectionProvider");
	}
	return context;
};