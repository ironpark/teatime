<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { Badge } from '$lib/components/ui/badge';
	import {
		Save,
		Trash2,
		Undo2,
		Redo2,
		ZoomIn,
		ZoomOut,
		Maximize,
		Play,
		Download,
		Upload,
		Settings
	} from 'lucide-svelte';
	import type { Node, Edge } from '@xyflow/svelte';

	interface Props {
		selectedNodes?: Node[];
		selectedEdges?: Edge[];
		onSave: () => void;
		onDelete: () => void;
	}

	let { selectedNodes = [], selectedEdges = [], onSave, onDelete }: Props = $props();

	function handleUndo() {
		// TODO: Implement undo functionality
		console.log('Undo');
	}

	function handleRedo() {
		// TODO: Implement redo functionality
		console.log('Redo');
	}

	function handleZoomIn() {
		// TODO: Implement zoom in
		console.log('Zoom in');
	}

	function handleZoomOut() {
		// TODO: Implement zoom out
		console.log('Zoom out');
	}

	function handleFitView() {
		// TODO: Implement fit view
		console.log('Fit view');
	}

	function handleTestRecipe() {
		// TODO: Implement recipe testing/simulation
		console.log('Test recipe');
	}

	function handleExport() {
		// TODO: Implement recipe export
		console.log('Export recipe');
	}

	function handleImport() {
		// TODO: Implement recipe import
		console.log('Import recipe');
	}

	function handleSettings() {
		// TODO: Implement settings
		console.log('Open settings');
	}

	let hasSelection = $derived(selectedNodes.length > 0 || selectedEdges.length > 0);
	let selectionCount = $derived(selectedNodes.length + selectedEdges.length);
</script>

<div class="flex items-center gap-2">
	<!-- File operations -->
	<div class="flex items-center gap-1">
		<Button variant="ghost" size="sm" onclick={onSave}>
			<Save class="mr-1 h-4 w-4" />
		</Button>
	</div>

	<Separator orientation="vertical" class="h-6" />

	<!-- Edit operations -->
	<div class="flex items-center gap-1">
		<Button variant="ghost" size="sm" onclick={handleUndo}>
			<Undo2 class="h-4 w-4" />
		</Button>

		<Button variant="ghost" size="sm" onclick={handleRedo}>
			<Redo2 class="h-4 w-4" />
		</Button>

		<Button
			variant="ghost"
			size="sm"
			onclick={onDelete}
			disabled={!hasSelection}
			class={hasSelection ? 'hover:bg-destructive hover:text-destructive-foreground' : ''}
		>
			<Trash2 class="h-4 w-4" />
		</Button>
	</div>

	<!-- View operations -->
	<div class="flex items-center gap-1">
		<Button variant="ghost" size="sm" onclick={handleZoomIn}>
			<ZoomIn class="h-4 w-4" />
		</Button>

		<Button variant="ghost" size="sm" onclick={handleZoomOut}>
			<ZoomOut class="h-4 w-4" />
		</Button>

		<Button variant="ghost" size="sm" onclick={handleFitView}>
			<Maximize class="h-4 w-4" />
		</Button>
	</div>

	<!-- Recipe operations -->
	<div class="flex items-center gap-1">
		<Button variant="secondary" size="sm" onclick={handleTestRecipe}>
			<Play class="mr-1 h-4 w-4" />
			Test
		</Button>

		<!-- <Button variant="outline" size="sm" onclick={handleSettings}>
			<Settings class="h-4 w-4" />
		</Button> -->
	</div>
</div>
