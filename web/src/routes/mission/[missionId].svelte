<script lang="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
  export async function load({ page, fetch }) {
    const missionId = page.params.missionId;
    let res = await fetch(`http://localhost:8080/api/v1/mission/${missionId}`, {
      credentials: "include",
    });
    const mission = await res.json();

    res = await fetch(`http://localhost:8080/api/v1/mission/${missionId}/step?offset=0&limit=5`, {
      credentials: "include",
    });
    const steps = await res.json();

    console.log(mission);
    console.log(steps);

    return {
      props: {
        mission,
        steps,
      },
    };
  }
</script>

<script lang="ts">
  // import { goto } from "$app/navigation";
  import Button from "$components/common/Button.svelte";
  import type { Step, Mission } from "./type";
  export let steps: Array<Step>;
  export let mission: Mission;

  let summary = "";

  async function submit() {
    let res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({
        summary,
        items: [],
      }),
    });

    console.log(await res.text());
    if (res.status !== 201) {
      console.error(await res.text());
    }

    res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step?offset=0&limit=5`, {
      credentials: "include",
    });
    const newSteps = await res.json();
    steps = newSteps;
    summary = "";
  }
</script>

<div class="m-2 p-2 summary empty:before:text-gray-400" contenteditable bind:innerHTML={summary}>
  {summary}
</div>
<Button onClick={submit} value="Submit" />
<ul>
  {#each steps as { id, summary, items }}
    <li>
      {id}
      {summary}
      {#each items as item}
        <li>{item.type}</li>
      {/each}
    </li>
  {/each}
</ul>

<style>
  .summary:empty::before {
    content: "What have you done?";
  }
</style>
