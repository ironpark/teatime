<script lang="ts">
	import { Switch } from '$lib/components/ui/switch';
	import { Label } from '$lib/components/ui/label';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, onUpdate }: Props = $props();
	
	const isChecked = $derived(prop.value === true || prop.value === 'true');
</script>

<div class="flex items-center space-x-2">
	<Switch
		id={`prop-${prop.key}`}
		checked={isChecked}
		onCheckedChange={(checked) => onUpdate(prop, checked.toString())}
	/>
	<Label for={`prop-${prop.key}`} class="text-sm">
		{isChecked ? 'Yes' : 'No'}
	</Label>
</div>