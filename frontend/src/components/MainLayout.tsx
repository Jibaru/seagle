import type React from "react";
import { useDatabaseStore } from "../store/DatabaseStore";
import { QueryInterface } from "./QueryInterface";
import { Sidebar } from "./Sidebar";
import { ThemeToggle } from "./ThemeToggle";
import { Button } from "./ui/button";

interface MainLayoutProps {
	onNewConnection: () => void;
}

export const MainLayout: React.FC<MainLayoutProps> = ({ onNewConnection }) => {
	const { state, selectDatabase, selectTable } = useDatabaseStore();

	return (
		<div className="flex h-screen overflow-hidden bg-gray-50 dark:bg-gray-900">
			<Sidebar />

			<div className="flex flex-1 flex-col overflow-hidden">
				<header className="flex-shrink-0 border-gray-200 border-b bg-white px-6 py-4 dark:border-gray-700 dark:bg-gray-800">
					<div className="flex items-center justify-between">
						<div className="flex items-center space-x-4">
							<h1 className="font-semibold text-gray-800 text-xl dark:text-gray-200">
								{state.selectedTable
									? `${state.selectedDatabase}.${state.selectedTable}`
									: state.selectedDatabase
										? `Database: ${state.selectedDatabase}`
										: "Select a database"}
							</h1>
						</div>
						<div className="flex items-center space-x-3">
							<ThemeToggle />
							<Button
								onClick={onNewConnection}
								variant="outline"
								className="border-gray-300 text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
							>
								New Connection
							</Button>
						</div>
					</div>
				</header>

				<main className="flex-1 overflow-hidden" style={{ minHeight: 0 }}>
					{state.selectedTable ? (
						<div className="flex flex-1 flex-col">
							<div className="border-gray-200 border-b bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
								<h2 className="font-medium text-gray-800 text-lg dark:text-gray-200">
									Table: {state.selectedDatabase}.{state.selectedTable}
								</h2>
								<div className="mt-1 text-gray-500 text-sm dark:text-gray-400">
									Table structure and data view - Coming soon
								</div>
							</div>
							<div className="flex-1 p-6">
								<div className="h-full rounded-lg border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800">
									<div className="text-gray-500 dark:text-gray-400">
										Table structure and data view will be implemented here for{" "}
										{state.selectedDatabase}.{state.selectedTable}
									</div>
								</div>
							</div>
						</div>
					) : state.selectedDatabase ? (
						<QueryInterface database={state.selectedDatabase} />
					) : (
						<div className="flex h-full items-center justify-center">
							<div className="text-center text-gray-500 dark:text-gray-400">
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
