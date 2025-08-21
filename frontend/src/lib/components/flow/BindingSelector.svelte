<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import { Badge } from '$lib/components/ui/badge';
	import { Link, Unlink, ArrowRight } from 'lucide-svelte';
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
		onToggleBinding: (prop: NodeProperty) => void;
	}

	let { prop, currentNodeId, nodes, edges, onToggleBinding }: Props = $props();

	const isCurrentlyBound = $derived(isBinding(prop.value));
	const bindingOptions = $derived(getBindingOptions(currentNodeId, prop.type, nodes, edges));

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

<!-- Inline bind button -->
{#if bindingOptions.length > 0 || isCurrentlyBound}
	<Button
		variant={isCurrentlyBound ? "default" : "outline"}
		size="sm"
		onclick={() => onToggleBinding(prop)}
		class="h-6 px-2 text-xs"
	>
		{#if isCurrentlyBound}
			<Unlink class="h-3 w-3 mr-1" />
			Unbind
		{:else}
			<Link class="h-3 w-3 mr-1" />
			Bind
		{/if}
	</Button>
{/if}