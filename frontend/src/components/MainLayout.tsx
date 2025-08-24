import type React from "react";
import { useState } from "react";
import { Sidebar } from "./Sidebar";
import { Button } from "./ui/button";

interface MainLayoutProps {
	databases: string[];
	onNewConnection: () => void;
}

export const MainLayout: React.FC<MainLayoutProps> = ({
	databases,
	onNewConnection,
}) => {
	const [selectedDatabase, setSelectedDatabase] = useState<string>();

	const handleDatabaseSelect = (database: string) => {
		setSelectedDatabase(database);
	};

	return (
		<div className="flex h-screen bg-gray-50">
			<Sidebar
				databases={databases}
				selectedDatabase={selectedDatabase}
				onDatabaseSelect={handleDatabaseSelect}
			/>

			<div className="flex flex-1 flex-col">
				<header className="border-gray-200 border-b bg-white px-6 py-4">
					<div className="flex items-center justify-between">
						<div className="flex items-center space-x-4">
							<h1 className="font-semibold text-gray-800 text-xl">
								{selectedDatabase
									? `Database: ${selectedDatabase}`
									: "Select a database"}
							</h1>
						</div>
						<Button
							onClick={onNewConnection}
							variant="outline"
							className="border-gray-300 text-gray-700 hover:bg-gray-50"
						>
							New Connection
						</Button>
					</div>
				</header>

				<main className="flex-1 p-6">
					{selectedDatabase ? (
						<div className="rounded-lg border border-gray-200 bg-white p-6">
							<h2 className="mb-4 font-medium text-gray-800 text-lg">
								Query Editor
							</h2>
							<div className="text-gray-500">
								Query editor will be implemented here for database:{" "}
								{selectedDatabase}
							</div>
						</div>
					) : (
						<div className="flex h-full items-center justify-center">
							<div className="text-center text-gray-500">
								<div className="mb-2 text-lg">
									Select a database from the sidebar to get started
								</div>
								<div className="text-sm">
									Available databases: {databases.length}
								</div>
							</div>
						</div>
					)}
				</main>
			</div>
		</div>
	);
};
