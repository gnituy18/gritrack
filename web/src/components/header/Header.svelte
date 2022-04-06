<script>
  import v1 from "$/apis/v1";
  import { session } from "$app/stores";
  import { missions } from "$stores/mission";
  import Button from "$components/common/Button.svelte";
  import Avatar from "$components/common/Avatar.svelte";
</script>

<header class="sticky top-0 flex-none box-border px-4 w-60 h-full border-r">
  <nav>
    <div class="flex justify-between items-center my-4">
      <Avatar alt={$session.currentUser.name} src={$session.currentUser.picture} />
      <Button
        size="s"
		theme="hidden"
        value="Logout"
        onClick={async () => {
          await fetch(v1("/auth/logout"), {
            method: "POST",
            credentials: "include",
          });

          window.location.href = "/";
        }}
      />
    </div>
    <div class="flex items-baseline">
      <h2 class="mr-auto">Missions</h2>
      <Button theme="hidden" size="s" href="/mission/create" value="Create" />
    </div>
    <ul class="my-4">
      {#each $missions as { id, name }}
        <li class="mt-1 rounded px-2 hover:cursor-pointer hover:bg-slate-50 ">
          <a href="/mission/{id}" class="block text-lg">{name}</a>
        </li>
      {/each}
    </ul>
  </nav>
</header>
