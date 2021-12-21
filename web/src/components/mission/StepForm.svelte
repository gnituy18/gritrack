<script lang="ts">
  import Button from "$components/common/Button.svelte";
  import type { Mission, Step } from "$types";

  export let mission: Mission;
  export let step: Step = { id: "", summary: "", items: [], createdAt: 0 };

  async function submit() {
    let res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({
        summary: step.summary,
        items: step.items,
      }),
    });

    if (res.status !== 201) {
      console.error(await res.text());
    }

    step.items = [];
    step.summary = "";
    duration = 0;
  }

  function addItem() {
    step.items = [{ type: 0, desc, time: { duration } }, ...step.items];
    duration = 0;
    desc = "";
  }

  let duration: number;
  let desc: string = "";
</script>

<div class="m-2 p-2 summary empty:before:text-gray-400" contenteditable bind:innerHTML={step.summary}>
  {step.summary}
</div>
<ul>
  <div class="flex items-center">
    <span on:click={addItem} class="m-2">+</span>
    <select>
      <option value="0" selected>⏰ time</option>
    </select>
    <input bind:value={duration} type="number" />
    <span>hours</span>
    <span bind:innerHTML={desc} class="mx-2 p-2 flex-grow" contenteditable />
  </div>
  <li class="p-2">
    {#each step.items as item}
      <li>⏰{item.time.duration} hours</li>
    {/each}
  </li>
</ul>
<Button onClick={submit} value="Submit" />

<style>
  .summary:empty::before {
    content: "What have you done?";
  }
</style>
