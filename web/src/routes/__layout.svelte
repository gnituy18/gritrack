<script>
  import "$/app.css";
  import Header from "$components/header/Header.svelte";
</script>

<div class="flex max-w-screen-lg h-screen  mx-auto">
  <Header />
  <div class="flex-grow">
    <slot />
  </div>
</div>
