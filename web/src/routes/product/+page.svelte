<script lang="ts">
  import type { PageServerData } from "./$types";
  import Button from "$lib/components/common/Button.svelte";
  import Logo from "$lib/components/header/Logo.svelte";

  export let data: PageServerData;
</script>

<header class="flex justify-between items-center m-2">
  <Logo href="/" />
  {#if data}
    <Button href="/login" value={`Continue with ${page.currentUser.name}`} reload />
  {:else}
    <Button href="/login" value="Try Gritter free" />
  {/if}
</header>
<main>
  <section class="m-4 text-center">
    <h2 class="text-3xl font-semibold">Stay focused in a world full of distractions</h2>
  </section>
</main>
