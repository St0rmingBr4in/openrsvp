<script lang="ts">
	import '../app.css';
	import '$lib/i18n';
	import { currentUser, isLoading } from '$lib/stores/auth';
	import { api } from '$lib/api/client';
	import { onMount } from 'svelte';
	import Toast from '$lib/components/ui/Toast.svelte';
	import { loadAppConfig } from '$lib/stores/config';

	// Gate rendering until the instance locale dictionary is loaded, so no
	// translation keys flash before the language is applied.
	let localeReady = $state(false);

	onMount(async () => {
		await loadAppConfig();
		localeReady = true;
		try {
			const user = await api.get<import('$lib/types').Organizer>('/auth/me');
			$currentUser = user;

			// Auto-save browser timezone to profile if not set yet.
			if (!user.timezone) {
				const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
				if (tz) {
					api.patch<import('$lib/types').Organizer>('/auth/me', { timezone: tz })
						.then((updated) => { $currentUser = updated; })
						.catch(() => {});
				}
			}
		} catch {
			$currentUser = null;
		} finally {
			$isLoading = false;
		}
	});

	let { children } = $props();
</script>

<div class="min-h-screen bg-neutral-50">
	{#if localeReady}
		{@render children()}
	{/if}
	<Toast />
</div>
