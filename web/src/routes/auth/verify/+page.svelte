<script lang="ts">
	import { goto, replaceState } from '$app/navigation';
	import { page } from '$app/stores';
	import { _ } from 'svelte-i18n';
	import { api } from '$lib/api/client';
	import { currentUser } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { Organizer } from '$lib/types';
	import { onMount } from 'svelte';

	let verifying = $state(true);
	let error = $state('');

	onMount(async () => {
		const token = $page.url.searchParams.get('token');

		if (!token) {
			error = $_('auth.verify.noToken');
			verifying = false;
			return;
		}

		try {
			const result = await api.post<{ token: string; organizer: Organizer }>('/auth/verify', { token });
			// Strip the raw token from the URL/history after use, so it does not
			// linger in browser history. This must run after the first await so
			// the SvelteKit client router is initialised; calling replaceState
			// synchronously in onMount throws and aborts verification (issue #5).
			replaceState('/auth/verify', {});
			$currentUser = result.organizer;
			toast.success($_('auth.verify.successToast'));
			goto('/events');
		} catch (err: unknown) {
			const apiErr = err as { message?: string };
			error = apiErr.message || $_('auth.verify.failedDefault');
			verifying = false;
		}
	});
</script>

<svelte:head>
	<title>{$_('auth.verify.pageTitle')}</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4">
	<div class="w-full max-w-md text-center">
		<a href="/" class="text-2xl font-bold text-primary">{$_('common.appName')}</a>

		{#if verifying}
			<h1 class="font-display mt-4 text-2xl font-semibold text-neutral-900">{$_('auth.verify.heading')}</h1>
			<div class="mt-6 flex flex-col items-center">
				<Spinner size="md" class="text-primary" />
				<p class="mt-4 text-neutral-600">{$_('auth.verify.pleaseWait')}</p>
			</div>
		{:else if error}
			<div class="mt-6">
				<div
					class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-error-light mb-4"
				>
					<svg class="h-6 w-6 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</div>
				<h2 class="font-display text-lg font-semibold text-neutral-900 mb-2">{$_('auth.verify.failedHeading')}</h2>
				<p class="text-sm text-neutral-600 mb-6">{error}</p>
				<Button href="/auth/login">{$_('auth.verify.tryAgain')}</Button>
			</div>
		{/if}
	</div>
</div>
