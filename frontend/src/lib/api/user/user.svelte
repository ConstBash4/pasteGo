<script lang="ts" context="module">
	import { apiClient } from '../api.svelte';
	import {
		APIResponseSchema,
		CredentialsSchema,
		type APIResponse,
		type Credentials
	} from '../types.svelte';

	export async function getUser(): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/user',
			method: 'GET',
			responseSchema: APIResponseSchema
		});
	}

	export async function updateUser(data: Credentials): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/user',
			method: 'PUT',
			requestData: data,
			requestSchema: CredentialsSchema,
			responseSchema: APIResponseSchema
		});
	}

	export async function deleteUser(): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/user',
			method: 'DELETE',
			responseSchema: APIResponseSchema
		});
	}
</script>
