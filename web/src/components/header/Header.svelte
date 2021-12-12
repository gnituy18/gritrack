<script>
  import { onMount } from "svelte";

  import Logo from "$components/header/Logo.svelte";
  import Button from "$components/common/Button.svelte";

  let missions = [];
  onMount(async () => {
    const resp = await fetch("http://localhost:8080/api/v1/mission", { credentials: "include" });
    missions = await resp.json();
  });
</script>

<header class="sticky top-0 flex-none box-border ml-4 w-60 h-full border-r">
  <nav>
    <div class="flex items-center my-4 mr-4">
      <Logo />
    </div>
    <div class="mt-8 mr-4">
      <Button href="/mission/create" value="Create a mission" />
    </div>
    <h2 class="mt-8">Missions</h2>
    <ul class="my-4 mr-4">
      {#each missions as { id, name }}
        <li class="mt-1 rounded-full px-4 hover:cursor-pointer hover:bg-gray-100 ">
          <a href="/mission/{id}" class="block text-lg">{name}</a>
        </li>
      {/each}
    </ul>
  </nav>
</header>
