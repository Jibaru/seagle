import type React from "react";
import { createContext, useContext, useEffect, useState } from "react";

type Theme = "light" | "dark" | "system";

interface ThemeContextType {
	theme: Theme;
	setTheme: (theme: Theme) => void;
	actualTheme: "light" | "dark";
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const useTheme = () => {
	const context = useContext(ThemeContext);
	if (!context) {
		throw new Error("useTheme must be used within a ThemeProvider");
	}
	return context;
};

interface ThemeProviderProps {
	children: React.ReactNode;
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
	const [theme, setTheme] = useState<Theme>(() => {
		// Get theme from localStorage or default to system
		const stored = localStorage.getItem("theme") as Theme;
		return stored || "system";
	});

	const [actualTheme, setActualTheme] = useState<"light" | "dark">("light");

	// Function to get system theme preference
	const getSystemTheme = (): "light" | "dark" => {
		return window.matchMedia("(prefers-color-scheme: dark)").matches
			? "dark"
			: "light";
	};

	// Update actual theme based on current theme setting
	useEffect(() => {
		const updateActualTheme = () => {
			const newActualTheme = theme === "system" ? getSystemTheme() : theme;
			setActualTheme(newActualTheme);

			// Update document class
			const root = document.documentElement;
			root.classList.remove("light", "dark");
			root.classList.add(newActualTheme);
		};

		updateActualTheme();

		// Listen for system theme changes if using system theme
		if (theme === "system") {
			const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
			const handleChange = () => updateActualTheme();
			mediaQuery.addEventListener("change", handleChange);

			return () => mediaQuery.removeEventListener("change", handleChange);
		}
	}, [theme]);

	// Save theme to localStorage when it changes
	useEffect(() => {
		localStorage.setItem("theme", theme);
	}, [theme]);

	return (
		<ThemeContext.Provider value={{ theme, setTheme, actualTheme }}>
			{children}
		</ThemeContext.Provider>
	);
};