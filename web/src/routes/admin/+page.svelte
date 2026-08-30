<script lang="ts">
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';
	import { api } from '$lib/api/client';
	import type { InstanceStats, ApiResponse } from '$lib/types';
	import AppShell from '$lib/components/layout/AppShell.svelte';
	import Spinner from '$lib/components/ui/Spinner.svelte';

	let stats = $state<InstanceStats | null>(null);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const res = await api.get<ApiResponse<InstanceStats>>('/admin/stats');
			stats = res.data;
		} catch (e: any) {
			error = e?.message || e?.error || $_('admin.loadFailed');
		} finally {
			loading = false;
		}
	});

	function pct(value: number, total: number): string {
		if (total === 0) return '0';
		return Math.round((value / total) * 100).toString();
	}

	function barWidth(value: number, total: number): string {
		if (total === 0) return '0%';
		return `${Math.max(2, Math.round((value / total) * 100))}%`;
	}

	const eventStatusItems = $derived(stats ? [
		{ label: $_('admin.evPublished'), value: stats.events.published, color: 'bg-success' },
		{ label: $_('admin.evDraft'), value: stats.events.draft, color: 'bg-neutral-400' },
		{ label: $_('admin.evCancelled'), value: stats.events.cancelled, color: 'bg-error' },
		{ label: $_('admin.evArchived'), value: stats.events.archived, color: 'bg-warning' },
	] : []);

	const rsvpItems = $derived(stats ? [
		{ label: $_('admin.rsvpAttending'), value: stats.attendees.attending, color: 'bg-success' },
		{ label: $_('admin.rsvpMaybe'), value: stats.attendees.maybe, color: 'bg-warning' },
		{ label: $_('admin.rsvpDeclined'), value: stats.attendees.declined, color: 'bg-error' },
		{ label: $_('admin.rsvpPending'), value: stats.attendees.pending, color: 'bg-info' },
		{ label: $_('admin.rsvpWaitlisted'), value: stats.attendees.waitlisted, color: 'bg-secondary' },
	] : []);

	const notifItems = $derived(stats ? [
		{ label: $_('admin.notifTotal'), value: stats.notifications.total, color: 'text-neutral-900' },
		{ label: $_('admin.notifSent'), value: stats.notifications.sent, color: 'text-info' },
		{ label: $_('admin.notifDelivered'), value: stats.notifications.delivered, color: 'text-success' },
		{ label: $_('admin.notifOpened'), value: stats.notifications.opened, color: 'text-primary' },
		{ label: $_('admin.notifBounced'), value: stats.notifications.bounced, color: 'text-error' },
		{ label: $_('admin.notifComplained'), value: stats.notifications.complained, color: 'text-warning' },
		{ label: $_('admin.notifFailed'), value: stats.notifications.failed, color: 'text-error' },
	] : []);

	const featureItems = $derived(stats ? [
		{ label: $_('admin.featWaitlist'), value: stats.features.waitlistEvents },
		{ label: $_('admin.featComments'), value: stats.features.commentsEnabledEvents },
		{ label: $_('admin.featCohosted'), value: stats.features.cohostedEvents },
		{ label: $_('admin.featQuestions'), value: stats.features.eventsWithQuestions },
		{ label: $_('admin.featCapacity'), value: stats.features.eventsWithCapacity },
		{ label: $_('admin.featSeries'), value: stats.features.seriesEvents },
	] : []);
</script>

<AppShell>
	<div class="space-y-8">
		<div>
			<h1 class="text-2xl font-bold font-display text-neutral-900">{$_('admin.title')}</h1>
			<p class="text-sm text-neutral-500 mt-1">{$_('admin.subtitle')}</p>
		</div>

		{#if loading}
			<div class="flex items-center justify-center py-20">
				<Spinner />
			</div>
		{:else if error}
			<div class="bg-error-light border border-error rounded-lg p-4 text-error">
				{error}
			</div>
		{:else if stats}
			<!-- Metric Cards -->
			<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
				<div class="bg-surface rounded-lg border border-neutral-200 p-5">
					<p class="text-sm font-medium text-neutral-500">{$_('admin.totalEvents')}</p>
					<p class="text-3xl font-bold font-mono text-neutral-900 mt-1">{stats.events.total}</p>
				</div>
				<div class="bg-surface rounded-lg border border-neutral-200 p-5">
					<p class="text-sm font-medium text-neutral-500">{$_('admin.totalGuests')}</p>
					<p class="text-3xl font-bold font-mono text-neutral-900 mt-1">{stats.attendees.totalHeadcount}</p>
					<p class="text-xs text-neutral-400 mt-1">{$_('admin.rsvpsPlusOnes', { values: { count: stats.attendees.total } })}</p>
				</div>
				<div class="bg-surface rounded-lg border border-neutral-200 p-5">
					<p class="text-sm font-medium text-neutral-500">{$_('admin.organizers')}</p>
					<p class="text-3xl font-bold font-mono text-neutral-900 mt-1">{stats.organizers.total}</p>
				</div>
				<div class="bg-surface rounded-lg border border-neutral-200 p-5">
					<p class="text-sm font-medium text-neutral-500">{$_('admin.avgGuests')}</p>
					<p class="text-3xl font-bold font-mono text-neutral-900 mt-1">{stats.attendees.avgPerEvent.toFixed(1)}</p>
				</div>
			</div>

			<!-- Charts Row -->
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<!-- Events by Status -->
				<div class="bg-surface rounded-lg border border-neutral-200 p-6">
					<h2 class="text-sm font-semibold font-display text-neutral-700 mb-4">{$_('admin.eventsByStatus')}</h2>
					<div class="space-y-3">
						{#each eventStatusItems as item}
							<div class="flex items-center gap-3">
								<span class="text-sm text-neutral-600 w-24 shrink-0">{item.label}</span>
								<div class="flex-1 bg-neutral-100 rounded-full h-5 overflow-hidden">
									<div
										class="{item.color} h-full rounded-full transition-all"
										style="width: {barWidth(item.value, stats.events.total)}"
									></div>
								</div>
								<span class="text-sm font-medium text-neutral-700 w-16 text-right">
									{item.value}
									<span class="text-neutral-400 text-xs">({pct(item.value, stats.events.total)}%)</span>
								</span>
							</div>
						{/each}
					</div>
				</div>

				<!-- RSVP Distribution -->
				<div class="bg-surface rounded-lg border border-neutral-200 p-6">
					<h2 class="text-sm font-semibold font-display text-neutral-700 mb-4">{$_('admin.rsvpDistribution')}</h2>
					<div class="space-y-3">
						{#each rsvpItems as item}
							<div class="flex items-center gap-3">
								<span class="text-sm text-neutral-600 w-24 shrink-0">{item.label}</span>
								<div class="flex-1 bg-neutral-100 rounded-full h-5 overflow-hidden">
									<div
										class="{item.color} h-full rounded-full transition-all"
										style="width: {barWidth(item.value, stats.attendees.total)}"
									></div>
								</div>
								<span class="text-sm font-medium text-neutral-700 w-16 text-right">
									{item.value}
									<span class="text-neutral-400 text-xs">({pct(item.value, stats.attendees.total)}%)</span>
								</span>
							</div>
						{/each}
					</div>
				</div>
			</div>

			<!-- Notification Health -->
			<div class="bg-surface rounded-lg border border-neutral-200 p-6">
				<h2 class="text-sm font-semibold font-display text-neutral-700 mb-4">{$_('admin.notifHealth')}</h2>
				{#if stats.notifications.total === 0}
					<p class="text-sm text-neutral-400">{$_('admin.noNotifs')}</p>
				{:else}
					<div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-7 gap-4">
						{#each notifItems as item}
							<div class="text-center">
								<p class="text-2xl font-bold {item.color}">{item.value}</p>
								<p class="text-xs text-neutral-500 mt-1">{item.label}</p>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Feature Adoption -->
			<div class="bg-surface rounded-lg border border-neutral-200 p-6">
				<h2 class="text-sm font-semibold font-display text-neutral-700 mb-4">{$_('admin.featureAdoption')}</h2>
				<div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
					{#each featureItems as item}
						<div class="text-center p-3 bg-neutral-50 rounded-lg">
							<p class="text-2xl font-bold text-neutral-900">{item.value}</p>
							<p class="text-xs text-neutral-500 mt-1">{item.label}</p>
							{#if stats.events.total > 0}
								<p class="text-xs text-primary mt-0.5">{$_('admin.pctOfEvents', { values: { pct: pct(item.value, stats.events.total) } })}</p>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
</AppShell>
