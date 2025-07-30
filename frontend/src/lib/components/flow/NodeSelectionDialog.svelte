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
		ChefHat,
		Utensils,
		CheckCircle,
		Timer,
		Thermometer,
		Scale,
		Workflow,
		Filter,
		Shuffle,
		Settings,
		Variable,
		Hash,
		Type,
		Calculator,
		BarChart3,
		Music,
		FileText,
		Languages,
		Bot,
		Save,
		FileImage,
		Terminal,
		Plug,
		Hand
	} from 'lucide-svelte';

	interface Props {
		open?: boolean;
		category?: 'trigger' | 'branch' | 'action' | 'utility' | null;
		onNodeSelect: (nodeType: string) => void;
		onClose: () => void;
	}

	let { open = false, category = null, onNodeSelect, onClose }: Props = $props();
	const nodeCategories = {
		trigger: {
			title: 'Triggers',
			description: 'Start your workflow with these trigger nodes',
			icon: Zap,
			color: 'text-blue-600',
			bgColor: 'bg-blue-50',
			nodes: [
				{
					type: 'command',
					title: 'Command Line',
					description: 'Execute command line programs',
					icon: Terminal,
					color: 'border-blue-200'
				},
				{
					type: 'mcp',
					title: 'MCP Server',
					description: 'Connect to MCP server resources',
					icon: Plug,
					color: 'border-blue-200'
				},
				{
					type: 'manual',
					title: 'Manual Trigger',
					description: 'Start workflow manually',
					icon: Hand,
					color: 'border-blue-200'
				},
				{
					type: 'file-input',
					title: 'File Input',
					description: 'Load any input file',
					icon: FileImage,
					color: 'border-blue-200'
				}
			]
		},
		branch: {
			title: 'Branches',
			description: 'Control flow with conditional branches',
			icon: GitBranch,
			color: 'text-yellow-600',
			bgColor: 'bg-yellow-50',
			nodes: [
				{
					type: 'condition',
					title: 'Condition',
					description: 'Branch based on conditions',
					icon: Filter,
					color: 'border-yellow-200'
				},
				{
					type: 'switch',
					title: 'Switch',
					description: 'Multiple path branching',
					icon: Shuffle,
					color: 'border-yellow-200'
				}
			]
		},
		action: {
			title: 'Actions',
			description: 'Perform actions and transformations',
			icon: Play,
			color: 'text-green-600',
			bgColor: 'bg-green-50',
			nodes: [
				{
					type: 'ffmpeg-extract',
					title: 'FFmpeg Audio Extract',
					description: 'Extract audio from video using FFmpeg',
					icon: Music,
					color: 'border-green-200'
				},
				{
					type: 'transcription',
					title: 'Speech Transcription',
					description: 'Convert audio to text',
					icon: FileText,
					color: 'border-green-200'
				},
				{
					type: 'translation',
					title: 'Text Translation',
					description: 'Translate text to another language',
					icon: Languages,
					color: 'border-green-200'
				},
				{
					type: 'llm-correction',
					title: 'LLM Text Correction',
					description: 'Improve text quality using AI',
					icon: Bot,
					color: 'border-green-200'
				},
				{
					type: 'file-save',
					title: 'Save to File',
					description: 'Save result to file system',
					icon: Save,
					color: 'border-green-200'
				}
			]
		},
		utility: {
			title: 'Utilities',
			description: 'Helper nodes for data processing and workflow management',
			icon: Settings,
			color: 'text-purple-600',
			bgColor: 'bg-purple-50',
			nodes: [
				{
					type: 'variable',
					title: 'Variable',
					description: 'Store and reference values',
					icon: Variable,
					color: 'border-purple-200'
				},
				{
					type: 'constant',
					title: 'Constant',
					description: 'Fixed constant value',
					icon: Hash,
					color: 'border-purple-200'
				},
				{
					type: 'formatter',
					title: 'Formatter',
					description: 'Format text and dates',
					icon: Type,
					color: 'border-purple-200'
				},
				{
					type: 'calculator',
					title: 'Calculator',
					description: 'Perform mathematical calculations',
					icon: Calculator,
					color: 'border-purple-200'
				},
				{
					type: 'aggregator',
					title: 'Aggregator',
					description: 'Aggregate and summarize data',
					icon: BarChart3,
					color: 'border-purple-200'
				}
			]
		}
	};

	let currentCategory = $derived(category ? nodeCategories[category] : null);
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
						onclick={() => onNodeSelect(node.type)}
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
			<div class="border-t pt-4">
				<div class="mb-3 flex items-center gap-2 text-sm font-medium">
					<Workflow class="h-4 w-4" />
					Quick Add
				</div>
				<div class="flex gap-2">
					<Badge
						variant="outline"
						class="hover:bg-primary hover:text-primary-foreground cursor-pointer"
						onclick={() => onNodeSelect('ingredient')}
					>
						+ Ingredient
					</Badge>
					<Badge
						variant="outline"
						class="hover:bg-primary hover:text-primary-foreground cursor-pointer"
						onclick={() => onNodeSelect('step')}
					>
						+ Step
					</Badge>
					<Badge
						variant="outline"
						class="hover:bg-primary hover:text-primary-foreground cursor-pointer"
						onclick={() => onNodeSelect('output')}
					>
						+ Output
					</Badge>
				</div>
			</div>
		{/if}
	</DialogContent>
</Dialog>
