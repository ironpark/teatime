<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Switch } from '$lib/components/ui/switch';
	import { Label } from '$lib/components/ui/label';
	import { Plus, X } from 'lucide-svelte';
	import type { NodeProperty } from '$bindings/internal/node';

	interface CommandArg {
		name: string;        // 옵션이름
		required: boolean;   // 필수여부
		list: boolean;      // 리스트여부
		description: string; // 설명
	}

	interface Props {
		prop: NodeProperty;
		onUpdate: (prop: NodeProperty, value: any) => void;
	}

	let { prop, onUpdate }: Props = $props();

	function parseArgList(value: any): CommandArg[] {
		console.log('parseArgList input:', value, typeof value);
		
		if (!value || (Array.isArray(value) && value.length === 0)) {
			console.log('empty value, returning default');
			return [{ name: '', required: false, list: false, description: '' }];
		}
		
		if (Array.isArray(value)) {
			const result = value.map(item => {
				if (typeof item === 'object' && item !== null) {
					return {
						name: item.name || '',
						required: Boolean(item.required),
						list: Boolean(item.list),
						description: item.description || ''
					};
				}
				return { name: String(item), required: false, list: false, description: '' };
			});
			console.log('parsed array result:', result);
			return result.length > 0 ? result : [{ name: '', required: false, list: false, description: '' }];
		}
		
		if (typeof value === 'string') {
			try {
				const parsed = JSON.parse(value);
				if (Array.isArray(parsed)) {
					return parseArgList(parsed);
				}
			} catch {
				// If parsing fails, treat as single argument name
				return [{ name: value, required: false, list: false, description: '' }];
			}
		}
		
		return [{ name: '', required: false, list: false, description: '' }];
	}

	function updateArgList(args: CommandArg[], keepEmpty = false) {
		console.log('updateArgList called with:', args, 'keepEmpty:', keepEmpty);
		if (keepEmpty) {
			// Keep all arguments when explicitly requested (e.g., when adding)
			console.log('keeping all args:', args);
			onUpdate(prop, args);
		} else {
			// Filter out completely empty arguments during normal updates
			const validArgs = args.filter(arg => 
				arg.name.trim() !== '' || 
				arg.description.trim() !== '' || 
				arg.required || 
				arg.list
			);
			
			console.log('filtered args:', validArgs);
			onUpdate(prop, validArgs.length > 0 ? validArgs : []);
		}
	}

	function addArg() {
		console.log('addArg called, current value:', prop.value);
		const currentArgs = parseArgList(prop.value);
		console.log('parsed args:', currentArgs);
		currentArgs.push({ name: '', required: false, list: false, description: '' });
		console.log('args after push:', currentArgs);
		updateArgList(currentArgs, true); // Keep empty arguments when adding
	}

	function removeArg(index: number) {
		const currentArgs = parseArgList(prop.value);
		currentArgs.splice(index, 1);
		if (currentArgs.length === 0) {
			currentArgs.push({ name: '', required: false, list: false, description: '' });
		}
		updateArgList(currentArgs);
	}

	function updateArg(index: number, field: keyof CommandArg, value: any) {
		const currentArgs = parseArgList(prop.value);
		if (currentArgs[index]) {
			currentArgs[index] = { ...currentArgs[index], [field]: value };
			updateArgList(currentArgs, true); // Keep empty arguments during editing
		}
	}

	const argList = $derived(parseArgList(prop.value));
</script>

<div class="space-y-2">
	{#each argList as arg, index}
		<div class="border rounded-md p-3 space-y-2 relative">
			<Button
				variant="ghost"
				size="icon"
				onclick={() => removeArg(index)}
				class="absolute top-2 right-2 h-5 w-5 text-red-500 hover:text-red-700"
			>
				<X class="h-3 w-3" />
			</Button>
			
			<div class="grid grid-cols-2 gap-2">
				<!-- Argument Name -->
				<div class="space-y-1">
					<Label class="text-xs text-muted-foreground">Option Name</Label>
					<Input
						value={arg.name}
						oninput={(e) => updateArg(index, 'name', e.currentTarget.value)}
						placeholder="--option or -o"
						class="h-7 text-xs"
						autocapitalize="none"
						spellcheck={false}
					/>
				</div>
				
				<!-- Flags -->
				<div class="space-y-1">
					<Label class="text-xs text-muted-foreground">Flags</Label>
					<div class="flex items-center gap-3">
						<div class="flex items-center space-x-1">
							<Switch
								id={`required-${index}`}
								checked={arg.required}
								onCheckedChange={(checked) => updateArg(index, 'required', checked)}
								class="scale-75"
							/>
							<Label for={`required-${index}`} class="text-xs">Req</Label>
						</div>
						<div class="flex items-center space-x-1">
							<Switch
								id={`list-${index}`}
								checked={arg.list}
								onCheckedChange={(checked) => updateArg(index, 'list', checked)}
								class="scale-75"
							/>
							<Label for={`list-${index}`} class="text-xs">List</Label>
						</div>
					</div>
				</div>
			</div>
			
			<!-- Description -->
			<div class="space-y-1 pr-6">
				<Label class="text-xs text-muted-foreground">Description</Label>
				<Textarea
					value={arg.description}
					oninput={(e) => updateArg(index, 'description', e.currentTarget.value)}
					placeholder="Describe what this argument does..."
					rows={2}
					class="text-xs resize-none"
					autocapitalize="none"
					spellcheck={false}
				/>
			</div>
		</div>
	{/each}
	
	<Button
		variant="outline"
		size="sm"
		onclick={addArg}
		class="w-full"
	>
		<Plus class="h-4 w-4 mr-2" />
		Add Argument
	</Button>
</div>