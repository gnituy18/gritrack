<script>
  import { session } from "$app/stores";

  import Avatar from "$components/common/Avatar.svelte";

  let user;
  session.subscribe(({ currentUser }) => {
    user = currentUser;
  });
</script>

<header class="flex-none box-border w-60 h-full border-r">
  <div class="flex justify- items-center m-4">
    <Avatar src={user.picture} alt={user.name} />
    <h1 class="ml-2">{user.name}</h1>
  </div>
  <hr />
  <nav>
    <div class="flex justify-between items-center m-4">
      <h2>Missions</h2>
      <a href="/mission/create" class="font-medium">+</a>
    </div>
  </nav>
</header>
