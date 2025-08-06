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

	// Load nodes on mount
	onMount(async () => {
		await nodeStore.loadNodes();
		
		// Create example workflow if nodes are available
		if (nodeStore.availableNodes.length > 0 && workflowStore.nodes.length === 0) {
			await workflowStore.createExampleWorkflow();
		}
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

	function saveRecipe() {
		const recipe = workflowStore.saveToLocal();
		console.log('Recipe saved:', recipe);
		
		// Show success message and redirect to recipe list
		alert('Recipe saved successfully!');
		goto('/recipes');
	}

	function executeRecipe() {
		const currentRecipe = {
			id: `temp-${Date.now()}`,
			name: recipeName,
			description: recipeDescription,
			nodes: nodes,
			edges: edges
		};

		// Store recipe for execution
		sessionStorage.setItem('recipe-to-execute', JSON.stringify(currentRecipe));
		goto('/execution');
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
