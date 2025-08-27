import { useEffect } from "react";
import { Database, Server, Loader2, AlertCircle, Trash2 } from "lucide-react";
import type React from "react";
import { Button } from "./ui/button";
import { useConnectionsStore } from "../store/ConnectionsStore";

interface SavedConnectionsProps {
	onConnectToSaved?: (connectionId: string) => void;
}

export const SavedConnections: React.FC<SavedConnectionsProps> = ({
	onConnectToSaved,
}) => {
	const { state, loadConnections, deleteConnection } = useConnectionsStore();
	const { connections, loading, error, connectingId, deletingId } = state;

	useEffect(() => {
		loadConnections();
	}, [loadConnections]);

	if (loading) {
		return (
			<div className="w-full max-w-4xl">
				<div className="flex items-center justify-center p-8">
					<Loader2 className="h-6 w-6 animate-spin text-blue-500" />
					<span className="ml-2 text-gray-600 dark:text-gray-400">
						Loading saved connections...
					</span>
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="w-full max-w-4xl">
				<div className="flex items-center justify-center p-8">
					<AlertCircle className="h-6 w-6 text-red-500" />
					<span className="ml-2 text-red-600 dark:text-red-400">{error}</span>
				</div>
			</div>
		);
	}

	if (connections.length === 0) {
		return (
			<div className="w-full max-w-4xl">
				<div className="text-center p-8">
					<Database className="h-12 w-12 text-gray-400 mx-auto mb-3" />
					<p className="text-gray-600 dark:text-gray-400">
						No saved connections found
					</p>
				</div>
			</div>
		);
	}

	return (
		<div className="w-full max-w-4xl">
			<div className="mb-4">
				<h2 className="text-xl font-semibold text-gray-900 dark:text-white">
					Recent Connections
				</h2>
				<p className="text-sm text-gray-600 dark:text-gray-400">
					{connections.length} saved connection{connections.length !== 1 ? "s" : ""}
				</p>
			</div>

			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{connections.map((connection) => (
					<div
						key={connection.id}
						className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow min-w-0"
					>
						<div className="flex items-start justify-between mb-3 min-w-0">
							<div className="flex items-center min-w-0 flex-1 mr-2">
								<Server className="h-5 w-5 text-blue-500 mr-2 flex-shrink-0" />
								<div className="min-w-0 flex-1">
									<h3 className="text-sm font-medium text-gray-900 dark:text-white truncate" title={connection.host}>
										{connection.host}
									</h3>
									<p className="text-xs text-gray-500 dark:text-gray-400">
										Port {connection.port}
									</p>
								</div>
							</div>
						</div>

						<div className="flex items-center justify-between min-w-0 gap-2">
							<div className="flex items-center text-xs text-gray-500 dark:text-gray-400 min-w-0 flex-1">
								<span className="truncate" title={connection.id}>
									{connection.id.substring(0, 8)}...
								</span>
							</div>
							<div className="flex items-center gap-2">
								<Button
									size="sm"
									variant="outline"
									onClick={() => deleteConnection(connection.id)}
									disabled={deletingId === connection.id || connectingId === connection.id}
									className="bg-red-50 hover:bg-red-100 text-red-600 border-red-200 dark:bg-red-900/20 dark:hover:bg-red-900/30 dark:text-red-400 dark:border-red-800 flex-shrink-0 disabled:opacity-50"
								>
									{deletingId === connection.id ? (
										<Loader2 className="h-3 w-3 animate-spin" />
									) : (
										<Trash2 className="h-3 w-3" />
									)}
								</Button>
								<Button
									size="sm"
									onClick={() => onConnectToSaved?.(connection.id)}
									disabled={connectingId === connection.id || deletingId === connection.id}
									className="bg-blue-600 hover:bg-blue-700 text-white dark:bg-blue-700 dark:hover:bg-blue-800 flex-shrink-0 disabled:opacity-50"
								>
									{connectingId === connection.id ? (
										<>
											<Loader2 className="h-3 w-3 animate-spin mr-1" />
											Connecting...
										</>
									) : (
										"Connect"
									)}
								</Button>
							</div>
						</div>
					</div>
				))}
			</div>
		</div>
	);
};