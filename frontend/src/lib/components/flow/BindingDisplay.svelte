<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Link, ArrowRight } from 'lucide-svelte';
	import type { Node, Edge } from '@xyflow/svelte';
	import type { NodeProperty } from '$bindings/internal/node';
	import { 
		getBindingOptions, 
		getBindingDisplay, 
		isBinding, 
		type BindingOption 
	} from './utils/binding';

	interface Props {
		prop: NodeProperty;
		currentNodeId: string;
		nodes: Node[];
		edges: Edge[];
		onUpdate: (prop: NodeProperty, value: any) => void;
		showBindingSelector: boolean;
		onSelectBinding: (option: BindingOption) => void;
		onCancelBinding: () => void;
	}

	let { 
		prop, 
		currentNodeId, 
		nodes, 
		edges, 
		onUpdate, 
		showBindingSelector,
		onSelectBinding,
		onCancelBinding 
	}: Props = $props();

	const isCurrentlyBound = $derived(isBinding(prop.value));
	const bindingOptions = $derived(getBindingOptions(currentNodeId, prop.type, nodes, edges));
	const currentBinding = $derived(
		isCurrentlyBound ? getBindingDisplay(prop.value as string, nodes) : null
	);

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

<div class="space-y-2">
	{#if isCurrentlyBound && currentBinding}
		<div class="flex items-center gap-2 p-2 bg-blue-50 border border-blue-200 rounded-md">
			<Link class="h-3 w-3 text-blue-600" />
			<span class="text-xs font-medium text-blue-800">
				{currentBinding.nodeLabel}
			</span>
			<ArrowRight class="h-3 w-3 text-blue-600" />
			<span class="text-xs text-blue-700">
				{currentBinding.propertyName}
			</span>
			<Badge variant="outline" class="text-[10px] px-1 py-0 bg-blue-100 text-blue-700 border-blue-300">
				{currentBinding.type}
			</Badge>
		</div>
	{/if}

	{#if showBindingSelector && bindingOptions.length > 0}
		<div class="space-y-2 p-3 bg-muted/50 rounded-md border">
			<div class="text-xs font-medium text-muted-foreground mb-2">
				Select binding source:
			</div>
			
			<div class="space-y-1 max-h-48 overflow-y-auto">
				{#each bindingOptions as option (option.bindingValue)}
					<Button
						variant="ghost"
						size="sm"
						onclick={() => onSelectBinding(option)}
						class="w-full justify-start h-auto p-2 text-xs"
					>
						<div class="flex items-center gap-2 w-full">
							<div class="flex items-center gap-1 flex-1">
								<span class="font-medium">{option.nodeLabel}</span>
								<ArrowRight class="h-3 w-3 text-muted-foreground" />
								<span>{option.propertyName}</span>
							</div>
							<div class="flex items-center gap-1">
								<Badge variant="outline" class="text-[9px] px-1 py-0">
									{option.type}
								</Badge>
								<Badge variant="secondary" class="text-[9px] px-1 py-0">
									{getPropertyTypeDisplay(option.propertyType)}
								</Badge>
							</div>
						</div>
					</Button>
				{/each}
			</div>
			
			<div class="flex justify-end pt-1 border-t">
				<Button
					variant="outline"
					size="sm"
					onclick={onCancelBinding}
					class="h-6 px-2 text-xs"
				>
					Cancel
				</Button>
			</div>
		</div>
	{/if}
</div>