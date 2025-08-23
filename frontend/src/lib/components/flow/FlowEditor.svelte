<script lang="ts">
	import {
		SvelteFlow,
		Controls,
		Background,
		BackgroundVariant,
		MiniMap,
		type Node,
		type Edge,
		SvelteFlowProvider,
		type NodeTargetEventWithPointer
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';

	import UnifiedNode from './node/UnifiedNode.svelte';
	import FloatingNodeMenu from './FloatingNodeMenu.svelte';
	import NodePropertiesSheet from './NodePropertiesSheet.svelte';
	import { setContext } from 'svelte';
	import type { RecipeStore } from '$lib/stores/recipe.svelte';
	import { settingsStore } from '$lib/stores/settings.svelte';
	interface Props {
		recipeStore: RecipeStore;
		onNodesChange?: (nodes: Node[]) => void;
		onEdgesChange?: (edges: Edge[]) => void;
		onSelectionChange?: (selection: { nodes: Node[]; edges: Edge[] }) => void;
	}

	let {
		onNodesChange,
		onEdgesChange,
		onSelectionChange,
		recipeStore
	}: Props = $props();

	// Function to handle edit button click
	function handleEditClick(nodeId: string) {
		const node = recipeStore.nodes.find(n => n.id === nodeId);
		if (node) {
			recipeStore.selectedNodes = [node];
			propertiesSheetOpen = true;
		}
	}

	// Function to handle double click
	function handleNodeDoubleClick(nodeId: string) {
		const node = recipeStore.nodes.find(n => n.id === nodeId);
		if (node) {
			recipeStore.selectedNodes = [node];
			propertiesSheetOpen = true;
		}
	}

	// Set context for node handlers
	setContext('nodeHandlers', {
		onEdit: handleEditClick,
		onDoubleClick: handleNodeDoubleClick
	});


	// All node types use the UnifiedNode component
	const nodeTypes = {
		trigger: UnifiedNode,
		branch: UnifiedNode,
		action: UnifiedNode,
		util: UnifiedNode
	};

	let propertiesSheetOpen = $state(false);

	async function addNode(nodeId: string) {
		await recipeStore.createNodeByRef(nodeId, 0, 0);
	}

	function handleConnect(connection: any) {
		// console.log('Connection attempt:', connection);
		// const newEdge = {
		// 	...connection,
		// 	id: `e${connection.source}-${connection.target}`,
		// 	type: 'smoothstep'
		// };
		// recipeStore.edges = [...recipeStore.edges, newEdge];
		// onEdgesChange?.(recipeStore.edges);
	}

	function handleError(event: Event) {
		console.error('SvelteFlow error:', event);
	}

	function handleInit() {
		console.log('SvelteFlow initialized');
		console.log('Nodes:', recipeStore.nodes);
		console.log('Edges:', recipeStore.edges);
	}
	
	function handleSelectionChange({ nodes: selectedNodesList, edges: selectedEdgesList }: { nodes: Node[], edges: Edge[] }) {
		console.log('Selection changed:', selectedNodesList);
		recipeStore.selectedNodes = selectedNodesList;
		recipeStore.selectedEdges = selectedEdgesList;
	}

	function handleNodeDrag(targetNode: Node, nodes: Node[], event: MouseEvent | TouchEvent) {
		console.log('Node dragged:', targetNode, nodes, event);
	}
</script>

<SvelteFlowProvider>
	<div class="relative flex-1">
		<SvelteFlow
			bind:nodes={recipeStore.nodes}
			bind:edges={recipeStore.edges}
			{nodeTypes}
			onconnect={handleConnect}
			oninit={handleInit}
			onselectionchange={handleSelectionChange}
			onerror={handleError}
			fitView
			snapGrid={[15, 15]}
			defaultEdgeOptions={{ type: 'smoothstep' }}
			proOptions={{ hideAttribution: true }}
			colorMode={settingsStore.theme}
		>
			<Controls />
			<Background variant={BackgroundVariant.Dots} />
			<MiniMap nodeStrokeWidth={3} pannable={true} zoomable={true} />
			<!-- <SelectionHandler /> -->
		</SvelteFlow>
	</div>

	<!-- Floating Node Menu -->
	<FloatingNodeMenu {addNode} />

	<!-- Node Properties Sheet -->
	<NodePropertiesSheet 
		recipeStore={recipeStore}
		bind:selectedNodes={recipeStore.selectedNodes} 
		bind:open={propertiesSheetOpen}
		onOpenChange={(open) => {
			propertiesSheetOpen = open;
		}}
	/>
</SvelteFlowProvider>
