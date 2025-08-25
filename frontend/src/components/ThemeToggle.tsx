import { Monitor, Moon, Sun } from "lucide-react";
import type React from "react";
import { useTheme } from "../contexts/ThemeContext";
import { Button } from "./ui/button";

export const ThemeToggle: React.FC = () => {
	const { theme, setTheme } = useTheme();

	const themes = [
		{
			value: "light" as const,
			label: "Light",
			icon: Sun,
		},
		{
			value: "dark" as const,
			label: "Dark", 
			icon: Moon,
		},
		{
			value: "system" as const,
			label: "System",
			icon: Monitor,
		},
	];

	const currentTheme = themes.find((t) => t.value === theme);
	const Icon = currentTheme?.icon || Sun;

	const cycleTheme = () => {
		const currentIndex = themes.findIndex((t) => t.value === theme);
		const nextIndex = (currentIndex + 1) % themes.length;
		setTheme(themes[nextIndex].value);
	};

	return (
		<Button
			variant="outline"
			size="sm"
			onClick={cycleTheme}
			className="flex items-center space-x-2"
			title={`Current: ${currentTheme?.label} theme. Click to cycle themes.`}
		>
			<Icon className="h-4 w-4" />
			<span className="text-xs">{currentTheme?.label}</span>
		</Button>
	);
};