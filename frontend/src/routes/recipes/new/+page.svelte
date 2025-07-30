<script lang="ts">
	import type { Node, Edge } from '@xyflow/svelte';
	import FlowEditor from '$lib/components/flow/FlowEditor.svelte';
	import RecipeToolbar from '$lib/components/flow/RecipeToolbar.svelte';
	import { Button } from '$lib/components/ui/button';
	import { SidebarTrigger } from '$lib/components/ui/sidebar';
	import { Separator } from '$lib/components/ui/separator';
	import { goto } from '$app/navigation';


	// Initial nodes for video transcription workflow
	let nodes = $state<Node[]>([
		// Trigger nodes
		{
			id: '1',
			type: 'trigger',
			position: { x: 50, y: 100 },
			data: {
				label: 'CLI Transcribe',
				nodeType: 'command',
				command: 'transcribe --input ${INPUT_FILE} --lang ${LANGUAGE} --output ${OUTPUT_PATH}',
				workingDir: '/workspace'
			}
		},
		{
			id: '2',
			type: 'trigger',
			position: { x: 50, y: 200 },
			data: {
				label: 'File Server MCP',
				nodeType: 'mcp',
				serverUrl: 'mcp://file-server/v1',
				method: 'process_video',
				params: '{"video_path": "${INPUT_FILE}"}'
			}
		},
		{
			id: '3',
			type: 'trigger',
			position: { x: 50, y: 350 },
			data: {
				label: 'Start Processing',
				nodeType: 'manual',
				description: 'Click to start video transcription workflow',
				requiresConfirmation: true
			}
		},
		// Processing nodes
		{
			id: '4',
			type: 'action',
			position: { x: 350, y: 150 },
			data: {
				label: 'Extract Audio',
				nodeType: 'ffmpeg-extract',
				outputFormat: 'wav',
				quality: 'high'
			}
		},
		{
			id: '5',
			type: 'action',
			position: { x: 600, y: 150 },
			data: {
				label: 'Speech to Text',
				nodeType: 'transcription',
				language: 'auto',
				model: 'whisper'
			}
		},
		{
			id: '6',
			type: 'action',
			position: { x: 850, y: 150 },
			data: {
				label: 'Translate Text',
				nodeType: 'translation',
				sourceLanguage: 'auto',
				targetLanguage: 'ko'
			}
		},
		{
			id: '7',
			type: 'action',
			position: { x: 1100, y: 150 },
			data: {
				label: 'AI Text Correction',
				nodeType: 'llm-correction',
				model: 'gpt-4',
				prompt: 'Improve grammar and fluency'
			}
		},
		{
			id: '8',
			type: 'action',
			position: { x: 1350, y: 150 },
			data: {
				label: 'Save Result',
				nodeType: 'file-save',
				outputPath: '/output/transcription.txt',
				format: 'txt'
			}
		}
	]);

	let edges = $state<Edge[]>([
		// Command trigger path
		{ id: 'e1-4', source: '1', target: '4', type: 'smoothstep' },
		// MCP trigger path
		{ id: 'e2-4', source: '2', target: '4', type: 'smoothstep' },
		// Manual trigger path
		{ id: 'e3-4', source: '3', target: '4', type: 'smoothstep' },
		// Processing pipeline
		{ id: 'e4-5', source: '4', target: '5', type: 'smoothstep' },
		{ id: 'e5-6', source: '5', target: '6', type: 'smoothstep' },
		{ id: 'e6-7', source: '6', target: '7', type: 'smoothstep' },
		{ id: 'e7-8', source: '7', target: '8', type: 'smoothstep' }
	]);

	let selectedNodes = $state<Node[]>([]);
	let selectedEdges = $state<Edge[]>([]);

	// Recipe metadata
	let recipeName = $state('Video Transcription Workflow');
	let recipeDescription = $state('Extract audio from video, transcribe, translate, and improve with AI');

	function handleSelectionChange(selection: { nodes: Node[]; edges: Edge[] }) {
		selectedNodes = selection.nodes;
		selectedEdges = selection.edges;
	}


	function deleteSelected() {
		const selectedNodeIds = new Set(selectedNodes.map((n) => n.id));
		const selectedEdgeIds = new Set(selectedEdges.map((e) => e.id));

		nodes = nodes.filter((node) => !selectedNodeIds.has(node.id));
		edges = edges.filter(
			(edge) =>
				!selectedEdgeIds.has(edge.id) &&
				!selectedNodeIds.has(edge.source) &&
				!selectedNodeIds.has(edge.target)
		);

		selectedNodes = [];
		selectedEdges = [];
	}

	function saveRecipe() {
		const recipe = {
			id: `recipe-${Date.now()}`,
			name: recipeName,
			description: recipeDescription,
			nodes: nodes,
			edges: edges,
			createdAt: new Date().toISOString()
		};

		// Save to localStorage for now (could be extended to save to backend)
		const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
		savedRecipes.push(recipe);
		localStorage.setItem('teatime-recipes', JSON.stringify(savedRecipes));

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
						bind:value={recipeName}
						class="border-none bg-transparent text-lg font-semibold outline-none"
						placeholder="Recipe name"
					/>
					<input
						bind:value={recipeDescription}
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
			bind:nodes
			bind:edges
			onSelectionChange={handleSelectionChange}
		/>
	</div>
</div>
