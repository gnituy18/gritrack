<script lang="ts">
  import type { Step } from "$lib/types";
  import { steps as storeSteps } from "$lib/stores/mission";
  import StepComp from "$lib/components/mission/Step.svelte";
  import Button from "$lib/components/common/Button.svelte";
  import type { PageData } from "./$types";
  import v1 from "$lib/apis/v1";

  export let data: PageData;

  let isOwner: boolean;
  let steps: Array<Step> = [];
  let count: number = 10;
  let hasMore: boolean = true;

  $: isOwner = data.user.id === data.mission.userId;
  $: steps = $storeSteps;
  $: steps = data.propSteps;

  async function fetchMoreStep() {
    const res = await fetch(
      v1(`/mission/${data.mission.id}/step?offset=${count}&limit=10`),
      {
        credentials: "include",
      }
    );
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
    <h1>{data.mission.name}</h1>
  {/if}
  {#key steps}
    {#if isOwner}
      <StepComp editing mission={data.mission} />
    {/if}
    {#each steps as step}
      <StepComp {step} mission={data.mission} />
    {/each}
  {/key}
</ul>
{#if hasMore}
  <div class="p-4 flex justify-center">
    <Button onClick={fetchMoreStep} value="show more" />
  </div>
{/if}
