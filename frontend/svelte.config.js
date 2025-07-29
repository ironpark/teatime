import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: vitePreprocess(),
	kit: {
		// adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		adapter: adapter({
            pages: 'dist',
			assets: 'dist',
			fallback: undefined,
			precompress: false,
			strict: true
        }),
		alias: {
			'$ipc': "./bindings/github.com/ironpark/teatime",
			'$ipc/*': "./bindings/github.com/ironpark/teatime/*",
			"$lib/*": "./src/lib/*",
			"$lib": "./src/lib",
			"$components/*": "./src/lib/components/*",
			"$components": "./src/lib/components",
			"$ui/*": "./src/lib/components/ui/*",
			"$ui": "./src/lib/components/ui",
		}
	}
};

export default config;
