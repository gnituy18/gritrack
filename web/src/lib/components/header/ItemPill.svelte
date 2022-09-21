<script lang="ts">
  export let classes: string = "";
</script>

<li
  class={`mt-4 rounded-full px-4 py-2 list-none hover:cursor-pointer group hover:bg-blue-50 bg-slate-50 ${classes}`}
>
  <slot />
</li>
