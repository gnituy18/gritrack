<script lang="ts">
	function load(_node: HTMLElement) {
		gapi.load('client:auth2', async () => {
			try {
				await gapi.client.init({
					clientId: import.meta.env.VITE_CLIENT_ID,
					scope: 'openid email profile'
				});
			} catch (err) {
				console.log(err);
			}
		});
	}

	async function signIn() {
		if (!gapi.auth2) {
			return;
		}

		await gapi.auth2.getAuthInstance().signIn();
		console.log('signed');
	}
</script>

<svelte:head>
	<script async defer src="https://apis.google.com/js/api.js" use:load></script>
</svelte:head>

<button on:click={signIn}> Log in </button>
