<script context="module" lang="ts">
import type { Load } from '@sveltejs/kit';
export const load: Load = ({ status }) => {
    return {
      props: {
        title: `${status}`,
      },
    };
  }
</script>

<script>
  export let title;
</script>

<h1>{title}</h1>
