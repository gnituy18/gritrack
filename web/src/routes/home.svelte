<script context="module" lang="ts">
  import type { Load } from "./__layout";

  export const load: Load = async ({ session }) => {
    return {
      status: 302,
      redirect: `/${session.currentUser.id}/`,
    };
  };
</script>

<script lang="ts">
  import type { Mission } from "$/types";
  import { missions as storeMissions } from "$stores/mission";
  import Header from "$components/header/Header.svelte";

  export let missions: Array<Mission> = [];
  $storeMissions = missions;
</script>

<div class="flex h-screen mx-auto overflow-y-auto">
  <Header />
  <div class="flex-grow">
    <slot />
  </div>
</div>
