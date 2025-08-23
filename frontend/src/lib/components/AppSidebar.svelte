<script lang="ts">
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
	import { Home, ChefHat, Play, Cog, Database } from 'lucide-svelte';
	import { page } from '$app/state';

	function isActive(url: string) {
		if (url === '/') {
			return page.url.pathname === '/';
		}
		return page.url.pathname.startsWith(url);
	}

	const data = {
		navMain: [
			{
				title: 'Recipes',
				url: '/',
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
				title: 'Application',
				url: '/settings',
				icon: Cog
			},
			{
				title: 'Configuration',
				url: '/secrets',
				icon: Database
			}
		]
	};
</script>

<Sidebar>
	<!-- mac os 에서는 35px 이상 패딩 필요 -->
	<SidebarHeader class="pt-[35px]">
		<a href="/" class="hover:bg-accent flex items-center gap-2 rounded-lg p-2 transition-colors">
			<div class="bg-primary flex h-8 w-8 items-center justify-center rounded-lg">
				<span class="text-primary-foreground text-sm font-bold">T</span>
			</div>
			<div class="flex flex-col">
				<span class="text-sm font-medium">Teatime</span>
				<span class="text-muted-foreground text-xs">Workflow Automation</span>
			</div>
		</a>
	</SidebarHeader>

	<SidebarContent>
		<SidebarGroup>
			<SidebarGroupLabel>Main</SidebarGroupLabel>
			<SidebarGroupContent>
				<SidebarMenu>
					{#each data.navMain as item}
						<SidebarMenuItem>
							<SidebarMenuButton isActive={isActive(item.url)}>
								<a
									href={item.url}
									class="flex w-full items-center gap-2"
									data-sveltekit-preload-data
								>
									<item.icon size={16} />
									<span>{item.title}</span>
								</a>
							</SidebarMenuButton>
						</SidebarMenuItem>
					{/each}
				</SidebarMenu>
			</SidebarGroupContent>
		</SidebarGroup>

		<SidebarGroup>
			<SidebarGroupLabel>Settings</SidebarGroupLabel>
			<SidebarGroupContent>
				<SidebarMenu>
					{#each data.navSecondary as item}
						<SidebarMenuItem>
							<SidebarMenuButton isActive={isActive(item.url)}>
								<a
									href={item.url}
									class="flex w-full items-center gap-2"
									data-sveltekit-preload-data
								>
									<item.icon size={16} />
									<span>{item.title}</span>
								</a>
							</SidebarMenuButton>
						</SidebarMenuItem>
					{/each}
				</SidebarMenu>
			</SidebarGroupContent>
		</SidebarGroup>
	</SidebarContent>

	<SidebarFooter>
		<div class="text-muted-foreground p-2 text-xs">
			<div>Version 0.1.0</div>
			<div class="text-xs opacity-60">Early Development</div>
		</div>
	</SidebarFooter>
</Sidebar>
