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
		Activity,
		Terminal,
		Plug,
		Hand,
		FileImage,
		Filter,
		Shuffle,
		Music,
		FileText,
		Languages,
		Bot,
		Save,
		Variable,
		Hash,
		Type,
		Calculator,
		BarChart3,
		Package
	} from 'lucide-svelte';
	import type { NodeInfo } from '$bindings/stores/models';

	interface Props {
		open?: boolean;
		category?: 'trigger' | 'branch' | 'action' | 'utility' | null;
		availableNodes?: NodeInfo[];
		onNodeSelect: (nodeId: string) => void;
		onClose: () => void;
	}

	let { open = false, category = null, availableNodes = [], onNodeSelect, onClose }: Props = $props();
	
	// Helper function to get icon based on node type
	function getNodeIcon(nodeType: string, nodeName: string): any {
		const type = nodeType.toLowerCase();
		const name = nodeName.toLowerCase();
		
		// Trigger icons
		if (name.includes('command') || name.includes('cli')) return Terminal;
		if (name.includes('webhook')) return Plug;
		if (name.includes('manual')) return Hand;
		if (name.includes('file')) return FileImage;
		
		// Branch icons
		if (name.includes('condition')) return Filter;
		if (name.includes('switch')) return Shuffle;
		if (name.includes('loop')) return GitBranch;
		
		// Action icons
		if (name.includes('llm') || name.includes('ai')) return Bot;
		if (name.includes('save')) return Save;
		if (name.includes('translate')) return Languages;
		if (name.includes('transcribe')) return FileText;
		if (name.includes('audio')) return Music;
		
		// Utility icons
		if (name.includes('variable')) return Variable;
		if (name.includes('constant')) return Hash;
		if (name.includes('format')) return Type;
		if (name.includes('calc')) return Calculator;
		if (name.includes('aggregate')) return BarChart3;
		
		// Default icons by category
		if (type.includes('trigger')) return Zap;
		if (type.includes('branch')) return GitBranch;
		if (type.includes('action')) return Activity;
		
		return Package;
	}
	
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
			const type = node.Type.toLowerCase();
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
			
			categories[categoryKey].nodes.push({
				id: node.Id,
				title: node.Name,
				description: node.Description || `${node.Type} node`,
				icon: getNodeIcon(node.Type, node.Name),
				color: borderColor
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
						onclick={() => onNodeSelect(node.id)}
					>
						<div class="flex w-full items-start gap-3 text-left">
							<div class="bg-muted rounded-lg p-2">
								<node.icon class="h-4 w-4" />
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

			<!-- Quick Actions -->
			{#if availableNodes.length > 0}
				<div class="border-t pt-4">
					<div class="mb-3 flex items-center gap-2 text-sm font-medium">
						<Workflow class="h-4 w-4" />
						Quick Add
					</div>
					<div class="flex gap-2 flex-wrap">
						{#each availableNodes.slice(0, 3) as node}
							<Badge
								variant="outline"
								class="hover:bg-primary hover:text-primary-foreground cursor-pointer"
								onclick={() => onNodeSelect(node.Id)}
							>
								+ {node.Name}
							</Badge>
						{/each}
					</div>
				</div>
			{/if}
		{/if}
	</DialogContent>
</Dialog>
