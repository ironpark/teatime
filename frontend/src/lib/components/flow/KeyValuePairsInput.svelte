<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Plus, X } from 'lucide-svelte';

	interface KeyValuePair {
		key: string;
		value: string;
	}

	interface Props {
		value: any;
		onUpdate: (value: any) => void;
	}

	let { value = $bindable(), onUpdate }: Props = $props();

	function parseKeyValuePairs(value: any): KeyValuePair[] {
		if (!value) return [{ key: '', value: '' }];
		
		if (typeof value === 'string') {
			try {
				const parsed = JSON.parse(value);
				return Object.entries(parsed).map(([key, value]) => ({ 
					key, 
					value: String(value) 
				}));
			} catch {
				return [{ key: '', value: '' }];
			}
		}
		
		if (typeof value === 'object') {
			const list = Object.entries(value).map(([key, value]) => ({ 
				key, 
				value: String(value) 
			}));
			if (list.length === 0) {
				return [{ key: '', value: '' }];
			}
			return list;
		}
		
		return [{ key: '', value: '' }];
	}

	function updateKeyValuePairs(pairs: KeyValuePair[]) {
		const obj: Record<string, string> = {};
		pairs.forEach(pair => {
			if (pair.key) {
				obj[pair.key] = pair.value;
			}
		});
		onUpdate(obj);
	}

	function addKeyValuePair() {
		const currentPairs = parseKeyValuePairs(value);
		currentPairs.push({ key: '', value: '' });
		updateKeyValuePairs(currentPairs);
	}

	function removeKeyValuePair(index: number) {
		const currentPairs = parseKeyValuePairs(value);
		currentPairs.splice(index, 1);
		if (currentPairs.length === 0) {
			currentPairs.push({ key: '', value: '' });
		}
		updateKeyValuePairs(currentPairs);
	}

	const keyValuePairs = $derived(parseKeyValuePairs(value));
</script>

<div class="space-y-2">
	{#each keyValuePairs as pair, index}
		<div class="flex items-center gap-2">
			<Input
				placeholder="Key"
				bind:value={pair.key}
				oninput={() => updateKeyValuePairs(keyValuePairs)}
				class="flex-1"
			/>
			<Input
				placeholder="Value"
				bind:value={pair.value}
				oninput={() => updateKeyValuePairs(keyValuePairs)}
				class="flex-1"
			/>
			<Button
				variant="ghost"
				size="icon"
				onclick={() => removeKeyValuePair(index)}
				class="h-8 w-8 text-red-500 hover:text-red-700"
			>
				<X class="h-4 w-4" />
			</Button>
		</div>
	{/each}
	<Button
		variant="outline"
		size="sm"
		onclick={addKeyValuePair}
		class="w-full"
	>
		<Plus class="h-4 w-4 mr-2" />
		Add Pair
	</Button>
</div>