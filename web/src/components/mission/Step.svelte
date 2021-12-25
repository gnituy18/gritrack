<script lang="ts">
  import type { Mission, Step, Item } from "$types";
  import { steps as storeSteps } from "$stores/mission";
  import ItemDisp from "$components/mission/ItemDisp.svelte";
  import Button from "$components/common/Button.svelte";
  import ItemForm from "$components/mission/ItemForm.svelte";

  export let editing: boolean = false;
  export let mission: Mission;
  export let step: Step = { id: undefined, summary: "", items: [], createdAt: undefined };
  const isNew = step.createdAt === undefined;
  const currentTs: number = (step as Step).createdAt || Math.floor(new Date().getTime() / 1000);
  let editingStep: Step = { ...step };

  console.log(step.items)

  async function submit() {
    if (isNew) {
      let res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step`, {
        method: "POST",
        credentials: "include",
        body: JSON.stringify({
          summary: editingStep.summary,
          items: editingStep.items,
        }),
      });

      if (res.status !== 201) {
        console.error(await res.text());
      }
    } else {
	  let res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step/${step.id}`, {
        method: "PUT",
        credentials: "include",
        body: JSON.stringify({
          summary: editingStep.summary,
          items: editingStep.items,
        }),
      });

      if (res.status !== 200) {
        console.error(await res.text());
      }
    }

    const res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step?offset=0&limit=10`, {
      credentials: "include",
    });
    $storeSteps = await res.json();
    editing = false;
  }

  function addItem(item: Item) {
    editingStep.items = [item, ...editingStep.items];
  }

  function getDisplayTime(ts: number): string {
    return new Date(ts * 1000).toLocaleDateString();
  }
</script>

<li class="mt-2 p-4 hover:bg-gray-100">
  {#if editing}
    <div class="m-2 p-2 summary empty:before:text-gray-400" contenteditable bind:innerHTML={editingStep.summary}>
      {editingStep.summary}
    </div>
    <ul>
      <ItemForm {addItem} />
      <li class="p-2">
        {#each editingStep.items as item}
          <li>⏰{item.time.duration} hours</li>
        {/each}
      </li>
    </ul>
    <Button onClick={submit} value="Submit" />
    <div
      on:click={() => {
        editing = false;
      }}
    >
      cancel
    </div>
  {:else}
    {getDisplayTime(currentTs)}
    {step.summary}
    {#each step.items as item}
      <ItemDisp {item} />
    {/each}
    <div
      on:click={() => {
        editing = true;
      }}
    >
      edit
    </div>
  {/if}
</li>

<style>
  .summary:empty::before {
    content: "What have you done?";
  }
</style>
