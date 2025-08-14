<script lang="ts">
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import {
		Edit,
		Trash2,
		Play,
		Copy,
		Zap,
		Database,
		Bot,
		FileText,
		Bell,
		Globe,
		Settings,
		Clock,
		Activity,
		GitBranch
	} from 'lucide-svelte';
	import type { RecipeInfo } from '$bindings/services';

	// 워크플로우 자동화를 위한 더미 데이터 생성
	type WorkflowData = {
		nodeCount: number;
		lastExecuted?: string;
		status: 'active' | 'inactive' | 'error';
		category: 'data' | 'ai' | 'file' | 'notification' | 'web' | 'system';
		triggerType: 'manual' | 'webhook' | 'schedule' | 'file-watch';
		executionCount: number;
	};

	interface Props {
		recipe: RecipeInfo;
		onEdit?: (recipe: RecipeInfo) => void;
		onDelete?: (recipe: RecipeInfo) => void;
		onDuplicate?: (recipe: RecipeInfo) => void;
		onExecute?: (recipe: RecipeInfo) => void;
	}

	let {
		recipe,
		onEdit = () => {},
		onDelete = () => {},
		onDuplicate = () => {},
		onExecute = () => {}
	}: Props = $props();

	// // 레시피 ID 기반 더미 워크플로우 데이터 생성
	// function generateWorkflowData(recipeId: string): WorkflowData {
	//   const hash = recipeId.split('').reduce((a, b) => {
	//     a = ((a << 5) - a) + b.charCodeAt(0);
	//     return a & a;
	//   }, 0);

	//   const categories: WorkflowData['category'][] = ['data', 'ai', 'file', 'notification', 'web', 'system'];
	//   const triggers: WorkflowData['triggerType'][] = ['manual', 'webhook', 'schedule', 'file-watch'];
	//   const statuses: WorkflowData['status'][] = ['active', 'inactive', 'error'];

	//   return {
	//     nodeCount: Math.abs(hash % 15) + 3, // 3-17 nodes
	//     lastExecuted: Math.abs(hash % 5) === 0 ? undefined : `${Math.abs(hash % 24) + 1} hours ago`,
	//     status: statuses[Math.abs(hash % statuses.length)],
	//     category: categories[Math.abs(hash % categories.length)],
	//     triggerType: triggers[Math.abs(hash % triggers.length)],
	//     executionCount: Math.abs(hash % 500) + 10
	//   };
	// }

	// let workflowData = $derived(generateWorkflowData(recipe.ID));

	const categoryIcons = {
		data: Database,
		ai: Bot,
		file: FileText,
		notification: Bell,
		web: Globe,
		system: Settings
	};

	const statusColors = {
		active: 'bg-green-100 text-green-800 border-green-200',
		inactive: 'bg-gray-100 text-gray-800 border-gray-200',
		error: 'bg-red-100 text-red-800 border-red-200'
	};

	const triggerIcons = {
		manual: Play,
		webhook: Zap,
		schedule: Clock,
		'file-watch': Activity
	};

	function handleEdit() {
		onEdit(recipe);
	}

	function handleDelete() {
		onDelete(recipe);
	}

	function handleDuplicate() {
		onDuplicate(recipe);
	}

	function handleExecute() {
		onExecute(recipe);
	}
	const obj: Record<string, string> = {
		webhook: 'webhook',
		schedule: 'schedule',
		command: 'cmd'
	};

	const triggers = $derived(
		recipe.Tags.filter((tag) => tag.includes('trigger:'))
			.map((tag) => obj[tag.split(':')[1].toLowerCase()])
			.filter((tag) => tag !== undefined)
	);
</script>

<Card class="group gap-1 transition-all duration-200 hover:shadow-lg">
	<CardHeader class="pb-3">
		<div class="flex items-start justify-between">
			<div class="flex-1">
				<div class="mb-2 flex items-center gap-2">
					<!-- <CategoryIcon class="w-4 h-4 text-muted-foreground" />
          <Badge variant="outline" class={`text-xs ${statusColors[workflowData.status]}`}>
            {workflowData.status}
          </Badge> -->
				</div>
				<CardTitle class="mb-1 text-lg">{recipe.Name}</CardTitle>
				<CardDescription class="line-clamp-2">
					{recipe.Description}
				</CardDescription>
			</div>

			<div class="flex gap-1 opacity-0 transition-opacity group-hover:opacity-100">
				<Button variant="ghost" size="icon" class="h-8 w-8" onclick={handleEdit}>
					<Edit class="h-4 w-4" />
				</Button>

				<Button variant="ghost" size="icon" class="h-8 w-8" onclick={handleDuplicate}>
					<Copy class="h-4 w-4" />
				</Button>

				<Button
					variant="ghost"
					size="icon"
					class="hover:bg-destructive hover:text-destructive-foreground h-8 w-8"
					onclick={handleDelete}
				>
					<Trash2 class="h-4 w-4" />
				</Button>
			</div>
		</div>
	</CardHeader>

	<CardContent class="gap-1">
		<div class="gap-1 space-y-3">
			<!-- Workflow stats -->
			<div class="flex flex-wrap gap-1">
				{#each triggers as trigger}
					<Badge variant="outline" class="bg-indigo-600 text-xs text-white">
						{trigger}
					</Badge>
				{/each}

				<Badge variant="outline" class="text-xs">
					<GitBranch class="mr-1 h-3 w-3" />
					nodes: {recipe.NodeCount}
				</Badge>

				<!-- <Badge variant="outline" class="text-xs">
					<Activity class="mr-1 h-3 w-3" />
					executions: {recipe.ExecutionCount}
				</Badge> -->
			</div>

			<!-- Last execution info -->
			{#if recipe.LastExecution}
				<div class="text-muted-foreground text-xs">
					Last executed: {recipe.LastExecution}
				</div>
			{:else}
				<div class="text-muted-foreground text-xs">Never executed</div>
			{/if}

			<!-- Actions -->
			<div class="flex gap-2 pt-2">
				<Button variant="default" size="sm" class="flex-1 gap-2" onclick={handleExecute}>
					<Play class="h-3 w-3" />
					Execute
				</Button>

				<Button variant="outline" size="sm" class="gap-2" onclick={handleEdit}>
					<Edit class="h-3 w-3" />
					Edit
				</Button>
			</div>
		</div>
	</CardContent>
</Card>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
