<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		type?: string;
		min?: number;
		max?: number;
		step?: number;
		placeholder?: string;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, type = 'text', min, max, step, placeholder, onUpdate }: Props = $props();
</script>

<Input
	id={`prop-${prop.key}`}
	{type}
	bind:value={prop.value}
	{min}
	{max}
	{step}
	oninput={(e) => onUpdate(prop, type === 'number' ? parseFloat(e.currentTarget.value) : e.currentTarget.value)}
	placeholder={placeholder || `Enter ${prop.name || prop.key}`}
/>