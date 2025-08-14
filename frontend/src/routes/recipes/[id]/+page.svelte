<script lang="ts">
	import type { Node, Edge } from '@xyflow/svelte';
	import FlowEditor from '$lib/components/flow/FlowEditor.svelte';
	import RecipeToolbar from '$lib/components/flow/RecipeToolbar.svelte';
	import { Button } from '$lib/components/ui/button';
	import { SidebarTrigger } from '$lib/components/ui/sidebar';
	import { Separator } from '$lib/components/ui/separator';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { nodeStore } from '$lib/stores/nodes.svelte';
	import { RecipesService } from '$bindings/services';
	import { Node as RecipeNode } from '$bindings/internal/recipe/models';
	import { RecipeStore } from '$lib/stores/recipe.svelte';
	// Get recipe ID from URL parameters
	const recipeId = page.params.id || '';
	const isNewRecipe = false;
	const recipeStore = new RecipeStore(recipeId);
	// Recipe loading state
	// 	let loading = $state(true);
	// let error = $state(null as string | null);
	// let currentRecipe = $state(null as any);

	// Load recipe data from the store

	function handleSelectionChange(selection: { nodes: Node[]; edges: Edge[] }) {
		recipeStore.selectedNodes = selection.nodes;
		recipeStore.selectedEdges = selection.edges;
	}

	function deleteSelected() {
		// recipeStore.deleteSelected();
	}

	// async function saveRecipe() {
	// 	try {
	// 		let recipe;
			
	// 		if (isNewRecipe) {
	// 			// Create new recipe
	// 			const result = await RecipesService.CreateRecipe(
	// 				recipeStore.info.name || 'Untitled Recipe',
	// 				recipeStore.info.description || 'No description'
	// 			);
				
	// 			if (!result || !result.Recipe) {
	// 				alert('Failed to create recipe. Please try again.');
	// 				return;
	// 			}
				
	// 			recipe = result.Recipe;
	// 			currentRecipe = recipe;
				
	// 			// Update URL to reflect the new recipe ID (remove 'new')
	// 			goto(`/recipes/${result.ID}`, { replaceState: true });
	// 		} else {
	// 			// Update existing recipe
	// 			if (!currentRecipe) {
	// 				alert('Recipe not found. Please try again.');
	// 				return;
	// 			}
				
	// 			recipe = { ...currentRecipe };
	// 		}
			
	// 		// Update the recipe with current workflow data
	// 		recipe.name = recipeStore.info.name || 'Untitled Recipe';
	// 		recipe.description = recipeStore.info.description || 'No description';
			
	// 		recipe.nodes = recipeStore.nodes;
			
	// 		recipe.edges = edges.map(edge => ({
	// 			id: edge.id,
	// 			source: edge.source,
	// 			target: edge.target,
	// 			sourceHandle: edge.sourceHandle || '',
	// 			targetHandle: edge.targetHandle || ''
	// 		}));
			
	// 		// Save the recipe
	// 		if (isNewRecipe) {
	// 			await RecipesService.SaveRecipe(recipe);
	// 		} else {
	// 			await RecipesService.UpdateRecipe(recipeId, recipe);
	// 		}
			
	// 		console.log('Recipe saved:', recipe);
	// 		alert(`Recipe ${isNewRecipe ? 'created' : 'updated'} successfully!`);
			
	// 		// Only redirect to recipes list if it was a new recipe
	// 		if (isNewRecipe) {
	// 			goto('/recipes');
	// 		}
			
	// 	} catch (error) {
	// 		console.error('Failed to save recipe:', error);
	// 		alert('Failed to save recipe. Please try again.');
	// 	}
	// }


</script>

<svelte:head>
	<title>{isNewRecipe ? 'New Recipe' : (recipeStore.info.name || 'Edit Recipe')} - Teatime</title>
</svelte:head>

{#if recipeStore.isLoading}
	<div class="flex h-screen w-full items-center justify-center">
		<div class="flex flex-col items-center gap-4">
			<div class="w-8 h-8 animate-spin rounded-full border-2 border-current border-t-transparent"></div>
			<p class="text-muted-foreground">Loading recipe...</p>
		</div>
	</div>
{:else if recipeStore.error}
	<div class="flex h-screen w-full items-center justify-center">
		<div class="flex flex-col items-center gap-4 text-center">
			<div class="text-destructive text-lg font-semibold">Error</div>
			<p class="text-muted-foreground">{recipeStore.error}</p>
			<div class="flex gap-2">
				<Button onclick={() => goto('/recipes')} variant="outline">
					Back to Recipes
				</Button>
				{#if !isNewRecipe}
					<Button onclick={() => recipeStore.load()}>
						Retry
					</Button>
				{/if}
			</div>
		</div>
	</div>
{:else}
	<div class="recipe-editor bg-background flex h-screen w-full flex-col">
		<!-- Header -->
		<header class="bg-card border-b">
			<div class="flex h-16 items-center gap-2 px-4">
				<SidebarTrigger />
				<Separator orientation="vertical" class="mr-2 h-4" />
				<div class="flex flex-1 items-center space-x-4">
					<div class="flex flex-col">
						<input
							bind:value={recipeStore.info.name}
							class="border-none bg-transparent text-lg font-semibold outline-none"
							placeholder="Recipe name"
						/>
						<input
							bind:value={recipeStore.info.description}
							class="text-muted-foreground border-none bg-transparent text-sm outline-none"
							placeholder="Recipe description"
						/>
					</div>
				</div>

				<div class="flex items-center gap-2">
					<RecipeToolbar
						selectedNodes={recipeStore.selectedNodes}
						selectedEdges={recipeStore.selectedEdges}
						onSave={() => {
							recipeStore.save();
						}}
						onDelete={deleteSelected}
					/>
				</div>
			</div>
		</header>

		<!-- Main content -->
		<div class="flex flex-1">
			<FlowEditor
				onSelectionChange={handleSelectionChange}
				recipeStore={recipeStore}
			/>
		</div>
	</div>
{/if}
