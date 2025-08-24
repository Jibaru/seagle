import type React from "react";
import { useDatabaseStore } from "../store/DatabaseStore";
import { Sidebar } from "./Sidebar";
import { Button } from "./ui/button";

interface MainLayoutProps {
	onNewConnection: () => void;
}

export const MainLayout: React.FC<MainLayoutProps> = ({ onNewConnection }) => {
	const { state, selectDatabase, selectTable } = useDatabaseStore();

	return (
		<div className="flex h-screen bg-gray-50">
			<Sidebar />

			<div className="flex flex-1 flex-col">
				<header className="border-gray-200 border-b bg-white px-6 py-4">
					<div className="flex items-center justify-between">
						<div className="flex items-center space-x-4">
							<h1 className="font-semibold text-gray-800 text-xl">
								{state.selectedTable
									? `${state.selectedDatabase}.${state.selectedTable}`
									: state.selectedDatabase
										? `Database: ${state.selectedDatabase}`
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
					{state.selectedTable ? (
						<div className="rounded-lg border border-gray-200 bg-white p-6">
							<h2 className="mb-4 font-medium text-gray-800 text-lg">
								Table: {state.selectedTable}
							</h2>
							<div className="text-gray-500">
								Table structure and data view will be implemented here for{" "}
								{state.selectedDatabase}.{state.selectedTable}
							</div>
						</div>
					) : state.selectedDatabase ? (
						<div className="rounded-lg border border-gray-200 bg-white p-6">
							<h2 className="mb-4 font-medium text-gray-800 text-lg">
								Query Editor - {state.selectedDatabase}
							</h2>
							<div className="text-gray-500">
								SQL query editor will be implemented here for database:{" "}
								{state.selectedDatabase}
								<br />
								<br />
								Select a table from the sidebar to view its structure.
							</div>
						</div>
					) : (
						<div className="flex h-full items-center justify-center">
							<div className="text-center text-gray-500">
								<div className="mb-2 text-lg">
									Select a database from the sidebar to get started
								</div>
								<div className="text-sm">
									Available databases: {state.databases.length}
								</div>
							</div>
						</div>
					)}
				</main>
			</div>
		</div>
	);
};
