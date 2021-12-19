<script lang="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
  export async function load({ page, fetch }) {
    const missionId = page.params.missionId;
    let res = await fetch(`http://localhost:8080/api/v1/mission/${missionId}`, {
      credentials: "include",
    });
    const mission = await res.json();

    res = await fetch(`http://localhost:8080/api/v1/mission/${missionId}/step?offset=0&limit=10`, {
      credentials: "include",
    });
    const steps = await res.json();

    return {
      props: {
        mission,
        steps,
      },
    };
  }
</script>

<script lang="ts">
  import type { Step, Mission } from "$types";
  import StepForm from "$components/mission/StepForm.svelte";
  import StepItem from "$components/mission/Step.svelte";

  export let mission: Mission;
  export let steps: Array<Step>;

  let editingStepId: string = "";
  function setEditingStepId(stepId: string) {
    editingStepId = stepId;
  }
</script>

<StepForm {mission} />
<ul class="mt-20">
  {#each steps as step}
    {#if step.id === editingStepId}
      <div on:click={() => setEditingStepId("")}>cancel</div>
      <StepForm {mission} />
    {:else}
      <div on:click={() => setEditingStepId(step.id)}>edit</div>
      <StepItem {step} />
    {/if}
  {/each}
</ul>
