<script lang="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
  export async function load({ page, fetch }) {
    try {
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
    } catch (error) {
      console.error(error);
      return {
        status: 404,
      };
    }
  }
</script>

<script lang="ts">
  import type { Step, Mission } from "$types";
  import { steps as storeSteps } from "$stores/mission";
  import StepComp from "$components/mission/Step.svelte";
  import Button from "$/components/common/Button.svelte";

  export let mission: Mission;
  export let propSteps: Array<Step>;

  let steps: Array<Step> = [];
  let noStepToday: boolean = true;
  let count: number = 10;
  let hasMore: boolean = true;

  $: steps = $storeSteps;
  $: steps = propSteps;
  $: noStepToday = steps.length === 0 || !isToday(steps[0].createdAt);

  function isToday(ts: number): boolean {
    return new Date().toLocaleDateString() === new Date(ts * 1000).toLocaleDateString();
  }

  async function fetchMoreStep() {
    const res = await fetch(`http://localhost:8080/api/v1/mission/${mission.id}/step?offset=${count}&limit=10`, {
      credentials: "include",
    });
    const moreSteps = await res.json();
    steps = [...steps, ...moreSteps];
    count += 10;
    if (moreSteps.length < 10) {
      hasMore = false;
    }
  }
</script>

<ul class="divide-y-2">
  {#key steps}
    {#if noStepToday}
      <StepComp editing {mission} />
    {/if}
    {#each steps as step}
      <StepComp {step} {mission} />
    {/each}
  {/key}
</ul>
{#if hasMore}
  <div class="p-4 flex justify-center">
    <Button onClick={fetchMoreStep} value="show more" />
  </div>
{/if}
