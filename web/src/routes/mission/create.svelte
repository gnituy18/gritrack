<script lang="ts">
  import { goto } from "$app/navigation";
  import Button from "$components/common/Button.svelte";

  let name = "";
  let readonly = false;

  async function handleSubmitClick() {
    readonly = true;
    await createMission();
    readonly = false;
  }

  async function createMission() {
    if (!name) {
      console.error("input invalid");
      return;
    }
    const resp = await fetch("http://localhost:8080/api/v1/mission", {
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
    await goto("/");
  }
</script>

<div class="m-4 w-80">
  <h2>Create a Mission</h2>
  <form>
    <label for="name" class="block mt-2">
      <div class="text-gray-500">Name</div>
      <input
        type="text"
        bind:value={name}
        {readonly}
        class="w-full rounded-full bg-gray-100 border-transparent focus:border-blue-300"
      />
    </label>
  </form>
  <div class="mt-8">
    <Button onClick={handleSubmitClick} value="Submit" />
  </div>
</div>
