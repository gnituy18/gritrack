<script context="module">
  /** @type {import('@sveltejs/kit').ErrorLoad} */
  export function load({ status }) {
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
