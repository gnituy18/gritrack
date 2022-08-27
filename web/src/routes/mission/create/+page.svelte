<script lang="ts">
  import type { Mission } from "$lib/types";
  import v1 from "$lib/apis/v1";
  import { missions } from "$lib/stores/mission";
  import { goto } from "$app/navigation";
  import Button from "$lib/components/common/Button.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;

  let name = "";
  let readonly = false;

  async function handleSubmitClick() {
    readonly = true;
    await createMission();
    readonly = false;
  }

  async function createMission() {
    console.log(data);
    if (!name) {
      console.error("input invalid");
      return;
    }
    const resp = await fetch(v1("/mission"), {
      method: "POST",
      credentials: "include",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({
        name: name,
      }),
    });
    if (resp.status !== 201) {
      console.error("create failed");
      return;
    }

    const newMissions: Array<Mission> = await (
      await fetch(v1("/mission"), { credentials: "include" })
    ).json();
    $missions = newMissions;

    await goto(`/${data.user.id}/${name}`);
  }
</script>

<div class="m-4 w-full">
  <h2>
    Set a mission that's really important to you and you want to spend more than
    a year on it.
  </h2>
  <form>
    <label for="name" class="block mt-2">
      <div class="text-gray-500">Name</div>
      <input
        type="text"
        bind:value={name}
        {readonly}
        class="w-80 rounded p-1 bg-gray-100 border-transparent focus:border-blue-300 focus:border"
      />
    </label>
  </form>
  <div class="mt-8">
    <Button onClick={handleSubmitClick} value="Submit" />
  </div>
</div>
