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
	<div class="profile">
		<Avatar src={user.picture} alt={user.name} />
		<span>{user.name}</span>
	</div>
	<nav>
		<ul>
			<li><Button href="/" value="home" /></li>
		</ul>
	</nav>
</header>

<style lang="scss">
	.profile {
		display: flex;
		align-items: center;
		padding: 16px 16px;
		border-bottom: 1px solid WhiteSmoke;

		&:hover {
			background-color: WhiteSmoke;
			cursor: pointer;
		}
	}

	span {
		color: SlateGray;
		font-size: 24px;
		margin: 0 8px;
	}


	header {
		height: 100%;
		width: 240px;
		border-right: 1px solid WhiteSmoke;
	}

	li {
		list-style-type: none;
	}
</style>
