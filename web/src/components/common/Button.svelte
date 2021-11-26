<script lang="ts">
	export let href: string = '';
	export let value: string = '';
	export let onClick: () => void = () => {};
</script>

<button on:click={onClick} class="button">
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
		border-radius: 4px;
		width: 100%;
		padding: 8px 24px;
		background-color: var(--blue);
		color: white;
		font-weight: 300;
		cursor: pointer;
		text-align: center;
		font-size: var(--font-size-normal);
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
