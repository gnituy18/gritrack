<script lang="ts">
  import type { Item } from "$lib/types";
  import { ItemType } from "$lib/types";

  export let item: Item;
  let duration: number = 0;
  $: duration = item.time.duration / 60;
</script>

<li class="flex my-1 rounded hover:bg-slate-200">
  {#if item.type === ItemType.Time}
    <p>⏰ {duration} {duration > 1 ? "minutes" : "minute"} : {item.desc}</p>
    <div class="ml-auto">
      <slot />
    </div>
  {/if}
</li>
