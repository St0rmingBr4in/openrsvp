<script lang="ts">
	import { goto } from '$app/navigation';
	import { _ } from 'svelte-i18n';
	import { api } from '$lib/api/client';
	import { currentUser } from '$lib/stores/auth';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			await api.post('/auth/logout');
		} catch {
			// Ignore errors on logout
		} finally {
			$currentUser = null;
			goto('/');
		}
	});
</script>

<svelte:head>
	<title>{$_('logout.pageTitle')}</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4">
	<div class="text-center">
		<Spinner size="md" class="text-primary mx-auto" />
		<p class="mt-4 text-neutral-600">{$_('logout.signingOut')}</p>
	</div>
</div>
