<script lang="ts">
	import { goto } from '$app/navigation';
	import { currentUser, isLoading } from '$lib/stores/auth';

	$effect(() => {
		if (!$isLoading) {
			goto($currentUser ? '/events' : '/auth/login', { replaceState: true });
		}
	});
</script>

<div class="flex items-center justify-center min-h-screen">
	<div class="loading-spinner"></div>
</div>
