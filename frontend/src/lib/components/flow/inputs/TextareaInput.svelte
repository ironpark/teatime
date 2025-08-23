<script lang="ts">
	import { Textarea } from '$lib/components/ui/textarea';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		rows?: number;
		placeholder?: string;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, rows = 3, placeholder, onUpdate }: Props = $props();
</script>

<Textarea
	id={`prop-${prop.key}`}
	bind:value={prop.value}
	oninput={(e) => onUpdate(prop, e.currentTarget.value)}
	{rows}
	placeholder={placeholder || `Enter ${prop.name || prop.key}`}
	autocapitalize="none"
	spellcheck={false}
/>