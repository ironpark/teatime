<script lang="ts">
	import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		options: { value: string; label: string }[];
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, options, onUpdate }: Props = $props();
</script>

<Select
	type="multiple"
	bind:value={prop.value}
	onValueChange={(value) => {
		onUpdate(prop, value);
	}}
>
	<SelectTrigger id={`prop-${prop.key}`}>
		{prop.value || 'Select option'}
	</SelectTrigger>
	<SelectContent>
		<SelectGroup>
		{#each options as option}
			<SelectItem value={option.value}>
				{option.label}
			</SelectItem>
		{/each}
		</SelectGroup>
	</SelectContent>
</Select>