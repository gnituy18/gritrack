<script>
	import { session } from '$app/stores';

	import Avatar from '$components/common/Avatar.svelte';
	import Button from '$components/common/Button.svelte';

	let user;
	session.subscribe(({ currentUser }) => {
		user = currentUser;
	});
</script>

<header>
	<nav>
		<div class="nav-content">
			<Avatar src={user.picture} alt={user.name} />
			<h1>{user.name}</h1>
		</div>
		<hr />
		<div class="nav-content">
			<Button href="/mission/create" value="Create a mission" />
		</div>
		<div class="nav-content">
			<h2>Your Missions</h2>
		</div>
	</nav>
</header>

<style lang="scss">
	.nav-content {
		display: flex;
		align-items: center;
		padding: 12px;
	}

	header {
		flex-shrink: 0;
		height: 100%;
		width: 240px;
		border-right: 1px solid WhiteSmoke;
	}

	li {
		margin: 16px 0;
		list-style-type: none;
	}
</style>
