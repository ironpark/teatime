<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Plus, X } from 'lucide-svelte';
	import type { NodeProperty } from '$bindings/internal/node';

	interface Props {
		prop: NodeProperty;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, onUpdate }: Props = $props();

	const itemType = $derived(prop.input?.placeholder === 'number' ? 'number' : prop.input?.placeholder === 'textarea' ? 'textarea' : 'text');
	const placeholder = $derived(prop.input?.placeholder || 'Enter item');
	const unique = $derived(prop.input?.unique || false);

	function parseArrayValue(value: any): string[] {
		if (!value) return [''];
		
		if (Array.isArray(value)) {
			const result = value.map(item => String(item));
			return result.length > 0 ? result : [''];
		}
		
		if (typeof value === 'string') {
			try {
				const parsed = JSON.parse(value);
				if (Array.isArray(parsed)) {
					return parsed.map(item => String(item));
				}
			} catch {
				// If parsing fails, treat as single item
				return [value];
			}
		}
		
		// Convert single value to array
		return [String(value)];
	}

	function updateArrayValue(items: string[], keepEmpty = false) {
		let filteredItems = keepEmpty ? items : items.filter(item => item.trim() !== '');
		
		// Apply unique constraint if enabled (always, regardless of keepEmpty)
		if (unique) {
			filteredItems = [...new Set(filteredItems)];
		}
		
		if (itemType === 'number') {
			const numberItems = filteredItems.map(item => {
				const num = parseFloat(item);
				return isNaN(num) ? 0 : num;
			});
			// Apply unique constraint for numbers too
			const finalNumbers = unique ? [...new Set(numberItems)] : numberItems;
			onUpdate(prop, finalNumbers);
		} else {
			onUpdate(prop, filteredItems);
		}
	}

	function addItem() {
		const currentItems = parseArrayValue(prop.value);
		currentItems.push('');
		// Don't filter when adding new items - keep the empty string
		onUpdate(prop, currentItems);
	}

	function removeItem(index: number) {
		const currentItems = parseArrayValue(prop.value);
		currentItems.splice(index, 1);
		if (currentItems.length === 0) {
			currentItems.push('');
		}
		updateArrayValue(currentItems);
	}

	function handleItemChange(index: number, newValue: string) {
		const currentItems = parseArrayValue(prop.value);
		currentItems[index] = newValue;
		// Keep empty strings during editing to allow users to type
		// Don't apply unique constraint during typing - only on blur
		if (unique) {
			// Just update without deduplication during typing
			onUpdate(prop, currentItems);
		} else {
			updateArrayValue(currentItems, true);
		}
	}
	
	function handleItemBlur() {
		// Apply unique constraint when focus leaves the input
		if (unique && hasDuplicates) {
			const currentItems = parseArrayValue(prop.value);
			updateArrayValue(currentItems, false); // This will apply deduplication
		}
	}

	const arrayItems = $derived(parseArrayValue(prop.value));
	const effectivePlaceholder = $derived(unique ? `${placeholder} (unique values only)` : placeholder);
	
	// Check for duplicates when unique constraint is enabled
	const hasDuplicates = $derived.by(() => {
		if (!unique) return false;
		const nonEmptyItems = arrayItems.filter(item => item.trim() !== '');
		return nonEmptyItems.length !== new Set(nonEmptyItems).size;
	});
	
	const duplicateItems = $derived.by(() => {
		if (!unique) return new Set();
		const itemCounts = new Map();
		arrayItems.forEach(item => {
			if (item.trim() !== '') {
				itemCounts.set(item, (itemCounts.get(item) || 0) + 1);
			}
		});
		return new Set([...itemCounts.entries()].filter(([, count]) => count > 1).map(([item]) => item));
	});
	
</script>

<div class="space-y-2">
	{#each arrayItems as item, index}
		<div class="flex items-center gap-2">
			{#if itemType === 'textarea'}
				<textarea
					value={item}
					oninput={(e) => handleItemChange(index, e.currentTarget.value)}
					onblur={handleItemBlur}
					placeholder={effectivePlaceholder}
					rows={2}
					class="flex-1 flex min-h-[80px] w-full rounded-md border {duplicateItems.has(item) ? 'border-red-500' : 'border-input'} bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					autocapitalize="none"
					autocorrect="off"
					spellcheck="false"
				></textarea>
			{:else}
				<input
					type={itemType === 'number' ? 'number' : 'text'}
					value={item}
					oninput={(e) => handleItemChange(index, e.currentTarget.value)}
					onblur={handleItemBlur}
					placeholder={effectivePlaceholder}
					class="flex-1 flex h-10 w-full rounded-md border {duplicateItems.has(item) ? 'border-red-500' : 'border-input'} bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					autocapitalize="none"
					autocorrect="off"
					spellcheck="false"
				/>
			{/if}
			<Button
				variant="ghost"
				size="icon"
				onclick={() => removeItem(index)}
				class="h-8 w-8 text-red-500 hover:text-red-700"
			>
				<X class="h-4 w-4" />
			</Button>
		</div>
	{/each}
	{#if unique && hasDuplicates}
		<div class="text-xs text-red-500 mt-2 flex items-center gap-1">
			<span class="inline-block w-1 h-1 bg-red-500 rounded-full"></span>
			Duplicate values detected. They will be removed when you finish editing.
		</div>
	{/if}
	<Button
		variant="outline"
		size="sm"
		onclick={addItem}
		class="w-full"
	>
		<Plus class="h-4 w-4 mr-2" />
		Add Item
	</Button>

</div>