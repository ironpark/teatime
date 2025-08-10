<script lang="ts">
	import { icons } from '@lucide/svelte';
	import type { IconProps } from '@lucide/svelte';
	type IconName = keyof typeof icons;
	interface Props extends IconProps {	
		name: IconName;
	}
	let { name, ...props }: Props = $props();

	const Icon = icons[name as IconName];
</script>

<Icon {...props} />
