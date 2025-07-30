<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { SvelteFlow, Controls, Background, BackgroundVariant, MiniMap, type Node, type Edge } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';
  
  import RecipeStepNode from '$lib/components/flow/node/RecipeStepNode.svelte';
  import IngredientNode from '$lib/components/flow/node/IngredientNode.svelte';
  import OutputNode from '$lib/components/flow/node/OutputNode.svelte';
  import VariableNode from '$lib/components/flow/node/VariableNode.svelte';
  import ConstantNode from '$lib/components/flow/node/ConstantNode.svelte';
  import FormatterNode from '$lib/components/flow/node/FormatterNode.svelte';
  import CalculatorNode from '$lib/components/flow/node/CalculatorNode.svelte';
  import AggregatorNode from '$lib/components/flow/node/AggregatorNode.svelte';
  import CommandNode from '$lib/components/flow/node/CommandNode.svelte';
  import MCPNode from '$lib/components/flow/node/MCPNode.svelte';
  import ManualNode from '$lib/components/flow/node/ManualNode.svelte';
  import FileInputNode from '$lib/components/flow/node/FileInputNode.svelte';
  import FFmpegExtractNode from '$lib/components/flow/node/FFmpegExtractNode.svelte';
  import TranscriptionNode from '$lib/components/flow/node/TranscriptionNode.svelte';
  import TranslationNode from '$lib/components/flow/node/TranslationNode.svelte';
  import LLMCorrectionNode from '$lib/components/flow/node/LLMCorrectionNode.svelte';
  import FileSaveNode from '$lib/components/flow/node/FileSaveNode.svelte';
  import FloatingNodeMenu from '$lib/components/flow/FloatingNodeMenu.svelte';
  import NodePropertiesSidebar from '$lib/components/flow/NodePropertiesSidebar.svelte';
  import RecipeToolbar from '$lib/components/flow/RecipeToolbar.svelte';
  import { Button } from '$lib/components/ui/button';
  import { SidebarTrigger } from '$lib/components/ui/sidebar';
  import { Separator } from '$lib/components/ui/separator';

  // Node types for different recipe components  
  const nodeTypes: any = {
    ingredient: IngredientNode,
    step: RecipeStepNode,
    output: OutputNode,
    variable: VariableNode,
    constant: ConstantNode,
    formatter: FormatterNode,
    calculator: CalculatorNode,
    aggregator: AggregatorNode,
    command: CommandNode,
    mcp: MCPNode,
    manual: ManualNode,
    'file-input': FileInputNode,
    'ffmpeg-extract': FFmpegExtractNode,
    transcription: TranscriptionNode,
    translation: TranslationNode,
    'llm-correction': LLMCorrectionNode,
    'file-save': FileSaveNode
  };

  let nodes = $state.raw<Node[]>([]);
  let edges = $state.raw<Edge[]>([]);
  let selectedNodes = $state<Node[]>([]);
  let selectedEdges = $state<Edge[]>([]);

  // Recipe metadata
  let recipeName = $state('Loading...');
  let recipeDescription = $state('');
  let recipeId = $state('');
  let originalRecipe: any = null;
  let loading = $state(true);
  let notFound = $state(false);

  onMount(() => {
    loadRecipe();
  });

  function loadRecipe() {
    const id = page.params.id;
    recipeId = id;
    
    try {
      const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
      const recipe = savedRecipes.find((r: any) => r.id === id);
      
      if (recipe) {
        originalRecipe = recipe;
        recipeName = recipe.name;
        recipeDescription = recipe.description;
        nodes = recipe.nodes || [];
        edges = recipe.edges || [];
        loading = false;
      } else {
        notFound = true;
        loading = false;
      }
    } catch (error) {
      console.error('Failed to load recipe:', error);
      notFound = true;
      loading = false;
    }
  }

  function onSelectionChange(selection: { nodes: Node[]; edges: Edge[] }) {
    selectedNodes = selection.nodes;
    selectedEdges = selection.edges;
  }

  function addNode(nodeType: string) {
    const newId = `${Date.now()}`;
    const newNode: Node = {
      id: newId,
      type: nodeType,
      position: { x: Math.random() * 500 + 100, y: Math.random() * 300 + 100 },
      data: getDefaultNodeData(nodeType)
    };
    
    nodes = [...nodes, newNode];
  }

  function getDefaultNodeData(nodeType: string) {
    switch (nodeType) {
      case 'ingredient':
        return { label: 'New Ingredient', amount: '1 unit', category: 'other' };
      case 'step':
        return { 
          label: 'New Step', 
          instruction: 'Add instruction here',
          duration: '1 minute',
          temperature: null
        };
      case 'output':
        return { label: 'New Output', description: 'Describe the result' };
      case 'command':
        return { label: 'Command Line', command: '', workingDir: '' };
      case 'mcp':
        return { label: 'MCP Server', serverUrl: '', method: '', params: '' };
      case 'manual':
        return { label: 'Manual Trigger', description: 'Click to start', requiresConfirmation: false };
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
    const selectedNodeIds = new Set(selectedNodes.map(n => n.id));
    const selectedEdgeIds = new Set(selectedEdges.map(e => e.id));
    
    nodes = nodes.filter(node => !selectedNodeIds.has(node.id));
    edges = edges.filter(edge => 
      !selectedEdgeIds.has(edge.id) && 
      !selectedNodeIds.has(edge.source) && 
      !selectedNodeIds.has(edge.target)
    );
    
    selectedNodes = [];
    selectedEdges = [];
  }

  function saveRecipe() {
    const updatedRecipe = {
      ...originalRecipe,
      id: recipeId,
      name: recipeName,
      description: recipeDescription,
      nodes: nodes,
      edges: edges,
      updatedAt: new Date().toISOString()
    };
    
    // Save to localStorage
    const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
    const index = savedRecipes.findIndex((r: any) => r.id === recipeId);
    
    if (index !== -1) {
      savedRecipes[index] = updatedRecipe;
    } else {
      savedRecipes.push(updatedRecipe);
    }
    
    localStorage.setItem('teatime-recipes', JSON.stringify(savedRecipes));
    
    console.log('Recipe updated:', updatedRecipe);
    
    // Show success message
    alert('Recipe updated successfully!');
  }

  function goBackToList() {
    goto('/recipes');
  }

  function executeRecipe() {
    const currentRecipe = {
      id: recipeId,
      name: recipeName,
      description: recipeDescription,
      nodes: nodes,
      edges: edges
    };
    
    // Store recipe for execution
    sessionStorage.setItem('recipe-to-execute', JSON.stringify(currentRecipe));
    goto('/execution');
  }

  function updateNodeData(nodeId: string, updates: any) {
    nodes = nodes.map(node => {
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

  function closeSidebar() {
    selectedNodes = [];
    selectedEdges = [];
  }
</script>

<svelte:head>
  <title>{recipeName} - Edit Recipe - Teatime</title>
</svelte:head>

{#if loading}
  <div class="h-screen w-full flex items-center justify-center">
    <div class="text-center">
      <div class="text-lg font-medium mb-2">Loading recipe...</div>
      <div class="text-sm text-muted-foreground">Please wait</div>
    </div>
  </div>
{:else if notFound}
  <div class="h-screen w-full flex items-center justify-center">
    <div class="text-center">
      <div class="text-lg font-medium mb-2">Recipe not found</div>
      <div class="text-sm text-muted-foreground mb-4">The recipe you're looking for doesn't exist.</div>
      <Button onclick={goBackToList}>
        ← Back to Recipes
      </Button>
    </div>
  </div>
{:else}
  <div class="recipe-editor h-screen w-full flex flex-col bg-background">
    <!-- Header -->
    <header class="border-b bg-card">
      <div class="flex h-16 items-center gap-2 px-4">
        <SidebarTrigger />
        <Separator orientation="vertical" class="mr-2 h-4" />
        <div class="flex items-center space-x-4 flex-1">
          <div class="flex flex-col">
            <input 
              bind:value={recipeName}
              class="text-lg font-semibold bg-transparent border-none outline-none"
              placeholder="Recipe name"
            />
            <input 
              bind:value={recipeDescription}
              class="text-sm text-muted-foreground bg-transparent border-none outline-none"
              placeholder="Recipe description"
            />
          </div>
        </div>
        
        <div class="flex items-center gap-2">
          <Button variant="outline" onclick={goBackToList}>
            ← Back to Recipes
          </Button>
          <Button variant="default" onclick={executeRecipe}>
            ▶ Execute Recipe
          </Button>
          <RecipeToolbar 
            {selectedNodes}
            {selectedEdges}
            onSave={saveRecipe}
            onDelete={deleteSelected}
          />
        </div>
      </div>
    </header>

    <!-- Main content -->
    <div class="flex-1 flex">
      <!-- Flow editor -->
      <main class="flex-1 relative">
        <SvelteFlow
          bind:nodes
          bind:edges
          {nodeTypes}
          onselectionchange={onSelectionChange}
          fitView
          snapGrid={[15, 15]}
          defaultEdgeOptions={{ type: 'smoothstep' }}
          proOptions={{ hideAttribution: true }}
        >
          <Controls />
          <Background variant={BackgroundVariant.Dots} />
          <MiniMap 
            nodeStrokeWidth={3}
            pannable={true}
            zoomable={true}
          />
        </SvelteFlow>
      </main>
    </div>

    <!-- Floating Node Menu -->
    <FloatingNodeMenu {addNode} />

    <!-- Node Properties Sidebar -->
    <NodePropertiesSidebar 
      {selectedNodes} 
      onNodeUpdate={updateNodeData}
      onClose={closeSidebar}
    />
  </div>
{/if}

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