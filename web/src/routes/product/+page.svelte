<script lang="ts">
  import Button from "$lib/components/common/Button.svelte";
  import Logo from "$lib/components/header/Logo.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<header class="flex justify-between items-center m-2">
  <Logo href="/" />
  {#if data.user}
    <Button href="/login" value={`Continue with ${data.user.name}`} reload />
  {:else}
    <Button href="/login" value="Try Gritrack free" />
  {/if}
</header>
<main>
  <section class="m-4 text-center">
    <h2 class="text-3xl font-semibold">
      Stay focused in a world full of distractions
    </h2>
  </section>
</main>
