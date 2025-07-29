<script>
	import {
		Sidebar,
		SidebarContent,
		SidebarFooter,
		SidebarGroup,
		SidebarGroupContent,
		SidebarGroupLabel,
		SidebarHeader,
		SidebarMenu,
		SidebarMenuButton,
		SidebarMenuItem,
		SidebarMenuSub,
		SidebarMenuSubButton,
		SidebarMenuSubItem
	} from '$ui/sidebar';
	import {
		Home,
		Settings,
		ChefHat,
		Play,
		FileText,
		Palette
	} from 'lucide-svelte';

	const data = {
		navMain: [
			{
				title: 'Dashboard',
				url: '/',
				icon: Home
			},
			{
				title: 'Recipes',
				url: '/recipes',
				icon: ChefHat
			},
			{
				title: 'Execution',
				url: '/execution',
				icon: Play
			}
		],
		navSecondary: [
			{
				title: 'Documentation',
				url: '/docs',
				icon: FileText
			},
			{
				title: 'Themes',
				url: '/themes',
				icon: Palette
			},
			{
				title: 'Settings',
				url: '/settings',
				icon: Settings
			}
		]
	};
</script>

<Sidebar>
	<SidebarHeader>
		<div class="flex items-center gap-2 p-2">
			<div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
				<span class="text-primary-foreground font-bold text-sm">T</span>
			</div>
			<div class="flex flex-col">
				<span class="text-sm font-medium">Teatime</span>
				<span class="text-xs text-muted-foreground">Workflow Automation</span>
			</div>
		</div>
	</SidebarHeader>

	<SidebarContent>
		<SidebarGroup>
			<SidebarGroupLabel>Main</SidebarGroupLabel>
			<SidebarGroupContent>
				<SidebarMenu>
					{#each data.navMain as item}
						<SidebarMenuItem>
							<SidebarMenuButton href={item.url}>
								<svelte:component this={item.icon} size={16} />
								<span>{item.title}</span>
							</SidebarMenuButton>
						</SidebarMenuItem>
					{/each}
				</SidebarMenu>
			</SidebarGroupContent>
		</SidebarGroup>

		<SidebarGroup>
			<SidebarGroupLabel>Tools</SidebarGroupLabel>
			<SidebarGroupContent>
				<SidebarMenu>
					{#each data.navSecondary as item}
						<SidebarMenuItem>
							<SidebarMenuButton href={item.url}>
								<svelte:component this={item.icon} size={16} />
								<span>{item.title}</span>
							</SidebarMenuButton>
						</SidebarMenuItem>
					{/each}
				</SidebarMenu>
			</SidebarGroupContent>
		</SidebarGroup>
	</SidebarContent>

	<SidebarFooter>
		<div class="p-2 text-xs text-muted-foreground">
			<div>Version 0.1.0</div>
			<div class="text-xs opacity-60">Early Development</div>
		</div>
	</SidebarFooter>
</Sidebar>