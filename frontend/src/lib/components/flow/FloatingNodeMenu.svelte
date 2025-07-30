<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Plus, Zap, GitBranch, Play, Settings } from 'lucide-svelte';
	import NodeSelectionDialog from './NodeSelectionDialog.svelte';

	const {
		addNode
	}: {
		addNode: (nodeType: string) => void;
	} = $props();

	let showDialog = $state(false);
	let selectedCategory = $state<'trigger' | 'branch' | 'action' | 'utility' | null>(null);

	function openDialog(category: 'trigger' | 'branch' | 'action' | 'utility') {
		selectedCategory = category;
		showDialog = true;
	}

	function handleNodeSelect(nodeType: string) {
		addNode(nodeType);
		showDialog = false;
	}

	function closeDialog() {
		showDialog = false;
	}
</script>

<!-- Floating Menu -->
<div class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 transform">
	<div class="bg-card flex items-center gap-2 rounded-full border p-2 shadow-lg">
		<!-- Trigger Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-12 w-12 rounded-full hover:bg-blue-100 hover:text-blue-600"
			onclick={() => openDialog('trigger')}
		>
			<Zap class="h-5 w-5" />
		</Button>

		<!-- Branch Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-12 w-12 rounded-full hover:bg-yellow-100 hover:text-yellow-600"
			onclick={() => openDialog('branch')}
		>
			<GitBranch class="h-5 w-5" />
		</Button>

		<!-- Action Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-12 w-12 rounded-full hover:bg-green-100 hover:text-green-600"
			onclick={() => openDialog('action')}
		>
			<Play class="h-5 w-5" />
		</Button>

		<!-- Utility Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-12 w-12 rounded-full hover:bg-purple-100 hover:text-purple-600"
			onclick={() => openDialog('utility')}
		>
			<Settings class="h-5 w-5" />
		</Button>

		<!-- Add Generic Node Button -->
		<div class="bg-border mx-1 h-8 w-px"></div>
		<Button
			variant="ghost"
			size="icon"
			class="hover:bg-primary hover:text-primary-foreground h-12 w-12 rounded-full"
			onclick={() => addNode('step')}
		>
			<Plus class="h-5 w-5" />
		</Button>
	</div>
</div>

<!-- Node Selection Dialog -->
<NodeSelectionDialog
	open={showDialog}
	category={selectedCategory}
	onNodeSelect={handleNodeSelect}
	onClose={closeDialog}
/>

<style>
</style>
