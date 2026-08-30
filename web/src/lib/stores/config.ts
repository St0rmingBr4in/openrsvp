import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import { setLocale } from '$lib/i18n';

export const smsEnabled = writable(false);

let loaded = false;

export async function loadAppConfig() {
	if (loaded) return;
	try {
		const result = await api.get<{ data: { smsEnabled: boolean; defaultLocale?: string } }>(
			'/config'
		);
		smsEnabled.set(result.data.smsEnabled);
		await setLocale(result.data.defaultLocale);
		loaded = true;
	} catch {
		smsEnabled.set(false);
		// Fall back to the default locale so the UI still renders.
		await setLocale(null);
	}
}
