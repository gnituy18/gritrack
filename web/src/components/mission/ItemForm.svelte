<script lang="ts">
  import type { Item } from "$types";

  export let handleAdd: (item: Item) => void;
  export let handleCancel: () => void;
  export let item: Item = { type: 1, desc: "", time: { duration: 1 } };

  let selected: number = 1;
</script>

<div class="flex items-center my-1">
  <select bind:value={selected} class="rounded border-slate-300 p-1 w-24">
    <option value={1}>⏰ time</option>
  </select>
  {#if item.type === 1}
    <input bind:value={item.time.duration} type="number" min={0} class="ml-2 border-slate-300 rounded p-1 w-12" />
    {#if item.time.duration > 1}
      <span class="ml-1">hours : </span>
    {:else}
      <span class="ml-1">hour : </span>
    {/if}
  {/if}
  <span
    bind:innerHTML={item.desc}
    class="flex-grow desc rounded mx-2 p-1 empty:before:text-gray-400 bg-white"
    contenteditable
  />
  <span on:click={handleCancel} class="mx-1 px-1 underline text-sm cursor-pointer">cancel</span>
  <span
    on:click={() => {
      handleAdd(item);
      handleCancel();
    }}
    class="mx-1 px-1 underline text-sm cursor-pointer">add</span
  >
</div>

<style>
  .desc:empty::before {
    content: "(What have you done?)";
  }
</style>
