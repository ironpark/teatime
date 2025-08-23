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
		Copy,
		GitBranch
	} from 'lucide-svelte';
	import type { RecipeInfo } from '$bindings/services';


	interface Props {
		recipe: RecipeInfo;
		onCardClick?: (recipe: RecipeInfo) => void;
		onEditDetails?: (recipe: RecipeInfo) => void;
		onDelete?: (recipe: RecipeInfo) => void;
		onDuplicate?: (recipe: RecipeInfo) => void;
	}

	let {
		recipe,
		onCardClick = () => {},
		onEditDetails = () => {},
		onDelete = () => {},
		onDuplicate = () => {}
	}: Props = $props();



	function handleCardClick() {
		onCardClick(recipe);
	}

	function handleEditDetails(event: Event) {
		event.stopPropagation(); // Prevent card click
		onEditDetails(recipe);
	}

	function handleDelete(event: Event) {
		event.stopPropagation(); // Prevent card click
		onDelete(recipe);
	}

	function handleDuplicate(event: Event) {
		event.stopPropagation(); // Prevent card click
		onDuplicate(recipe);
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

<Card class="recipe-card group gap-1 transition-all duration-150 cursor-pointer relative" onclick={handleCardClick}>
	<!-- Floating Action Buttons -->
	<div class="absolute top-2 right-2 flex gap-1 opacity-0 transition-opacity group-hover:opacity-100 z-10">
		<Button variant="ghost" size="icon" class="h-8 w-8 bg-background/80 backdrop-blur-sm hover:bg-background shadow-sm" onclick={handleEditDetails}>
			<Edit class="h-4 w-4" />
		</Button>

		<Button variant="ghost" size="icon" class="h-8 w-8 bg-background/80 backdrop-blur-sm hover:bg-background shadow-sm" onclick={handleDuplicate}>
			<Copy class="h-4 w-4" />
		</Button>

		<Button
			variant="ghost"
			size="icon"
			class="h-8 w-8 bg-background/80 backdrop-blur-sm hover:bg-destructive hover:text-destructive-foreground shadow-sm"
			onclick={handleDelete}
		>
			<Trash2 class="h-4 w-4" />
		</Button>
	</div>

	<CardHeader class="pb-3">
		<div class="w-full">
			<div class="mb-2 flex items-center gap-2">
				<!-- <CategoryIcon class="w-4 h-4 text-muted-foreground" />
          <Badge variant="outline" class={`text-xs ${statusColors[workflowData.status]}`}>
            {workflowData.status}
          </Badge> -->
			</div>
			<CardTitle class="mb-1 text-lg pr-28">{recipe.Name}</CardTitle>
			<CardDescription class="line-clamp-2 pr-28">
				{recipe.Description}
			</CardDescription>
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

	:global(.recipe-card) {
		transition: all 0.15s ease-out;
	}

	:global(.recipe-card:hover) {
		box-shadow: 
			0 20px 25px -5px rgba(30, 64, 175, 0.12),
			0 10px 10px -5px rgba(29, 78, 216, 0.08),
			0 0 0 1px rgba(37, 99, 235, 0.06);
		transform: translateY(-4px);
	}

	/* Dark mode glow */
	:global(.dark .recipe-card:hover) {
		box-shadow: 
			0 20px 25px -5px rgba(30, 58, 138, 0.18),
			0 10px 10px -5px rgba(29, 78, 216, 0.12),
			0 0 0 1px rgba(37, 99, 235, 0.08);
	}
</style>
