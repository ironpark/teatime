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

	function addNode(nodeType: string) {
		const newId = `${Date.now()}`;
		// Determine category type
		const categoryType =
			nodeType === 'trigger' ||
			nodeType === 'branch' ||
			nodeType === 'action' ||
			nodeType === 'util'
				? nodeType
				: getCategoryType(nodeType);

		const newNode: Node = {
			id: newId,
			type: categoryType,
			position: { x: Math.random() * 500 + 100, y: Math.random() * 300 + 100 },
			data: {
				...getDefaultNodeData(nodeType),
				nodeType: nodeType
			}
		};

		nodes = [...nodes, newNode];
		onNodesChange?.(nodes);
	}

	function getCategoryType(nodeType: string): 'trigger' | 'branch' | 'action' | 'util' {
		const triggerTypes = ['command', 'mcp', 'manual', 'file-input', 'ingredient'];
		const actionTypes = [
			'ffmpeg-extract',
			'transcription',
			'translation',
			'llm-correction',
			'file-save',
			'output',
			'step'
		];
		const utilTypes = ['variable', 'constant', 'formatter', 'calculator', 'aggregator'];

		if (triggerTypes.includes(nodeType)) return 'trigger';
		if (actionTypes.includes(nodeType)) return 'action';
		if (utilTypes.includes(nodeType)) return 'util';
		return 'util';
	}

	function getDefaultNodeData(nodeType: string) {
		switch (nodeType) {
			case 'command':
				return { label: 'Command Line', command: '', workingDir: '' };
			case 'mcp':
				return { label: 'MCP Server', serverUrl: '', method: '', params: '' };
			case 'manual':
				return {
					label: 'Manual Trigger',
					description: 'Click to start',
					requiresConfirmation: false
				};
			case 'file-input':
				return { label: 'File Input', filePath: '', fileType: 'any' };
			case 'ffmpeg-extract':
				return { label: 'Extract Audio', outputFormat: 'wav', quality: 'high' };
			case 'transcription':
				return { label: 'Speech to Text', language: 'auto', model: 'whisper' };
			case 'translation':
				return { label: 'Translate Text', sourceLanguage: 'auto', targetLanguage: 'en' };
			case 'llm-correction':
				return { label: 'AI Text Correction', model: 'gpt-4', prompt: 'Improve text quality' };
			case 'file-save':
				return { label: 'Save Result', outputPath: '/output/result.txt', format: 'txt' };
			case 'variable':
				return { label: 'New Variable', value: '', type: 'string' };
			case 'constant':
				return { label: 'New Constant', value: '', type: 'string' };
			case 'formatter':
				return { label: 'New Formatter', format: 'YYYY-MM-DD', type: 'date' };
			case 'calculator':
				return { label: 'New Calculator', expression: 'a + b', operation: 'add' };
			case 'aggregator':
				return { label: 'New Aggregator', operation: 'sum', field: '' };
			default:
				return { label: 'New Node' };
		}
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

		onNodesChange?.(nodes);
		onEdgesChange?.(edges);

		selectedNodes = [];
		selectedEdges = [];
	}

	function updateNodeData(nodeId: string, updates: any) {
		nodes = nodes.map((node) => {
			if (node.id === nodeId) {
				return {
					...node,
					data: {
						...node.data,
						...updates
					}
				};
			}
			return node;
		});
		onNodesChange?.(nodes);
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
