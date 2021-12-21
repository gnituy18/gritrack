<script lang="ts">
  import { Item, ItemType } from "$types";

  export let item: Item;
</script>

<li>
  {#if item.type === ItemType.Time}
    ⏰{item.time.duration} hours
  {/if}
</li>
