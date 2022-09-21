<script lang="ts">
  import v1 from "$lib/apis/v1";
  import { page } from "$app/stores";
  import missions from "$lib/stores/mission";
  import Button from "$lib/components/common/Button.svelte";
  import Avatar from "$lib/components/common/Avatar.svelte";
  import Dropdown from "$lib/components/common/Dropdown.svelte";
  import ItemPill from "$lib/components/header/ItemPill.svelte";
  import { goto } from "$app/navigation";

  async function deleteMission(missionId: string) {
    try {
      await missions.delete(missionId);
    } catch (error) {
      console.error(error);
    }

    await goto(`/${$page.data.user.id}`);
  }
</script>

<header
  class="sticky top-0 flex flex-none flex-col justify-between box-border px-4 w-60 h-full border-r"
>
  <nav>
    <div class="m-4">
      <Avatar alt={$page.data.user.name} src={$page.data.user.picture} />
      <h2 class="m-4 text-center">{$page.data.user.name}</h2>
    </div>
    <hr />
    <h2 class="m-4 text-center">Missions</h2>
    <ul class="my-4">
      {#each $missions as { id, name }}
        <ItemPill>
          <a href="/{$page.data.user.id}/{name}" class="block text-lg"
            >{name}
            <Dropdown
              classes="float-right invisible group-hover:visible"
              items={[{ label: "delete", action: () => deleteMission(id) }]}
            >
              ⋯
            </Dropdown>
          </a>
        </ItemPill>
      {/each}
    </ul>
    <a href="/mission/create">
      <ItemPill classes="text-blue-400 text-center">New Mission</ItemPill>
    </a>
  </nav>

  <div class="text-center">
    <Button
      size="s"
      theme="hidden"
      value="Logout"
      onClick={async () => {
        await fetch(v1("/auth/logout"), {
          method: "POST",
          credentials: "include",
        });

        await goto("/");
      }}
    />
  </div>
</header>
