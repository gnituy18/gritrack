<script lang="ts">
  import type { Step } from "$types";

  export let step: Step;
</script>

<li class="m-2">
  {step.createdAt}
  {step.summary}
  {#each step.items as item}
    <li>{item.type}</li>
  {/each}
</li>
<hr />
