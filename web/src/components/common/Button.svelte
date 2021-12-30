<script lang="ts">
  export let href: string = "";
  export let value: string = "";
  export let onClick: () => void = () => {};
</script>

<a {href}>
  <button
    class="border border-slate-500 rounded-full px-3 text-slate-500 bg-slate-300 hover:bg-slate-200"
    on:click={onClick}
  >
    {value}
  </button>
</a>
