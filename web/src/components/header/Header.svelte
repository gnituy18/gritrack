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
	<Avatar src={user.picture} alt={user.name} />
	<nav>
		<ul>
			<li><Button href="/" value="home" /></li>
		</ul>
	</nav>
</header>

<style>
	header {
		height: 100%;
		width: 240px;
	}

	li {
		list-style-type: none;
	}
</style>
