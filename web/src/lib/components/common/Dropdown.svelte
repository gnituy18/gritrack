<script lang="ts">
  import type { DropdownItem } from "$lib/types";
  export let classes: string = "";
  export let items: Array<DropdownItem> = [];
</script>

<div class={`relative ${classes} group-dropdown`}>
  <slot />
  <div
    class="absolute bg-gray-100 p-2 rounded invisible z-10 [.group-dropdown:hover_&]:visible"
  >
    {#each items as { label, action }}
      <li
        class="p-1 rounded-2 text-red-500 hover:text-red-700"
        on:click={() => action(label)}
      >
        {label}
      </li>
    {/each}
  </div>
</div>
