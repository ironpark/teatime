import { RecipesService } from '$bindings/services';
import { Recipe, NodeData } from '$bindings/internal/recipe/models';
import type { Edge, Node } from '@xyflow/svelte';

type CustomNode = Node & {
	data: NodeData;
}
// Recipe store for managing the current workflow
export class RecipeStore {
	nodes = $state.raw<CustomNode[]>([]);
	edges = $state.raw<Edge[]>([]);
	selectedNodes = $state.raw<Node[]>([]);
	selectedEdges = $state.raw<Edge[]>([]);
    recipeId: string;
	// Recipe metadata
	info = $state({
		name: '',
		path: '',
		description: ''
	});
    private loading = $state(false);
    error = $state<string | null>(null);

    constructor(recipeId: string) {
        this.recipeId = recipeId;
        this.load();
    }

    get isLoading() {
        return this.loading;
    }
	
    async load() {
        this.loading = true; // TODO: use a promise instead
        this.error = null;
        const recipe = await RecipesService.GetRecipe(this.recipeId);
        if (!recipe) {
            this.error = 'Failed to load recipe';
            this.loading = false;
            return false;
        }
        this.info = {
            name: recipe.name,
            path: recipe.path,
            description: recipe.description
        };
        this.nodes = recipe.nodes as CustomNode[];
        this.edges = recipe.edges;
        this.loading = false;
        return true;
    }
	
    async save() {
        console.log('Saving recipe', this.nodes, this.edges);
        const data = {
            name: this.info.name,
            path: this.info.path,
            description: this.info.description,
            nodes: this.nodes,
            edges: this.edges
        };
        const recipe = await RecipesService.SaveRecipe(Recipe.createFrom(data));
        if (!recipe) {
            return false;
        }
        return true;
    }

    async createNodeByRef(nodeId: string, x: number, y: number) {
        const newNode = await RecipesService.CreateNode(nodeId, x, y);
        if (!newNode) {
            return false;
        }
        this.nodes = [...this.nodes, newNode as CustomNode];
        return true;
    }

	// Update node data
	updateNode(nodeId: string, updates: Record<string, unknown>) {
        const index = this.nodes.findIndex(node => node.id === nodeId);
        if (index === -1) {
            return;
        }
        this.nodes[index] = {
            ...this.nodes[index],
            data: {
                ...this.nodes[index].data,
                ...updates
            }
        };
    }

    // Update node data
	updateNodeData(nodeId: string, updates: Record<string, unknown>) {
        const index = this.nodes.findIndex(node => node.id === nodeId);
        if (index === -1) {
            return;
        }
        this.nodes[index] = {
            ...this.nodes[index],
            data: {
                ...this.nodes[index].data,
                ...updates
            }
        };
	}

    deleteNode(nodeId: string) {
        this.nodes = this.nodes.filter(node => node.id !== nodeId);
    }
	
	// Clear workflow
	clear() {
		this.nodes = [];
		this.edges = [];
		this.selectedNodes = [];
		this.selectedEdges = [];
		this.info = {
			name: '',
			path: '',
			description: ''
		};
	}	
}
