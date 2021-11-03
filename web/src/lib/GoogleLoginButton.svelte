<script lang="ts">
	import Button from '$lib/Button.svelte';
	function signIn() {
		gapi.load('client:auth2', async () => {
			try {
				await gapi.client.init({
					clientId: import.meta.env.VITE_CLIENT_ID,
					scope: 'openid email profile'
				});
				await gapi.auth2.getAuthInstance().signIn();
				console.log('signed');
			} catch (err) {
				console.log(err);
			}
		});
	}
</script>

<svelte:head>
	<script async defer src="https://apis.google.com/js/api.js"></script>
</svelte:head>

<Button on:click={signIn} value="Log in with Google" />
