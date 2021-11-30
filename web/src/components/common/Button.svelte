<script lang="ts">
  export let href: string = "";
  export let value: string = "";
  export let onClick: () => void = () => {};
</script>

<button class="w-full rounded pt-2 pb-2 text-gray-500 bg-gray-300 hover:bg-gray-200"  on:click={onClick}>
  {#if href !== ""}
    <a {href}>
      {value}
    </a>
  {:else}
    {value}
  {/if}
</button>
