import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { resetMode, setMode } from 'mode-watcher';

export interface Settings {
	language: 'en' | 'ko' | 'ja' | 'zh';
	theme: 'light' | 'dark' | 'system';
}

const defaultSettings: Settings = {
	language: 'en',
	theme: 'system'
};

function createSettingsStore() {
	const { subscribe, set, update } = writable<Settings>(defaultSettings);

	return {
		subscribe,
		set,
		update,
		load: () => {
			if (!browser) return;
			
			try {
				const saved = localStorage.getItem('teatime-settings');
				if (saved) {
					const parsed = JSON.parse(saved) as Partial<Settings>;
					update(current => ({ ...current, ...parsed }));
				}
			} catch (error) {
				console.error('Failed to load settings:', error);
			}
		},
		save: (settings: Settings) => {
			if (!browser) return;
			
			try {
				localStorage.setItem('teatime-settings', JSON.stringify(settings));
				set(settings);
				applyTheme(settings.theme);
			} catch (error) {
				console.error('Failed to save settings:', error);
			}
		},
		updateSetting: <K extends keyof Settings>(key: K, value: Settings[K]) => {
			update(current => {
				const updated = { ...current, [key]: value };
				if (browser) {
					try {
						localStorage.setItem('teatime-settings', JSON.stringify(updated));
						if (key === 'theme') {
							applyTheme(value as Settings['theme']);
						}
					} catch (error) {
						console.error('Failed to save setting:', error);
					}
				}
				return updated;
			});
		}
	};
}

export function applyTheme(theme: Settings['theme']) {
	if (!browser) return;
	
	if (theme === 'system') {
		resetMode();
	} else {
		setMode(theme);
	}
}

export const settingsStore = createSettingsStore();