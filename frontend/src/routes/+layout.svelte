<script lang="ts">
	import '../app.css';
	import { SidebarProvider, SidebarInset } from '$ui/sidebar';
	import AppSidebar from '$components/AppSidebar.svelte';
	import { ModeWatcher } from 'mode-watcher';

	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

<ModeWatcher />
<SidebarProvider>
	<AppSidebar />
	<SidebarInset>
		{@render children?.()}
	</SidebarInset>
</SidebarProvider>
