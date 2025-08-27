import { Eye, EyeOff, Settings, X } from "lucide-react";
import type React from "react";
import { useEffect, useState } from "react";
import { GetConfig } from "../../wailsjs/go/handlers/GetConfigHandler";
import { SetConfig } from "../../wailsjs/go/handlers/SetConfigHandler";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";

interface SettingsFormProps {
	isOpen: boolean;
	onClose: () => void;
}

export const SettingsForm: React.FC<SettingsFormProps> = ({
	isOpen,
	onClose,
}) => {
	const [openAIAPIKey, setOpenAIAPIKey] = useState("");
	const [originalKey, setOriginalKey] = useState("");
	const [loading, setLoading] = useState(false);
	const [saving, setSaving] = useState(false);
	const [message, setMessage] = useState("");
	const [showPassword, setShowPassword] = useState(false);

	useEffect(() => {
		if (isOpen) {
			loadConfig();
		}
	}, [isOpen]);

	const loadConfig = async () => {
		try {
			setLoading(true);
			const result = await GetConfig();
			if (result.success && result.Config) {
				const key = result.Config.openAIAPIKey || "";
				setOpenAIAPIKey(key);
				setOriginalKey(key);
			}
		} catch (error) {
			console.error("Failed to load config:", error);
			setMessage("Failed to load configuration");
		} finally {
			setLoading(false);
		}
	};

	const handleSave = async () => {
		try {
			setSaving(true);
			setMessage("");

			const result = await SetConfig({
				openAIAPIKey: openAIAPIKey,
			});

			if (result.success) {
				setMessage("Configuration saved successfully!");
				setOriginalKey(openAIAPIKey);
				setTimeout(() => {
					setMessage("");
					onClose();
				}, 2000);
			} else {
				setMessage(result.message || "Failed to save configuration");
			}
		} catch (error) {
			console.error("Failed to save config:", error);
			setMessage("Failed to save configuration");
		} finally {
			setSaving(false);
		}
	};

	const handleCancel = () => {
		setOpenAIAPIKey(originalKey);
		setMessage("");
		onClose();
	};

	const hasChanges = openAIAPIKey !== originalKey;

	if (!isOpen) return null;

	return (
		<div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
			<div className="w-full max-w-md rounded-lg border border-gray-200 bg-white p-6 shadow-lg dark:border-gray-700 dark:bg-gray-800">
				<div className="mb-4 flex items-center justify-between">
					<div className="flex items-center space-x-2">
						<Settings className="h-5 w-5 text-gray-700 dark:text-gray-300" />
						<h2 className="font-semibold text-gray-900 text-lg dark:text-white">
							Settings
						</h2>
					</div>
					<Button
						variant="outline"
						size="sm"
						onClick={handleCancel}
						className="h-8 w-8 p-0"
					>
						<X className="h-4 w-4" />
					</Button>
				</div>

				{loading ? (
					<div className="py-8 text-center">
						<div className="text-gray-600 dark:text-gray-400">
							Loading configuration...
						</div>
					</div>
				) : (
					<div className="space-y-4">
						<div className="space-y-2">
							<Label
								htmlFor="openai-key"
								className="text-gray-700 dark:text-gray-300"
							>
								OpenAI API Key
							</Label>
							<div className="relative">
								<Input
									id="openai-key"
									type={showPassword ? "text" : "password"}
									value={openAIAPIKey}
									onChange={(e) => setOpenAIAPIKey(e.target.value)}
									placeholder="Enter your OpenAI API key"
									className="pr-10 border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
								/>
								<Button
									type="button"
									variant="ghost"
									size="sm"
									onClick={() => setShowPassword(!showPassword)}
									className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
								>
									{showPassword ? (
										<EyeOff className="h-4 w-4 text-gray-500 dark:text-gray-400" />
									) : (
										<Eye className="h-4 w-4 text-gray-500 dark:text-gray-400" />
									)}
								</Button>
							</div>
							<p className="text-gray-500 text-sm dark:text-gray-400">
								Required for AI-powered query generation features
							</p>
						</div>

						{message && (
							<div
								className={`rounded-md border p-3 text-sm ${
									message.includes("success")
										? "border-green-200 bg-green-50 text-green-700 dark:border-green-700 dark:bg-green-900/20 dark:text-green-400"
										: "border-red-200 bg-red-50 text-red-700 dark:border-red-700 dark:bg-red-900/20 dark:text-red-400"
								}`}
							>
								{message}
							</div>
						)}

						<div className="flex justify-end space-x-2 pt-4">
							<Button
								variant="outline"
								onClick={handleCancel}
								disabled={saving}
								className="border-gray-300 text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
							>
								Cancel
							</Button>
							<Button
								onClick={handleSave}
								disabled={saving || !hasChanges}
								className="bg-blue-600 text-white hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-800"
							>
								{saving ? "Saving..." : "Save"}
							</Button>
						</div>
					</div>
				)}
			</div>
		</div>
	);
};