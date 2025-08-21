<script lang="ts">
	import { setContext } from 'svelte';
	import { RecipeStore } from '$lib/stores/recipe.svelte';

	interface Props {
		recipeStore: RecipeStore;
		children: import('svelte').Snippet;
	}

	let { recipeStore, children }: Props = $props();

	// Create recipe context with store and recipe ID
	const recipeContext = {
		recipeStore: recipeStore,
		recipeID: recipeStore.recipeId
	};

	// Provide context to child components
	setContext('recipe', recipeContext);
</script>

{@render children()}