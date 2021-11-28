<script lang="ts">
  import Button from "$components/common/Button.svelte";
  function signIn() {
    gapi.load("client:auth2", async () => {
      try {
        await gapi.client.init({
          clientId: import.meta.env.ENV_GOOGLE_CLIENT_ID,
          scope: "openid email profile",
        });
        const user = await gapi.auth2.getAuthInstance().signIn();
        const accessToken = user.getAuthResponse().access_token;
        const resp = await fetch("http://localhost:8080/api/v1/auth", {
          method: "POST",
          redirect: "manual",
          credentials: "include",
          headers: {
            "content-type": "application/json",
            "Access-Control-Request-Method": "POST",
            "Access-Control-Request-Headers": "Content-Type",
          },
          body: JSON.stringify({
            type: 1,
            google: { accessToken },
          }),
        });
        if (resp.type === "opaqueredirect") {
          window.location.href = "/";
        }
      } catch (err) {
        console.log(err);
      }
    });
  }
</script>

<svelte:head>
  <script async defer src="https://apis.google.com/js/api.js"></script>
</svelte:head>

<Button onClick={signIn} value="Log in with Google" />
