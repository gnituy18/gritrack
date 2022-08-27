<script lang="ts">
  import type { Mission } from "$lib/types";
  import { missions as storeMissions } from "$lib/stores/mission";
  import Header from "$lib/components/header/Header.svelte";

  export let missions: Array<Mission> = [];
  $storeMissions = missions;
</script>

<div class="flex h-screen mx-auto overflow-y-auto">
  <Header />
  <div class="flex-grow">
    <slot />
  </div>
</div>
