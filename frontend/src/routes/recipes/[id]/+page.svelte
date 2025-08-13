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
	import { nodeStore, workflowStore } from '$lib/stores/nodes.svelte';
	import { RecipesService } from '$bindings/services';
	import { Node as RecipeNode } from '$bindings/internal/recipe/models';

	// Get recipe ID from URL parameters
	const recipeId = page.params.id || '';
	const isNewRecipe = false;

	// Recipe loading state
	let loading = $state(true);
	let error = $state(null as string | null);
	let currentRecipe = $state(null as any);

	// Load nodes and recipe on mount
	onMount(async () => {
		await nodeStore.loadNodes();
		if (!isNewRecipe) {
			await loadRecipe(recipeId);
		} else {
			// Reset workflow store for new recipes
			workflowStore.clear();
			loading = false;
		}
	});

	// Use store values reactively
	let nodes = $derived(workflowStore.nodes);
	let edges = $derived(workflowStore.edges);
	let selectedNodes = $derived(workflowStore.selectedNodes);
	let selectedEdges = $derived(workflowStore.selectedEdges);
	let recipeName = $derived(workflowStore.recipeName);
	let recipeDescription = $derived(workflowStore.recipeDescription);

	// Load recipe data from the store
	async function loadRecipe(id: string) {
		try {
			loading = true;
			error = null;
			
			// Get recipe from backend
			const recipe = await RecipesService.GetRecipe(id);
			
			if (!recipe) {
				error = 'Recipe not found';
				return;
			}
			
			currentRecipe = recipe;
			
			// Load recipe data into workflow store
			workflowStore.recipeName = recipe.name || '';
			workflowStore.recipeDescription = recipe.description || '';
			// workflowStore.nodes = recipe.nodes || [];
			// Convert recipe nodes to workflow nodes
			if (recipe.nodes && recipe.nodes.length > 0) {
				workflowStore.nodes = recipe.nodes.map((node: RecipeNode) => ({
					id: node.id,
					type: node.type,
					position: { x: node.pos[0], y: node.pos[1] },
					data: {
						icon: node.icon || '',
						label: node.name || '',
						name: node.name || '',
						nodeType: node.type || '',
						description: node.description || '',
						properties: node.properties || [],
						outputs: node.output || [],
						backendNodeRef: node.ref || ''
					}
				}));
				// workflowStore.nodes = workflowNodes;
			}
			
			// Convert recipe edges to workflow edges  
			if (recipe.edges && recipe.edges.length > 0) {
				const workflowEdges = recipe.edges.map((edge: any) => ({
					id: edge.id,
					source: edge.source,
					target: edge.target,
					sourceHandle: edge.sourceHandle || null,
					targetHandle: edge.targetHandle || null
				}));
				workflowStore.edges = workflowEdges;
			}
			
		} catch (err) {
			console.error('Failed to load recipe:', err);
			error = 'Failed to load recipe. Please try again.';
		} finally {
			loading = false;
		}
	}

	function handleSelectionChange(selection: { nodes: Node[]; edges: Edge[] }) {
		workflowStore.selectedNodes = selection.nodes;
		workflowStore.selectedEdges = selection.edges;
	}

	function deleteSelected() {
		workflowStore.deleteSelected();
	}

	async function saveRecipe() {
		try {
			let recipe;
			
			if (isNewRecipe) {
				// Create new recipe
				const result = await RecipesService.CreateRecipe(
					workflowStore.recipeName || 'Untitled Recipe',
					workflowStore.recipeDescription || 'No description'
				);
				
				if (!result || !result.Recipe) {
					alert('Failed to create recipe. Please try again.');
					return;
				}
				
				recipe = result.Recipe;
				currentRecipe = recipe;
				
				// Update URL to reflect the new recipe ID (remove 'new')
				goto(`/recipes/${result.ID}`, { replaceState: true });
			} else {
				// Update existing recipe
				if (!currentRecipe) {
					alert('Recipe not found. Please try again.');
					return;
				}
				
				recipe = { ...currentRecipe };
			}
			
			// Update the recipe with current workflow data
			recipe.name = workflowStore.recipeName || 'Untitled Recipe';
			recipe.description = workflowStore.recipeDescription || 'No description';
			
			recipe.nodes = nodes.map(node => ({
				id: node.id,
				ref: node.data?.backendNodeRef || '',
				position: { x: node.position.x, y: node.position.y },
				properties: node.data?.properties || [],
				output: node.data?.outputs || [],
				name: node.data?.label || '',
				description: node.data?.description || '',
				type: node.data?.nodeType || ''
			}));
			
			recipe.edges = edges.map(edge => ({
				id: edge.id,
				source: edge.source,
				target: edge.target,
				sourceHandle: edge.sourceHandle || '',
				targetHandle: edge.targetHandle || ''
			}));
			
			// Save the recipe
			if (isNewRecipe) {
				await RecipesService.SaveRecipe(recipe);
			} else {
				await RecipesService.UpdateRecipe(recipeId, recipe);
			}
			
			console.log('Recipe saved:', recipe);
			alert(`Recipe ${isNewRecipe ? 'created' : 'updated'} successfully!`);
			
			// Only redirect to recipes list if it was a new recipe
			if (isNewRecipe) {
				goto('/recipes');
			}
			
		} catch (error) {
			console.error('Failed to save recipe:', error);
			alert('Failed to save recipe. Please try again.');
		}
	}


</script>

<svelte:head>
	<title>{isNewRecipe ? 'New Recipe' : (recipeName || 'Edit Recipe')} - Teatime</title>
</svelte:head>

{#if loading}
	<div class="flex h-screen w-full items-center justify-center">
		<div class="flex flex-col items-center gap-4">
			<div class="w-8 h-8 animate-spin rounded-full border-2 border-current border-t-transparent"></div>
			<p class="text-muted-foreground">Loading recipe...</p>
		</div>
	</div>
{:else if error}
	<div class="flex h-screen w-full items-center justify-center">
		<div class="flex flex-col items-center gap-4 text-center">
			<div class="text-destructive text-lg font-semibold">Error</div>
			<p class="text-muted-foreground">{error}</p>
			<div class="flex gap-2">
				<Button onclick={() => goto('/recipes')} variant="outline">
					Back to Recipes
				</Button>
				{#if !isNewRecipe}
					<Button onclick={() => loadRecipe(recipeId)}>
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
							bind:value={workflowStore.recipeName}
							class="border-none bg-transparent text-lg font-semibold outline-none"
							placeholder="Recipe name"
						/>
						<input
							bind:value={workflowStore.recipeDescription}
							class="text-muted-foreground border-none bg-transparent text-sm outline-none"
							placeholder="Recipe description"
						/>
					</div>
				</div>

				<div class="flex items-center gap-2">
					<RecipeToolbar
						{selectedNodes}
						{selectedEdges}
						onSave={saveRecipe}
						onDelete={deleteSelected}
					/>
				</div>
			</div>
		</header>

		<!-- Main content -->
		<div class="flex flex-1">
			<FlowEditor
				bind:nodes={workflowStore.nodes}
				bind:edges={workflowStore.edges}
				onSelectionChange={handleSelectionChange}
			/>
		</div>
	</div>
{/if}
