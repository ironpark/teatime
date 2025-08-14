<script lang="ts">
  import { onMount } from 'svelte';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { SidebarTrigger } from '$lib/components/ui/sidebar';
  import { Separator } from '$lib/components/ui/separator';
  import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '$lib/components/ui/dialog';
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';
  import { 
    Plus, 
    Search, 
    ChefHat
  } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { RecipesService } from '$bindings/services';
  import { RecipeInfo } from '$bindings/services';
  import RecipeCard from '$lib/components/RecipeCard.svelte';

  let recipes: RecipeInfo[] = $state([]);
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
      recipes = await RecipesService.ListRecipes();
      console.log(recipes);
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

  function editRecipe(recipe: RecipeInfo) {
    goto(`/recipes/${recipe.ID}`);
  }

  async function deleteRecipe(recipe: RecipeInfo) {
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

  async function duplicateRecipe(recipe: RecipeInfo) {
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

  async function executeRecipe(recipe: RecipeInfo) {
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
</script>

<svelte:head>
  <title>Teatime</title>
</svelte:head>

<div class="flex flex-col h-full">
  <header class="flex h-16 items-center gap-2 border-b px-4">
    <SidebarTrigger />
    <Separator orientation="vertical" class="mr-2 h-4" />
    <h1 class="text-lg font-semibold">Teatime</h1>
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
          <RecipeCard
            {recipe}
            onEdit={editRecipe}
            onDelete={deleteRecipe}
            onDuplicate={duplicateRecipe}
            onExecute={executeRecipe}
          />
        {/each}
      </div>
    {/if}
  </main>
</div>

