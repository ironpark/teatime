<script lang="ts">
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogHeader,
		DialogTitle
	} from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import {
		Zap,
		GitBranch,
		Play,
		Settings,
		Workflow,
		Package
	} from 'lucide-svelte';
	import LucideIcon from '$lib/components/LucideIcon.svelte';
	import type { NodeInfo } from '$bindings/internal/node';

	interface Props {
		open?: boolean;
		category?: 'trigger' | 'branch' | 'action' | 'utility' | null;
		availableNodes?: NodeInfo[];
		onNodeSelect: (nodeRef: string) => void;
		onClose: () => void;
	}

	let { open = false, category = null, availableNodes = [], onNodeSelect, onClose }: Props = $props();
	
	// Build dynamic categories from available nodes
	let nodeCategories = $derived(() => {
		const categories = {
			trigger: {
				title: 'Triggers',
				description: 'Start your workflow with these trigger nodes',
				icon: Zap,
				color: 'text-blue-600',
				bgColor: 'bg-blue-50',
				nodes: [] as any[]
			},
			branch: {
				title: 'Branches',
				description: 'Control flow with conditional branches',
				icon: GitBranch,
				color: 'text-yellow-600',
				bgColor: 'bg-yellow-50',
				nodes: [] as any[]
			},
			action: {
				title: 'Actions',
				description: 'Perform actions and transformations',
				icon: Play,
				color: 'text-green-600',
				bgColor: 'bg-green-50',
				nodes: [] as any[]
			},
			utility: {
				title: 'Utilities',
				description: 'Helper nodes for data processing and workflow management',
				icon: Settings,
				color: 'text-purple-600',
				bgColor: 'bg-purple-50',
				nodes: [] as any[]
			}
		};
		
		// Categorize nodes
		availableNodes.forEach(node => {
			const type = String(node.type).toLowerCase();
			let categoryKey: 'trigger' | 'branch' | 'action' | 'utility' = 'utility';
			let borderColor = 'border-purple-200';
			
			if (type.includes('trigger')) {
				categoryKey = 'trigger';
				borderColor = 'border-blue-200';
			} else if (type.includes('branch') || type.includes('conditional') || type.includes('loop') || type.includes('switch')) {
				categoryKey = 'branch';
				borderColor = 'border-yellow-200';
			} else if (type.includes('action')) {
				categoryKey = 'action';
				borderColor = 'border-green-200';
			}
			console.log(node);
			categories[categoryKey].nodes.push({
				ref: node.ref,
				title: node.name,
				description: node.description || `${node.type} node`,
				color: borderColor,
				iconString: node.icon
			});
		});
		
		return categories;
	});

	let currentCategory = $derived(category ? nodeCategories()[category] : null);
</script>

<Dialog {open} onOpenChange={(isOpen) => !isOpen && onClose()}>
	<DialogContent class="min-w-[500px] max-w-2xl">
		{#if currentCategory}
			<DialogHeader>
				<DialogTitle class="flex items-center gap-2">
					<div class={`rounded-lg p-2 ${currentCategory.bgColor}`}>
						<currentCategory.icon class={`h-5 w-5 ${currentCategory.color}`} />
					</div>
					{currentCategory.title}
				</DialogTitle>
				<DialogDescription>
					{currentCategory.description}
				</DialogDescription>
			</DialogHeader>

			<div class="grid grid-cols-1 gap-4 py-4 md:grid-cols-2">
				{#each currentCategory.nodes as node}
					<Button
						variant="outline"
						class={`h-auto justify-start p-4 transition-all hover:shadow-md ${node.color}`}
						onclick={() => onNodeSelect(node.ref)}
					>
						<div class="flex w-full items-start gap-3 text-left">
							<div class="bg-muted rounded-lg p-2">
								{#if node.iconString}
									<LucideIcon name={node.iconString as any} class="h-4 w-4" />
								{:else}
									<Package class="h-4 w-4" />
								{/if}
							</div>
							<div class="flex-1">
								<div class="text-sm font-medium">{node.title}</div>
								<div class="text-muted-foreground mt-1 text-xs">
									{node.description}
								</div>
							</div>
						</div>
					</Button>
				{/each}
			</div>

		{/if}
	</DialogContent>
</Dialog>
