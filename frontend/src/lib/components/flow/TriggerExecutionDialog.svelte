<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import { Play, Settings } from 'lucide-svelte';
	import type { NodeProperty } from '$bindings/internal/node';
	import { ExecuteTriggerNode } from '$bindings/services/recipesservice';

	interface CommandArg {
		name: string;
		required: boolean;
		list: boolean;
		description: string;
	}

	interface Props {
		open: boolean;
		onOpenChange: (open: boolean) => void;
		nodeLabel: string;
		args: CommandArg[];
		properties: NodeProperty[];
		recipeID: string;
		nodeID: string;
	}

	let { open, onOpenChange, nodeLabel, args, properties, recipeID, nodeID }: Props = $props();

	// Store argument values
	let argValues = $state<Record<string, string | string[]>>({});
	
	// Execution state
	let isExecuting = $state(false);
	let errorMessage = $state<string | null>(null);

	// Initialize argument values when dialog opens
	$effect(() => {
		if (open) {
			const newValues: Record<string, string | string[]> = {};
			for (const arg of args) {
				if (arg.list) {
					newValues[arg.name] = [];
				} else {
					newValues[arg.name] = '';
				}
			}
			argValues = newValues;
			errorMessage = null; // Clear error when dialog opens
		}
	});

	async function handleExecute() {
		if (isExecuting) return;
		
		isExecuting = true;
		errorMessage = null; // Clear previous error
		
		try {
			// Convert argValues to the format expected by the backend
			const formattedArgs: Record<string, any> = {};
			for (const [key, value] of Object.entries(argValues)) {
				if (value !== '' && !(Array.isArray(value) && value.length === 0)) {
					formattedArgs[key] = value;
				}
			}
			
			// Convert properties to the format expected by the backend
			const formattedProperties: Record<string, any> = {};
			for (const prop of properties) {
				if (prop.value !== undefined && prop.value !== null) {
					formattedProperties[prop.key] = prop.value;
				}
			}
			
			console.log('Executing trigger with properties:', formattedProperties, 'args:', formattedArgs);
			await ExecuteTriggerNode(recipeID, nodeID, formattedProperties, formattedArgs);
			console.log('Trigger executed successfully');
			onOpenChange(false);
		} catch (error) {
			console.error('Failed to execute trigger:', error);
			errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
		} finally {
			isExecuting = false;
		}
	}

	function handleCancel() {
		onOpenChange(false);
	}

	function updateArgValue(argName: string, value: string) {
		argValues = { ...argValues, [argName]: value };
	}

	function updateListArgValue(argName: string, value: string) {
		// Split by comma and trim whitespace
		const values = value.split(',').map(v => v.trim()).filter(v => v.length > 0);
		argValues = { ...argValues, [argName]: values };
	}

	function getArgDisplayValue(argName: string): string {
		const value = argValues[argName];
		if (Array.isArray(value)) {
			return value.join(', ');
		}
		return value || '';
	}

	// Check if all required arguments are filled and not currently executing
	const isValidForm = $derived.by(() => {
		if (isExecuting) return false;
		
		for (const arg of args) {
			if (arg.required) {
				const value = argValues[arg.name];
				if (!value || (Array.isArray(value) && value.length === 0)) {
					return false;
				}
			}
		}
		return true;
	});

	// Get property display value
	function getPropertyDisplayValue(prop: NodeProperty): string {
		if (!prop.value) return 'Not set';
		
		// Special formatting for ArgList
		if (prop.type === 6 && Array.isArray(prop.value)) {
			const argNames: string[] = [];
			for (const arg of prop.value) {
				if (typeof arg === 'object' && arg !== null && arg.name) {
					argNames.push(arg.name.startsWith('--') ? arg.name : `--${arg.name}`);
				}
			}
			return argNames.join(', ') || 'No arguments';
		}
		
		if (typeof prop.value === 'object') {
			return JSON.stringify(prop.value);
		}
		
		return String(prop.value);
	}

	// Filter out args property from displayed properties
	const displayProperties: NodeProperty[] = $derived.by(() => {
		return properties.filter((prop: NodeProperty) => 
			prop.key !== 'args' && 
			(!prop.optional || prop.value)
		);
	});
</script>

<Dialog.Root {open} onOpenChange={onOpenChange}>
	<Dialog.Content class="max-w-md">
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-2">
				<Play class="h-4 w-4" />
				Run Trigger: {nodeLabel}
			</Dialog.Title>
			<Dialog.Description>
				Enter the required arguments to execute this trigger.
			</Dialog.Description>
		</Dialog.Header>

		<div class="space-y-6 py-4">
			<!-- Properties Section -->
			{#if displayProperties.length > 0}
				<div class="space-y-3">
					<div class="flex items-center gap-2">
						<Settings class="h-4 w-4" />
						<h3 class="text-sm font-medium">Current Properties</h3>
					</div>
					<div class="space-y-2 max-h-32 overflow-y-auto">
						{#each displayProperties as prop (prop.key)}
							<div class="flex items-center justify-between p-2 bg-muted/50 rounded-md text-xs">
								<div class="flex items-center gap-2">
									<span class="font-medium text-muted-foreground">
										{prop.name || prop.key}:
									</span>
									<span class="text-gray-900 break-words max-w-[200px] truncate">
										{getPropertyDisplayValue(prop)}
									</span>
								</div>
								<Badge variant="outline" class="text-[9px] px-1 py-0">
									{prop.type === 1 ? 'bool' : prop.type === 5 ? 'string' : prop.type === 6 ? 'json' : 'other'}
								</Badge>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			<!-- Arguments Section -->
			{#if args.length === 0}
				<div class="text-center text-muted-foreground py-8">
					<p>This trigger has no configurable arguments.</p>
					<p class="text-sm mt-1">Click "Execute" to run the trigger.</p>
				</div>
			{:else}
				<div class="space-y-3">
					<div class="flex items-center gap-2">
						<Play class="h-4 w-4" />
						<h3 class="text-sm font-medium">Runtime Arguments</h3>
					</div>
					<div class="space-y-4">
						{#each args as arg (arg.name)}
					<div class="space-y-2">
						<div class="flex items-center justify-between">
							<Label for={`arg-${arg.name}`} class="text-sm font-medium">
								{arg.name}
								{#if arg.required}
									<span class="text-red-500 ml-1">*</span>
								{/if}
							</Label>
							<div class="flex items-center gap-1">
								{#if arg.required}
									<Badge variant="destructive" class="text-[10px] px-1 py-0">
										Required
									</Badge>
								{:else}
									<Badge variant="secondary" class="text-[10px] px-1 py-0">
										Optional
									</Badge>
								{/if}
								{#if arg.list}
									<Badge variant="outline" class="text-[10px] px-1 py-0">
										List
									</Badge>
								{/if}
							</div>
						</div>
						
						{#if arg.description}
							<p class="text-xs text-muted-foreground">{arg.description}</p>
						{/if}

						<Input
							id={`arg-${arg.name}`}
							value={getArgDisplayValue(arg.name)}
							oninput={(e) => {
								if (arg.list) {
									updateListArgValue(arg.name, e.currentTarget.value);
								} else {
									updateArgValue(arg.name, e.currentTarget.value);
								}
							}}
							placeholder={arg.list ? 'value1, value2, value3' : arg.name.startsWith('--') ? arg.name : `--${arg.name}`}
							autocapitalize="none"
							autocorrect="off"
							spellcheck="false"
							class="text-sm"
						/>

						{#if arg.list}
							<p class="text-xs text-muted-foreground">
								Separate multiple values with commas
							</p>
						{/if}
					</div>
				{/each}
				</div>
			</div>
		{/if}
		</div>

		<!-- Error Message -->
		{#if errorMessage}
			<div class="mx-6 mb-4 p-3 bg-destructive/10 border border-destructive/20 rounded-md">
				<div class="flex items-start gap-2">
					<div class="w-4 h-4 text-destructive mt-0.5">⚠️</div>
					<div class="flex-1">
						<p class="text-sm font-medium text-destructive mb-1">Execution Failed</p>
						<p class="text-xs text-destructive/80">{errorMessage}</p>
					</div>
				</div>
			</div>
		{/if}

		<Dialog.Footer class="flex items-center gap-2">
			<Button 
				variant="outline" 
				onclick={handleCancel}
				disabled={isExecuting}
			>
				Cancel
			</Button>
			<Button 
				onclick={handleExecute}
				disabled={!isValidForm || isExecuting}
				class="flex items-center gap-2"
			>
				{#if isExecuting}
					<div class="w-3 h-3 animate-spin rounded-full border border-current border-t-transparent"></div>
					Executing...
				{:else}
					<Play class="h-3 w-3" />
					Execute
				{/if}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>