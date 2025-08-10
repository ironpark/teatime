import {RecipesService, SettingsService} from '$bindings/services';
import type { Node, Edge } from '@xyflow/svelte';
import type { NodeInfo } from '$bindings/internal/node';
import { NodeType } from '$bindings/internal/node/models';

// Store for available node types from backend
class NodeStore {
	availableNodes = $state<NodeInfo[]>([]);
	isLoading = $state(false);
	error = $state<string | null>(null);
	triggerNodes = $derived(this.availableNodes.filter(n => n.type == NodeType.NodeTypeTrigger));
	actionNodes = $derived(this.availableNodes.filter(n => n.type == NodeType.NodeTypeAction));
	branchNodes = $derived(this.availableNodes.filter(n => n.type == NodeType.NodeTypeBranch));
	utilityNodes = $derived(this.availableNodes.filter(n => n.type == NodeType.NodeTypeUtil));
	
	// Load nodes from backend
	async loadNodes() {
		if (this.isLoading) return;
		
		this.isLoading = true;
		this.error = null;
		
		try {
			const nodes = await RecipesService.GetAvailableNodes();
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
		return this.availableNodes.filter(n => n.type == type);
	}
	
	// Create a flow node for the editor
	async createFlowNode(nodeRef: string, position: { x: number; y: number }): Promise<Node | null> {
		const node = await RecipesService.CreateNode(nodeRef, Math.round(position.x), Math.round(position.y));
		if (!node) return null;
		
		const categoryType = this.getCategoryType(node.type);
		
		return {
			id: node.id, // Use the backend-generated ID
			type: categoryType,
			position,
			data: {
				icon: node.icon,
				label: node.name,
				name: node.name,
				nodeType: node.type,
				description: node.description,
				properties: node.properties,
				outputs: node.output,
				backendNodeRef: node.ref
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
	async addNode(nodeRef: string) {
		const position = {
			x: Math.random() * 500 + 100,
			y: Math.random() * 300 + 100
		};
		
		const newNode = await nodeStore.createFlowNode(nodeRef, position);
		console.log('newNode', newNode);
		if (newNode) {
			this.nodes = [...this.nodes, newNode];
		}
		return newNode;
	}
	
	// Add a node at specific position
	async addNodeAt(nodeRef: string, position: { x: number; y: number }) {
		const newNode = await nodeStore.createFlowNode(nodeRef, position);
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
			const node = await this.addNodeAt(triggers[0].ref, { x: xPos, y: yPos });
			if (node) {
				newNodes.push(node);
				xPos += spacing;
			}
		}
		
		// Add a few action nodes if available
		for (let i = 0; i < Math.min(3, actions.length); i++) {
			const node = await this.addNodeAt(actions[i].ref, { x: xPos, y: yPos });
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