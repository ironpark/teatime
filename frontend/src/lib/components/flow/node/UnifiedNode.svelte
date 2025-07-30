<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import { nodeConfigs, nodeIcons } from './nodeConfigs';

	let { id, data, type, selected = false, isConnectable = true }: NodeProps = $props();

	// Get the specific node type from data or use the type directly
	const specificNodeType = $derived.by(() => data.nodeType || type);

	// Get node configuration
	const config = $derived.by(
		() => nodeConfigs[type as keyof typeof nodeConfigs] || nodeConfigs.util
	);

	// Get icon - use specific icon if available, otherwise use type icon
	const Icon = $derived.by(() => config?.icon || nodeIcons[specificNodeType]);

	// Convert data object to array of key-value pairs for rendering
	const dataEntries = $derived.by(() =>
		Object.entries(data).filter(([key]) => key !== 'label' && key !== 'nodeType')
	);

	// Determine handle type based on node type
	const handleType = $derived.by(() => {
		if (type === 'trigger') return 'source';
		if (type === 'action' && (specificNodeType === 'output' || specificNodeType === 'file-save'))
			return 'target';
		return 'both';
	});
</script>

{#if handleType === 'target' || handleType === 'both'}
	<Handle type="target" position={Position.Left} {isConnectable} />
{/if}

<div
	class={`transition-all duration-200 ${config?.minWidth || 'w-48'} 
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
		{#if dataEntries.length > 0}
			<div class="space-y-1">
				{#each dataEntries as [key, value]}
					<div class="flex items-start gap-2 text-sm">
						<span class="text-muted-foreground min-w-[80px] capitalize"
							>{key.replace(/_/g, ' ')}:</span
						>
						<span class="flex-1 break-words">
							{#if typeof value === 'boolean'}
								<span class={value ? 'text-green-600' : 'text-red-600'}>
									{value ? 'Yes' : 'No'}
								</span>
							{:else if typeof value === 'object' && value !== null}
								<code class="bg-muted rounded px-1 py-0.5 text-xs">
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
