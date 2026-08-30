<script lang="ts">
	import { _ } from 'svelte-i18n';
	import Modal from '$lib/components/ui/Modal.svelte';
	import { api } from '$lib/api/client';

	interface Props {
		/** Page path recorded as the feedback source. */
		source?: string;
	}

	let { source = '' }: Props = $props();

	let open = $state(false);
	let message = $state('');
	let contact = $state('');
	let submitting = $state(false);
	let error = $state('');
	let sent = $state(false);

	async function handleSubmit() {
		if (!message.trim()) {
			error = $_('feedback.enterMessage');
			return;
		}
		submitting = true;
		error = '';
		try {
			await api.post('/feedback/public', {
				message: message.trim(),
				contact: contact.trim() || undefined,
				source: source || undefined
			});
			sent = true;
			message = '';
			contact = '';
		} catch (err) {
			const apiErr = err as { message?: string };
			error = apiErr.message || $_('feedback.sendFailed');
		} finally {
			submitting = false;
		}
	}

	function reset() {
		open = false;
		// Reset transient state after the modal closes.
		setTimeout(() => {
			sent = false;
			error = '';
			message = '';
			contact = '';
		}, 150);
	}
</script>

<button
	type="button"
	onclick={() => (open = true)}
	class="text-xs text-neutral-400 hover:text-neutral-600 underline underline-offset-2 transition-colors"
>
	{$_('feedback.reportProblem')}
</button>

<Modal bind:open title={$_('feedback.title')}>
	{#if sent}
		<div class="text-center py-2">
			<div class="w-12 h-12 rounded-full bg-success-light flex items-center justify-center mx-auto mb-3">
				<svg class="w-6 h-6 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
			</div>
			<h3 class="font-display text-lg font-semibold text-neutral-900 mb-1">{$_('feedback.thankYou')}</h3>
			<p class="text-sm text-neutral-600 mb-4">{$_('feedback.sent')}</p>
			<button
				type="button"
				onclick={reset}
				class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-hover transition-colors"
			>
				{$_('feedback.close')}
			</button>
		</div>
	{:else}
		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
			<div class="space-y-4">
				<p class="text-sm text-neutral-500">
					{$_('feedback.intro')}
				</p>
				<div>
					<label for="guest-feedback-message" class="block text-sm font-medium text-neutral-700 mb-1.5">
						{$_('feedback.messageLabel')} <span class="text-error">*</span>
					</label>
					<textarea
						id="guest-feedback-message"
						bind:value={message}
						rows="4"
						maxlength="2000"
						required
						placeholder={$_('feedback.messagePlaceholder')}
						class="w-full rounded-md border border-neutral-300 px-4 py-2.5 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary transition-colors resize-none"
					></textarea>
					<p class="mt-1 text-xs text-neutral-400">{message.length}/2000</p>
				</div>
				<div>
					<label for="guest-feedback-contact" class="block text-sm font-medium text-neutral-700 mb-1.5">
						{$_('feedback.emailLabel')} <span class="text-neutral-400 font-normal">{$_('feedback.optional')}</span>
					</label>
					<input
						id="guest-feedback-contact"
						type="email"
						bind:value={contact}
						placeholder="you@example.com"
						class="w-full rounded-md border border-neutral-300 px-4 py-2.5 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary transition-colors"
					/>
				</div>
				{#if error}
					<div class="rounded-md bg-error-light border border-error/20 px-4 py-3 text-sm text-error">
						{error}
					</div>
				{/if}
			</div>

			<div class="mt-4 flex justify-end gap-3">
				<button
					type="button"
					onclick={() => (open = false)}
					class="rounded-md border border-neutral-300 px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-neutral-50 transition-colors"
				>
					{$_('feedback.cancel')}
				</button>
				<button
					type="submit"
					disabled={submitting || !message.trim()}
					class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				>
					{submitting ? $_('feedback.sending') : $_('feedback.send')}
				</button>
			</div>
		</form>
	{/if}
</Modal>
