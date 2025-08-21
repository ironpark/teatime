import type { Node, Edge } from '@xyflow/svelte';
import type { NodeProperty } from '$bindings/internal/node';

export interface BindingOption {
	nodeId: string;
	nodeLabel: string;
	type: 'input' | 'output';
	key: string;
	propertyName: string;
	propertyType: number;
	bindingValue: string; // @[nodeid.(input|output).key] format
}

/**
 * Parses a binding value to extract node ID, type, and key
 * @param bindingValue - Format: "@[nodeid.(input|output).key]"
 * @returns Parsed binding info or null if invalid
 */
export function parseBindingValue(bindingValue: string): {
	nodeId: string;
	type: 'input' | 'output';
	key: string;
} | null {
	const match = bindingValue.match(/^@\[([^.]+)\.(input|output)\.(.+)\]$/);
	if (!match) return null;
	
	return {
		nodeId: match[1],
		type: match[2] as 'input' | 'output',
		key: match[3]
	};
}

/**
 * Creates a binding value in the required format
 */
export function createBindingValue(nodeId: string, type: 'input' | 'output', key: string): string {
	return `@[${nodeId}.${type}.${key}]`;
}

/**
 * Checks if a value is a binding
 */
export function isBinding(value: any): boolean {
	return typeof value === 'string' && value.startsWith('@[') && value.endsWith(']');
}

/**
 * Gets all upstream nodes from a given node by traversing edges
 */
export function getUpstreamNodes(nodeId: string, nodes: Node[], edges: Edge[]): Node[] {
	const upstream = new Set<string>();
	const toVisit = [nodeId];
	
	while (toVisit.length > 0) {
		const currentId = toVisit.shift()!;
		
		// Find all edges that target the current node
		const incomingEdges = edges.filter(edge => edge.target === currentId);
		
		for (const edge of incomingEdges) {
			if (!upstream.has(edge.source)) {
				upstream.add(edge.source);
				toVisit.push(edge.source);
			}
		}
	}
	
	// Return nodes excluding the current node itself
	return nodes.filter(node => upstream.has(node.id));
}

/**
 * Gets all available binding options for a property type from upstream nodes
 */
export function getBindingOptions(
	currentNodeId: string,
	propertyType: number,
	nodes: Node[],
	edges: Edge[]
): BindingOption[] {
	const upstreamNodes = getUpstreamNodes(currentNodeId, nodes, edges);
	const options: BindingOption[] = [];
	
	for (const node of upstreamNodes) {
		const nodeData = node.data;
		
		// Check input properties
		if (nodeData.properties && Array.isArray(nodeData.properties)) {
			for (const prop of nodeData.properties as NodeProperty[]) {
				if (prop.type === propertyType) {
					options.push({
						nodeId: node.id,
						nodeLabel: (nodeData as any).label || (nodeData as any).name || node.id,
						type: 'input',
						key: prop.key,
						propertyName: prop.name || prop.key,
						propertyType: prop.type,
						bindingValue: createBindingValue(node.id, 'input', prop.key)
					});
				}
			}
		}
		
		// Check output properties
		if (nodeData.outputs && Array.isArray(nodeData.outputs)) {
			for (const output of nodeData.outputs as NodeProperty[]) {
				if (output.type === propertyType) {
					options.push({
						nodeId: node.id,
						nodeLabel: (nodeData as any).label || (nodeData as any).name || node.id,
						type: 'output',
						key: output.key,
						propertyName: output.name || output.key,
						propertyType: output.type,
						bindingValue: createBindingValue(node.id, 'output', output.key)
					});
				}
			}
		}
	}
	
	return options;
}

/**
 * Gets display information for a binding value
 */
export function getBindingDisplay(bindingValue: string, nodes: Node[]): {
	nodeLabel: string;
	propertyName: string;
	type: 'input' | 'output';
} | null {
	const parsed = parseBindingValue(bindingValue);
	if (!parsed) return null;
	
	const node = nodes.find(n => n.id === parsed.nodeId);
	if (!node) return null;
	
	const nodeData = node.data as any;
	const nodeLabel = nodeData.label || nodeData.name || node.id;
	
	// Find the property name
	let propertyName = parsed.key;
	
	if (parsed.type === 'input' && nodeData.properties) {
		const prop = nodeData.properties.find((p: NodeProperty) => p.key === parsed.key);
		if (prop) propertyName = prop.name || prop.key;
	} else if (parsed.type === 'output' && nodeData.outputs) {
		const output = nodeData.outputs.find((o: NodeProperty) => o.key === parsed.key);
		if (output) propertyName = output.name || output.key;
	}
	
	return {
		nodeLabel,
		propertyName,
		type: parsed.type
	};
}