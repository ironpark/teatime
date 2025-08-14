import { RecipesService } from '$bindings/services';
import { Recipe, NodeData } from '$bindings/internal/recipe/models';
import type { Edge, Node } from '@xyflow/svelte';

// Type definitions
interface RecipeInfo {
    name: string;
    path: string;
    description: string;
}

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
	info = $state<RecipeInfo>({
		name: '',
		path: '',
		description: ''
	});
    private _loading = $state(false);
    private _hasChanges = $state(false);
    
    error = $state<string | null>(null);

    constructor(recipeId: string) {
        this.recipeId = recipeId;
        this.load();
    }

    get isLoading() {
        return this._loading;
    }

    get hasChanges() {
        return this._hasChanges;
    }
	
    async load(): Promise<boolean> {
        return await this.withLoading(async () => {
            const recipe = await RecipesService.GetRecipe(this.recipeId);
            if (!recipe) {
                throw new Error('Recipe not found');
            }
            
            this.info = {
                name: recipe.name,
                path: recipe.path,
                description: recipe.description
            };
            this.nodes = recipe.nodes as CustomNode[];
            this.edges = recipe.edges;
            return true;
        }) ?? false;
    }
	
    async save(): Promise<boolean> {
        return await this.withLoading(async () => {
            const data = this.serializeRecipeData();
            const recipe = await RecipesService.SaveRecipe(Recipe.createFrom(data));
            if (!recipe) {
                throw new Error('Failed to save recipe');
            }
            
            this._hasChanges = false;
            return true;
        }) ?? false;
    }

    async createNodeByRef(nodeId: string, x: number, y: number): Promise<boolean> {
        return await this.withLoading(async () => {
            const newNode = await RecipesService.CreateNode(nodeId, x, y);
            if (!newNode) {
                throw new Error('Failed to create node');
            }
            
            this.nodes = [...this.nodes, newNode as CustomNode];
            this.markChanged();
            return true;
        }) ?? false;
    }


    // Update node data
    updateNodeData(nodeId: string, updates: Record<string, unknown>): boolean {
        const nodeIndex = this.findNodeIndex(nodeId);
        if (nodeIndex === -1) return false;

        this.nodes[nodeIndex] = {
            ...this.nodes[nodeIndex],
            data: { ...this.nodes[nodeIndex].data, ...updates }
        };
        this.markChanged();
        return true;
    }

    deleteNode(nodeId: string): boolean {
        const initialLength = this.nodes.length;
        this.nodes = this.nodes.filter(node => node.id !== nodeId);
        
        if (this.nodes.length < initialLength) {
            this.markChanged();
            return true;
        }
        return false;
    }
	
	// Clear workflow
	clear(): void {
		this.nodes = [];
		this.edges = [];
		this.selectedNodes = [];
		this.selectedEdges = [];
		this.info = {
			name: '',
			path: '',
			description: ''
		};
		this._hasChanges = false;
	}

    // Helper methods
    private findNodeIndex(nodeId: string): number {
        return this.nodes.findIndex(node => node.id === nodeId);
    }

    hasNode(nodeId: string): boolean {
        return this.findNodeIndex(nodeId) !== -1;
    }

    getNode(nodeId: string): CustomNode | undefined {
        return this.nodes.find(node => node.id === nodeId);
    }

    private markChanged(): void {
        this._hasChanges = true;
    }

    private serializeRecipeData() {
        return {
            name: this.info.name,
            path: this.info.path,
            description: this.info.description,
            nodes: this.nodes,
            edges: this.edges
        };
    }
    
    private async withLoading<T>(operation: () => Promise<T>): Promise<T | null> {
        try {
            this._loading = true;
            this.error = null;
            return await operation();
        } catch (error) {
            this.error = error instanceof Error ? error.message : 'Unknown error';
            return null;
        } finally {
            this._loading = false;
        }
    }
}
