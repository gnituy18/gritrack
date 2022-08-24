<script lang="ts">
  import "app.css"
  import { appName } from "$lib/const";
</script>

<svelte:head>
  <title>{appName}</title>
</svelte:head>

<slot />
