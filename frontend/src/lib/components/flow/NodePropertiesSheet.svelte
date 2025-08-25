<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import { Settings } from 'lucide-svelte';
	import type { Node } from '@xyflow/svelte';
	import { InputType, type NodeProperty } from '$bindings/internal/node';
	import { Badge } from '$lib/components/ui/badge';
	import { useSvelteFlow } from '@xyflow/svelte';
	import { 
		KeyValuePairsInput, 
		DynamicListInput, 
		SelectInput, 
		TextInput, 
		TextareaInput, 
		RangeInput, 
		SwitchInput, 
		ExpressionInput,
		MultiSelectInput,
		ArgListInput
	} from './inputs';
	import BindingSelector from './BindingSelector.svelte';
	import BindingDisplay from './BindingDisplay.svelte';
	import { isBinding, getBindingOptions, type BindingOption } from './utils/binding';
	import type { RecipeStore } from '$lib/stores/recipe.svelte';

	const { updateNodeData } = useSvelteFlow();
	
	interface Props {
		recipeStore: RecipeStore;
		editNode?: Node;
		open?: boolean;
		onOpenChange?: (open: boolean) => void;
	}

	let { editNode, open = $bindable(false), onOpenChange, recipeStore }: Props = $props();
	const selectedNodeId = $derived(editNode?.id);
	const selectedNode = $derived(recipeStore.nodes.find(n => n.id === selectedNodeId));
	
	// Binding state management
	let showBindingSelector = $state<{[key: string]: boolean}>({});
	async function updateNodeProperty(field: string, value: any) {
		if (selectedNode) {
			// For basic node properties like label
			if (field === 'label') {
				await recipeStore.updateNodeData(selectedNode.id, value, getPropertiesMap());
			} else {
				// For other data properties, update via updateNodeData (local change)
				updateNodeData(selectedNode.id, { [field]: value });
			}
		}
	}

	async function updateProperty(prop: NodeProperty, value: any) {
		if (selectedNode && selectedNode.data.properties && Array.isArray(selectedNode.data.properties)) {
			const properties = [...selectedNode.data.properties] as NodeProperty[];
			const index = properties.findIndex((p: NodeProperty) => p.key === prop.key);
			if (index !== -1) {
				properties[index] = { ...properties[index], value };
				
				// Update via backend API
				const propertiesMap = properties.reduce((acc, prop) => {
					acc[prop.key] = prop.value;
					return acc;
				}, {} as Record<string, any>);
				
				await recipeStore.updateNodeData(selectedNode.id, selectedNode.data.label, propertiesMap);
				// updateNodeData(selectedNode.id, { [prop.key]: value });
			}
		}
	}

	function getPropertiesMap(): Record<string, any> {
		if (!selectedNode?.data.properties || !Array.isArray(selectedNode.data.properties)) {
			return {};
		}
		return selectedNode.data.properties.reduce((acc, prop) => {
			acc[prop.key] = prop.value;
			return acc;
		}, {} as Record<string, any>);
	}

	// Binding functions
	function toggleBindingMode(prop: NodeProperty) {
		if (isBinding(prop.value)) {
			// Clear binding - restore to default value
			const defaultValue = getDefaultValueForType(prop.type);
			updateProperty(prop, defaultValue);
		} else {
			// Enter binding mode
			showBindingSelector = { ...showBindingSelector, [prop.key]: true };
		}
	}

	function selectBinding(prop: NodeProperty, option: BindingOption) {
		updateProperty(prop, option.bindingValue);
		showBindingSelector = { ...showBindingSelector, [prop.key]: false };
	}

	function cancelBinding(propKey: string) {
		showBindingSelector = { ...showBindingSelector, [propKey]: false };
	}

	function getDefaultValueForType(type: number): any {
		// Based on PropertyType enum
		switch (type) {
			case 1: return false; // Bool
			case 2: case 3: case 4: return 0; // Int64, Uint64, Float64
			case 5: case 6: case 7: case 8: return ''; // String, JSON, XML, Date
			case 9: case 10: case 11: return []; // Arrays
			default: return '';
		}
	}

	// Helper function to get property type display
	function getPropertyTypeDisplay(propType: number): string {
		const typeMap: Record<number, string> = {
			0: 'invalid',
			1: 'boolean',
			2: 'int64',
			3: 'uint64', 
			4: 'float64',
			5: 'string',
			6: 'json',
			7: 'xml',
			8: 'date',
			9: 'string[]',
			10: 'number[]',
			11: 'boolean[]'
		};
		return typeMap[propType] || 'unknown';
	}

	const titleCase = (str: any) => {
		if (!str) return '';
		if (typeof str === 'string') {
			return str.charAt(0).toUpperCase() + str.slice(1);
		}
		return str;
	}

</script>

<Sheet.Root open={open} onOpenChange={(val) => {
	console.log('Sheet onOpenChange:', val);
	open = val;
	onOpenChange?.(val);
}}>
	<Sheet.Portal>
		<Sheet.Content class="w-[450px] flex flex-col">
		{#if selectedNode}
			<Sheet.Header>
				<Sheet.Title class="flex items-center gap-2">
					<Settings class="h-4 w-4" />
					{titleCase(selectedNode.data.name)} {titleCase(selectedNode.data.nodeType)} Properties
				</Sheet.Title>
			</Sheet.Header>

			<ScrollArea class="flex-1 overflow-y-auto">
				<div class="space-y-6 px-6">
					<!-- Node Info -->
					<div class="space-y-3">
						<div class="space-y-2">
							<div class="flex items-center gap-2">
								<span class="text-sm text-muted-foreground">ID:</span>
								<code class="text-xs font-mono bg-muted px-1.5 py-0.5 rounded">
									{selectedNode.id}
								</code>
							</div>
							{#if selectedNode.data.name}
								<div class="flex items-center gap-2">
									<span class="text-sm text-muted-foreground">Name:</span>
									<span class="text-sm font-medium">
										{selectedNode.data.name}
									</span>
								</div>
							{/if}
							{#if selectedNode.data.description}
								<div class="space-y-1">
									<p class="text-sm text-muted-foreground/80">
										{selectedNode.data.description}
									</p>
								</div>
							{/if}
						</div>
					</div>

					<!-- Node Properties -->
					<Separator />
					<div class="space-y-3">
						<div class="text-sm font-medium">Node Properties</div>
						<div class="space-y-3">
							<!-- Label field -->
							<div class="space-y-1.5">
								<div class="flex items-center justify-between">
									<Label for="node-label">
										Label
										<span class="text-red-500 ml-1">*</span>
									</Label>
									<Badge variant="outline" class="text-[10px] px-1 py-0">
										string
									</Badge>
								</div>
								<Input
									id="node-label"
									value={selectedNode.data.label || ''}
									oninput={(e) => updateNodeProperty('label', (e.target as HTMLInputElement).value)}
									placeholder="Enter node label"
								/>
							</div>
							
							<!-- Dynamic Properties -->
							{#if selectedNode.data.properties && Array.isArray(selectedNode.data.properties)}
								{#each selectedNode.data.properties as prop (prop.key)}
										<div class="space-y-1.5">
											<div class="flex items-center justify-between">
												<Label for={`prop-${prop.key}`}>
													{prop.name || prop.key}
													<!-- <Badge variant="outline" class="text-[12px] px-2 rounded-sm">
														{getPropertyTypeDisplay(prop.type)}
													</Badge> -->
													{#if !prop.optional}
														<span class="text-red-500 ml-1">*</span>
													{/if}
													
												</Label>
												<div class="flex items-center gap-2">
													<!-- Binding Button -->
													<BindingSelector 
														{prop} 
														currentNodeId={selectedNode.id} 
														nodes={recipeStore.nodes} 
														edges={recipeStore.edges} 
														onToggleBinding={toggleBindingMode} 
													/>
												</div>
											</div>
											{#if prop.description}
												<p class="text-xs text-muted-foreground">{prop.description}</p>
											{/if}
											
											<!-- Binding Display and Selector -->
											<BindingDisplay 
												{prop} 
												currentNodeId={selectedNode.id} 
												nodes={recipeStore.nodes} 
												edges={recipeStore.edges} 
												onUpdate={updateProperty}
												showBindingSelector={showBindingSelector[prop.key] || false}
												onSelectBinding={(option) => selectBinding(prop, option)}
												onCancelBinding={() => cancelBinding(prop.key)}
											/>
											
											<!-- Show input controls only if not bound -->
											{#if !isBinding(prop.value) && prop.input}
												{#if prop.input.type === InputType.InputTypeSwitch}
													<SwitchInput {prop} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeRange}
													<RangeInput {prop} min={prop.input.min ?? 0} max={prop.input.max ?? 0} step={prop.input.step ?? 0} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeTextarea}
													<TextareaInput {prop} rows={prop.input.rows} placeholder={prop.input.placeholder} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeNumber}
													<TextInput {prop} type="number" min={prop.input.min ?? 0} max={prop.input.max ?? 0} step={prop.input.step ?? 0} placeholder={prop.input.placeholder} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeSelect}
													<SelectInput {prop} options={prop.options?.map((option) => ({ value: option, label: option })) ?? []} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeMultiSelect}
													<MultiSelectInput {prop} options={prop.options?.map((option) => ({ value: option, label: option })) ?? []} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeExpression}
													<ExpressionInput {prop} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeKeyValue}
													<KeyValuePairsInput {prop} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeDynamicList}
													<DynamicListInput {prop} onUpdate={updateProperty} />
												{:else if prop.input.type === InputType.InputTypeArgList}
													<ArgListInput {prop} onUpdate={updateProperty} />
												{:else}
													<!-- Default to text input -->
													<TextInput {prop} type="text" placeholder={prop.input.placeholder} onUpdate={updateProperty} />
												{/if}
											{:else if !isBinding(prop.value) && prop.options && prop.options.length > 0}
												<!-- Fallback to select if options are provided -->
												<SelectInput {prop} options={prop.options.map((option) => ({ value: option, label: option }))} onUpdate={updateProperty} />
											{:else if !isBinding(prop.value)}
												<!-- Fallback based on property type -->
												{#if prop.type === 1}
													<!-- Bool -->
													<SwitchInput {prop} onUpdate={updateProperty} />
												{:else if prop.type === 2 || prop.type === 3 || prop.type === 4}
													<!-- Int64, Uint64, Float64 -->
													<TextInput {prop} type="number" onUpdate={updateProperty} />
												{:else if prop.type === 5 || prop.type === 6 || prop.type === 7 || prop.type === 8 || prop.type === 9 || prop.type === 10 || prop.type === 11}
													<!-- String, JSON, XML, Date, Arrays -->
													<TextareaInput {prop} rows={3} onUpdate={updateProperty} />
												{:else}
													<!-- Default -->
													<TextInput {prop} type="text" onUpdate={updateProperty} />
												{/if}
											{/if}
										</div>
								{/each}
							{/if}
						</div>
					</div>

					<!-- Output Properties -->
					{#if selectedNode.data.outputs && Array.isArray(selectedNode.data.outputs)}
						<Separator />
						<div class="space-y-3">
							<div class="text-sm font-medium">Outputs</div>
							<div class="space-y-2">
								{#each selectedNode.data.outputs as output (output.key)}
									<div class="flex items-center justify-between p-2 bg-muted/50 rounded-md">
										<div class="flex items-center gap-2">
											<span class="text-sm font-medium">
												{output.name || output.key}
											</span>
											{#if output.optional}
												<Badge variant="outline" class="text-[10px] px-1 py-0">
													optional
												</Badge>
											{/if}
										</div>
										<Badge variant="secondary" class="text-[10px] px-1 py-0">
											{getPropertyTypeDisplay(output.type)}
										</Badge>
									</div>
								{/each}
							</div>
						</div>
					{/if}

				</div>
			</ScrollArea>

			<div class="border-t px-6 py-4 bg-background flex-shrink-0">
				<div class="flex justify-end">
					<Button variant="outline" onclick={() => {
						open = false;
						onOpenChange?.(false);
					}}>
						Close
					</Button>
				</div>
			</div>
		{:else}
			<div class="flex items-center justify-center h-full">
				<p class="text-muted-foreground">Select a node to view properties</p>
			</div>
		{/if}
		</Sheet.Content>
	</Sheet.Portal>
</Sheet.Root>