<script lang="ts">
	import { onMount } from 'svelte';
	import AppBase from '$lib/layouts/AppBase.svelte';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Globe, Palette, Monitor, Sun, Moon, Cog } from 'lucide-svelte';
	import { settingsStore, type UISettings } from '$lib/stores/settings.svelte';
	import { mode } from 'mode-watcher';

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
		// Settings are loaded automatically via constructor
	});

	async function updateLanguage(language: UISettings['language']) {
		await settingsStore.updateLanguage(language);
	}

	async function updateTheme(theme: UISettings['theme']) {
		await settingsStore.updateTheme(theme);
	}


	function getLanguageLabel(value: string) {
		return languages.find((lang) => lang.value === value)?.label || value;
	}

	function getThemeLabel(value: string) {
		return themes.find((theme) => theme.value === value)?.label || value;
	}
</script>

<svelte:head>
	<title>Application - Teatime</title>
</svelte:head>

<AppBase title="Application" icon={Cog}>
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
						value={settingsStore.language}
						onValueChange={(value) =>
							value && updateLanguage(value as UISettings['language'])}
					>
						<Select.Trigger id="language-select" class="w-full max-w-sm">
							<div class="flex items-center gap-2">
								<span>{languages.find((lang) => lang.value === settingsStore.language)?.flag}</span>
								<span>{getLanguageLabel(settingsStore.language)}</span>
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
							value={settingsStore.theme}
							onValueChange={(value) =>
								value && updateTheme(value as UISettings['theme'])}
						>
							<Select.Trigger id="theme-select" class="w-full max-w-sm">
								<div class="flex items-center gap-2">
									{#if settingsStore.theme === 'light'}
										<Sun class="h-4 w-4" />
									{:else if settingsStore.theme === 'dark'}
										<Moon class="h-4 w-4" />
									{:else}
										<Monitor class="h-4 w-4" />
									{/if}
									<span>{getThemeLabel(settingsStore.theme)}</span>
									{#if settingsStore.theme === 'system'}
										<span class="text-muted-foreground text-xs">
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
</AppBase>
