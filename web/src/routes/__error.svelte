<script context="module" lang="ts">
  import type { ErrorLoad } from "@sveltejs/kit";

  export const load: ErrorLoad = ({ status }) => {
    return {
      props: {
        title: `${status}`,
      },
    };
  };
</script>

<script lang="ts">
  export let title: string;
</script>

<h1>{title}</h1>
