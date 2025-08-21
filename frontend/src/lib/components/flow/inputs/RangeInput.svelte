<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		min?: number;
		max?: number;
		step?: number;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, min = 0, max = 100, step = 1, onUpdate }: Props = $props();
</script>

<div class="space-y-2">
	<div class="flex items-center justify-between">
		<Input
			id={`prop-${prop.key}`}
			type="range"
			bind:value={prop.value}
			{min}
			{max}
			{step}
			oninput={(e) => onUpdate(prop, parseFloat(e.currentTarget.value))}
			class="flex-1"
		/>
		<span class="ml-3 text-sm font-medium w-12 text-right">
			{prop.value ?? 0}
		</span>
	</div>
</div>