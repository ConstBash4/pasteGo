<script lang="ts" context="module">
	import { z } from 'zod';

	export const APIResponseSchema = z.object({
		code: z.number(),
		explanation: z.string(),
		message: z.record(z.unknown()).optional()
	});
	export type APIResponse = z.infer<typeof APIResponseSchema>;

	export const CredentialsSchema = z.object({
		username: z.string(),
		password: z.string()
	});
	export type Credentials = z.infer<typeof CredentialsSchema>;

	export const UserDataSchema = z.object({
		username: z.string()
	});
	export type UserData = z.infer<typeof UserDataSchema>;

	export const PasteInfoSchema = z.object({
		id: z.string(),
		author: z.string(),
		created: z.number(),
		updated: z.number().optional(),
		expTime: z.number(),
		lifetime: z.string().optional(),
		text: z.string(),
		password: z.string().optional(),
		hasPassword: z.boolean(),
		public: z.boolean()
	});
	export type PasteInfo = z.infer<typeof PasteInfoSchema>;

	export const PasteInfoListSchema = z.object({
		pastes: z.array(PasteInfoSchema)
	});
	export type PasteInfoList = z.infer<typeof PasteInfoListSchema>;

	export const PastePasswordSchema = z.object({
		password: z.string()
	});
	export type PastePassword = z.infer<typeof PastePasswordSchema>;
</script>
