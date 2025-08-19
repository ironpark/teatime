import { SettingsService } from '$bindings/services';
import { Setting } from '$bindings/internal/database/models';
import { browser } from '$app/environment';
import { resetMode, setMode } from 'mode-watcher';

// Frontend settings interface (for UI compatibility)
export interface UISettings {
    language: 'en' | 'ko' | 'ja' | 'zh';
    theme: 'light' | 'dark' | 'system';
    autoStart: boolean;
}

// Default settings
const defaultSettings: UISettings = {
    language: 'en',
    theme: 'system',
    autoStart: false
};

export class SettingsStore {
    private _settings = $state<UISettings>(defaultSettings);
    private _loading = $state(false);
    private _error = $state<string | null>(null);

    constructor() {
        this.load();
    }

    get settings(): UISettings {
        return this._settings;
    }

    get isLoading(): boolean {
        return this._loading;
    }

    get error(): string | null {
        return this._error;
    }

    get language(): UISettings['language'] {
        return this._settings.language;
    }

    get theme(): UISettings['theme'] {
        return this._settings.theme;
    }

    get autoStart(): boolean {
        return this._settings.autoStart;
    }

    async load(): Promise<boolean> {
        return await this.withLoading(async () => {
            try {
                const dbSettings = await SettingsService.GetSettings();
                if (dbSettings) {
                    this._settings = this.mapFromDatabase(dbSettings);
                    this.applyTheme(this._settings.theme);
                }
                return true;
            } catch (error) {
                // Fallback to localStorage if backend fails
                this.loadFromLocalStorage();
                throw new Error('Failed to load settings from backend, using local storage');
            }
        }) ?? false;
    }

    async updateSettings(settings: Partial<UISettings>): Promise<boolean> {
        return await this.withLoading(async () => {
            const updatedSettings = { ...this._settings, ...settings };
            
            // Update backend
            const dbSettings = this.mapToDatabase(updatedSettings);
            await SettingsService.UpdateSettings(dbSettings);
            
            // Update local state
            this._settings = updatedSettings;
            
            // Apply theme changes immediately
            if (settings.theme) {
                this.applyTheme(settings.theme);
            }
            
            // Persist to localStorage as backup
            this.saveToLocalStorage(updatedSettings);
            
            return true;
        }) ?? false;
    }

    async updateTheme(theme: UISettings['theme']): Promise<boolean> {
        return await this.withLoading(async () => {
            await SettingsService.UpdateTheme(theme);
            this._settings.theme = theme;
            this.applyTheme(theme);
            this.saveToLocalStorage(this._settings);
            return true;
        }) ?? false;
    }

    async updateLanguage(language: UISettings['language']): Promise<boolean> {
        return await this.updateSettings({ language });
    }

    async updateAutoStart(autoStart: boolean): Promise<boolean> {
        return await this.updateSettings({ autoStart });
    }

    // Helper methods
    private mapFromDatabase(dbSettings: Setting): UISettings {
        return {
            language: this.validateLanguage(dbSettings.Language),
            theme: this.validateTheme(dbSettings.Theme),
            autoStart: dbSettings.AutoStart === 1
        };
    }

    private mapToDatabase(uiSettings: UISettings): Setting {
        return Setting.createFrom({
            Theme: uiSettings.theme,
            Language: uiSettings.language,
            AutoStart: uiSettings.autoStart
        });
    }

    private validateLanguage(language: string): UISettings['language'] {
        const validLanguages: UISettings['language'][] = ['en', 'ko', 'ja', 'zh'];
        return validLanguages.includes(language as UISettings['language']) 
            ? language as UISettings['language'] 
            : 'en';
    }

    private validateTheme(theme: string): UISettings['theme'] {
        const validThemes: UISettings['theme'][] = ['light', 'dark', 'system'];
        return validThemes.includes(theme as UISettings['theme']) 
            ? theme as UISettings['theme'] 
            : 'system';
    }

    private applyTheme(theme: UISettings['theme']): void {
        if (!browser) return;
        
        if (theme === 'system') {
            resetMode();
        } else {
            setMode(theme);
        }
    }

    private loadFromLocalStorage(): void {
        if (!browser) return;
        
        try {
            const saved = localStorage.getItem('teatime-settings');
            if (saved) {
                const parsed = JSON.parse(saved) as Partial<UISettings>;
                this._settings = { ...defaultSettings, ...parsed };
                this.applyTheme(this._settings.theme);
            }
        } catch (error) {
            console.error('Failed to load settings from localStorage:', error);
        }
    }

    private saveToLocalStorage(settings: UISettings): void {
        if (!browser) return;
        
        try {
            localStorage.setItem('teatime-settings', JSON.stringify(settings));
        } catch (error) {
            console.error('Failed to save settings to localStorage:', error);
        }
    }

    private async withLoading<T>(operation: () => Promise<T>): Promise<T | null> {
        try {
            this._loading = true;
            this._error = null;
            return await operation();
        } catch (error) {
            this._error = error instanceof Error ? error.message : 'Unknown error';
            return null;
        } finally {
            this._loading = false;
        }
    }

    // Reset to defaults
    async reset(): Promise<boolean> {
        return await this.updateSettings(defaultSettings);
    }
}

// Export singleton instance
export const settingsStore = new SettingsStore();

// Export theme apply function for backward compatibility
export function applyTheme(theme: UISettings['theme']): void {
    if (!browser) return;
    
    if (theme === 'system') {
        resetMode();
    } else {
        setMode(theme);
    }
}