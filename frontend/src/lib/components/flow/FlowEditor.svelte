<script lang="ts">
	import {
		SvelteFlow,
		Controls,
		Background,
		BackgroundVariant,
		MiniMap,
		type Node,
		type Edge,
		SvelteFlowProvider
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';

	import UnifiedNode from './node/UnifiedNode.svelte';
	import FloatingNodeMenu from './FloatingNodeMenu.svelte';
	import NodePropertiesSheet from './NodePropertiesSheet.svelte';
	import { workflowStore } from '$lib/stores/nodes.svelte';
	import { setContext } from 'svelte';

	interface Props {
		nodes: Node[];
		edges: Edge[];
		onNodesChange?: (nodes: Node[]) => void;
		onEdgesChange?: (edges: Edge[]) => void;
		onSelectionChange?: (selection: { nodes: Node[]; edges: Edge[] }) => void;
	}

	let {
		nodes = $bindable([]),
		edges = $bindable([]),
		onNodesChange,
		onEdgesChange,
		onSelectionChange
	}: Props = $props();

	// Function to handle edit button click
	function handleEditClick(nodeId: string) {
		const node = nodes.find(n => n.id === nodeId);
		if (node) {
			selectedNodes = [node];
			propertiesSheetOpen = true;
		}
	}

	// Function to handle double click
	function handleNodeDoubleClick(nodeId: string) {
		const node = nodes.find(n => n.id === nodeId);
		if (node) {
			selectedNodes = [node];
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

	let selectedNodes = $state<Node[]>([]);
	let selectedEdges = $state<Edge[]>([]);
	let propertiesSheetOpen = $state(false);

	async function addNode(nodeId: string) {
		const newNode = await workflowStore.addNode(nodeId);
		if (newNode) {
			nodes = [...workflowStore.nodes];
			if (onNodesChange) {
				onNodesChange(workflowStore.nodes);
			}
		}
	}

	function updateNodeData(nodeId: string, updates: Record<string, unknown>) {
		workflowStore.updateNode(nodeId, updates);
		if (onNodesChange) {
			onNodesChange(workflowStore.nodes);
		}
	}

	function handlePropertiesSheetChange(open: boolean) {
		propertiesSheetOpen = open;
		if (!open) {
			// Optionally clear selection when sheet closes
			// selectedNodes = [];
			// selectedEdges = [];
		}
	}

	function handleConnect(connection: any) {
		console.log('Connection attempt:', connection);
		const newEdge = {
			...connection,
			id: `e${connection.source}-${connection.target}`,
			type: 'smoothstep'
		};
		edges = [...edges, newEdge];
		onEdgesChange?.(edges);
	}

	function handleError(event: CustomEvent) {
		console.error('SvelteFlow error:', event.detail);
	}

	function handleInit() {
		console.log('SvelteFlow initialized');
		console.log('Nodes:', nodes);
		console.log('Edges:', edges);
	}
	
	function handleSelectionChange({ nodes: selectedNodesList, edges: selectedEdgesList }: { nodes: Node[], edges: Edge[] }) {
		console.log('Selection changed:', selectedNodesList);
		selectedNodes = selectedNodesList;
		selectedEdges = selectedEdgesList;
		
		// Don't auto-open sheet on selection anymore - only open via edit button
	}
</script>

<SvelteFlowProvider>
	<div class="relative flex-1">
		<SvelteFlow
			bind:nodes
			bind:edges
			{nodeTypes}
			onconnect={handleConnect}
			onerror={handleError}
			oninit={handleInit}
			onselectionchange={handleSelectionChange}
			fitView
			snapGrid={[15, 15]}
			defaultEdgeOptions={{ type: 'smoothstep' }}
			proOptions={{ hideAttribution: true }}
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
		bind:selectedNodes 
		onNodeUpdate={updateNodeData} 
		bind:open={propertiesSheetOpen}
		onOpenChange={handlePropertiesSheetChange}
	/>
</SvelteFlowProvider>

<style>
	:global(.svelte-flow) {
		background-color: hsl(var(--background));
	}

	:global(.svelte-flow__node) {
		font-family: inherit;
	}

	:global(.svelte-flow__edge) {
		stroke: hsl(var(--border));
	}

	:global(.svelte-flow__edge.selected) {
		stroke: hsl(var(--primary));
	}

	:global(.svelte-flow__connection-line) {
		stroke: hsl(var(--primary));
	}
</style>
