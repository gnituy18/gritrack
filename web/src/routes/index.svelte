<script context="module">
	export async function load() {
		return {
			status: 302,
			redirect: '/login'
		};
	}
</script>

<script>
	import Header from '../lib/Header.svelte';
</script>

<Header />
<h1>Welcome to Gritter</h1>
