<script>
	import Header from '$components/header/Header.svelte'
	import '../app.scss';
</script>

<div>
	<Header />
	<slot />
</div>

<style>
	div {
		display: flex;
		height: 100vh;
		width: 100vw;
	}
</style>
