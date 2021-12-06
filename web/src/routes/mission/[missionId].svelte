<script lang="ts" context="module">
  /** @type {import('@sveltejs/kit').Load} */
  export async function load({ page, fetch }) {
    const missionId = page.params.missionId;
    let res = await fetch(`http://localhost:8080/api/v1/mission/${missionId}`, {
      credentials: "include",
    });
    const mission = await res.json();

    res = await fetch(`http://localhost:8080/api/v1/mission/${missionId}/step?offset=0&limit=5`, {
      credentials: "include",
    });
    const steps = await res.json();

    console.log(mission);
    console.log(steps);

    return {
      props: {
        mission,
        steps,
      },
    };
  }
</script>

<script lang="ts">
  import type { Step } from "./type";
  export let steps: Array<Step>;
  console.log(steps)
</script>

<ul>
  {#each steps as { id, summary, items }}
    <li>
      {id}
      {summary}
      {#each items as item}
        <li>{item.type}</li>
      {/each}
    </li>
  {/each}
</ul>
