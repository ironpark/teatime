<script lang="ts">
	import { onMount } from 'svelte';
	import { SidebarTrigger } from '$lib/components/ui/sidebar';
	import { Separator } from '$lib/components/ui/separator';
	import { cn } from '$lib/utils';
    type Props = {
        title?: string;
        icon?: any;
        children: any;
        actions?: () => any;
        class?: string;
        mainClass?: string;
    };

    let {
        title,
        icon,
        children,
        actions,
        class: className,
        mainClass: mainClassName
    }: Props = $props();

    const classes = cn('settings-page bg-background flex h-screen w-full flex-col', className);
    const mainClasses = cn('flex-1 overflow-y-auto p-6', mainClassName);
</script>


<div class={classes}>
	<!-- Header -->
	<header class="bg-card border-b">
		<div class="flex h-16 items-center gap-2 px-4">
			<SidebarTrigger />
			<Separator orientation="vertical" class="mr-2 h-4" />
            {#if icon}
                <div class="flex items-center gap-2">
                    {#if icon}
                        {@const Icon = icon}
                        <Icon class="text-muted-foreground h-5 w-5" />
                    {/if}
                    <h1 class="text-lg font-semibold">{title || ''}</h1>
                </div>
            {/if}
            <div class="flex-1"></div>
            {#if actions}
                {@render actions?.()}
            {/if}
		</div>
	</header>

	<!-- Main content -->
	<main class={mainClasses}>
        {@render children?.()}
	</main>
</div>


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
