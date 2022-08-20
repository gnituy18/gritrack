<script>
  import Button from "$/components/common/Button.svelte";
  import Logo from "$components/header/Logo.svelte";
</script>

<header class="flex justify-between items-center m-2">
  <Logo href="/" />
  <Button href="/login" value="Try Gritter free" />
</header>
<main>
  <section class="m-4 text-center">
    <h2 class="text-3xl font-semibold">Stay focused in a world full of distractions</h2>
  </section>
</main>
