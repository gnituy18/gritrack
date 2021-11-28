import preprocess from 'svelte-preprocess';
import path from 'path';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: [
		preprocess({
			postcss: true
		})
	],
	kit: {
		target: '#root',
		vite: {
			envPrefix: 'ENV_',
			envDir: '../',
			resolve: {
				alias: {
					$: path.resolve('./src'),
					$components: path.resolve('./src/components'),
					$stores: path.resolve('./src/stores')
				}
			}
		}
	}
};

export default config;
