import { useState } from "react";
import { DatabaseConnectionForm } from "./components/DatabaseConnectionForm";
import { MainLayout } from "./components/MainLayout";
import { WelcomeScreen } from "./components/WelcomeScreen";
import { Button } from "./components/ui/button";
import { ThemeProvider } from "./contexts/ThemeContext";
import { DatabaseProvider, useDatabaseStore } from "./store/DatabaseStore";
import "./App.css";

function AppContent() {
	const [currentScreen, setCurrentScreen] = useState<
		"welcome" | "connection" | "connected"
	>("welcome");
	const [isConnected, setIsConnected] = useState(false);
	const { setDatabases, resetState } = useDatabaseStore();

	const handleNewConnection = () => {
		setCurrentScreen("connection");
	};

	const handleConnectionChange = (
		connected: boolean,
		databaseList?: string[],
	) => {
		setIsConnected(connected);
		if (connected) {
			setCurrentScreen("connected");
			if (databaseList) {
				setDatabases(databaseList);
			}
		} else {
			resetState();
		}
	};

	const handleBackToWelcome = () => {
		setCurrentScreen("welcome");
		setIsConnected(false);
		resetState();
	};

	if (currentScreen === "welcome") {
		return <WelcomeScreen onNewConnection={handleNewConnection} />;
	}

	if (currentScreen === "connection") {
		return (
			<div className="flex min-h-screen items-center justify-center bg-gray-900 p-4 dark:bg-gray-950">
				<div className="w-full max-w-4xl space-y-4">
					<div className="mb-4 flex justify-start">
						<Button
							onClick={handleBackToWelcome}
							variant="outline"
							className="border-gray-600 text-gray-300 hover:bg-gray-800 dark:border-gray-500 dark:text-gray-300 dark:hover:bg-gray-800"
						>
							← Back
						</Button>
					</div>
					<DatabaseConnectionForm onConnectionChange={handleConnectionChange} />
				</div>
			</div>
		);
	}

	if (currentScreen === "connected") {
		return <MainLayout onNewConnection={handleBackToWelcome} />;
	}

	return null;
}

function App() {
	return (
		<ThemeProvider>
			<DatabaseProvider>
				<AppContent />
			</DatabaseProvider>
		</ThemeProvider>
	);
}

export default App;
