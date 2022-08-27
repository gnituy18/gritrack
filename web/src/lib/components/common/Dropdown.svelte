<script lang="ts">
  import type { DropdownItem } from "$lib/types";
  export let classes: string = "";
  export let items: Array<DropdownItem> = [];
</script>

<div class={`group relative ${classes}`}>
  <slot />
  <div class="absolute bg-gray-100 p-2 rounded hidden z-10 group-hover:block">
    {#each items as { label, action }}
      <li
        class="p-1 rounded text-red-500 hover:text-red-700"
        on:click={() => action(label)}
      >
        {label}
      </li>
    {/each}
  </div>
</div>
