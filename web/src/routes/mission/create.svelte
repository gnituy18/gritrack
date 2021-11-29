<script lang="ts">
  import Button from "$components/common/Button.svelte";

  let name = "";
  let description = "";
  let readonly = false

  async function handleSubmitClick() {
    readonly = true;
    await createMission();
  }

  async function createMission() {
    const resp = await fetch("http://localhost:8080/api/v1/mission", {
      method: "POST",
      credentials: "include",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        description: description,
      }),
    });
    console.log(resp);
  }
</script>

<h2>Create a Mission</h2>
<form>
  <div>
    <label for="name">Name</label>
	<input type="text" bind:value={name} {readonly} />
  </div>
  <div>
    <label for="description">Description</label>
	<input type="text" bind:value={description} {readonly} />
  </div>
</form>
<Button onClick={handleSubmitClick} value="Submit" />
