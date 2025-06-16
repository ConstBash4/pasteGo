<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { onMount } from 'svelte';
	import { getCookie } from '$lib/utils/cookies.svelte';
	import { goto } from '$app/navigation';
	import { z } from 'zod';
	import { logout } from '$lib/api/auth/auth.svelte';
	import Frame from '$lib/components/Frame.svelte';
	import { deleteUser, updateUser } from '$lib/api/user/user.svelte';

	onMount(UpdateUsername);

	function UpdateUsername() {
		const usernameCookie = getCookie('username');
		if (!usernameCookie) {
			goto('/auth/login');
		} else {
			username = usernameCookie;
		}
	}

	let error: string | null = null;
	let showForm = false;
	let isLoading = false;
	let login: string = '';
	let password: string = '';
	let username = '';
	$: filledForm = login.trim() !== '' || password.trim() !== '';

	const RegisterSchema = z
		.object({
			login: z.string().refine((value) => value === '' || value.length >= 4, {
				message: 'Имя пользователя должно содержать минимум 4 символа'
			}),
			password: z.string().refine((value) => value === '' || value.length >= 6, {
				message: 'Пароль должен содержать минимум 6 символов'
			})
		})
		.refine((data) => data.login !== '' || data.password !== '', {
			message: 'Имя пользователя или пароль должен быть заполнен'
		});

	// Функция для отображения заметок
	function showPastes() {
		goto('/paste');
	}

	// Функция для удаления пользователя
	async function handleDeleteUser() {
		error = null;
		try {
			await deleteUser();
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при удалении пользователя:' + err;
				console.log(error);
			}
		} finally {
			goto('/');
		}
	}

	// Функция для выхода
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
			goto('/');
		}
	}

	// Функция для обновления пользователя
	async function handleUpdateUser() {
		error = null;
		try {
			// Валидация данных перед отправкой
			const validatedData = RegisterSchema.parse({
				login,
				password
			});

			let response = await updateUser({
				username: validatedData.login,
				password: validatedData.password
			});

			if (response.code != 0) {
				error =
					response.code + ': ' + response.explanation ||
					'Произошла ошибка при обновлении пользователя';
			}
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при выходе:' + err;
				console.log(error);
			}
		} finally {
			UpdateUsername();
			showForm = false;
			password = login = '';
		}
	}
</script>

<svelte:head>
	<title>Profile</title>
</svelte:head>

<Header mode={''} />
<Frame {content} />
<Footer />

{#snippet content()}
	<div class="container">
		<h2>{username}</h2>
		<div class="actions">
			<button on:click={showPastes}>Список вставок</button>
			<button on:click={handleDeleteUser}
				>{isLoading ? 'Загрузка...' : 'Удалить пользователя'}</button
			>
			<button on:click={handleLogout}>{isLoading ? 'Загрузка...' : 'Выйти'}</button>
		</div>
		{#if error}
			<div class="error">{error}</div>
		{/if}
	</div>
	<button class="showForm" on:click={() => (showForm = !showForm)}> Изменить пользователя </button>
	{#if showForm}
		<form on:submit|preventDefault={handleUpdateUser} class="container">
			<div class="form-group">
				<label for="login">Логин</label>
				<input type="text" id="login" bind:value={login} placeholder="Новый логин" />
			</div>

			<div class="form-group">
				<label for="password">Пароль</label>
				<input type="password" id="password" bind:value={password} placeholder="Новый пароль" />
			</div>

			<button type="submit" disabled={!filledForm}>Изменить</button>
		</form>
	{/if}
{/snippet}

<style>
	.container {
		background-color: #1e1e1e;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
		max-width: 400px;
		width: 100%;
		margin: 10px auto;
	}

	.error {
		color: #ff4757;
		font-size: 1rem;
		margin-top: 1rem;
		text-align: center;
	}

	.container h2 {
		margin-top: 0;
		color: #00ffcc;
		text-align: center;
		font-size: 1.5rem;
		margin-bottom: 1.5rem;
	}

	.container input {
		width: 100%;
		padding: 0.5rem;
		font-size: 1rem;
		margin-top: 0.5rem;
	}

	.actions {
		margin-top: 1rem;
	}

	button {
		width: 100%;
		padding: 0.75rem;
		background-color: #00ffcc;
		margin-top: 10px;
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

	button:disabled {
		background-color: #00e0b8;
		cursor: not-allowed;
	}

	.showForm {
		margin: 10px auto;
		max-width: max-content;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group input:focus {
		border-color: #00ffcc;
		outline: none;
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
		box-sizing: border-box;
	}
</style>
