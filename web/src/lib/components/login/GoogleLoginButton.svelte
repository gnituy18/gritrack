<script lang="ts">
  import { onMount } from "svelte";
  import { env } from "$env/dynamic/public";
  import v1 from "$lib/apis/v1";
  import { goto } from "$app/navigation";

  async function handleCredentialResponse(response) {

    try {
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
          google: { accessToken: response.credential },
        }),
      });
      if (resp.type === "opaqueredirect") {
        await goto("/");
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(() => {
    google.accounts.id.initialize({
      client_id: env.PUBLIC_GOOGLE_CLIENT_ID,
      callback: handleCredentialResponse,
    });
    google.accounts.id.renderButton(
      document.getElementById("googleSignInButton"),
      { theme: "outline", size: "large" } // customization attributes
    );
  });
</script>

<svelte:head>
  <script src="https://accounts.google.com/gsi/client" async defer></script>
</svelte:head>

<div id="googleSignInButton" />
