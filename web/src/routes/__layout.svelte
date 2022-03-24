<script context="module" lang="ts">
  import type { Load } from "./__layout";
  import v1 from "$/apis/v1";

  export const load: Load = async ({ fetch }) => {
    const res = await fetch(v1("/mission"), { credentials: "include" });
    if (res.status !== 200) {
      return {
        status: res.status,
      };
    }

    const missions = await res.json();

    return {
      props: {
        missions,
      },
    };
  };
</script>

<script lang="ts">
  import type { Mission } from "$/types";
  import "$/app.css";
  import { missions as storeMissions } from "$stores/mission";
  import Header from "$components/header/Header.svelte";

  export let missions: Array<Mission> = [];
  $storeMissions = missions;
</script>

<div class="flex max-w-screen-lg h-screen mx-auto overflow-y-auto">
  <Header />
  <div class="flex-grow">
    <slot />
  </div>
</div>
