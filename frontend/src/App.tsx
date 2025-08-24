import { useState } from "react";
import { DatabaseConnectionForm } from "./components/DatabaseConnectionForm";
import { WelcomeScreen } from "./components/WelcomeScreen";
import { Button } from "./components/ui/button";
import "./App.css";

function App() {
	const [currentScreen, setCurrentScreen] = useState<
		"welcome" | "connection" | "connected"
	>("welcome");
	const [isConnected, setIsConnected] = useState(false);

	const handleNewConnection = () => {
		setCurrentScreen("connection");
	};

	const handleConnectionChange = (connected: boolean) => {
		setIsConnected(connected);
		if (connected) {
			setCurrentScreen("connected");
		}
	};

	const handleBackToWelcome = () => {
		setCurrentScreen("welcome");
		setIsConnected(false);
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
			<div className="flex min-h-screen items-center justify-center bg-gray-900">
				<div className="space-y-4 text-center">
					<div className="mb-4 rounded border border-green-400 bg-green-100 p-6 text-green-700">
						Successfully connected to the database!
					</div>
					<Button
						onClick={handleBackToWelcome}
						variant="outline"
						className="border-gray-600 text-gray-300 hover:bg-gray-800"
					>
						New Connection
					</Button>
				</div>
			</div>
		);
	}

	return null;
}

export default App;
