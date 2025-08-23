<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Plus, Zap, GitBranch, Play, Settings } from 'lucide-svelte';
	import NodeSelectionDialog from './NodeSelectionDialog.svelte';
	import { onMount } from 'svelte';
	import { nodeStore } from '$lib/stores/nodes.svelte';

	const {
		addNode
	}: {
		addNode: (nodeRef: string) => void;
	} = $props();

	let showDialog = $state(false);
	let selectedCategory = $state<'trigger' | 'branch' | 'action' | 'utility' | null>(null);

	onMount(async () => {
		// Ensure nodes are loaded
		if (nodeStore.availableNodes.length === 0) {
			await nodeStore.loadNodes();
		}
	});

	function openDialog(category: 'trigger' | 'branch' | 'action' | 'utility') {
		selectedCategory = category;
		showDialog = true;
	}

	function handleNodeSelect(nodeRef: string) {
		addNode(nodeRef);
		showDialog = false;
	}

	function closeDialog() {
		showDialog = false;
	}
</script>

<!-- Floating Menu -->
<div class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 transform">
	<div class="bg-card flex items-center gap-2 rounded-full border p-1 shadow-lg">
		<!-- Trigger Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-10 w-10 rounded-full hover:bg-blue-100 hover:text-blue-600"
			onclick={() => openDialog('trigger')}
		>
			<Zap class="h-4 w-4" />
		</Button>

		<!-- Branch Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-10 w-10 rounded-full hover:bg-yellow-100 hover:text-yellow-600"
			onclick={() => openDialog('branch')}
		>
			<GitBranch class="h-4 w-4" />
		</Button>

		<!-- Action Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-10 w-10 rounded-full hover:bg-green-100 hover:text-green-600"
			onclick={() => openDialog('action')}
		>
			<Play class="h-4 w-4" />
		</Button>

		<!-- Utility Button -->
		<Button
			variant="ghost"
			size="icon"
			class="h-10	 w-10 rounded-full hover:bg-purple-100 hover:text-purple-600"
			onclick={() => openDialog('utility')}
		>
			<Settings class="h-4 w-4" />
		</Button>
	</div>
</div>

<!-- Node Selection Dialog -->
<NodeSelectionDialog
	open={showDialog}
	category={selectedCategory}
	availableNodes={nodeStore.availableNodes}
	onNodeSelect={handleNodeSelect}
	onClose={closeDialog}
/>

<style>
</style>
