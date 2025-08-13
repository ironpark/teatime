<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Badge } from '$lib/components/ui/badge';
  import { Input } from '$lib/components/ui/input';
  import { SidebarTrigger } from '$lib/components/ui/sidebar';
  import { Separator } from '$lib/components/ui/separator';
  import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '$lib/components/ui/dialog';
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';
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
  import { RecipesService } from '$bindings/services';
  import type { Recipe as DBRecipe } from '$bindings/internal/database';

  type Recipe = DBRecipe & {
    estimatedTime?: string;
    difficulty?: 'easy' | 'medium' | 'hard';
    servings?: number;
  };

  let recipes: Recipe[] = $state([]);
  let searchTerm = $state('');
  let loading = $state(true);
  
  // Dialog state
  let dialogOpen = $state(false);
  let newRecipeName = $state('');
  let newRecipeDescription = $state('');
  let creating = $state(false);

  onMount(() => {
    loadRecipes();
  });

  async function loadRecipes() {
    try {
      loading = true;
      const dbRecipes = await RecipesService.ListRecipes();
      recipes = dbRecipes.map((recipe) => ({
        ...recipe,
        estimatedTime: '30 min', // Default for now
        difficulty: 'medium' as const, // Default for now
        servings: 4 // Default for now
      }));
    } catch (error) {
      console.error('Failed to load recipes:', error);
      recipes = [];
    } finally {
      loading = false;
    }
  }


  function createNewRecipe() {
    dialogOpen = true;
  }

  async function handleCreateRecipe() {
    if (!newRecipeName.trim()) return;
    
    try {
      creating = true;
      const result = await RecipesService.CreateRecipe(
        newRecipeName.trim(),
        newRecipeDescription.trim()
      );
      
      if (result) {
        // Reset form and close dialog
        newRecipeName = '';
        newRecipeDescription = '';
        dialogOpen = false;
        
        // Reload recipes to show the new one
        await loadRecipes();
        
        // Navigate to edit the new recipe
        goto(`/recipes/${result.ID}`);
      }
    } catch (error) {
      console.error('Failed to create recipe:', error);
      alert('Failed to create recipe. Please try again.');
    } finally {
      creating = false;
    }
  }

  function cancelCreateRecipe() {
    newRecipeName = '';
    newRecipeDescription = '';
    dialogOpen = false;
  }

  function editRecipe(recipe: Recipe) {
    goto(`/recipes/${recipe.ID}`);
  }

  async function deleteRecipe(recipe: Recipe) {
    if (confirm(`Are you sure you want to delete "${recipe.Name}"?`)) {
      try {
        await RecipesService.DeleteRecipe(recipe.ID);
        await loadRecipes();
      } catch (error) {
        console.error('Failed to delete recipe:', error);
        alert('Failed to delete recipe. Please try again.');
      }
    }
  }

  async function duplicateRecipe(recipe: Recipe) {
    try {
      // Get the full recipe data
      const fullRecipe = await RecipesService.GetRecipe(recipe.ID);
      if (!fullRecipe) {
        alert('Failed to load recipe for duplication');
        return;
      }
      
      // Create a new recipe with copied data
      const result = await RecipesService.CreateRecipe(
        `${recipe.Name} (Copy)`,
        recipe.Description
      );
      
      if (result) {
        // Update the new recipe with the copied workflow data
        fullRecipe.name = `${recipe.Name} (Copy)`;
        await RecipesService.UpdateRecipe(result.ID, fullRecipe);
        await loadRecipes();
      }
    } catch (error) {
      console.error('Failed to duplicate recipe:', error);
      alert('Failed to duplicate recipe. Please try again.');
    }
  }

  async function executeRecipe(recipe: Recipe) {
    try {
      // Get the full recipe data
      const fullRecipe = await RecipesService.GetRecipe(recipe.ID);
      if (fullRecipe) {
        sessionStorage.setItem('recipe-to-execute', JSON.stringify(fullRecipe));
        goto('/execution');
      } else {
        alert('Failed to load recipe for execution');
      }
    } catch (error) {
      console.error('Failed to load recipe:', error);
      alert('Failed to load recipe. Please try again.');
    }
  }

  let filteredRecipes = $derived(
    recipes.filter(recipe =>
      recipe.Name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      recipe.Description.toLowerCase().includes(searchTerm.toLowerCase())
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
      
      <Dialog bind:open={dialogOpen}>
        <DialogTrigger>
          <Button onclick={createNewRecipe} class="gap-2">
            <Plus class="w-4 h-4" />
            New Recipe
          </Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Create New Recipe</DialogTitle>
            <DialogDescription>
              Create a new automation recipe. You can add steps and configure it after creation.
            </DialogDescription>
          </DialogHeader>
          
          <div class="space-y-4 py-4">
            <div class="space-y-2">
              <Label for="recipe-name">Recipe Name</Label>
              <Input
                id="recipe-name"
                bind:value={newRecipeName}
                placeholder="Enter recipe name..."
                disabled={creating}
                onkeydown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleCreateRecipe();
                  }
                }}
              />
            </div>
            
            <div class="space-y-2">
              <Label for="recipe-description">Description (optional)</Label>
              <Textarea
                id="recipe-description"
                bind:value={newRecipeDescription}
                placeholder="Describe what this recipe does..."
                disabled={creating}
                rows={3}
              />
            </div>
          </div>
          
          <div class="flex justify-end gap-2">
            <Button variant="outline" onclick={cancelCreateRecipe} disabled={creating}>
              Cancel
            </Button>
            <Button 
              onclick={handleCreateRecipe} 
              disabled={creating || !newRecipeName.trim()}
              class="gap-2"
            >
              {#if creating}
                <div class="w-4 h-4 animate-spin rounded-full border-2 border-current border-t-transparent"></div>
                Creating...
              {:else}
                <Plus class="w-4 h-4" />
                Create Recipe
              {/if}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
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
        {#each filteredRecipes as recipe (recipe.ID)}
          <Card class="group hover:shadow-lg transition-all duration-200">
            <CardHeader class="pb-3">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <CardTitle class="text-lg mb-1">{recipe.Name}</CardTitle>
                  <CardDescription class="line-clamp-2">
                    {recipe.Description}
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

                <!-- Recipe info -->
                <div class="text-xs text-muted-foreground">
                  Recipe ID: {recipe.ID}
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
                  Created {new Date(recipe.CreatedAt).toLocaleDateString()}
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
