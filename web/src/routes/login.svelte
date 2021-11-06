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
