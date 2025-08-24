import type React from "react";
import { useState } from "react";
import { Connect } from "../../wailsjs/go/handlers/ConnectHandler";
import { Disconnect } from "../../wailsjs/go/handlers/DisconnectHandler";
import { TestConnection } from "../../wailsjs/go/handlers/TestConnectionHandler";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";

interface DatabaseConnectionFormProps {
	onConnectionChange: (connected: boolean, databases?: string[]) => void;
}

interface DatabaseConfig {
	host: string;
	port: number;
	database: string;
	username: string;
	password: string;
	sslmode: "disable" | "require" | "verify-ca" | "verify-full";
	connectionString: string;
	useConnectionString: boolean;
}

export const DatabaseConnectionForm: React.FC<DatabaseConnectionFormProps> = ({
	onConnectionChange,
}) => {
	const [config, setConfig] = useState<DatabaseConfig>({
		host: "localhost",
		port: 5432,
		database: "",
		username: "",
		password: "",
		sslmode: "require",
		connectionString: "",
		useConnectionString: false,
	});
	const [loading, setLoading] = useState(false);
	const [connected, setConnected] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleInputChange = (
		field: keyof DatabaseConfig,
		value: string | number | boolean,
	) => {
		setConfig((prev) => ({
			...prev,
			[field]: value,
		}));
	};

	const handleTestConnection = async () => {
		setLoading(true);
		setError(null);

		try {
			await TestConnection(config);
			setError(null);
			alert("Connection test successful!");
		} catch (err) {
			setError(err as string);
		} finally {
			setLoading(false);
		}
	};

	const handleConnect = async () => {
		setLoading(true);
		setError(null);

		try {
			const result = await Connect(config);

			setConnected(true);
			onConnectionChange(true, result?.databases || []);
			setError(null);

			// Log the databases received from the connection
			if (result?.databases) {
				console.log("Available databases:", result.databases);
			}
		} catch (err) {
			setError(err as string);
		} finally {
			setLoading(false);
		}
	};

	const handleDisconnect = async () => {
		setLoading(true);

		try {
			await Disconnect();
			setConnected(false);
			onConnectionChange(false);
			setError(null);
		} catch (err) {
			setError(err as string);
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="mx-auto w-full max-w-2xl rounded-lg bg-white p-8 shadow-lg">
			<h2 className="mb-4 font-bold text-gray-800 text-xl">
				Database Connection
			</h2>

			<div className="mb-6">
				<div className="flex items-center space-x-4">
					<label className="flex items-center">
						<input
							type="radio"
							name="connectionType"
							checked={!config.useConnectionString}
							onChange={() => handleInputChange("useConnectionString", false)}
							disabled={connected || loading}
							className="mr-2"
						/>
						<span className="text-gray-700">Connection Form</span>
					</label>
					<label className="flex items-center">
						<input
							type="radio"
							name="connectionType"
							checked={config.useConnectionString}
							onChange={() => handleInputChange("useConnectionString", true)}
							disabled={connected || loading}
							className="mr-2"
						/>
						<span className="text-gray-700">Connection String</span>
					</label>
				</div>
			</div>

			{error && (
				<div className="mb-4 rounded border border-red-400 bg-red-100 p-3 text-red-700">
					{error}
				</div>
			)}

			{connected && (
				<div className="mb-4 rounded border border-green-400 bg-green-100 p-3 text-green-700">
					Connected to database successfully
				</div>
			)}

			{config.useConnectionString ? (
				<div className="space-y-4">
					<div>
						<Label htmlFor="connectionString" className="text-gray-700">
							Connection String
						</Label>
						<textarea
							id="connectionString"
							value={config.connectionString}
							onChange={(e) =>
								handleInputChange("connectionString", e.target.value)
							}
							disabled={connected || loading}
							placeholder="postgresql://username:password@localhost:5432/database_name?sslmode=require"
							className="flex min-h-[80px] w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-gray-800 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
							rows={3}
						/>
						<p className="mt-1 text-gray-500 text-sm">
							Example:
							postgresql://user:pass@localhost:5432/mydb?sslmode=require
						</p>
					</div>
				</div>
			) : (
				<div className="grid grid-cols-2 gap-4">
					<div>
						<Label htmlFor="host" className="text-gray-700">
							Host
						</Label>
						<Input
							id="host"
							type="text"
							value={config.host}
							onChange={(e) => handleInputChange("host", e.target.value)}
							disabled={connected || loading}
							className="text-gray-800"
						/>
					</div>

					<div>
						<Label htmlFor="port" className="text-gray-700">
							Port
						</Label>
						<Input
							id="port"
							type="number"
							value={config.port}
							onChange={(e) =>
								handleInputChange("port", Number.parseInt(e.target.value))
							}
							disabled={connected || loading}
							className="text-gray-800"
						/>
					</div>

					<div>
						<Label htmlFor="database" className="text-gray-700">
							Database
						</Label>
						<Input
							id="database"
							type="text"
							value={config.database}
							onChange={(e) => handleInputChange("database", e.target.value)}
							disabled={connected || loading}
							className="text-gray-800"
						/>
					</div>

					<div>
						<Label htmlFor="username" className="text-gray-700">
							Username
						</Label>
						<Input
							id="username"
							type="text"
							value={config.username}
							onChange={(e) => handleInputChange("username", e.target.value)}
							disabled={connected || loading}
							className="text-gray-800"
						/>
					</div>

					<div>
						<Label htmlFor="password" className="text-gray-700">
							Password
						</Label>
						<Input
							id="password"
							type="password"
							value={config.password}
							onChange={(e) => handleInputChange("password", e.target.value)}
							disabled={connected || loading}
							className="text-gray-800"
						/>
					</div>

					<div className="col-span-2">
						<Label htmlFor="sslmode" className="text-gray-700">
							SSL Mode
						</Label>
						<select
							id="sslmode"
							value={config.sslmode}
							onChange={(e) => handleInputChange("sslmode", e.target.value)}
							disabled={connected || loading}
							className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-gray-800 text-sm ring-offset-background file:border-0 file:bg-transparent file:font-medium file:text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value="disable">Disable</option>
							<option value="require">Require</option>
							<option value="verify-ca">Verify CA</option>
							<option value="verify-full">Verify Full</option>
						</select>
					</div>
				</div>
			)}

			<div className="mt-6 flex space-x-2">
				{!connected ? (
					<>
						<Button
							onClick={handleTestConnection}
							disabled={loading}
							variant="secondary"
							className="flex-1 border-gray-300 bg-gray-200 text-gray-800 hover:bg-gray-300"
						>
							{loading ? "Testing..." : "Test"}
						</Button>
						<Button
							onClick={handleConnect}
							disabled={loading}
							className="flex-1 bg-blue-600 text-white hover:bg-blue-700"
						>
							{loading ? "Connecting..." : "Connect"}
						</Button>
					</>
				) : (
					<Button
						onClick={handleDisconnect}
						disabled={loading}
						variant="destructive"
						className="w-full bg-red-600 text-white hover:bg-red-700"
					>
						{loading ? "Disconnecting..." : "Disconnect"}
					</Button>
				)}
			</div>
		</div>
	);
};
