import { RecipesService } from '$bindings/services';
import { Recipe, NodeData } from '$bindings/internal/recipe/models';
import type { EditSession } from '$bindings/services/models';
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
    private _editSession = $state<EditSession | null>(null);
    
    error = $state<string | null>(null);

    constructor(recipeId: string) {
        this.recipeId = recipeId;
        this.load();
    }

    get isLoading() {
        return this._loading;
    }

    get hasChanges() {
        return this._editSession?.needs_save ?? false;
    }

    get editSession() {
        return this._editSession;
    }
	
    async load(): Promise<boolean> {
        return await this.withLoading(async () => {
            const editSession = await RecipesService.GetRecipe(this.recipeId);
            if (!editSession || !editSession.recipe) {
                throw new Error('Recipe not found');
            }
            
            this._editSession = editSession;
            const recipe = editSession.recipe;
            
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
            // Update the current edit session with latest recipe data
            const data = this.serializeRecipeData();
            await RecipesService.UpdateRecipe(this.recipeId, Recipe.createFrom(data));
            
            // Save the session
            const savedSession = await RecipesService.SaveRecipe(this.recipeId);
            if (!savedSession) {
                throw new Error('Failed to save recipe');
            }
            
            this._editSession = savedSession;
            return true;
        }) ?? false;
    }

    async createNodeByRef(nodeRef: string, x: number, y: number): Promise<boolean> {
        return await this.withLoading(async () => {
            const newNode = await RecipesService.CreateNode(this.recipeId, nodeRef, x, y);
            if (!newNode) {
                throw new Error('Failed to create node');
            }
            
            this.nodes = [...this.nodes, newNode as CustomNode];
            // No need to call markChanged() - the backend will handle session updates
            return true;
        }) ?? false;
    }


    // Update node data - now calls backend API
    async updateNodeData(nodeId: string, label: string, properties: Record<string, any>, x?: number, y?: number): Promise<boolean> {
        const node = this.getNode(nodeId);
        if (!node) return false;

        const currentX = x ?? node.position.x;
        const currentY = y ?? node.position.y;
        
        const updatedNode = await RecipesService.UpdateNode(
            this.recipeId, 
            nodeId, 
            currentX, 
            currentY, 
            label, 
            properties
        );
        
        if (!updatedNode) {
            throw new Error('Failed to update node');
        }
        
        // Update local state
        const nodeIndex = this.findNodeIndex(nodeId);
        if (nodeIndex !== -1) {
            // this.nodes[nodeIndex] = updatedNode as CustomNode;
            this.nodes = this.nodes.map(node => node.id === nodeId ? updatedNode as CustomNode : node);
        }
        return true;
    }

    async deleteNode(nodeId: string): Promise<boolean> {
        return await this.withLoading(async () => {
            await RecipesService.DeleteNode(this.recipeId, nodeId);
            
            // Update local state
            const initialLength = this.nodes.length;
            this.nodes = this.nodes.filter(node => node.id !== nodeId);
            
            return this.nodes.length < initialLength;
        }) ?? false;
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
		this._editSession = null;
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

    // markChanged is no longer needed - EditSession handles this automatically

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
