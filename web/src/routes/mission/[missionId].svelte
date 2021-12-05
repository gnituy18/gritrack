<script>
  import { page } from "$app/stores";

  let missionId = "";
  page.subscribe(({ params }) => {
    missionId = params.missionId;
  });
</script>

<h1>
  {missionId}
</h1>
