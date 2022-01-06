<script lang="ts" context="module">
  import type { Load } from "@sveltejs/kit";
  import v1 from "$apis/v1";

  export const load: Load = async ({ params, fetch }) => {
    try {
      const missionId = params.missionId;
      let res = await fetch(v1(`/mission/${missionId}`), {
        credentials: "include",
      });
      if (res.status !== 200) {
        return {
          status: res.status,
        };
      }
      const mission = await res.json();

      res = await fetch(v1(`/mission/${missionId}/step?offset=0&limit=10`), {
        credentials: "include",
      });
      if (res.status !== 200) {
        return {
          status: res.status,
        };
      }
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
        status: 500,
      };
    }
  };
</script>

<script lang="ts">
  import type { Step, Mission } from "$types";
  import { steps as storeSteps } from "$stores/mission";
  import { session } from "$app/stores";
  import StepComp from "$components/mission/Step.svelte";
  import Button from "$/components/common/Button.svelte";

  export let mission: Mission;
  export let propSteps: Array<Step>;

  let isOwner: boolean;
  let steps: Array<Step> = [];
  let noStepToday: boolean = true;
  let count: number = 10;
  let hasMore: boolean = true;

  $: isOwner = $session.currentUser.id === mission.userId;
  $: steps = $storeSteps;
  $: steps = propSteps;
  $: noStepToday = steps.length === 0 || !isToday(steps[0].createdAt);

  function isToday(ts: number): boolean {
    return new Date().toLocaleDateString() === new Date(ts * 1000).toLocaleDateString();
  }

  async function fetchMoreStep() {
    const res = await fetch(v1(`/mission/${mission.id}/step?offset=${count}&limit=10`), {
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
  {#if !isOwner}
    <h1>{mission.name}</h1>
  {/if}
  {#key steps}
    {#if noStepToday && isOwner}
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
