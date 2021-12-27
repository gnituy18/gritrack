<script lang="ts">
  import { Item, ItemType } from "$types";

  export let item: Item;
</script>

<li class="flex my-1 rounded hover:bg-slate-200">
  {#if item.type === ItemType.Time}
    <p>⏰ {item.time.duration} {item.time.duration > 1 ? "hours" : "hour"} : {item.desc}</p>
    <div class="ml-auto">
      <slot />
    </div>
  {/if}
</li>
