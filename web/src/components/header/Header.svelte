<script>
	import { session } from '$app/stores';

	import Avatar from '$components/common/Avatar.svelte';

	let user;
	session.subscribe(({ currentUser }) => {
		user = currentUser;
	});
</script>

<header>
	<div class="profile">
		<Avatar src={user.picture} alt={user.name} />
		<span>{user.name}</span>
	</div>
	<hr />
	<nav>
		<ul>
			<li><a href="/mission">Missions</a></li>
		</ul>
	</nav>
</header>

<style lang="scss">
	.profile {
		display: flex;
		align-items: center;
		padding: 12px;

		&:hover {
			background-color: WhiteSmoke;
			cursor: pointer;
		}
	}

	span {
		color: var(--color-text);
		font-size: var(--font-size-large);
		margin: 0 8px;
	}

	header {
		height: 100%;
		width: 240px;
		border-right: 1px solid WhiteSmoke;
	}

	li {
		margin: 16px 0;
		list-style-type: none;
	}
</style>
