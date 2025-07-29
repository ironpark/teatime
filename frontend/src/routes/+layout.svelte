<script>
	import '../app.css';
	import { SidebarProvider, SidebarInset } from '$ui/sidebar';
	import AppSidebar from '$components/AppSidebar.svelte';
</script>

<SidebarProvider>
	<AppSidebar />
	<SidebarInset>
		<slot />
	</SidebarInset>
</SidebarProvider>
