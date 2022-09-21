<script lang="ts">
  import { afterUpdate, onMount } from "svelte";
  import { env } from "$env/dynamic/public";
  import v1 from "$lib/apis/v1";
  import { goto } from "$app/navigation";

  onMount(() => {
    google.accounts.id.initialize({
      client_id: env.PUBLIC_GOOGLE_CLIENT_ID,
      callback: async ({ credential }) => {
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
              google: { idToken: credential },
            }),
          });
          if (resp.type === "opaqueredirect") {
            await goto("/");
          }
        } catch (e) {
          console.error(e);
        }
      },
    });
  });

  afterUpdate(() => {
    google.accounts.id.renderButton(
      document.getElementById("googleSignInButton"),
      { type: "standard", theme: "outline", size: "large" }
    );
  });
</script>

<div id="googleSignInButton" />
