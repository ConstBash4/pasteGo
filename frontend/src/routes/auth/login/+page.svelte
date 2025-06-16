<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Frame from '$lib/components/Frame.svelte';
	import { authenticate } from '$lib/api/auth/auth.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { z } from 'zod';
	import { checkAndRefreshTokens, getCookie } from '$lib/utils/cookies.svelte';

	onMount(() => {
		const usernameCookie = getCookie('username');
		if (usernameCookie) {
			goto('/profile');
		}
	});

	// Схема валидации данных
	const RegisterSchema = z.object({
		login: z.string().min(4, 'Имя пользователя должно содержать минимум 4 символа'),
		password: z.string().min(6, 'Пароль должен содержать минимум 6 символов')
	});

	let login: string = '';
	let password: string = '';
	let error: string | null = null;
	let isLoading = false;

	// Обработка отправки формы
	async function handleRegister() {
		error = null;
		isLoading = true;
		try {
			// Валидация данных перед отправкой
			const validatedData = RegisterSchema.parse({
				login,
				password
			});

			let response = await authenticate({
				username: validatedData.login,
				password: validatedData.password
			});

			if (response.code == 0) {
				goto('/profile');
				checkAndRefreshTokens();
			} else {
				error = response.code + ': ' + response.explanation || 'Произошла ошибка при регистрации';
			}
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при регистрации:' + err;
			}
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<Header />
<Frame {content} />
<Footer />

{#snippet content()}
	<form on:submit|preventDefault={handleRegister} class="login-form">
		<h2>Вход</h2>
		<div class="form-group">
			<label for="login">Логин</label>
			<input type="text" id="login" bind:value={login} placeholder="Введите логин" required />
		</div>
		<div class="form-group">
			<label for="password">Пароль</label>
			<input
				type="password"
				id="password"
				bind:value={password}
				placeholder="Введите пароль"
				required
			/>
		</div>
		<button type="submit">{isLoading ? 'Загрузка...' : 'Войти'}</button>
		{#if error}
			<div class="error">{error}</div>
		{/if}
	</form>
{/snippet}

<style>
	.login-form {
		background-color: #1e1e1e;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
		max-width: 400px;
		width: 100%;
		margin: 10px auto;
	}

	.login-form h2 {
		margin-top: 0;
		color: #00ffcc;
		text-align: center;
		font-size: 1.5rem;
		margin-bottom: 1.5rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		color: #ccc;
	}

	.form-group input {
		width: 100%;
		padding: 0.75rem;
		background-color: #2a2a2a;
		border: 1px solid #444;
		color: #fff;
		border-radius: 5px;
		font-size: 1rem;
		transition: border-color 0.3s;
	}

	.form-group input:focus {
		border-color: #00ffcc;
		outline: none;
	}

	.error {
		color: #ff4757;
		font-size: 1rem;
		margin-top: 1rem;
		text-align: center;
	}

	button {
		width: 100%;
		padding: 0.75rem;
		background-color: #00ffcc;
		color: #000;
		border: none;
		border-radius: 5px;
		font-size: 1rem;
		cursor: pointer;
		transition: background-color 0.3s;
	}

	button:hover {
		background-color: #00e0b8;
	}

	input {
		box-sizing: border-box;
	}
</style>
