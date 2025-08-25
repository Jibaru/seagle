import { Bird } from "lucide-react";
import type React from "react";
import { Button } from "./ui/button";

interface WelcomeScreenProps {
	onNewConnection: () => void;
}

export const WelcomeScreen: React.FC<WelcomeScreenProps> = ({
	onNewConnection,
}) => {
	return (
		<div className="flex min-h-screen items-center justify-center bg-gray-50 dark:bg-gray-950">
			<div className="space-y-8 text-center">
				<div className="mb-6 flex justify-center">
					<Bird className="h-24 w-24 text-blue-500 dark:text-blue-400" />
				</div>

				<div className="space-y-4">
					<h1 className="font-bold text-5xl text-gray-900 dark:text-white">Seagle</h1>
					<p className="max-w-md text-gray-700 text-xl dark:text-gray-300">
						AI-powered database management tool.
					</p>
				</div>

				<div className="space-y-4">
					<p className="text-gray-600 dark:text-gray-400">
						Connect to your PostgreSQL database to get started
					</p>

					<Button
						onClick={onNewConnection}
						size="lg"
						className="bg-blue-600 px-8 py-3 text-lg text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-800"
					>
						New Connection
					</Button>
				</div>
			</div>
		</div>
	);
};
