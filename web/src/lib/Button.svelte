<script lang="ts">
	export let href: string = '';
	export let value: string;
</script>

<button class="button">
	{#if href !== ''}
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
		border: 0;
		border-radius: 8px;
		padding: 8px 24px;
		background-color: transparent;
		color: var(--blue);
		cursor: pointer;
		text-align: center;
		transition: background-color 0.1s;
		user-select: none;

		&:hover {
			background-color: var(--blue-light);
		}
	}

	a {
		color: inherit;
	}
</style>
