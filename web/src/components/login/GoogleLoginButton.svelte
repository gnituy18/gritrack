<script lang="ts">
  import { env } from "$env/dynamic/public";
  import Button from "$components/common/Button.svelte";
  import v1 from "$apis/v1";

  function signIn() {
    gapi.load("client:auth2", async () => {
      try {
        await gapi.client.init({
          clientId: env.PUBLIC_GOOGLE_CLIENT_ID,
          scope: "openid email profile",
        });
        const user = await gapi.auth2.getAuthInstance().signIn();
        const accessToken = user.getAuthResponse().access_token;
        const resp = await fetch(v1("/auth"), {
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
          window.location.href = "/home";
        }
      } catch (err) {
        console.error(err);
      }
    });
  }
</script>

<svelte:head>
  <script async defer src="https://apis.google.com/js/api.js"></script>
</svelte:head>

<Button onClick={signIn} value="Log in with Google" />
