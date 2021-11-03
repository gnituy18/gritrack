<script lang="ts">
	export let href: string = '';
	export let value: string;
</script>

<button class="button">
	{#if href !== undefined}
		<a {href}>
			{value}
		</a>
	{:else}
		{value}
	{/if}
</button>

<style lang="scss">
	.button {
		display: inline-block;
		border-radius: 8px;
		border: 0;
		padding: 8px 24px;
		background-color: transparent;
		color: black;
		cursor: pointer;
		text-align: center;

		&:hover {
			background-color: rgba(0, 0, 0, 0.05);
			text-decoration: none;
		}
	}

	a {
		&:visited,
		&:hover {
			color: black;
			text-decoration: none;
		}
	}
</style>
