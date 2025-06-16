<script lang="ts">
	import { goto } from '$app/navigation';
	import { z } from 'zod';
	import { logout } from '$lib/api/auth/auth.svelte';

	export let mode: 'auth' | 'profile' | '' = '';
	export let username: string = '';
	let isMenuOpen = false;
	let error: string | null = null;

	async function handleLogout() {
		error = null;
		try {
			await logout();
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при выходе:' + err;
				console.log(error);
			}
		} finally {
			mode = 'auth';
			username = '';
			goto('/');
		}
	}
</script>

<header class="header">
	<nav class="nav">
		<a class="logo" href="/">PasteGo</a>
		<ul class="nav-links"></ul>
		{#if mode === 'auth'}
			{@render authBlock()}
		{:else if mode === 'profile'}
			{@render profileBlock()}
		{/if}
	</nav>
</header>

{#snippet authBlock()}
	<div class="auth-button" on:click={() => (isMenuOpen = !isMenuOpen)}>
		Вход/Регистрация
		<div class="dropdown-arrow" style:transform={isMenuOpen ? 'rotate(180deg)' : 'rotate(0deg)'}>
			▼
		</div>
	</div>
	<div class="dropdown-menu" style:display={isMenuOpen ? 'block' : 'none'}>
		<a href="/auth/login">Войти</a>
		<a href="/auth/registration">Зарегистрироваться</a>
	</div>
{/snippet}

{#snippet profileBlock()}
	<div class="auth-button" on:click={() => (isMenuOpen = !isMenuOpen)}>
		{username}
		<div class="dropdown-arrow" style:transform={isMenuOpen ? 'rotate(180deg)' : 'rotate(0deg)'}>
			▼
		</div>
	</div>
	<div class="dropdown-menu" style:display={isMenuOpen ? 'block' : 'none'}>
		<a href="/profile">Профиль</a>
		<button on:click={handleLogout}>Выйти</button>
	</div>
{/snippet}

<style>
	.header {
		background-color: #1e1e1e;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
		position: fixed;
		z-index: 1000;
		top: 0;
		left: 0;
		right: 0;
		padding: 1rem 2rem;
		height: 32px;
	}

	.nav {
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
		position: relative;
	}

	.logo {
		font-size: 1.5rem;
		font-weight: bold;
		color: #00ffcc;
		text-decoration: none;
	}

	.nav-links {
		list-style: none;
		display: flex;
		gap: 1.5rem;
		margin: 0;
		padding: 0;
	}

	.auth-button {
		position: relative;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #ccc;
		font-weight: bold;
		user-select: none;
	}

	.dropdown-arrow {
		transition: transform 0.3s;
		font-size: 0.8rem;
		user-select: none;
	}

	.dropdown-menu {
		position: absolute;
		top: 100%;
		right: 0;
		background-color: #1e1e1e;
		border: 1px solid #333;
		border-radius: 4px;
		padding: 0.5rem 0;
		z-index: 1001;
		min-width: 120px;
		white-space: nowrap;
		user-select: none;
	}

	.dropdown-menu a,
	.dropdown-menu button {
		display: block;
		padding: 0.5rem 1rem;
		color: #ccc;
		text-decoration: none;
		transition: background-color 0.3s;
		font-size: 1rem;
	}

	.dropdown-menu button {
		color: #a01010;
		background: none;
		border: none;
		width: 100%;
		text-align: left;
		cursor: pointer;
	}

	.dropdown-menu a:hover {
		background-color: #2a2a2a;
	}

	.dropdown-menu button:hover {
		background-color: #2b0c0c;
	}
</style>
