<script lang="ts">
	export let href: string = '';
	export let value: string = '';
	export let onClick: () => void = () => {};
</script>

<button on:click={onClick}>
	{#if href !== ''}
		<a {href}>
			{value}
		</a>
	{:else}
		{value}
	{/if}
</button>
