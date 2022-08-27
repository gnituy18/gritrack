<script lang="ts">
  import type { Mission, Step, Item } from "$lib/types";
  import v1 from "$lib/apis/v1";
  import { steps as storeSteps } from "$lib/stores/mission";
  import { page } from "$app/stores";
  import ItemDisp from "$lib/components/mission/ItemDisp.svelte";
  import Button from "$lib/components/common/Button.svelte";
  import ItemForm from "$lib/components/mission/ItemForm.svelte";

  export let editing: boolean = false;
  export let mission: Mission;
  export let step: Step = {
    date: undefined,
    id: undefined,
    summary: "",
    items: [],
    createdAt: undefined,
  };

  const isNew = step.createdAt === undefined;
  let date = isNew
    ? formatDate(new Date())
    : formatDate(new Date(step.time * 1000));

  let isOwner: boolean;
  let editingStep: Step = { ...step };
  let showItemForm: boolean = false;
  $: isOwner = $page.data.user.id === mission.userId;
  $: editingStep = { ...editingStep, date };

  function formatDate(date: Date) {
    return `${date.getFullYear()}-${addZero(date.getMonth() + 1)}-${addZero(
      date.getDate()
    )}`;
  }

  function addZero(num: number) {
    return num < 10 ? `0${num}` : num;
  }

  async function submit() {
    if (isNew) {
      let res = await fetch(v1(`/mission/${mission.id}/step`), {
        method: "POST",
        credentials: "include",
        body: JSON.stringify({
          date: editingStep.date,
          summary: editingStep.summary,
          items: editingStep.items,
        }),
      });

      if (!res.ok) {
        console.error(await res.text());
      }
    } else {
      let res = await fetch(v1(`/mission/${mission.id}/step/${step.id}`), {
        method: "PUT",
        credentials: "include",
        body: JSON.stringify({
          summary: editingStep.summary,
          items: editingStep.items,
        }),
      });

      if (!res.ok) {
        console.error(await res.text());
      }
    }

    const res = await fetch(
      v1(`/mission/${mission.id}/step?offset=0&limit=10`),
      {
        credentials: "include",
      }
    );
    $storeSteps = await res.json();
    editing = false;
  }

  function addItem(item: Item) {
    editingStep.items = [item, ...editingStep.items];
  }

  function removeItem(i: number) {
    editingStep.items = [
      ...editingStep.items.slice(0, i),
      ...editingStep.items.slice(i + 1, editingStep.items.length),
    ];
  }
</script>

<li class="border-gray-100 p-4 hover:bg-slate-50">
  <div class="flex">
    {#if isNew}
      <input type="date" bind:value={date} />
    {:else}
      <time
        on:click={() => {
          editing = true;
          editingStep = { ...step };
        }}
        class="inlint-block border border-slate-300 rounded-full px-2 text-sm bg-slate-200"
      >
        {date}
      </time>
    {/if}
    {#if isOwner}
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
    {/if}
  </div>
  <div class="mt-2 ml-2">
    {#if editing}
      <div
        class="my-2 rounded p-1 summary empty:before:text-gray-400 bg-white"
        contenteditable
        bind:innerHTML={editingStep.summary}
      >
        {editingStep.summary}
      </div>
    {:else}
      <div class="my-1 p-1">{@html step.summary}</div>
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
            <span
              on:click={() => removeItem(i)}
              class="underline cursor-pointer">remove</span
            >
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
