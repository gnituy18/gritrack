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
        propSteps: steps,
      },
    };
  }
</script>

<script lang="ts">
  import type { Step, Mission } from "$types";
  import { steps as storeSteps } from "$stores/mission";
  import StepComp from "$components/mission/Step.svelte";

  export let mission: Mission;
  export let propSteps: Array<Step>;

  let steps: Array<Step> = [];
  let noStepToday: boolean = true;

  $: steps = $storeSteps;
  $: steps = propSteps;
  $: noStepToday = steps.length === 0 || !isToday(steps[0].createdAt);

  function isToday(ts: number): boolean {
    return new Date().toLocaleDateString() === new Date(ts * 1000).toLocaleDateString();
  }
</script>

<ul>
  {#key steps}
    {#if noStepToday}
      <StepComp editing {mission} />
    {/if}
    {#each steps as step}
      <StepComp {step} {mission} />
    {/each}
  {/key}
</ul>
