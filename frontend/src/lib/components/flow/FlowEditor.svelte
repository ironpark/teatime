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
	import NodePropertiesSidebar from './NodePropertiesSidebar.svelte';
	import { workflowStore } from '$lib/stores/nodes.svelte';

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

	// All node types use the UnifiedNode component
	const nodeTypes = {
		trigger: UnifiedNode,
		branch: UnifiedNode,
		action: UnifiedNode,
		util: UnifiedNode
	};

	let selectedNodes = $state<Node[]>([]);
	let selectedEdges = $state<Edge[]>([]);

	async function addNode(nodeId: string) {
		const newNode = await workflowStore.addNode(nodeId);
		if (newNode && onNodesChange) {
			onNodesChange(workflowStore.nodes);
		}
	}

	function updateNodeData(nodeId: string, updates: Record<string, unknown>) {
		workflowStore.updateNode(nodeId, updates);
		if (onNodesChange) {
			onNodesChange(workflowStore.nodes);
		}
	}

	function closeSidebar() {
		selectedNodes = [];
		selectedEdges = [];
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

	function handleError(error: Error, el: Element) {
		console.error('SvelteFlow error:', error);
	}

	function handleInit() {
		console.log('SvelteFlow initialized');
		console.log('Nodes:', nodes);
		console.log('Edges:', edges);
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

	<!-- Node Properties Sidebar -->
	<NodePropertiesSidebar {selectedNodes} onNodeUpdate={updateNodeData} onClose={closeSidebar} />
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
