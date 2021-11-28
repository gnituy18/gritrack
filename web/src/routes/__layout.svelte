<script>
  import "$/app.css";
  import Header from "$components/header/Header.svelte";
</script>

<div class="flex w-screen h-screen">
  <Header />
  <div class="flex-grow">
    <slot />
  </div>
</div>
