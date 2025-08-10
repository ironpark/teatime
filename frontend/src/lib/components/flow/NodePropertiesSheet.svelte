<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import { Settings } from 'lucide-svelte';
	import type { Node } from '@xyflow/svelte';
	interface NodeProperty {
		name?: string;
		description?: string;
		key: string;
		value?: any;
		type: number;
		optional?: boolean;
		options?: string[];
	}
	import { Switch } from '$lib/components/ui/switch';
	import { Badge } from '$lib/components/ui/badge';

	interface Props {
		selectedNodes: Node[];
		onNodeUpdate: (nodeId: string, updates: any) => void;
		open?: boolean;
		onOpenChange?: (open: boolean) => void;
	}

	let { selectedNodes = $bindable([]), onNodeUpdate, open = $bindable(false), onOpenChange }: Props = $props();

	let selectedNode = $derived(selectedNodes?.[0]);
	
	$effect(() => {
		console.log('NodePropertiesSheet - open:', open);
		console.log('NodePropertiesSheet - selectedNode:', selectedNode);
	});

	function updateNodeData(field: string, value: any) {
		if (selectedNode) {
			onNodeUpdate(selectedNode.id, { [field]: value });
		}
	}

	function updateProperty(prop: NodeProperty, value: any) {
		if (selectedNode && selectedNode.data.properties && Array.isArray(selectedNode.data.properties)) {
			const properties = [...selectedNode.data.properties] as NodeProperty[];
			const index = properties.findIndex((p: NodeProperty) => p.key === prop.key);
			if (index !== -1) {
				properties[index] = { ...properties[index], value };
				onNodeUpdate(selectedNode.id, { properties });
			}
		}
	}

	function getPropertyTypeDisplay(propType: number): string {
		const typeMap: Record<number, string> = {
			1: 'boolean',
			2: 'number',
			3: 'int32',
			4: 'int64',
			5: 'string',
			6: 'json',
			7: 'xml',
			8: 'date',
			9: 'text',
			10: 'string[]',
			11: 'number[]',
			12: 'boolean[]',
			13: 'json[]',
			14: 'xml[]'
		};
		return typeMap[propType] || 'unknown';
	}

	// Dynamic property input component based on type
	function renderPropertyInput(prop: NodeProperty): any {
		if (prop.options && prop.options.length > 0) {
			// Render select for properties with options
			return {
				component: 'select',
				value: prop.value || '',
				options: prop.options.map((option: string) => ({
					value: option,
					label: option
				}))
			};
		}

		switch (prop.type) {
			case 1: // Bool
				return {
					component: 'switch',
					value: prop.value === 'true' || prop.value === true
				};
			case 2: // Float32
			case 3: // Float64
			case 4: // Int64
				return {
					component: 'input',
					type: 'number',
					value: prop.value || ''
				};
			case 6: // JSON
			case 7: // XML
			case 9: // Text
				return {
					component: 'textarea',
					value: prop.value || '',
					rows: 3
				};
			default: // String and others
				return {
					component: 'input',
					type: 'text',
					value: prop.value || ''
				};
		}
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
									oninput={(e) => updateNodeData('label', (e.target as HTMLInputElement).value)}
									placeholder="Enter node label"
								/>
							</div>
							
							<!-- Dynamic Properties -->
							{#if selectedNode.data.properties && Array.isArray(selectedNode.data.properties)}
								{#each selectedNode.data.properties as prop (prop.key)}
									{#if !prop.optional || prop.value}
										{@const inputConfig = renderPropertyInput(prop)}
										<div class="space-y-1.5">
											<div class="flex items-center justify-between">
												<Label for={`prop-${prop.key}`}>
													{prop.name || prop.key}
													{#if !prop.optional}
														<span class="text-red-500 ml-1">*</span>
													{/if}
												</Label>
												<Badge variant="outline" class="text-[10px] px-1 py-0">
													{getPropertyTypeDisplay(prop.type)}
												</Badge>
											</div>
											{#if prop.description}
												<p class="text-xs text-muted-foreground">{prop.description}</p>
											{/if}
											
											{#if inputConfig.component === 'switch'}
												<div class="flex items-center space-x-2">
													<Switch
														id={`prop-${prop.key}`}
														checked={inputConfig.value}
														onCheckedChange={(checked) => updateProperty(prop, checked.toString())}
													/>
													<Label for={`prop-${prop.key}`} class="text-sm">
														{inputConfig.value ? 'Yes' : 'No'}
													</Label>
												</div>
											{:else if inputConfig.component === 'textarea'}
												<Textarea
													id={`prop-${prop.key}`}
													value={inputConfig.value}
													oninput={(e) => updateProperty(prop, e.currentTarget.value)}
													rows={inputConfig.rows}
													placeholder={`Enter ${prop.name || prop.key}`}
												/>
											{:else if inputConfig.component === 'input'}
												<Input
													id={`prop-${prop.key}`}
													type={inputConfig.type}
													value={inputConfig.value}
													oninput={(e) => updateProperty(prop, e.currentTarget.value)}
													placeholder={`Enter ${prop.name || prop.key}`}
												/>
											{:else if inputConfig.component === 'select'}
												<Select
													type="single"
													value={inputConfig.value}
													onValueChange={(value) => updateProperty(prop, value)}
												>
													<SelectTrigger id={`prop-${prop.key}`}>
														{inputConfig.value || 'Select option'}
													</SelectTrigger>
													<SelectContent>
														{#each inputConfig.options as option}
															<SelectItem value={option.value}>
																{option.label}
															</SelectItem>
														{/each}
													</SelectContent>
												</Select>
											{/if}
										</div>
									{/if}
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

					<!-- Position Info -->
					<Separator />
					<div class="space-y-3 p-4 bg-muted/30 rounded-lg">
						<div class="text-sm font-medium">Position</div>
						<div class="grid grid-cols-2 gap-3">
							<div class="space-y-1.5">
								<Label for="node-x">X</Label>
								<Input
									id="node-x"
									type="number"
									value={Math.round(selectedNode.position.x)}
									readonly
									class="text-xs"
								/>
							</div>
							<div class="space-y-1.5">
								<Label for="node-y">Y</Label>
								<Input
									id="node-y"
									type="number"
									value={Math.round(selectedNode.position.y)}
									readonly
									class="text-xs"
								/>
							</div>
						</div>
					</div>
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