<script lang="ts">
  import { missions as storeMissions } from "$lib/stores/mission";
  import Header from "$lib/components/header/Header.svelte";
  import type { LayoutData } from "./$types";

  export let data: LayoutData;
  $storeMissions = data.missions;
</script>

<div class="flex h-screen mx-auto overflow-y-auto">
  <Header />
  <div class="flex-grow">
    <slot />
  </div>
</div>
