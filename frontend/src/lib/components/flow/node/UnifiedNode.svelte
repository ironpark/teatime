<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import { nodeConfigs } from './nodeConfigs';
	import type { NodeProperty } from '$bindings/internal/node/models';
	import { HelpCircle, Package, Edit } from 'lucide-svelte';
	import LucideIcon from '$lib/components/LucideIcon.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Button } from '$lib/components/ui/button';
	import { getContext } from 'svelte';

	let { data, type: nodeType, selected = false, isConnectable = true, id }: NodeProps = $props();

	const typeName = $derived.by(() => {
		if (nodeType === 'trigger') return 'Trigger';
		if (nodeType === 'action') return 'Action';
		if (nodeType === 'branch') return 'Branch';
		if (nodeType === 'util') return 'Utility';
		return 'Unknown';
	});
	// Get handlers from context
	const handlers = getContext<{
		onEdit: (nodeId: string) => void;
		onDoubleClick: (nodeId: string) => void;
	}>('nodeHandlers');

	// Get node configuration
	const config = nodeConfigs[nodeType as keyof typeof nodeConfigs] || nodeConfigs.util
	// Get icon from data
	const iconString = $derived.by(() => data.nodeIcon || data.icon || null);
	console.log(data);
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

	// Get output handles array if available
	const outputHandles = $derived.by(() => {
		if (Array.isArray(data.outputHandles)) {
			return data.outputHandles as Array<{id: string, label?: string, description?: string}>;
		}
		// Return default handle if no output handles are specified
		return [{id: 'default', label: 'Output'}];
	});

	// Determine handle type based on node type
	const handleType = $derived.by(() => {
		if (nodeType === 'trigger') return 'source';
		if (nodeType === 'action')
			return 'target';
		return 'both';
	});

	// Helper function to get property type display
	function getPropertyTypeDisplay(propType: number): string {
		const typeMap: Record<number, string> = {
			0: 'invalid',
			1: 'boolean',
			2: 'int64',
			3: 'uint64', 
			4: 'float64',
			5: 'string',
			6: 'json',
			7: 'xml',
			8: 'date',
			9: 'string[]',
			10: 'number[]',
			11: 'boolean[]'
		};
		return typeMap[propType] || 'unknown';
	}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class={`transition-all duration-200 ${properties.length > 0 || outputs.length > 0 ? 'min-w-[240px]' : config?.minWidth || 'w-48'} 
		${selected ? 'ring-primary shadow-lg ring-2' : 'shadow-sm'} 
		${config?.color?.border || 'border-gray-200'} 
		rounded-lg border bg-white p-3`}
	ondblclick={() => handlers?.onDoubleClick?.(id)}
>
	<div class="mb-2 flex items-center justify-between">
		<div class="flex items-center gap-2">
			<div class={`rounded p-1 ${config?.color?.iconBg || 'bg-gray-100'}`}>
				{#if iconString}
					<LucideIcon name={iconString as any} class={`h-4 w-4 ${config?.color?.iconText || 'text-gray-600'}`} />
				{:else}
					<Package class={`h-4 w-4 ${config?.color?.iconText || 'text-gray-600'}`} />
				{/if}
			</div>
			<span class="text-sm font-medium">{data.name} {typeName}</span>
		</div>
		<div class="flex items-center gap-1">
			{#if config?.badge}
				<span class={`rounded px-2 py-1 text-xs ${config?.badge?.className || ''}`}>
				 {config?.badge?.text || ''}
				</span>
			{/if}
			<Button 
				variant="ghost" 
				size="icon" 
				class="h-6 w-6 hover:bg-gray-100"
				onclick={(e) => {
					e.stopPropagation();
					handlers?.onEdit?.(id);
				}}
			>
				<Edit class="h-3 w-3" />
			</Button>
		</div>
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



	</div>
</div>

{#if handleType === 'target' || handleType === 'both'}
	<Handle type="target" position={Position.Left} {isConnectable} />
{/if}

<!-- Multiple output handles - stack them vertically -->

{#each outputHandles as handle, index}
	{#if handle.label == 'Else'}
		<Handle 
			id={handle.id}	
			type="source" 
			position={Position.Bottom}>
			<div class="flex flex-col items-center justify-center absolute top-3 left-0">
				<span class="text-xs text-zinc-800 bg-gray-100 rounded-md px-1 py-0.5">{handle.label || handle.id}</span>
			</div>
		</Handle>
	{:else}
		<Handle 
			id={handle.id}
			type="source" 
			position={Position.Right} 
			style="top: {40 + (index * 55)}px; width: 12px; height: 12px; background-color: rgba(255, 255, 255, 0.8); border-radius: 3px; border: 1px solid black; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06);">
			<div class="flex flex-col items-center justify-center absolute top-3 left-0">
				<span class="text-xs text-zinc-800 bg-gray-100 rounded-md px-1 py-0.5">{handle.label || handle.id}</span>
			</div>
		</Handle>
	{/if}
{/each}