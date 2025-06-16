<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Frame from '$lib/components/Frame.svelte';
	import { getCookie } from '$lib/utils/cookies.svelte';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	let mode: 'auth' | 'profile' | '' = 'auth';
	let username: string = '';

	onMount(() => {
		const usernameCookie = getCookie('username');
		if (usernameCookie) {
			mode = 'profile';
			username = usernameCookie;
		}
	});
</script>

<svelte:head>
	<title>{page.status}</title>
</svelte:head>

<Header {mode} {username} />
<Frame {content} />
<Footer />

{#snippet content()}
	<div class="container mx-auto p-4">
		<h1 class="text-3xl font-bold">{page.status}</h1>
		<p>{page.error?.message}</p>
	</div>
{/snippet}

<style>
	.container {
		max-width: 800px;
		margin: 10px auto;
		text-align: center;
	}
</style>
