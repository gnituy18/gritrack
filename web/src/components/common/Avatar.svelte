<script lang="ts">
  export let src: string;
  export let alt: string;
</script>

<figure>
  <img class="w-12 rounded-full" {src} {alt} />
</figure>
