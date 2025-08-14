import {RecipesService} from '$bindings/services';
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
		const data = {
			...node,
			type: categoryType
		};
		return data as unknown as Node;
	}
	
	// Determine visual category from node type
	getCategoryType(nodeType: string) {
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

export const nodeStore = new NodeStore();