<script>
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$ui/card';
	import { Button } from '$ui/button';
	import { Badge } from '$ui/badge';
	import { SidebarTrigger } from '$ui/sidebar';
	import { Separator } from '$ui/separator';
	import { 
		Activity, 
		Clock, 
		PlayCircle, 
		Settings, 
		ChefHat,
		TrendingUp
	} from 'lucide-svelte';

	const stats = [
		{
			title: 'Active Recipes',
			value: '12',
			change: '+2 from last week',
			icon: ChefHat
		},
		{
			title: 'Total Executions',
			value: '2,847',
			change: '+15% from last month',
			icon: Activity
		},
		{
			title: 'Success Rate',
			value: '98.2%',
			change: '+0.5% from last week',
			icon: TrendingUp
		},
		{
			title: 'Avg. Execution Time',
			value: '1.2s',
			change: '-0.3s from last week',
			icon: Clock
		}
	];

	const recentRecipes = [
		{
			name: 'Data Processing Recipe',
			status: 'running',
			lastRun: '2 minutes ago',
			executions: 156
		},
		{
			name: 'Image Enhancement Recipe',
			status: 'idle',
			lastRun: '1 hour ago',
			executions: 89
		},
		{
			name: 'API Integration Recipe',
			status: 'running',
			lastRun: '5 minutes ago',
			executions: 234
		},
		{
			name: 'Backup Automation Recipe',
			status: 'scheduled',
			lastRun: '12 hours ago',
			executions: 45
		}
	];
</script>

<div class="flex flex-col h-full">
	<header class="flex h-16 items-center gap-2 border-b px-4">
		<SidebarTrigger />
		<Separator orientation="vertical" class="mr-2 h-4" />
		<h1 class="text-lg font-semibold">Dashboard</h1>
	</header>

	<main class="flex-1 p-6 space-y-6">
		<div class="flex items-center justify-between">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Welcome to Teatime</h2>
				<p class="text-muted-foreground">
					Monitor and manage your AI workflow automation
				</p>
			</div>
			<Button>
				<PlayCircle class="mr-2 h-4 w-4" />
				New Recipe
			</Button>
		</div>

		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
			{#each stats as stat}
				<Card>
					<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
						<CardTitle class="text-sm font-medium">{stat.title}</CardTitle>
						<svelte:component this={stat.icon} class="h-4 w-4 text-muted-foreground" />
					</CardHeader>
					<CardContent>
						<div class="text-2xl font-bold">{stat.value}</div>
						<p class="text-xs text-muted-foreground">{stat.change}</p>
					</CardContent>
				</Card>
			{/each}
		</div>

		<div class="grid gap-4 md:grid-cols-2">
			<Card>
				<CardHeader>
					<CardTitle>Recent Recipes</CardTitle>
					<CardDescription>
						Your most recently active recipes
					</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="space-y-4">
						{#each recentRecipes as recipe}
							<div class="flex items-center justify-between">
								<div class="space-y-1">
									<p class="text-sm font-medium leading-none">{recipe.name}</p>
									<p class="text-sm text-muted-foreground">
										{recipe.executions} executions • {recipe.lastRun}
									</p>
								</div>
								<Badge 
									variant={recipe.status === 'running' ? 'default' : 
											recipe.status === 'scheduled' ? 'secondary' : 'outline'}
								>
									{recipe.status}
								</Badge>
							</div>
						{/each}
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Quick Actions</CardTitle>
					<CardDescription>
						Common tasks and shortcuts
					</CardDescription>
				</CardHeader>
				<CardContent class="space-y-3">
					<Button variant="outline" class="w-full justify-start">
						<ChefHat class="mr-2 h-4 w-4" />
						Create New Recipe
					</Button>
					<Button variant="outline" class="w-full justify-start">
						<PlayCircle class="mr-2 h-4 w-4" />
						Execute Recipe
					</Button>
					<Button variant="outline" class="w-full justify-start">
						<Settings class="mr-2 h-4 w-4" />
						System Settings
					</Button>
				</CardContent>
			</Card>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>System Status</CardTitle>
				<CardDescription>
					Current system health and performance
				</CardDescription>
			</CardHeader>
			<CardContent>
				<div class="grid gap-4 md:grid-cols-3">
					<div class="space-y-2">
						<div class="flex items-center justify-between">
							<span class="text-sm">CPU Usage</span>
							<span class="text-sm font-medium">34%</span>
						</div>
						<div class="w-full bg-secondary rounded-full h-2">
							<div class="bg-primary h-2 rounded-full" style="width: 34%"></div>
						</div>
					</div>
					<div class="space-y-2">
						<div class="flex items-center justify-between">
							<span class="text-sm">Memory Usage</span>
							<span class="text-sm font-medium">67%</span>
						</div>
						<div class="w-full bg-secondary rounded-full h-2">
							<div class="bg-primary h-2 rounded-full" style="width: 67%"></div>
						</div>
					</div>
					<div class="space-y-2">
						<div class="flex items-center justify-between">
							<span class="text-sm">Active Connections</span>
							<span class="text-sm font-medium">8/20</span>
						</div>
						<div class="w-full bg-secondary rounded-full h-2">
							<div class="bg-primary h-2 rounded-full" style="width: 40%"></div>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</main>
</div>
