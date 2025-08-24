import { useState } from "react";
import { DatabaseConnectionForm } from "./components/DatabaseConnectionForm";
import { MainLayout } from "./components/MainLayout";
import { WelcomeScreen } from "./components/WelcomeScreen";
import { Button } from "./components/ui/button";
import "./App.css";

function App() {
	const [currentScreen, setCurrentScreen] = useState<
		"welcome" | "connection" | "connected"
	>("welcome");
	const [isConnected, setIsConnected] = useState(false);
	const [databases, setDatabases] = useState<string[]>([]);

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
			setDatabases([]);
		}
	};

	const handleBackToWelcome = () => {
		setCurrentScreen("welcome");
		setIsConnected(false);
		setDatabases([]);
	};

	if (currentScreen === "welcome") {
		return <WelcomeScreen onNewConnection={handleNewConnection} />;
	}

	if (currentScreen === "connection") {
		return (
			<div className="flex min-h-screen items-center justify-center bg-gray-900 p-4">
				<div className="w-full max-w-4xl space-y-4">
					<div className="mb-4 flex justify-start">
						<Button
							onClick={handleBackToWelcome}
							variant="outline"
							className="border-gray-600 text-gray-300 hover:bg-gray-800"
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
		return (
			<MainLayout databases={databases} onNewConnection={handleBackToWelcome} />
		);
	}

	return null;
}

export default App;
