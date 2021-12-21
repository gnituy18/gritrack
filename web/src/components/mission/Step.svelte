<script lang="ts">
  import type { Step } from "$types";
  import Item from "$components/mission/Item.svelte";

  export let step: Step;
</script>

<li class="m-2">
  {step.createdAt}
  {step.summary}
  {#each step.items as item}
    <Item {item} />
  {/each}
</li>
<hr />
