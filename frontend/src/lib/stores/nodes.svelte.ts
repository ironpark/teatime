import * as TeaTime from '$bindings/services/teatime';
import type { NodeInfo, Node as BackendNode } from '$bindings/stores/models';
import type { Node, Edge } from '@xyflow/svelte';

// Store for available node types from backend
class NodeStore {
	availableNodes = $state<NodeInfo[]>([]);
	isLoading = $state(false);
	error = $state<string | null>(null);
	
	// Categorized nodes for easier access
	get triggerNodes() {
		return this.availableNodes.filter(n => 
			n.Type.toLowerCase().includes('trigger')
		);
	}
	
	get actionNodes() {
		return this.availableNodes.filter(n => 
			n.Type.toLowerCase().includes('action') && 
			!n.Type.toLowerCase().includes('trigger')
		);
	}
	
	get branchNodes() {
		return this.availableNodes.filter(n => 
			n.Type.toLowerCase().includes('branch') || 
			n.Type.toLowerCase().includes('conditional') || 
			n.Type.toLowerCase().includes('loop') ||
			n.Type.toLowerCase().includes('switch')
		);
	}
	
	get utilityNodes() {
		return this.availableNodes.filter(n => 
			!n.Type.toLowerCase().includes('trigger') &&
			!n.Type.toLowerCase().includes('action') &&
			!n.Type.toLowerCase().includes('branch') &&
			!n.Type.toLowerCase().includes('conditional') &&
			!n.Type.toLowerCase().includes('loop') &&
			!n.Type.toLowerCase().includes('switch')
		);
	}
	
	// Load nodes from backend
	async loadNodes() {
		if (this.isLoading) return;
		
		this.isLoading = true;
		this.error = null;
		
		try {
			const nodes = await TeaTime.GetNodeInfos();
			this.availableNodes = nodes;
			console.log('Loaded nodes:', nodes);
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to load nodes';
			console.error('Failed to load nodes:', err);
		} finally {
			this.isLoading = false;
		}
	}
	
	// Get nodes by type
	getNodesByType(type: string) {
		return TeaTime.GetNodeInfosByType(type);
	}
	
	// Create a new node instance from backend
	async createNode(nodeId: string): Promise<BackendNode | null> {
		try {
			const node = await TeaTime.CreateNode(nodeId);
			return node;
		} catch (err) {
			console.error('Failed to create node:', err);
			return null;
		}
	}
	
	// Create a flow node for the editor
	async createFlowNode(nodeId: string, position: { x: number; y: number }): Promise<Node | null> {
		const backendNode = await this.createNode(nodeId);
		if (!backendNode) return null;
		
		const categoryType = this.getCategoryType(backendNode.Type);
		
		return {
			id: `node-${Date.now()}-${Math.random()}`,
			type: categoryType,
			position,
			data: {
				label: backendNode.Name,
				nodeType: backendNode.Type,
				description: backendNode.Description,
				properties: backendNode.Properties,
				outputs: backendNode.Output,
				backendNodeId: backendNode.Id
			}
		};
	}
	
	// Determine visual category from node type
	getCategoryType(nodeType: string): 'trigger' | 'branch' | 'action' | 'util' {
		const type = nodeType.toLowerCase();
		if (type.includes('trigger')) return 'trigger';
		if (type.includes('branch') || type.includes('conditional') || type.includes('loop')) return 'branch';
		if (type.includes('action')) return 'action';
		return 'util';
	}
	
	// Get icon for a node type
	getNodeIcon(nodeType: string, nodeName: string): string {
		const type = nodeType.toLowerCase();
		const name = nodeName.toLowerCase();
		
		// Trigger icons
		if (name.includes('command') || name.includes('cli')) return 'Terminal';
		if (name.includes('webhook')) return 'Plug';
		if (name.includes('manual')) return 'Hand';
		if (name.includes('file')) return 'FileImage';
		
		// Branch icons
		if (name.includes('condition')) return 'Filter';
		if (name.includes('switch')) return 'Shuffle';
		if (name.includes('loop')) return 'GitBranch';
		
		// Action icons
		if (name.includes('llm') || name.includes('ai')) return 'Bot';
		if (name.includes('save')) return 'Save';
		if (name.includes('translate')) return 'Languages';
		if (name.includes('transcribe')) return 'FileText';
		if (name.includes('audio')) return 'Music';
		
		// Utility icons
		if (name.includes('variable')) return 'Variable';
		if (name.includes('constant')) return 'Hash';
		if (name.includes('format')) return 'Type';
		if (name.includes('calc')) return 'Calculator';
		if (name.includes('aggregate')) return 'BarChart3';
		
		// Default icons by category
		if (type.includes('trigger')) return 'Zap';
		if (type.includes('branch')) return 'GitBranch';
		if (type.includes('action')) return 'Activity';
		
		return 'Package';
	}
}

// Create singleton instance
export const nodeStore = new NodeStore();

// Workflow store for managing the current workflow
class WorkflowStore {
	nodes = $state<Node[]>([]);
	edges = $state<Edge[]>([]);
	selectedNodes = $state<Node[]>([]);
	selectedEdges = $state<Edge[]>([]);
	
	// Recipe metadata
	recipeName = $state('');
	recipeDescription = $state('');
	
	// Add a node to the workflow
	async addNode(nodeId: string) {
		const position = {
			x: Math.random() * 500 + 100,
			y: Math.random() * 300 + 100
		};
		
		const newNode = await nodeStore.createFlowNode(nodeId, position);
		if (newNode) {
			this.nodes = [...this.nodes, newNode];
		}
		return newNode;
	}
	
	// Add a node at specific position
	async addNodeAt(nodeId: string, position: { x: number; y: number }) {
		const newNode = await nodeStore.createFlowNode(nodeId, position);
		if (newNode) {
			this.nodes = [...this.nodes, newNode];
		}
		return newNode;
	}
	
	// Update node data
	updateNode(nodeId: string, updates: Record<string, unknown>) {
		this.nodes = this.nodes.map(node => {
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
	}
	
	// Delete selected nodes and edges
	deleteSelected() {
		const selectedNodeIds = new Set(this.selectedNodes.map(n => n.id));
		const selectedEdgeIds = new Set(this.selectedEdges.map(e => e.id));
		
		this.nodes = this.nodes.filter(node => !selectedNodeIds.has(node.id));
		this.edges = this.edges.filter(edge => 
			!selectedEdgeIds.has(edge.id) &&
			!selectedNodeIds.has(edge.source) &&
			!selectedNodeIds.has(edge.target)
		);
		
		this.selectedNodes = [];
		this.selectedEdges = [];
	}
	
	// Clear workflow
	clear() {
		this.nodes = [];
		this.edges = [];
		this.selectedNodes = [];
		this.selectedEdges = [];
		this.recipeName = '';
		this.recipeDescription = '';
	}
	
	// Load workflow from saved data
	loadWorkflow(data: {
		nodes: Node[];
		edges: Edge[];
		name?: string;
		description?: string;
	}) {
		this.nodes = data.nodes;
		this.edges = data.edges;
		this.recipeName = data.name || '';
		this.recipeDescription = data.description || '';
	}
	
	// Get workflow as JSON
	toJSON() {
		return {
			id: `recipe-${Date.now()}`,
			name: this.recipeName,
			description: this.recipeDescription,
			nodes: this.nodes,
			edges: this.edges,
			createdAt: new Date().toISOString()
		};
	}
	
	// Save workflow to localStorage
	saveToLocal() {
		const recipe = this.toJSON();
		const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
		savedRecipes.push(recipe);
		localStorage.setItem('teatime-recipes', JSON.stringify(savedRecipes));
		return recipe;
	}
	
	// Create example workflow
	async createExampleWorkflow() {
		this.clear();
		
		const triggers = nodeStore.triggerNodes;
		const actions = nodeStore.actionNodes;
		
		let xPos = 50;
		const yPos = 150;
		const spacing = 250;
		const newNodes: Node[] = [];
		
		// Add one trigger if available
		if (triggers.length > 0) {
			const node = await this.addNodeAt(triggers[0].Id, { x: xPos, y: yPos });
			if (node) {
				newNodes.push(node);
				xPos += spacing;
			}
		}
		
		// Add a few action nodes if available
		for (let i = 0; i < Math.min(3, actions.length); i++) {
			const node = await this.addNodeAt(actions[i].Id, { x: xPos, y: yPos });
			if (node) {
				newNodes.push(node);
				xPos += spacing;
			}
		}
		
		// Create edges to connect the nodes
		const newEdges: Edge[] = [];
		for (let i = 0; i < newNodes.length - 1; i++) {
			newEdges.push({
				id: `edge-${i}`,
				source: newNodes[i].id,
				target: newNodes[i + 1].id,
				type: 'smoothstep'
			});
		}
		
		this.edges = newEdges;
		this.recipeName = 'Example Workflow';
		this.recipeDescription = 'Auto-generated example workflow';
	}
}

// Create singleton instance
export const workflowStore = new WorkflowStore();