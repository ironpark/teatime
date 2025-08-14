import {RecipesService} from '$bindings/services';
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
	getNodesByType(type: string): NodeInfo[] {
		return this.availableNodes.filter(node => node.type === type);
	}
}



export const nodeStore = new NodeStore();