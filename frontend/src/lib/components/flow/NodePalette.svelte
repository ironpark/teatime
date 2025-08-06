<script lang="ts">
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Separator } from '$lib/components/ui/separator';
  import { ChefHat, Plus, Zap, Activity, Save } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { nodeStore } from '$lib/stores/nodes.svelte';

  interface Props {
    addNode: (nodeId: string) => void;
  }

  let { addNode }: Props = $props();

  onMount(async () => {
    // Ensure nodes are loaded
    if (nodeStore.availableNodes.length === 0) {
      await nodeStore.loadNodes();
    }
  });

  // Dynamically group nodes by their type
  let nodeCategories = $derived(() => {
    const categories: Record<string, { title: string; icon: unknown; items: Array<{
      type: string;
      label: string;
      description: string;
      color: string;
    }> }> = {};
    
    nodeStore.availableNodes.forEach(node => {
      const nodeType = node.Type.toLowerCase();
      let category = 'Actions';
      let icon = Activity;
      let color = 'text-blue-600';
      
      if (nodeType.includes('trigger')) {
        category = 'Triggers';
        icon = Zap;
        color = 'text-yellow-600';
      } else if (nodeType.includes('branch') || nodeType.includes('conditional') || nodeType.includes('loop')) {
        category = 'Control Flow';
        icon = ChefHat;
        color = 'text-purple-600';
      } else if (nodeType.includes('save') || nodeType.includes('output')) {
        category = 'Output';
        icon = Save;
        color = 'text-green-600';
      } else if (nodeType.includes('action')) {
        category = 'Actions';
        icon = Activity;
        color = 'text-blue-600';
      }
      
      if (!categories[category]) {
        categories[category] = {
          title: category,
          icon,
          items: []
        };
      }
      
      categories[category].items.push({
        type: node.Id,
        label: node.Name,
        description: node.Description || `${node.Type} node`,
        color
      });
    });
    
    return Object.values(categories);
  });

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

    {#if nodeStore.isLoading}
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-center">
            <p class="text-sm text-muted-foreground">Loading nodes...</p>
          </div>
        </CardContent>
      </Card>
    {:else if nodeCategories().length === 0}
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-center">
            <p class="text-sm text-muted-foreground">No nodes available</p>
          </div>
        </CardContent>
      </Card>
    {:else}
      {#each nodeCategories() as category}
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
      
        {#if category !== nodeCategories()[nodeCategories().length - 1]}
          <Separator />
        {/if}
      {/each}
    {/if}

    <!-- Quick Actions -->
    {#if !nodeStore.isLoading && nodeStore.availableNodes.length > 0}
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
                // Add a sample flow with available nodes
                const triggers = nodeStore.triggerNodes;
                const actions = nodeStore.actionNodes;
                
                if (triggers.length > 0) handleAddNode(triggers[0].Id);
                if (actions.length > 0) {
                  handleAddNode(actions[0].Id);
                  if (actions.length > 1) handleAddNode(actions[1].Id);
                }
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
                // Add a more complex template
                const triggers = nodeStore.triggerNodes;
                const actions = nodeStore.actionNodes;
                
                if (triggers.length > 0) handleAddNode(triggers[0].Id);
                actions.slice(0, 4).forEach(action => handleAddNode(action.Id));
              }}
            >
              <Activity class="w-4 h-4 mr-2" />
              Add Complex Workflow
            </Button>
          </div>
        </CardContent>
      </Card>
    {/if}
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