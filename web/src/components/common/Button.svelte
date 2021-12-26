<script lang="ts">
  export let href: string = "";
  export let value: string = "";
  export let onClick: () => void = () => {};
</script>

<button
  class="border border-slate-500 rounded-full px-3 text-slate-500 bg-slate-300 hover:bg-slate-200"
  on:click={onClick}
>
  {#if href !== ""}
    <a {href}>
      {value}
    </a>
  {:else}
    {value}
  {/if}
</button>
