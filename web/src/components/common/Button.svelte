<script lang="ts">
  export let href: string = "";
  export let value: string = "";
  export let onClick: () => Promise<void> | void = () => {};
</script>

<a class="shrink-0" {href}>
  <button class="border rounded-md px-5 py-2 text-white bg-blue-500 font-medium hover:bg-blue-600" on:click={onClick}>
    {value}
  </button>
</a>
