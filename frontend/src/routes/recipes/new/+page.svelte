<script lang="ts">
	import type { Node, Edge } from '@xyflow/svelte';
	import FlowEditor from '$lib/components/flow/FlowEditor.svelte';
	import RecipeToolbar from '$lib/components/flow/RecipeToolbar.svelte';
	import { Button } from '$lib/components/ui/button';
	import { SidebarTrigger } from '$lib/components/ui/sidebar';
	import { Separator } from '$lib/components/ui/separator';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { nodeStore, workflowStore } from '$lib/stores/nodes.svelte';
	import { RecipesService } from '$bindings/services';

	// Load nodes on mount
	onMount(async () => {
		await nodeStore.loadNodes();
	});

	// Use store values reactively
	let nodes = $derived(workflowStore.nodes);
	let edges = $derived(workflowStore.edges);
	let selectedNodes = $derived(workflowStore.selectedNodes);
	let selectedEdges = $derived(workflowStore.selectedEdges);
	let recipeName = $derived(workflowStore.recipeName);
	let recipeDescription = $derived(workflowStore.recipeDescription);

	function handleSelectionChange(selection: { nodes: Node[]; edges: Edge[] }) {
		workflowStore.selectedNodes = selection.nodes;
		workflowStore.selectedEdges = selection.edges;
	}

	function deleteSelected() {
		workflowStore.deleteSelected();
	}

	async function saveRecipe() {
		try {
			// Create the recipe in the backend
			const result = await RecipesService.CreateRecipe(
				workflowStore.recipeName || 'Untitled Recipe',
				workflowStore.recipeDescription || 'No description'
			);
			
			if (result && result.Recipe) {
				// Update the recipe with the workflow data
				result.Recipe.nodes = nodes.map(node => ({
					id: node.id,
					ref: node.data?.backendNodeRef || '',
					position: { x: node.position.x, y: node.position.y },
					properties: node.data?.properties || [],
					output: node.data?.outputs || [],
					name: node.data?.label || '',
					description: node.data?.description || '',
					type: node.data?.nodeType || ''
				}));
				
				result.Recipe.edges = edges.map(edge => ({
					id: edge.id,
					source: edge.source,
					target: edge.target,
					sourceHandle: edge.sourceHandle || '',
					targetHandle: edge.targetHandle || ''
				}));
				
				// Save the updated recipe
				await RecipesService.SaveRecipe(result.Recipe);
				
				console.log('Recipe saved:', result);
				alert('Recipe saved successfully!');
				goto('/recipes');
			} else {
				alert('Failed to create recipe. Please try again.');
			}
		} catch (error) {
			console.error('Failed to save recipe:', error);
			alert('Failed to save recipe. Please try again.');
		}
	}


</script>

<svelte:head>
	<title>New Recipe - Teatime</title>
</svelte:head>

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
