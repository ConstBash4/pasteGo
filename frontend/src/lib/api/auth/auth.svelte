<script lang="ts" context="module">
	import { apiClient } from '../api.svelte';
	import {
		APIResponseSchema,
		CredentialsSchema,
		type APIResponse,
		type Credentials
	} from '../types.svelte';

	export async function authenticate(data: Credentials): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/auth',
			method: 'POST',
			requestData: data,
			requestSchema: CredentialsSchema,
			responseSchema: APIResponseSchema
		});
	}

	export async function registration(data: Credentials): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/registration',
			method: 'POST',
			requestData: data,
			requestSchema: CredentialsSchema,
			responseSchema: APIResponseSchema
		});
	}

	export async function refresh(): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/update_tokens',
			method: 'POST',
			responseSchema: APIResponseSchema
		});
	}

	export async function logout(): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/logout',
			method: 'DELETE',
			responseSchema: APIResponseSchema
		});
	}
</script>
