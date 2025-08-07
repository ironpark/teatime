<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Badge } from '$lib/components/ui/badge';
  import { Input } from '$lib/components/ui/input';
  import { SidebarTrigger } from '$lib/components/ui/sidebar';
  import { Separator } from '$lib/components/ui/separator';
  import { 
    Plus, 
    Search, 
    Clock, 
    Users, 
    ChefHat, 
    Edit, 
    Trash2,
    Play,
    Copy
  } from 'lucide-svelte';
  import { goto } from '$app/navigation';

  type Recipe = {
    id?: string;
    name: string;
    description: string;
    nodes: any[];
    edges: any[];
    createdAt: string;
    updatedAt?: string;
    estimatedTime?: string;
    difficulty?: 'easy' | 'medium' | 'hard';
    servings?: number;
  };

  let recipes: Recipe[] = $state([]);
  let searchTerm = $state('');
  let loading = $state(true);

  onMount(() => {
    loadRecipes();
  });

  function loadRecipes() {
    try {
      const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
      recipes = savedRecipes.map((recipe: any, index: number) => ({
        ...recipe,
        id: recipe.id || `recipe-${index}`,
        estimatedTime: calculateEstimatedTime(recipe.nodes),
        difficulty: calculateDifficulty(recipe.nodes),
        servings: extractServings(recipe.nodes)
      }));
    } catch (error) {
      console.error('Failed to load recipes:', error);
      recipes = [];
    } finally {
      loading = false;
    }
  }

  function calculateEstimatedTime(nodes: any[]): string {
    const steps = nodes.filter(node => node.type === 'step');
    let totalMinutes = 0;
    
    steps.forEach(step => {
      if (step.data.duration) {
        const match = step.data.duration.match(/(\d+)/);
        if (match) {
          totalMinutes += parseInt(match[1]);
        }
      }
    });
    
    if (totalMinutes === 0) return 'Unknown';
    if (totalMinutes < 60) return `${totalMinutes} min`;
    
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;
    return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`;
  }

  function calculateDifficulty(nodes: any[]): 'easy' | 'medium' | 'hard' {
    const stepCount = nodes.filter(node => node.type === 'step').length;
    const ingredientCount = nodes.filter(node => node.type === 'ingredient').length;
    
    const complexity = stepCount + (ingredientCount / 3);
    
    if (complexity <= 3) return 'easy';
    if (complexity <= 6) return 'medium';
    return 'hard';
  }

  function extractServings(nodes: any[]): number {
    const output = nodes.find(node => node.type === 'output');
    return output?.data.servings || 4;
  }

  function createNewRecipe() {
    goto('/recipes/new');
  }

  function editRecipe(recipe: Recipe) {
    goto(`/recipes/${recipe.id}`);
  }

  function deleteRecipe(recipe: Recipe) {
    if (confirm(`Are you sure you want to delete "${recipe.name}"?`)) {
      const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
      const filteredRecipes = savedRecipes.filter((r: any, index: number) => 
        (r.id || `recipe-${index}`) !== recipe.id
      );
      localStorage.setItem('teatime-recipes', JSON.stringify(filteredRecipes));
      loadRecipes();
    }
  }

  function duplicateRecipe(recipe: Recipe) {
    const newRecipe = {
      ...recipe,
      id: `recipe-${Date.now()}`,
      name: `${recipe.name} (Copy)`,
      createdAt: new Date().toISOString()
    };
    
    const savedRecipes = JSON.parse(localStorage.getItem('teatime-recipes') || '[]');
    savedRecipes.push(newRecipe);
    localStorage.setItem('teatime-recipes', JSON.stringify(savedRecipes));
    loadRecipes();
  }

  function executeRecipe(recipe: Recipe) {
    // Store recipe for execution
    sessionStorage.setItem('recipe-to-execute', JSON.stringify(recipe));
    goto('/execution');
  }

  let filteredRecipes = $derived(
    recipes.filter(recipe =>
      recipe.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      recipe.description.toLowerCase().includes(searchTerm.toLowerCase())
    )
  );

  const difficultyColors = {
    easy: 'bg-green-100 text-green-800 border-green-200',
    medium: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    hard: 'bg-red-100 text-red-800 border-red-200'
  };
</script>

<svelte:head>
  <title>Recipes - Teatime</title>
</svelte:head>

<div class="flex flex-col h-full">
  <header class="flex h-16 items-center gap-2 border-b px-4">
    <SidebarTrigger />
    <Separator orientation="vertical" class="mr-2 h-4" />
    <h1 class="text-lg font-semibold">Recipe Library</h1>
  </header>

  <main class="flex-1 p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold tracking-tight">Recipe Library</h2>
        <p class="text-muted-foreground">
          Manage and organize your automation recipes
        </p>
      </div>
      
      <Button onclick={createNewRecipe} class="gap-2">
        <Plus class="w-4 h-4" />
        New Recipe
      </Button>
    </div>

    <!-- Search and filters -->
    <div class="flex items-center gap-4">
      <div class="relative flex-1 max-w-md">
        <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          bind:value={searchTerm}
          placeholder="Search recipes..."
          class="pl-10"
        />
      </div>
      
      <div class="flex items-center gap-2 text-sm text-muted-foreground">
        <span>{filteredRecipes.length} recipe{filteredRecipes.length === 1 ? '' : 's'}</span>
      </div>
    </div>

    <!-- Recipe grid -->
    {#if loading}
      <div class="flex items-center justify-center h-64">
        <div class="text-muted-foreground">Loading recipes...</div>
      </div>
    {:else if filteredRecipes.length === 0}
      <div class="text-center py-12">
        <ChefHat class="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 class="text-lg font-semibold mb-2">
          {searchTerm ? 'No recipes found' : 'No recipes yet'}
        </h3>
        <p class="text-muted-foreground mb-4">
          {searchTerm 
            ? 'Try adjusting your search terms'
            : 'Get started by creating your first recipe'
          }
        </p>
        {#if !searchTerm}
          <Button onclick={createNewRecipe} class="gap-2">
            <Plus class="w-4 h-4" />
            Create Your First Recipe
          </Button>
        {/if}
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each filteredRecipes as recipe (recipe.id)}
          <Card class="group hover:shadow-lg transition-all duration-200">
            <CardHeader class="pb-3">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <CardTitle class="text-lg mb-1">{recipe.name}</CardTitle>
                  <CardDescription class="line-clamp-2">
                    {recipe.description}
                  </CardDescription>
                </div>
                
                <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-8 w-8"
                    onclick={() => editRecipe(recipe)}
                  >
                    <Edit class="w-4 h-4" />
                  </Button>
                  
                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-8 w-8"
                    onclick={() => duplicateRecipe(recipe)}
                  >
                    <Copy class="w-4 h-4" />
                  </Button>
                  
                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-8 w-8 hover:bg-destructive hover:text-destructive-foreground"
                    onclick={() => deleteRecipe(recipe)}
                  >
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </CardHeader>
            
            <CardContent>
              <div class="space-y-3">
                <!-- Recipe stats -->
                <div class="flex flex-wrap gap-2">
                  <Badge variant="outline" class="text-xs">
                    <Clock class="w-3 h-3 mr-1" />
                    {recipe.estimatedTime}
                  </Badge>
                  
                  <Badge variant="outline" class="text-xs">
                    <Users class="w-3 h-3 mr-1" />
                    {recipe.servings} servings
                  </Badge>
                  
                  <Badge variant="outline" class={`text-xs ${difficultyColors[recipe.difficulty!]}`}>
                    {recipe.difficulty}
                  </Badge>
                </div>

                <!-- Node count info -->
                <div class="text-xs text-muted-foreground">
                  {recipe.nodes.filter(n => n.type === 'ingredient').length} ingredients • 
                  {recipe.nodes.filter(n => n.type === 'step').length} steps
                </div>

                <!-- Actions -->
                <div class="flex gap-2 pt-2">
                  <Button
                    variant="default"
                    size="sm"
                    class="flex-1 gap-2"
                    onclick={() => executeRecipe(recipe)}
                  >
                    <Play class="w-3 h-3" />
                    Execute
                  </Button>
                  
                  <Button
                    variant="outline"
                    size="sm"
                    class="gap-2"
                    onclick={() => editRecipe(recipe)}
                  >
                    <Edit class="w-3 h-3" />
                    Edit
                  </Button>
                </div>

                <!-- Creation date -->
                <div class="text-xs text-muted-foreground pt-2 border-t">
                  Created {new Date(recipe.createdAt).toLocaleDateString()}
                </div>
              </div>
            </CardContent>
          </Card>
        {/each}
      </div>
    {/if}
  </main>
</div>

<style>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
