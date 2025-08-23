<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import AppBase from '$lib/layouts/AppBase.svelte';
	import { 
		Plus, 
		Edit, 
		Trash2, 
		Eye, 
		EyeOff, 
		Key, 
		Shield,
		Settings
	} from 'lucide-svelte';
	import { 
		ListSecrets, 
		CreateSecret, 
		UpdateSecret, 
		DeleteSecret,
		GetSecretData
	} from '$bindings/services/secretsservice';
	import {
		ListEnvironmentVariables,
		CreateEnvironmentVariable,
		UpdateEnvironmentVariable,
		DeleteEnvironmentVariable
	} from '$bindings/services/environmentvariablesservice';
	import type { 
		SecretInfo, 
		SecretCreateRequest, 
		SecretUpdateRequest,
		EnvironmentVariableInfo,
		EnvironmentVariableCreateRequest,
		EnvironmentVariableUpdateRequest
	} from '$bindings/services/models';
		
	// State
	let secrets = $state<SecretInfo[]>([]);
	let environments = $state<EnvironmentVariableInfo[]>([]);
	let loading = $state(false);
	let showCreateDialog = $state(false);
	let showEditDialog = $state(false);
	let showDeleteDialog = $state(false);
	let showCreateEnvDialog = $state(false);
	let showEditEnvDialog = $state(false);
	let showDeleteEnvDialog = $state(false);
	let editingSecret = $state<SecretInfo | null>(null);
	let deletingSecret = $state<SecretInfo | null>(null);
	let editingEnv = $state<EnvironmentVariableInfo | null>(null);
	let deletingEnv = $state<EnvironmentVariableInfo | null>(null);
	let showValue = $state(false);
	let activeTab = $state<'secrets' | 'environments'>('secrets');
	
	// Form state
	let formData = $state({
		name: '',
		description: '',
		value: ''
	});
	
	// Environment form state
	let envFormData = $state({
		name: '',
		description: '',
		value: ''
	});
	
	onMount(async () => {
		await loadSecrets();
		await loadEnvironments();
	});
	
	async function loadSecrets() {
		try {
			loading = true;
			secrets = await ListSecrets() || [];
		} catch (error) {
			console.error('Failed to load secrets:', error);
		} finally {
			loading = false;
		}
	}
	
	
	function openCreateDialog() {
		formData = {
			name: '',
			description: '',
			value: ''
		};
		showValue = false;
		showCreateDialog = true;
	}
	
	async function openEditDialog(secret: SecretInfo) {
		try {
			editingSecret = secret;
			const secretValue = await GetSecretData(secret.id);
			formData = {
				name: secret.name,
				description: secret.description,
				value: secretValue || ''
			};
			showValue = false;
			showEditDialog = true;
		} catch (error) {
			console.error('Failed to load secret data:', error);
		}
	}
	
	async function handleCreate() {
		try {
			const request: SecretCreateRequest = {
				name: formData.name,
				description: formData.description,
				value: formData.value
			};
			
			await CreateSecret(request);
			await loadSecrets();
			showCreateDialog = false;
		} catch (error) {
			console.error('Failed to create secret:', error);
			alert(`Failed to create secret: ${error instanceof Error ? error.message : String(error)}`);
		}
	}
	
	async function handleUpdate() {
		if (!editingSecret) return;
		
		try {
			const request: SecretUpdateRequest = {
				id: editingSecret.id,
				name: formData.name,
				description: formData.description,
				value: formData.value
			};
			
			await UpdateSecret(request);
			await loadSecrets();
			showEditDialog = false;
			editingSecret = null;
		} catch (error) {
			console.error('Failed to update secret:', error);
			alert(`Failed to update secret: ${error instanceof Error ? error.message : String(error)}`);
		}
	}
	
	function openDeleteDialog(secret: SecretInfo) {
		deletingSecret = secret;
		showDeleteDialog = true;
	}
	
	async function confirmDelete() {
		if (!deletingSecret) return;
		
		try {
			await DeleteSecret(deletingSecret.id);
			await loadSecrets();
			showDeleteDialog = false;
			deletingSecret = null;
		} catch (error) {
			console.error('Failed to delete secret:', error);
			alert(`Failed to delete secret: ${error instanceof Error ? error.message : String(error)}`);
		}
	}
	
	function toggleValueVisibility() {
		showValue = !showValue;
	}
	
	
	function formatDate(dateStr: string): string {
		if (!dateStr) return 'Never';
		return new Date(dateStr).toLocaleDateString();
	}
	
	function validateSecretName(name: string): { isValid: boolean; error: string } {
		if (!name) {
			return { isValid: false, error: 'Secret name cannot be empty' };
		}
		
		// Check for invalid characters (only allow letters, numbers, underscores, and hyphens)
		const validPattern = /^[a-zA-Z0-9_-]+$/;
		if (!validPattern.test(name)) {
			return { isValid: false, error: 'Name can only contain letters, numbers, underscores, and hyphens' };
		}
		
		return { isValid: true, error: '' };
	}
	
	// Reactive validation
	let nameValidation = $derived(validateSecretName(formData.name));
	let envNameValidation = $derived(validateSecretName(envFormData.name));
	
	// Environment functions
	async function loadEnvironments() {
		try {
			environments = await ListEnvironmentVariables() || [];
		} catch (error) {
			console.error('Failed to load environments:', error);
			environments = [];
		}
	}
	
	function openCreateEnvDialog() {
		envFormData = {
			name: '',
			description: '',
			value: ''
		};
		showCreateEnvDialog = true;
	}
	
	function openEditEnvDialog(env: EnvironmentVariableInfo) {
		editingEnv = env;
		envFormData = {
			name: env.name,
			description: env.description || '',
			value: env.value
		};
		showEditEnvDialog = true;
	}
	
	function openDeleteEnvDialog(env: EnvironmentVariableInfo) {
		deletingEnv = env;
		showDeleteEnvDialog = true;
	}
	
	async function handleCreateEnv() {
		try {
			const request: EnvironmentVariableCreateRequest = {
				name: envFormData.name,
				value: envFormData.value,
				description: envFormData.description
			};
			
			await CreateEnvironmentVariable(request);
			await loadEnvironments();
			showCreateEnvDialog = false;
		} catch (error) {
			console.error('Failed to create environment variable:', error);
			alert(`Failed to create environment variable: ${error instanceof Error ? error.message : String(error)}`);
		}
	}
	
	async function handleUpdateEnv() {
		if (!editingEnv) return;
		
		try {
			const request: EnvironmentVariableUpdateRequest = {
				id: editingEnv.id,
				name: envFormData.name,
				value: envFormData.value,
				description: envFormData.description
			};
			
			await UpdateEnvironmentVariable(request);
			await loadEnvironments();
			showEditEnvDialog = false;
			editingEnv = null;
		} catch (error) {
			console.error('Failed to update environment variable:', error);
			alert(`Failed to update environment variable: ${error instanceof Error ? error.message : String(error)}`);
		}
	}
	
	async function confirmDeleteEnv() {
		if (!deletingEnv) return;
		
		try {
			await DeleteEnvironmentVariable(deletingEnv.id);
			await loadEnvironments();
			showDeleteEnvDialog = false;
			deletingEnv = null;
		} catch (error) {
			console.error('Failed to delete environment variable:', error);
			alert(`Failed to delete environment variable: ${error instanceof Error ? error.message : String(error)}`);
		}
	}
</script>

<svelte:head>
	<title>Secrets & Environments - Teatime</title>
</svelte:head>

<AppBase title="Secrets & Environments" icon={Key}>
	{#snippet actions()}
		{#if activeTab === 'secrets'}
			<Button onclick={openCreateDialog} class="gap-2">
				<Plus class="h-4 w-4" />
				Add Secret
			</Button>
		{:else}
			<Button onclick={openCreateEnvDialog} class="gap-2">
				<Plus class="h-4 w-4" />
				Add Environment
			</Button>
		{/if}
	{/snippet}
	
	<Tabs.Root value={activeTab} onValueChange={(value) => activeTab = value as 'secrets' | 'environments'}>
		<Tabs.List class="grid w-full grid-cols-2">
			<Tabs.Trigger value="secrets" class="flex items-center gap-2">
				<Shield class="w-4 h-4" />
				Secrets
			</Tabs.Trigger>
			<Tabs.Trigger value="environments" class="flex items-center gap-2">
				<Settings class="w-4 h-4" />
				Environment Variables
			</Tabs.Trigger>
		</Tabs.List>
		
		<Tabs.Content value="secrets" class="space-y-4">
			<div class="space-y-2">
				<p class="text-muted-foreground">Securely store API keys, tokens, and other sensitive data in the system keychain.</p>
			</div>
			
			<div class="rounded-md border">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Name</Table.Head>
					<Table.Head>Last Used</Table.Head>
					<Table.Head>Created</Table.Head>
					<Table.Head class="w-[100px]">Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if loading}
					<Table.Row>
						<Table.Cell colspan={4} class="text-center py-8 text-muted-foreground">
							Loading secrets...
						</Table.Cell>
					</Table.Row>
				{:else if secrets.length === 0}
					<Table.Row>
						<Table.Cell colspan={4} class="text-center py-8 text-muted-foreground">
							No secrets found. Create your first secret to get started.
						</Table.Cell>
					</Table.Row>
				{:else}
					{#each secrets as secret}
						<Table.Row>
							<Table.Cell class="font-medium">
								<div class="space-y-1">
									<div class="flex items-center gap-2">
										<Shield class="w-4 h-4 text-green-500" />
										<span>{secret.name}</span>
									</div>
									{#if secret.description}
										<div class="text-xs text-muted-foreground ml-6">{secret.description}</div>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell class="text-muted-foreground">
								{formatDate(secret.last_used_at)}
							</Table.Cell>
							<Table.Cell class="text-muted-foreground">
								{formatDate(secret.created_at)}
							</Table.Cell>
							<Table.Cell>
								<div class="flex items-center gap-1">
									<Button
										variant="ghost"
										size="sm"
										onclick={() => openEditDialog(secret)}
										class="h-8 w-8 p-0"
									>
										<Edit class="h-4 w-4" />
									</Button>
									<button
										class="h-8 w-8 p-0 text-destructive hover:text-destructive hover:bg-gray-100 rounded flex items-center justify-center"
										onclick={() => openDeleteDialog(secret)}
									>
										<Trash2 class="h-4 w-4" />
									</button>
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</div>
		</Tabs.Content>
		
		<Tabs.Content value="environments" class="space-y-4">
			<div class="space-y-2">
				<p class="text-muted-foreground">Manage environment variables for non-sensitive configuration data.</p>
			</div>
			
			<div class="rounded-md border">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Name</Table.Head>
							<Table.Head>Value</Table.Head>
							<Table.Head class="w-[100px]">Actions</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#if environments.length === 0}
							<Table.Row>
								<Table.Cell colspan={3} class="text-center py-8 text-muted-foreground">
									No environment variables found. Create your first environment variable to get started.
								</Table.Cell>
							</Table.Row>
						{:else}
							{#each environments as env}
								<Table.Row>
									<Table.Cell class="font-medium">
										<div class="space-y-1">
											<div class="flex items-center gap-2">
												<Settings class="w-4 h-4 text-blue-500" />
												<span>{env.name}</span>
											</div>
											{#if env.description}
												<div class="text-xs text-muted-foreground ml-6">{env.description}</div>
											{/if}
										</div>
									</Table.Cell>
									<Table.Cell class="text-muted-foreground font-mono text-sm">
										{env.value}
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1">
											<Button
												variant="ghost"
												size="sm"
												onclick={() => openEditEnvDialog(env)}
												class="h-8 w-8 p-0"
											>
												<Edit class="h-4 w-4" />
											</Button>
											<button
												class="h-8 w-8 p-0 text-destructive hover:text-destructive hover:bg-gray-100 rounded flex items-center justify-center"
												onclick={() => openDeleteEnvDialog(env)}
											>
												<Trash2 class="h-4 w-4" />
											</button>
										</div>
									</Table.Cell>
								</Table.Row>
							{/each}
						{/if}
					</Table.Body>
				</Table.Root>
			</div>
		</Tabs.Content>
	</Tabs.Root>
</AppBase>
<!-- Create Dialog -->
<Dialog.Root open={showCreateDialog} onOpenChange={(open) => showCreateDialog = open}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Create New Secret</Dialog.Title>
			<Dialog.Description>
				Add a new secret to securely store sensitive information like API keys and tokens.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="name">Name</Label>
				<Input
					id="name"
					bind:value={formData.name}
					placeholder="e.g., openai_api_key"
					class={nameValidation.isValid ? '' : 'border-red-500'}
					required
				/>
				{#if !nameValidation.isValid && formData.name}
					<p class="text-sm text-red-500">{nameValidation.error}</p>
				{/if}
				<p class="text-xs text-muted-foreground">Only letters, numbers, underscores, and hyphens are allowed</p>
			</div>
			
			<div class="space-y-2">
				<Label for="description">Description (optional)</Label>
				<Textarea
					id="description"
					bind:value={formData.description}
					placeholder="Brief description of this secret"
					rows={2}
				/>
			</div>
			
			<div class="space-y-2">
				<Label for="value">Secret Value</Label>
				<div class="relative">
					<Input
						id="value"
						type={showValue ? 'text' : 'password'}
						bind:value={formData.value}
						placeholder="Enter your API key, token, or secret value"
						class="pr-10"
						required
					/>
					<Button
						type="button"
						variant="ghost"
						size="sm"
						class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
						onclick={toggleValueVisibility}
					>
						{#if showValue}
							<EyeOff class="h-4 w-4" />
						{:else}
							<Eye class="h-4 w-4" />
						{/if}
					</Button>
				</div>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showCreateDialog = false)}>Cancel</Button>
			<Button onclick={handleCreate} disabled={!formData.name || !formData.value || !nameValidation.isValid}>
				Create Secret
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Edit Dialog -->
<Dialog.Root open={showEditDialog} onOpenChange={(open) => showEditDialog = open}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Edit Secret</Dialog.Title>
			<Dialog.Description>
				Update the secret information and value.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="edit-name">Name</Label>
				<Input
					id="edit-name"
					bind:value={formData.name}
					class={nameValidation.isValid ? '' : 'border-red-500'}
					required
				/>
				{#if !nameValidation.isValid && formData.name}
					<p class="text-sm text-red-500">{nameValidation.error}</p>
				{/if}
				<p class="text-xs text-muted-foreground">Only letters, numbers, underscores, and hyphens are allowed</p>
			</div>
			
			<div class="space-y-2">
				<Label for="edit-description">Description</Label>
				<Textarea
					id="edit-description"
					bind:value={formData.description}
					rows={2}
				/>
			</div>
			
			<div class="space-y-2">
				<Label for="edit-value">Secret Value</Label>
				<div class="relative">
					<Input
						id="edit-value"
						type={showValue ? 'text' : 'password'}
						bind:value={formData.value}
						placeholder="Enter your API key, token, or secret value"
						class="pr-10"
						required
					/>
					<Button
						type="button"
						variant="ghost"
						size="sm"
						class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
						onclick={toggleValueVisibility}
					>
						{#if showValue}
							<EyeOff class="h-4 w-4" />
						{:else}
							<Eye class="h-4 w-4" />
						{/if}
					</Button>
				</div>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showEditDialog = false)}>Cancel</Button>
			<Button onclick={handleUpdate} disabled={!formData.name || !formData.value || !nameValidation.isValid}>
				Update
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Delete Confirmation Dialog -->
<AlertDialog.Root bind:open={showDeleteDialog} >
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete Secret</AlertDialog.Title>
			<AlertDialog.Description>
				Are you sure you want to delete "{deletingSecret?.name}"? This action cannot be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={confirmDelete} variant="destructive">Delete</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<!-- Create Environment Dialog -->
<Dialog.Root open={showCreateEnvDialog} onOpenChange={(open) => showCreateEnvDialog = open}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Create New Environment Variable</Dialog.Title>
			<Dialog.Description>
				Add a new environment variable for configuration data.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="env-name">Name</Label>
				<Input
					id="env-name"
					bind:value={envFormData.name}
					placeholder="e.g., api_base_url"
					class={envNameValidation.isValid ? '' : 'border-red-500'}
					required
				/>
				{#if !envNameValidation.isValid && envFormData.name}
					<p class="text-sm text-red-500">{envNameValidation.error}</p>
				{/if}
				<p class="text-xs text-muted-foreground">Only letters, numbers, underscores, and hyphens are allowed</p>
			</div>
			
			<div class="space-y-2">
				<Label for="env-description">Description (optional)</Label>
				<Textarea
					id="env-description"
					bind:value={envFormData.description}
					placeholder="Brief description of this environment variable"
					rows={2}
				/>
			</div>
			
			<div class="space-y-2">
				<Label for="env-value">Value</Label>
				<Input
					id="env-value"
					bind:value={envFormData.value}
					placeholder="Enter the environment variable value"
					required
				/>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showCreateEnvDialog = false)}>Cancel</Button>
			<Button onclick={handleCreateEnv} disabled={!envFormData.name || !envFormData.value || !envNameValidation.isValid}>
				Create Environment Variable
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Edit Environment Dialog -->
<Dialog.Root open={showEditEnvDialog} onOpenChange={(open) => showEditEnvDialog = open}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Edit Environment Variable</Dialog.Title>
			<Dialog.Description>
				Update the environment variable information and value.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="edit-env-name">Name</Label>
				<Input
					id="edit-env-name"
					bind:value={envFormData.name}
					class={envNameValidation.isValid ? '' : 'border-red-500'}
					required
				/>
				{#if !envNameValidation.isValid && envFormData.name}
					<p class="text-sm text-red-500">{envNameValidation.error}</p>
				{/if}
				<p class="text-xs text-muted-foreground">Only letters, numbers, underscores, and hyphens are allowed</p>
			</div>
			
			<div class="space-y-2">
				<Label for="edit-env-description">Description</Label>
				<Textarea
					id="edit-env-description"
					bind:value={envFormData.description}
					rows={2}
				/>
			</div>
			
			<div class="space-y-2">
				<Label for="edit-env-value">Value</Label>
				<Input
					id="edit-env-value"
					bind:value={envFormData.value}
					placeholder="Enter the environment variable value"
					required
				/>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showEditEnvDialog = false)}>Cancel</Button>
			<Button onclick={handleUpdateEnv} disabled={!envFormData.name || !envFormData.value || !envNameValidation.isValid}>
				Update
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Delete Environment Confirmation Dialog -->
<AlertDialog.Root open={showDeleteEnvDialog} onOpenChange={(open) => showDeleteEnvDialog = open}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete Environment Variable</AlertDialog.Title>
			<AlertDialog.Description>
				Are you sure you want to delete "{deletingEnv?.name}"? This action cannot be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={confirmDeleteEnv} variant="destructive">Delete</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<style>
	:global(.secrets-manager .lucide) {
		flex-shrink: 0;
	}
</style>