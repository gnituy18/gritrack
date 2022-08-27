<script lang="ts">
  import { appName } from "$lib/const";
  import "../app.css";
</script>

<svelte:head>
  <title>{appName}</title>
</svelte:head>

<slot />
