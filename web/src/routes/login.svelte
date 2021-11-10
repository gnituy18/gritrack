<script context="module">
	export const ssr = false;
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

		if (res.ok) {
			return {
				status: 302,
				redirect: '/'
			};
		}

		return {};
	}
</script>

<script>
	import GoogleLoginButton from '$lib/components/GoogleLoginButton.svelte';
</script>

<div>
	<h1>Welcome to Gritter</h1>
	<GoogleLoginButton />
</div>

<style>
	div {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		height: 100%;
		width: 100%;
	}
</style>
