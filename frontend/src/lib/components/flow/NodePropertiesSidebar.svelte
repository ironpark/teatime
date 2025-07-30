<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import { Separator } from '$lib/components/ui/separator';
	import { X, Settings } from 'lucide-svelte';
	import type { Node } from '@xyflow/svelte';
	import { useOnSelectionChange } from '@xyflow/svelte';

	interface Props {
		selectedNodes: Node[];
		onNodeUpdate: (nodeId: string, updates: any) => void;
		onClose: () => void;
	}

	let { selectedNodes = $bindable([]), onNodeUpdate, onClose }: Props = $props();

	useOnSelectionChange((selection) => {
		selectedNodes = [...selection.nodes];
		// selectedEdges = [...selection.edges];
	});
	let isVisible = $derived.by(() => selectedNodes && selectedNodes.length > 0);
	let selectedNode = $derived.by(() => selectedNodes && selectedNodes[0]);

	$effect(() => {
		console.log('NodePropertiesSidebar - selectedNodes:', selectedNodes);
		console.log('NodePropertiesSidebar - isVisible:', isVisible);
		console.log('NodePropertiesSidebar - selectedNode:', selectedNode);
	});

	function updateNodeData(field: string, value: any) {
		if (selectedNode) {
			onNodeUpdate(selectedNode.id, { [field]: value });
		}
	}

	function handleInputChange(field: string, event: any) {
		updateNodeData(field, event.target.value);
	}
</script>

{#if isVisible}
	<div
		class="bg-card fixed bottom-4 right-4 top-20 z-40 flex w-80 flex-col rounded-lg border shadow-lg"
	>
		<!-- Header -->
		<div class="flex items-center justify-between border-b p-4">
			<div class="flex items-center gap-2">
				<Settings class="text-muted-foreground h-4 w-4" />
				<h3 class="text-sm font-semibold">Node Properties</h3>
			</div>
			<Button variant="ghost" size="icon" class="h-6 w-6" onclick={onClose}>
				<X class="h-4 w-4" />
			</Button>
		</div>

		<!-- Content -->
		<div class="flex-1 space-y-4 overflow-y-auto p-4">
			{#if selectedNode}
				<!-- Node Info -->
				<div class="space-y-2">
					<div class="text-muted-foreground text-xs uppercase tracking-wide">Node Info</div>
					<div class="space-y-1">
						<div class="text-sm">ID: <span class="font-mono text-xs">{selectedNode.id}</span></div>
						<div class="text-sm">Type: <span class="capitalize">{selectedNode.type}</span></div>
					</div>
				</div>

				<Separator />

				<!-- Node Properties based on type -->
				{#if selectedNode.data.nodeType === 'ingredient'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Ingredient Properties
						</div>

						<div class="space-y-1">
							<Label for="ingredient-name">Name</Label>
							<Input
								id="ingredient-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Ingredient name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="ingredient-amount">Amount</Label>
							<Input
								id="ingredient-amount"
								value={selectedNode.data.amount}
								oninput={(e) => handleInputChange('amount', e)}
								placeholder="e.g., 2 cups, 500g"
							/>
						</div>

						<div class="space-y-1">
							<Label for="ingredient-category">Category</Label>
							<Select
								value={selectedNode.data.category}
								onValueChange={(value) => updateNodeData('category', value)}
							>
								<SelectTrigger>
									{selectedNode.data.category || 'Select category'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="dry">Dry ingredients</SelectItem>
									<SelectItem value="wet">Wet ingredients</SelectItem>
									<SelectItem value="protein">Protein</SelectItem>
									<SelectItem value="vegetable">Vegetables</SelectItem>
									<SelectItem value="spice">Spices</SelectItem>
									<SelectItem value="other">Other</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'step'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">Step Properties</div>

						<div class="space-y-1">
							<Label for="step-name">Step Name</Label>
							<Input
								id="step-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Step name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="step-instruction">Instruction</Label>
							<Textarea
								id="step-instruction"
								value={selectedNode.data.instruction}
								oninput={(e) => handleInputChange('instruction', e)}
								placeholder="Detailed cooking instruction"
								rows={3}
							/>
						</div>

						<div class="space-y-1">
							<Label for="step-duration">Duration</Label>
							<Input
								id="step-duration"
								value={selectedNode.data.duration}
								oninput={(e) => handleInputChange('duration', e)}
								placeholder="e.g., 10 minutes"
							/>
						</div>

						<div class="space-y-1">
							<Label for="step-temperature">Temperature (optional)</Label>
							<Input
								id="step-temperature"
								value={selectedNode.data.temperature || ''}
								oninput={(e) => handleInputChange('temperature', e.target.value || null)}
								placeholder="e.g., 180°C, 350°F"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'output'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Output Properties
						</div>

						<div class="space-y-1">
							<Label for="output-name">Output Name</Label>
							<Input
								id="output-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Output name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="output-description">Description</Label>
							<Textarea
								id="output-description"
								value={selectedNode.data.description}
								oninput={(e) => handleInputChange('description', e)}
								placeholder="Describe the final result"
								rows={3}
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'variable'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Variable Properties
						</div>

						<div class="space-y-1">
							<Label for="variable-name">Variable Name</Label>
							<Input
								id="variable-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Variable name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="variable-type">Type</Label>
							<Select
								value={selectedNode.data.type}
								onValueChange={(value) => updateNodeData('type', value)}
							>
								<SelectTrigger>
									{selectedNode.data.type || 'Select type'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="string">String</SelectItem>
									<SelectItem value="number">Number</SelectItem>
									<SelectItem value="boolean">Boolean</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="variable-value">Value</Label>
							<Input
								id="variable-value"
								value={selectedNode.data.value}
								oninput={(e) => handleInputChange('value', e)}
								placeholder="Variable value"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'constant'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Constant Properties
						</div>

						<div class="space-y-1">
							<Label for="constant-name">Constant Name</Label>
							<Input
								id="constant-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Constant name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="constant-type">Type</Label>
							<Select
								value={selectedNode.data.type}
								onValueChange={(value) => updateNodeData('type', value)}
							>
								<SelectTrigger>
									{selectedNode.data.type || 'Select type'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="string">String</SelectItem>
									<SelectItem value="number">Number</SelectItem>
									<SelectItem value="boolean">Boolean</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="constant-value">Value</Label>
							<Input
								id="constant-value"
								value={selectedNode.data.value}
								oninput={(e) => handleInputChange('value', e)}
								placeholder="Constant value"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'formatter'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Formatter Properties
						</div>

						<div class="space-y-1">
							<Label for="formatter-name">Formatter Name</Label>
							<Input
								id="formatter-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Formatter name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="formatter-type">Format Type</Label>
							<Select
								value={selectedNode.data.type}
								onValueChange={(value) => updateNodeData('type', value)}
							>
								<SelectTrigger>
									{selectedNode.data.type || 'Select format type'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="text">Text</SelectItem>
									<SelectItem value="date">Date</SelectItem>
									<SelectItem value="number">Number</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="formatter-format">Format Pattern</Label>
							<Input
								id="formatter-format"
								value={selectedNode.data.format}
								oninput={(e) => handleInputChange('format', e)}
								placeholder="e.g., YYYY-MM-DD, #,##0.00"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'calculator'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Calculator Properties
						</div>

						<div class="space-y-1">
							<Label for="calculator-name">Calculator Name</Label>
							<Input
								id="calculator-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Calculator name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="calculator-operation">Operation</Label>
							<Select
								value={selectedNode.data.operation}
								onValueChange={(value) => updateNodeData('operation', value)}
							>
								<SelectTrigger>
									{selectedNode.data.operation || 'Select operation'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="add">Addition</SelectItem>
									<SelectItem value="subtract">Subtraction</SelectItem>
									<SelectItem value="multiply">Multiplication</SelectItem>
									<SelectItem value="divide">Division</SelectItem>
									<SelectItem value="custom">Custom Expression</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="calculator-expression">Expression</Label>
							<Input
								id="calculator-expression"
								value={selectedNode.data.expression}
								oninput={(e) => handleInputChange('expression', e)}
								placeholder="e.g., a + b, (x * y) / 2"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'aggregator'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Aggregator Properties
						</div>

						<div class="space-y-1">
							<Label for="aggregator-name">Aggregator Name</Label>
							<Input
								id="aggregator-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Aggregator name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="aggregator-operation">Operation</Label>
							<Select
								value={selectedNode.data.operation}
								onValueChange={(value) => updateNodeData('operation', value)}
							>
								<SelectTrigger>
									{selectedNode.data.operation || 'Select operation'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="sum">Sum</SelectItem>
									<SelectItem value="average">Average</SelectItem>
									<SelectItem value="count">Count</SelectItem>
									<SelectItem value="min">Minimum</SelectItem>
									<SelectItem value="max">Maximum</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="aggregator-field">Field (optional)</Label>
							<Input
								id="aggregator-field"
								value={selectedNode.data.field || ''}
								oninput={(e) => handleInputChange('field', e)}
								placeholder="Field to aggregate"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'command'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Command Properties
						</div>

						<div class="space-y-1">
							<Label for="command-name">Command Name</Label>
							<Input
								id="command-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Command name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="command-cmd">Command</Label>
							<Textarea
								id="command-cmd"
								value={selectedNode.data.command}
								oninput={(e) => handleInputChange('command', e)}
								placeholder="e.g., transcribe --input INPUT_FILE --lang $LANGUAGE"
								rows={3}
							/>
						</div>

						<div class="space-y-1">
							<Label for="command-workingdir">Working Directory</Label>
							<Input
								id="command-workingdir"
								value={selectedNode.data.workingDir || ''}
								oninput={(e) => handleInputChange('workingDir', e)}
								placeholder="/workspace"
							/>
						</div>

						<div class="space-y-1">
							<Label for="command-timeout">Timeout (seconds)</Label>
							<Input
								id="command-timeout"
								type="number"
								value={selectedNode.data.timeout || ''}
								oninput={(e) =>
									handleInputChange('timeout', e.target.value ? parseInt(e.target.value) : null)}
								placeholder="30"
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'mcp'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							MCP Server Properties
						</div>

						<div class="space-y-1">
							<Label for="mcp-name">Server Name</Label>
							<Input
								id="mcp-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="MCP Server name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="mcp-url">Server URL</Label>
							<Input
								id="mcp-url"
								value={selectedNode.data.serverUrl}
								oninput={(e) => handleInputChange('serverUrl', e)}
								placeholder="mcp://server-name/v1"
							/>
						</div>

						<div class="space-y-1">
							<Label for="mcp-method">Method</Label>
							<Input
								id="mcp-method"
								value={selectedNode.data.method || ''}
								oninput={(e) => handleInputChange('method', e)}
								placeholder="process_video"
							/>
						</div>

						<div class="space-y-1">
							<Label for="mcp-params">Parameters (JSON)</Label>
							<Textarea
								id="mcp-params"
								value={selectedNode.data.params || ''}
								oninput={(e) => handleInputChange('params', e)}
								placeholder={`{"video_path": "INPUT_FILE"}`}
								rows={3}
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'manual'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Manual Trigger Properties
						</div>

						<div class="space-y-1">
							<Label for="manual-name">Trigger Name</Label>
							<Input
								id="manual-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Manual trigger name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="manual-description">Description</Label>
							<Textarea
								id="manual-description"
								value={selectedNode.data.description || ''}
								oninput={(e) => handleInputChange('description', e)}
								placeholder="Describe what this trigger does"
								rows={2}
							/>
						</div>

						<div class="space-y-1">
							<Label for="manual-confirmation">Requires Confirmation</Label>
							<Select
								value={selectedNode.data.requiresConfirmation ? 'true' : 'false'}
								onValueChange={(value) => updateNodeData('requiresConfirmation', value === 'true')}
							>
								<SelectTrigger>
									{selectedNode.data.requiresConfirmation ? 'Yes' : 'No'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="false">No</SelectItem>
									<SelectItem value="true">Yes</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'file-input'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							File Input Properties
						</div>

						<div class="space-y-1">
							<Label for="file-input-name">Input Name</Label>
							<Input
								id="file-input-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="File input name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="file-input-path">File Path</Label>
							<Input
								id="file-input-path"
								value={selectedNode.data.filePath}
								oninput={(e) => handleInputChange('filePath', e)}
								placeholder="/path/to/input/file"
							/>
						</div>

						<div class="space-y-1">
							<Label for="file-input-type">File Type</Label>
							<Select
								value={selectedNode.data.fileType}
								onValueChange={(value) => updateNodeData('fileType', value)}
							>
								<SelectTrigger>
									{selectedNode.data.fileType || 'Select file type'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="any">Any file</SelectItem>
									<SelectItem value="video">Video files</SelectItem>
									<SelectItem value="audio">Audio files</SelectItem>
									<SelectItem value="image">Image files</SelectItem>
									<SelectItem value="text">Text files</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'ffmpeg-extract'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							FFmpeg Extract Properties
						</div>

						<div class="space-y-1">
							<Label for="ffmpeg-name">Extract Name</Label>
							<Input
								id="ffmpeg-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Extract operation name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="ffmpeg-format">Output Format</Label>
							<Select
								value={selectedNode.data.outputFormat}
								onValueChange={(value) => updateNodeData('outputFormat', value)}
							>
								<SelectTrigger>
									{selectedNode.data.outputFormat || 'Select format'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="wav">WAV</SelectItem>
									<SelectItem value="mp3">MP3</SelectItem>
									<SelectItem value="flac">FLAC</SelectItem>
									<SelectItem value="aac">AAC</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="ffmpeg-quality">Quality</Label>
							<Select
								value={selectedNode.data.quality}
								onValueChange={(value) => updateNodeData('quality', value)}
							>
								<SelectTrigger>
									{selectedNode.data.quality || 'Select quality'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="low">Low</SelectItem>
									<SelectItem value="medium">Medium</SelectItem>
									<SelectItem value="high">High</SelectItem>
									<SelectItem value="lossless">Lossless</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'transcription'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Transcription Properties
						</div>

						<div class="space-y-1">
							<Label for="transcription-name">Transcription Name</Label>
							<Input
								id="transcription-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Transcription operation name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="transcription-language">Language</Label>
							<Select
								value={selectedNode.data.language}
								onValueChange={(value) => updateNodeData('language', value)}
							>
								<SelectTrigger>
									{selectedNode.data.language || 'Select language'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="auto">Auto-detect</SelectItem>
									<SelectItem value="en">English</SelectItem>
									<SelectItem value="ko">Korean</SelectItem>
									<SelectItem value="ja">Japanese</SelectItem>
									<SelectItem value="zh">Chinese</SelectItem>
									<SelectItem value="es">Spanish</SelectItem>
									<SelectItem value="fr">French</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="transcription-model">Model</Label>
							<Select
								value={selectedNode.data.model}
								onValueChange={(value) => updateNodeData('model', value)}
							>
								<SelectTrigger>
									{selectedNode.data.model || 'Select model'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="whisper">Whisper</SelectItem>
									<SelectItem value="whisper-large">Whisper Large</SelectItem>
									<SelectItem value="wav2vec2">Wav2Vec2</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'translation'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							Translation Properties
						</div>

						<div class="space-y-1">
							<Label for="translation-name">Translation Name</Label>
							<Input
								id="translation-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Translation operation name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="translation-source">Source Language</Label>
							<Select
								value={selectedNode.data.sourceLanguage}
								onValueChange={(value) => updateNodeData('sourceLanguage', value)}
							>
								<SelectTrigger>
									{selectedNode.data.sourceLanguage || 'Select source language'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="auto">Auto-detect</SelectItem>
									<SelectItem value="en">English</SelectItem>
									<SelectItem value="ko">Korean</SelectItem>
									<SelectItem value="ja">Japanese</SelectItem>
									<SelectItem value="zh">Chinese</SelectItem>
									<SelectItem value="es">Spanish</SelectItem>
									<SelectItem value="fr">French</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="translation-target">Target Language</Label>
							<Select
								value={selectedNode.data.targetLanguage}
								onValueChange={(value) => updateNodeData('targetLanguage', value)}
							>
								<SelectTrigger>
									{selectedNode.data.targetLanguage || 'Select target language'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="en">English</SelectItem>
									<SelectItem value="ko">Korean</SelectItem>
									<SelectItem value="ja">Japanese</SelectItem>
									<SelectItem value="zh">Chinese</SelectItem>
									<SelectItem value="es">Spanish</SelectItem>
									<SelectItem value="fr">French</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'llm-correction'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							LLM Correction Properties
						</div>

						<div class="space-y-1">
							<Label for="llm-name">Correction Name</Label>
							<Input
								id="llm-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="LLM correction operation name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="llm-model">Model</Label>
							<Select
								value={selectedNode.data.model}
								onValueChange={(value) => updateNodeData('model', value)}
							>
								<SelectTrigger>
									{selectedNode.data.model || 'Select model'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="gpt-4">GPT-4</SelectItem>
									<SelectItem value="gpt-3.5-turbo">GPT-3.5 Turbo</SelectItem>
									<SelectItem value="claude-3">Claude 3</SelectItem>
									<SelectItem value="gemini-pro">Gemini Pro</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div class="space-y-1">
							<Label for="llm-prompt">Correction Prompt</Label>
							<Textarea
								id="llm-prompt"
								value={selectedNode.data.prompt}
								oninput={(e) => handleInputChange('prompt', e)}
								placeholder="Improve grammar and fluency"
								rows={3}
							/>
						</div>
					</div>
				{:else if selectedNode.data.nodeType === 'file-save'}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							File Save Properties
						</div>

						<div class="space-y-1">
							<Label for="file-save-name">Save Name</Label>
							<Input
								id="file-save-name"
								value={selectedNode.data.label}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="File save operation name"
							/>
						</div>

						<div class="space-y-1">
							<Label for="file-save-path">Output Path</Label>
							<Input
								id="file-save-path"
								value={selectedNode.data.outputPath}
								oninput={(e) => handleInputChange('outputPath', e)}
								placeholder="/output/result.txt"
							/>
						</div>

						<div class="space-y-1">
							<Label for="file-save-format">Format</Label>
							<Select
								value={selectedNode.data.format}
								onValueChange={(value) => updateNodeData('format', value)}
							>
								<SelectTrigger>
									{selectedNode.data.format || 'Select format'}
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="txt">Text (.txt)</SelectItem>
									<SelectItem value="json">JSON (.json)</SelectItem>
									<SelectItem value="csv">CSV (.csv)</SelectItem>
									<SelectItem value="pdf">PDF (.pdf)</SelectItem>
									<SelectItem value="docx">Word (.docx)</SelectItem>
								</SelectContent>
							</Select>
						</div>
					</div>
				{:else}
					<div class="space-y-3">
						<div class="text-muted-foreground text-xs uppercase tracking-wide">
							General Properties
						</div>

						<div class="space-y-1">
							<Label for="node-label">Label</Label>
							<Input
								id="node-label"
								value={selectedNode.data.label || ''}
								oninput={(e) => handleInputChange('label', e)}
								placeholder="Node label"
							/>
						</div>
					</div>
				{/if}

				<Separator />

				<!-- Position Info -->
				<div class="space-y-2">
					<div class="text-muted-foreground text-xs uppercase tracking-wide">Position</div>
					<div class="grid grid-cols-2 gap-2">
						<div class="space-y-1">
							<Label for="node-x">X</Label>
							<Input
								id="node-x"
								type="number"
								value={Math.round(selectedNode.position.x)}
								readonly
								class="text-xs"
							/>
						</div>
						<div class="space-y-1">
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
			{/if}
		</div>
	</div>
{/if}

<style>
	/* Custom scrollbar */
	.overflow-y-auto::-webkit-scrollbar {
		width: 6px;
	}

	.overflow-y-auto::-webkit-scrollbar-track {
		background: transparent;
	}

	.overflow-y-auto::-webkit-scrollbar-thumb {
		background-color: hsl(var(--border));
		border-radius: 3px;
	}

	.overflow-y-auto::-webkit-scrollbar-thumb:hover {
		background-color: hsl(var(--muted-foreground));
	}
</style>
