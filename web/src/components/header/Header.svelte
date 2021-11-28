<script>
  import { session } from "$app/stores";

  import Avatar from "$components/common/Avatar.svelte";
  import Button from "$components/common/Button.svelte";

  let user;
  session.subscribe(({ currentUser }) => {
    user = currentUser;
  });
</script>

<header class="flex-none box-border w-60 h-full border-r">
  <div class="flex items-center p-2">
    <Avatar src={user.picture} alt={user.name} />
    <h1 class="ml-4">{user.name}</h1>
  </div>
  <hr />
  <nav>
    <div class="flex justify-between p-4">
      <h2>Missions</h2>
      <a href="/mission/create">+</a>
    </div>
  </nav>
</header>
