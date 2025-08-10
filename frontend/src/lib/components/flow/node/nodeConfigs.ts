
export type NodeType = 'trigger' | 'branch' | 'action' | 'util';

export interface NodeConfig {
	label: string;
	color: {
		border: string;
		bg: string;
		text: string;
		handle: string;
		hover: string;
		iconBg: string;
		iconText: string;
	};
	minWidth?: string;
	badge?: {
		text: string;
		className: string;
	};
}

// Type colors
const typeColors = {
	trigger: {
		border: 'border-blue-200',
		bg: 'bg-blue-50',
		text: 'text-blue-800',
		handle: 'rgb(59 130 246)', // blue-500
		hover: 'hover:border-blue-300',
		iconBg: 'bg-blue-100',
		iconText: 'text-blue-600'
	},
	branch: {
		border: 'border-purple-200',
		bg: 'bg-purple-50',
		text: 'text-purple-800',
		handle: 'rgb(147 51 234)', // purple-500
		hover: 'hover:border-purple-300',
		iconBg: 'bg-purple-100',
		iconText: 'text-purple-600'
	},
	action: {
		border: 'border-green-200',
		bg: 'bg-green-50',
		text: 'text-green-800',
		handle: 'rgb(34 197 94)', // green-500
		hover: 'hover:border-green-300',
		iconBg: 'bg-green-100',
		iconText: 'text-green-600'
	},
	util: {
		border: 'border-gray-200',
		bg: 'bg-gray-50',
		text: 'text-gray-800',
		handle: 'rgb(107 114 128)', // gray-500
		hover: 'hover:border-gray-300',
		iconBg: 'bg-gray-100',
		iconText: 'text-gray-600'
	}
};

export const nodeConfigs: Record<'trigger' | 'branch' | 'action' | 'util', NodeConfig> = {
	trigger: {
		label: 'Trigger',
		color: typeColors.trigger,
		minWidth: 'min-w-[200px]'
	},
	branch: {
		label: 'Branch',
		color: typeColors.branch,
		minWidth: 'min-w-[200px]'
	},
	action: {
		label: 'Action',
		color: typeColors.action,
		minWidth: 'min-w-[200px]'
	},
	util: {
		label: 'Utility',
		color: typeColors.util,
		minWidth: 'min-w-[200px]'
	}
};