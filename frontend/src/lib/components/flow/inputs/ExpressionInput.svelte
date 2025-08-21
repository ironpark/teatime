<script lang="ts">
	import { Textarea } from '$lib/components/ui/textarea';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, onUpdate }: Props = $props();
</script>

<div class="space-y-2">
	<Textarea
		id={`prop-${prop.key}`}
		bind:value={prop.value}
		oninput={(e) => onUpdate(prop, e.currentTarget.value)}
		rows={3}
		placeholder="{'{'}expression{'}'} or @[binding]"
		class="font-mono text-sm"
		autocapitalize="none"
		autocorrect="off"
		spellcheck="false"
	/>
	<div class="text-xs text-muted-foreground">
		Use {'{'}expression{'}'} for calculations or @[nodeName.output.field] for bindings
	</div>
</div>