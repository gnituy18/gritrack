<script lang="ts">
  import steps from "$lib/stores/step";
  import Step from "$lib/components/mission/Step.svelte";
  import Button from "$lib/components/common/Button.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;

  let isOwner: boolean = data.user.id === data.mission.userId;

  $: steps.set(data.mission.id, data.steps, data.more, data.currentOffset);
</script>

<ul class="divide-y-2">
  {#if !isOwner}
    <h1>{data.mission.name}</h1>
  {/if}
  {#key $steps}
    {#if isOwner}
      <Step editing mission={data.mission} />
    {/if}
    {#each $steps as step}
      <Step {step} mission={data.mission} />
    {/each}
  {/key}
</ul>
{#if steps.hasMore()}
  <div class="p-4 flex justify-center">
    <Button onClick={() => steps.updateMore(7)} value="show more" />
  </div>
{/if}
