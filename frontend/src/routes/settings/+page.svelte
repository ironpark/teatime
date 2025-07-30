<script lang="ts">
	import { onMount } from 'svelte';
	import { SidebarTrigger } from '$lib/components/ui/sidebar';
	import { Separator } from '$lib/components/ui/separator';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Settings, Globe, Palette, Monitor, Sun, Moon } from 'lucide-svelte';
	import { settingsStore, type Settings as SettingsType } from '$lib/stores/settings';
	import { mode } from 'mode-watcher';

	let settings = $state<SettingsType>({
		language: 'en',
		theme: 'system'
	});


	const languages = [
		{ value: 'en' as const, label: 'English', flag: '🇺🇸' },
		{ value: 'ko' as const, label: '한국어', flag: '🇰🇷' },
		{ value: 'ja' as const, label: '日本語', flag: '🇯🇵' },
		{ value: 'zh' as const, label: '中文', flag: '🇨🇳' }
	];

	const themes = [
		{ value: 'light' as const, label: 'Light', icon: Sun, description: 'Light theme' },
		{ value: 'dark' as const, label: 'Dark', icon: Moon, description: 'Dark theme' },
		{
			value: 'system' as const,
			label: 'System',
			icon: Monitor,
			description: 'Follow system preference'
		}
	];

	onMount(() => {
		// Load settings from store
		settingsStore.load();
		const unsubscribe = settingsStore.subscribe((value) => {
			settings = { ...value };
		});

		return unsubscribe;
	});

	function updateSetting<K extends keyof SettingsType>(key: K, value: SettingsType[K]) {
		// Update settings and save immediately
		settingsStore.updateSetting(key, value);
	}

	function getLanguageLabel(value: string) {
		return languages.find((lang) => lang.value === value)?.label || value;
	}

	function getThemeLabel(value: string) {
		return themes.find((theme) => theme.value === value)?.label || value;
	}
</script>

<svelte:head>
	<title>Settings - Teatime</title>
</svelte:head>

<div class="settings-page bg-background flex h-screen w-full flex-col">
	<!-- Header -->
	<header class="bg-card border-b">
		<div class="flex h-16 items-center gap-2 px-4">
			<SidebarTrigger />
			<Separator orientation="vertical" class="mr-2 h-4" />
			<div class="flex items-center gap-2">
				<Settings class="text-muted-foreground h-5 w-5" />
				<h1 class="text-lg font-semibold">Settings</h1>
			</div>
		</div>
	</header>

	<!-- Main content -->
	<main class="flex-1 overflow-y-auto p-6">
		<div class="space-y-8">
			<!-- Language Settings -->
			<section class="space-y-6">
				<div class="space-y-2">
					<h2 class="flex items-center gap-2 text-xl font-semibold">
						<Globe class="h-5 w-5" />
						Language & Region
					</h2>
					<p class="text-muted-foreground text-sm">
						Choose your preferred language for the interface
					</p>
				</div>

				<div class="space-y-3">
					<Label for="language-select">Interface Language</Label>
					<Select.Root
						type="single"
						value={settings.language}
						onValueChange={(value) =>
							value && updateSetting('language', value as SettingsType['language'])}
					>
						<Select.Trigger id="language-select" class="w-full max-w-sm">
							<div class="flex items-center gap-2">
								<span>{languages.find((lang) => lang.value === settings.language)?.flag}</span>
								<span>{getLanguageLabel(settings.language)}</span>
							</div>
						</Select.Trigger>
						<Select.Content>
							{#each languages as language}
								<Select.Item value={language.value}>
									<div class="flex items-center gap-2">
										<span>{language.flag}</span>
										<span>{language.label}</span>
									</div>
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</section>

			<!-- Theme Settings -->
			<section class="space-y-6">
				<div class="space-y-2">
					<h2 class="flex items-center gap-2 text-xl font-semibold">
						<Palette class="h-5 w-5" />
						Appearance
					</h2>
					<p class="text-muted-foreground text-sm">
						Customize the look and feel of the application
					</p>
				</div>

				<div class="space-y-4">
					<div class="space-y-3">
						<Label for="theme-select">Theme</Label>
						<Select.Root
							type="single"
							value={settings.theme}
							onValueChange={(value) =>
								value && updateSetting('theme', value as SettingsType['theme'])}
						>
							<Select.Trigger id="theme-select" class="w-full max-w-sm">
								<div class="flex items-center gap-2">
									{#if settings.theme === 'light'}
										<Sun class="h-4 w-4" />
									{:else if settings.theme === 'dark'}
										<Moon class="h-4 w-4" />
									{:else}
										<Monitor class="h-4 w-4" />
									{/if}
									<span>{getThemeLabel(settings.theme)}</span>
									{#if settings.theme === 'system'}
										<span class="text-xs text-muted-foreground">
											(currently {mode.current})
										</span>
									{/if}
								</div>
							</Select.Trigger>
							<Select.Content>
								{#each themes as theme}
									<Select.Item value={theme.value}>
										<div class="flex items-center gap-2">
											<theme.icon class="h-4 w-4" />
											<div class="flex flex-col">
												<span>{theme.label}</span>
												<span class="text-muted-foreground text-xs">{theme.description}</span>
											</div>
										</div>
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				</div>
			</section>

			<!-- Additional Settings -->
			<section class="space-y-6">
				<div class="space-y-2">
					<h2 class="text-xl font-semibold">Advanced</h2>
					<p class="text-muted-foreground text-sm">Additional application settings</p>
				</div>

				<div class="rounded-lg border border-dashed p-6 text-center">
					<p class="text-muted-foreground text-sm">
						More settings will be available in future updates.
					</p>
				</div>
			</section>
		</div>
	</main>
</div>

<style>
	/* Custom scrollbar */
	.overflow-y-auto::-webkit-scrollbar {
		width: 6px;
	}

	.overflow-y-auto::-webkit-scrollbar-track {
		background: transparent;
	}

	.overflow-y-auto::-webkit-scrollbar-thumb {
		background-color: hsl(var(--border));
		border-radius: 3px;
	}

	.overflow-y-auto::-webkit-scrollbar-thumb:hover {
		background-color: hsl(var(--muted-foreground));
	}
</style>
