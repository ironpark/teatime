<script lang="ts">
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Separator } from '$lib/components/ui/separator';
  import { ChefHat, Utensils, CheckCircle, Plus } from 'lucide-svelte';

  interface Props {
    addNode: (nodeType: string) => void;
  }

  let { addNode }: Props = $props();

  const nodeCategories = [
    {
      title: 'Ingredients',
      icon: ChefHat,
      items: [
        { type: 'ingredient', label: 'Ingredient', description: 'Add recipe ingredients', color: 'text-red-600' }
      ]
    },
    {
      title: 'Processing',
      icon: Utensils,
      items: [
        { type: 'step', label: 'Recipe Step', description: 'Cooking or preparation step', color: 'text-blue-600' }
      ]
    },
    {
      title: 'Output',
      icon: CheckCircle,
      items: [
        { type: 'output', label: 'Final Result', description: 'Recipe end result', color: 'text-green-600' }
      ]
    }
  ];

  function handleAddNode(nodeType: string) {
    addNode(nodeType);
  }
</script>

<div class="h-full p-4 overflow-y-auto">
  <div class="space-y-6">
    <div>
      <h2 class="text-lg font-semibold mb-2">Node Palette</h2>
      <p class="text-sm text-muted-foreground">
        Drag and drop or click to add components to your recipe
      </p>
    </div>

    {#each nodeCategories as category}
      <Card>
        <CardHeader class="pb-3">
          <CardTitle class="flex items-center gap-2 text-base">
            <category.icon class="w-4 h-4" />
            {category.title}
          </CardTitle>
        </CardHeader>
        
        <CardContent class="pt-0">
          <div class="space-y-2">
            {#each category.items as item}
              <Button
                variant="ghost"
                class="w-full justify-start h-auto p-3 hover:bg-muted/50"
                onclick={() => handleAddNode(item.type)}
              >
                <div class="flex items-start gap-3 w-full">
                  <Plus class={`w-4 h-4 mt-0.5 ${item.color}`} />
                  <div class="text-left flex-1">
                    <div class="font-medium text-sm">{item.label}</div>
                    <div class="text-xs text-muted-foreground">{item.description}</div>
                  </div>
                </div>
              </Button>
            {/each}
          </div>
        </CardContent>
      </Card>
      
      {#if category !== nodeCategories[nodeCategories.length - 1]}
        <Separator />
      {/if}
    {/each}

    <!-- Quick Actions -->
    <Card>
      <CardHeader class="pb-3">
        <CardTitle class="text-base">Quick Actions</CardTitle>
      </CardHeader>
      
      <CardContent class="pt-0">
        <div class="space-y-2">
          <Button
            variant="outline"
            size="sm"
            class="w-full justify-start"
            onclick={() => {
              handleAddNode('ingredient');
              handleAddNode('step');
              handleAddNode('output');
            }}
          >
            <Plus class="w-4 h-4 mr-2" />
            Add Basic Flow
          </Button>
          
          <Button
            variant="outline"
            size="sm"
            class="w-full justify-start"
            onclick={() => {
              for (let i = 0; i < 3; i++) {
                handleAddNode('ingredient');
              }
              handleAddNode('step');
              handleAddNode('output');
            }}
          >
            <ChefHat class="w-4 h-4 mr-2" />
            Add Recipe Template
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</div>

<style>
  /* Custom scrollbar */
  .overflow-y-auto {
    scrollbar-width: thin;
    scrollbar-color: hsl(var(--border)) transparent;
  }
  
  .overflow-y-auto::-webkit-scrollbar {
    width: 6px;
  }
  
  .overflow-y-auto::-webkit-scrollbar-track {
    background: transparent;
  }
  
  .overflow-y-auto::-webkit-scrollbar-thumb {
    background-color: hsl(var(--border));
    border-radius: 3px;
  }
  
  .overflow-y-auto::-webkit-scrollbar-thumb:hover {
    background-color: hsl(var(--border) / 0.8);
  }
</style>