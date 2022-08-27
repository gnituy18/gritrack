<script lang="ts">
  import v1 from "$lib/apis/v1";
  import { page } from "$app/stores";
  import { missions } from "$lib/stores/mission";
  import Button from "$lib/components/common/Button.svelte";
  import Avatar from "$lib/components/common/Avatar.svelte";
  import Dropdown from "$lib/components/common/Dropdown.svelte";
  import type { Mission } from "$lib/types";
  import { goto } from "$app/navigation";

  async function deleteMission(missionId: string) {
    const res = await fetch(v1(`/mission/${missionId}`), {
      method: "DELETE",
      credentials: "include",
    });

    if (!res.ok) {
      console.error("create failed");
      return;
    }

    const newMissions: Array<Mission> = await (
      await fetch(v1("/mission"), { credentials: "include" })
    ).json();
    $missions = newMissions;

    await goto(`/${$page.data.user.id}`);
  }
</script>

<header class="sticky top-0 flex-none box-border px-4 w-60 h-full border-r">
  <nav>
    <div class="flex justify-between items-center my-4">
      <Avatar alt={$page.data.user.name} src={$page.data.user.picture} />
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
          <a href="/{$page.data.user.id}/{name}" class="block text-lg"
            >{name}
            <Dropdown
              classes="float-right"
              items={[{ label: "delete", action: () => deleteMission(id) }]}
            >
              ...
            </Dropdown>
          </a>
        </li>
      {/each}
    </ul>
  </nav>
</header>
