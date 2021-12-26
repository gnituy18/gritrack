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
  const displayDate: Date = step.createdAt ? new Date(step.createdAt * 1000) : new Date();
  let editingStep: Step = { ...step };

  let showItemForm: boolean = false;

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

  function removeItem(i: number) {
    editingStep.items = [...editingStep.items.slice(0, i), ...editingStep.items.slice(i + 1, editingStep.items.length)];
  }
</script>

<li class="border-gray-100 p-4 hover:bg-gray-100">
  <div class="flex">
    <time
      class="inlint-block border border-slate-300 rounded-full px-2 text-sm bg-slate-200"
      datetime={displayDate.toISOString()}>{displayDate.toLocaleDateString()}</time
    >
    <div class="ml-auto underline cursor-pointer">
      {#if editing}
        <span
          on:click={() => {
            editing = false;
          }}
        >
          cancel
        </span>
      {:else}
        <span
          on:click={() => {
            editing = true;
          }}
        >
          edit
        </span>
      {/if}
    </div>
  </div>
  <div class="mt-2 ml-2">
    {#if editing}
      <div
        class="my-2 rounded-md p-1 summary empty:before:text-gray-400 bg-white"
        contenteditable
        bind:innerHTML={editingStep.summary}
      >
        {editingStep.summary}
      </div>
    {:else}
      <div class="my-1 p-1">{step.summary}</div>
    {/if}
    {#if editing}
      <span
        on:click={() => {
          showItemForm = true;
        }}
        class="underline cursor-pointer">Add an item</span
      >
      {#if showItemForm}
        <ItemForm
          handleAdd={addItem}
          handleCancel={() => {
            showItemForm = false;
          }}
        />
      {/if}
    {/if}
    <ul>
      {#if editing}
        {#each editingStep.items as item, i}
          <ItemDisp {item}>
            <span on:click={() => removeItem(i)} class="underline cursor-pointer">remove</span>
          </ItemDisp>
        {/each}
      {:else}
        {#each step.items as item}
          <ItemDisp {item} />
        {/each}
      {/if}
    </ul>
    {#if editing}
      <div class="flex mt-4 justify-end">
        <Button onClick={submit} value="Submit" />
      </div>
    {/if}
  </div>
</li>

<style>
  .summary:empty::before {
    content: "Today's summary";
  }
</style>
