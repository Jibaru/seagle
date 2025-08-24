import { Database } from "lucide-react";
import type React from "react";

interface SidebarProps {
	databases: string[];
	onDatabaseSelect?: (database: string) => void;
	selectedDatabase?: string;
}

export const Sidebar: React.FC<SidebarProps> = ({
	databases,
	onDatabaseSelect,
	selectedDatabase,
}) => {
	return (
		<div className="w-64 border-gray-200 border-r bg-white shadow-sm">
			<div className="border-gray-200 border-b p-4">
				<h2 className="font-semibold text-gray-800 text-lg">Databases</h2>
			</div>
			<div className="p-2">
				{databases.length === 0 ? (
					<div className="p-4 text-center text-gray-500">
						No databases available
					</div>
				) : (
					<ul className="space-y-1">
						{databases.map((database) => (
							<li key={database}>
								<button
									type="button"
									onClick={() => onDatabaseSelect?.(database)}
									className={`flex w-full items-center rounded-md px-3 py-2 text-left text-sm transition-colors ${
										selectedDatabase === database
											? "bg-blue-100 text-blue-700"
											: "text-gray-700 hover:bg-gray-100"
									}`}
								>
									<Database className="mr-2 h-4 w-4 flex-shrink-0" />
									<span className="truncate">{database}</span>
								</button>
							</li>
						))}
					</ul>
				)}
			</div>
		</div>
	);
};
