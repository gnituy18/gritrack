<script lang="ts">
  throw new Error("@migration task: Add data prop (https://github.com/sveltejs/kit/discussions/5774#discussioncomment-3292707)");

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
    {#if isOwner}
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
