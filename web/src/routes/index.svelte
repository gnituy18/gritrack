<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ fetch }) {
		const res = await fetch('http://localhost:8080/api/v1/user/current', {
			credentials: 'include'
		});

		res.headers.forEach((a, b) => {
			console.log(a, b);
		});

		if (!res.ok) {
			return {
				status: 302,
				redirect: '/login'
			};
		}

		return {};
	}
</script>

<script>
	import Header from '$lib/components/Header.svelte';
</script>

<Header />
<h1>Welcome to Gritter</h1>
