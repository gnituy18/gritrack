<script lang="ts">
  import type { Item } from "$types";

  export let addItem: (item: Item) => void = () => {};

  export let item: Item = { type: 1, desc: "", time: { duration: 1 } };

  let selected: number = 1;
</script>

<div class="flex items-center">
  <select bind:value={selected}>
    <option value={1}>⏰ time</option>
  </select>
  {#if item.type === 1}
    <input bind:value={item.time.duration} type="number" />
  {/if}
  <span>hours</span>
  <span bind:innerHTML={item.desc} class="mx-2 p-2 flex-grow" contenteditable />
  <span on:click={() => addItem(item)} class="m-2">+</span>
</div>
