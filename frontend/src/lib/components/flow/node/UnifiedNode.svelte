<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import { nodeConfigs, nodeIcons } from './nodeConfigs';
	import type { NodeProperty } from '$bindings/node/types/models';
	import { HelpCircle } from 'lucide-svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';

	let { id, data, type, selected = false, isConnectable = true }: NodeProps = $props();

	// Get the specific node type from data or use the type directly
	const specificNodeType = $derived.by(() => data.nodeType || type);

	// Get node configuration
	const config = $derived.by(
		() => nodeConfigs[type as keyof typeof nodeConfigs] || nodeConfigs.util
	);

	// Get icon - use specific icon if available, otherwise use type icon
	const Icon = $derived.by(() => config?.icon || nodeIcons[specificNodeType]);

	// Get properties array if available (filter out optional properties without values)
	const properties = $derived.by(() => {
		if (Array.isArray(data.properties)) {
			const props = data.properties as NodeProperty[];
			// Filter out optional properties that have no value
			return props.filter(prop => !prop.optional || prop.value);
		}
		return [];
	});

	// Get outputs array if available
	const outputs = $derived.by(() => {
		if (Array.isArray(data.outputs)) {
			return data.outputs as NodeProperty[];
		}
		return [];
	});

	// Convert other data fields to array for rendering (excluding special fields)
	const dataEntries = $derived.by(() =>
		Object.entries(data).filter(
			([key]) => !['label', 'nodeType', 'description', 'properties', 'outputs', 'backendNodeId'].includes(key)
		)
	);

	// Determine handle type based on node type
	const handleType = $derived.by(() => {
		if (type === 'trigger') return 'source';
		if (type === 'action' && (specificNodeType === 'output' || specificNodeType === 'file-save'))
			return 'target';
		return 'both';
	});

	// Helper function to get property type display
	function getPropertyTypeDisplay(propType: number): string {
		const typeMap: Record<number, string> = {
			1: 'boolean',
			2: 'number',
			3: 'number',
			4: 'number',
			5: 'string',
			6: 'json',
			7: 'xml',
			8: 'date',
			9: 'text',
			10: 'string[]',
			11: 'number[]',
			12: 'boolean[]',
			13: 'json[]',
			14: 'xml[]'
		};
		return typeMap[propType] || 'unknown';
	}
</script>

{#if handleType === 'target' || handleType === 'both'}
	<Handle type="target" position={Position.Left} {isConnectable} />
{/if}

<div
	class={`transition-all duration-200 ${properties.length > 0 || outputs.length > 0 ? 'min-w-[240px]' : config?.minWidth || 'w-48'} 
		${selected ? 'ring-primary shadow-lg ring-2' : 'shadow-sm'} 
		${config?.color?.border || 'border-gray-200'} 
		rounded-lg border bg-white p-3`}
>
	<div class="mb-2 flex items-center justify-between">
		<div class="flex items-center gap-2">
			<div class={`rounded p-1 ${config?.color?.iconBg || 'bg-gray-100'}`}>
				<Icon class={`h-4 w-4 ${config?.color?.iconText || 'text-gray-600'}`} />
			</div>
			<span class="text-sm font-medium">{config?.label || 'Node'}</span>
		</div>
		{#if config?.badge}
			<span class={`rounded px-2 py-1 text-xs ${config?.badge?.className || ''}`}>
				{config?.badge?.text || ''}
			</span>
		{/if}
	</div>

	<div class="space-y-2">
		<h3 class="text-base font-semibold">{data.label || 'Untitled'}</h3>
		
		{#if data.description}
			<p class="text-xs text-muted-foreground">{data.description}</p>
		{/if}

		{#if properties.length > 0}
			<div class="space-y-2">
				<div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">Properties</div>
				<div class="space-y-1">
					{#each properties as prop}
						<div class="rounded bg-gray-50 p-1.5 text-xs">
							<div class="flex items-center justify-between gap-2">
								<div class="flex items-center gap-2 flex-1">
									{#if prop.description}
									<Tooltip.Root>
										<Tooltip.Trigger>
											<HelpCircle class="h-3 w-3 text-gray-400 hover:text-gray-600 cursor-help" />
										</Tooltip.Trigger>
										<Tooltip.Content>
											<p class="text-xs max-w-[200px]">{prop.description}</p>
										</Tooltip.Content>
									</Tooltip.Root>
									{/if}
									<span class="font-medium text-gray-700 min-w-[60px]">
										{prop.name || prop.key}:
									</span>
									<div class="flex-1">
										{#if prop.value}
											<span class="text-gray-900 break-words">
												{prop.value}
											</span>
										{:else if prop.options && prop.options.length > 0}
											<span class="text-gray-500 italic">
												[{prop.options.join(', ')}]
											</span>
										{:else}
											<span class="text-gray-400 italic">
												Not set
											</span>
										{/if}
									</div>
								</div>
								<span class="text-muted-foreground text-[10px] whitespace-nowrap">
									{getPropertyTypeDisplay(prop.type)}
								</span>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		{#if outputs.length > 0}
			<div class="space-y-2">
				<div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">Outputs</div>
				<div class="space-y-1">
					{#each outputs as output}
						<div class="flex items-center justify-between rounded bg-blue-50 px-1.5 py-1 text-xs">
							<span class="font-medium text-blue-900">{output.name || output.key}</span>
							<span class="text-blue-600 text-[10px]">
								{getPropertyTypeDisplay(output.type)}
							</span>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		{#if dataEntries.length > 0}
			<div class="space-y-1">
				{#each dataEntries as [key, value]}
					<div class="flex items-start gap-2 text-xs">
						<span class="text-muted-foreground min-w-[60px] capitalize">
							{key.replace(/_/g, ' ')}:
						</span>
						<span class="flex-1 break-words">
							{#if typeof value === 'boolean'}
								<span class={value ? 'text-green-600' : 'text-red-600'}>
									{value ? 'Yes' : 'No'}
								</span>
							{:else if typeof value === 'object' && value !== null}
								<code class="bg-muted rounded px-1 py-0.5 text-[10px]">
									{JSON.stringify(value, null, 2)}
								</code>
							{:else if value === null || value === undefined}
								<span class="text-muted-foreground italic">Not set</span>
							{:else}
								<span>{String(value)}</span>
							{/if}
						</span>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

{#if handleType === 'source' || handleType === 'both'}
	<Handle type="source" position={Position.Right} {isConnectable} />
{/if}
