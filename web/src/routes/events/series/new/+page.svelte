<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';
	import { api } from '$lib/api/client';
	import { currentUser } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import { smsEnabled, loadAppConfig } from '$lib/stores/config';
	import { datetimeLocalToUTC } from '$lib/utils/dates';
	import { getTimezoneOptions, getTimezoneLabel } from '$lib/utils/timezones';
	import type { EventSeries } from '$lib/types';
	import AppShell from '$lib/components/layout/AppShell.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Card from '$lib/components/ui/Card.svelte';

	const defaultTz = $currentUser?.timezone
		|| Intl.DateTimeFormat().resolvedOptions().timeZone
		|| '';

	let submitting = $state(false);

	// Form fields
	let title = $state('');
	let description = $state('');
	let location = $state('');
	let timezone = $state(defaultTz);
	let startDate = $state('');
	let eventTime = $state('');
	let durationMinutes = $state('');
	let recurrenceRule = $state('weekly');
	let endCondition = $state<'none' | 'count' | 'date'>('none');
	let maxOccurrences = $state('');
	let recurrenceEnd = $state('');
	let contactRequirement = $state('email');
	let showHeadcount = $state(false);
	let showGuestList = $state(false);
	let rsvpDeadlineOffsetHours = $state('');
	let maxCapacity = $state('');
	let retentionDays = $state('30');
	let showRetention = $state(false);

	// Validation
	let errors: Record<string, string> = $state({});

	const tzOptions = getTimezoneOptions(defaultTz);

	const recurrenceOptions = $derived([
		{ value: 'weekly', label: $_('series.recWeekly') },
		{ value: 'biweekly', label: $_('series.recBiweekly') },
		{ value: 'monthly', label: $_('series.recMonthly') }
	]);

	const contactRequirementOptions = $derived([
		{ value: 'email_or_phone', label: $_('eventForm.contactEmailOrPhone') },
		{ value: 'email', label: $_('eventForm.contactEmail') },
		{ value: 'phone', label: $_('eventForm.contactPhone') },
		{ value: 'email_and_phone', label: $_('eventForm.contactBoth') }
	]);

	const filteredContactOptions = $derived(
		$smsEnabled
			? contactRequirementOptions
			: contactRequirementOptions.filter(o => o.value !== 'phone')
	);

	// Today's date for the min attribute on the date picker
	const today = new Date().toISOString().split('T')[0];

	onMount(() => {
		if (!$currentUser?.isAdmin) {
			goto('/events');
			return;
		}
		loadAppConfig();
	});

	function validate(): boolean {
		errors = {};
		if (!title.trim()) errors.title = $_('eventForm.titleRequired');
		if (!startDate) errors.startDate = $_('series.startDateRequired');
		if (!eventTime) errors.eventTime = $_('series.eventTimeRequired');
		if (!timezone) errors.timezone = $_('eventForm.timezoneRequired');

		if (endCondition === 'count') {
			const n = parseInt(maxOccurrences);
			if (isNaN(n) || n < 1) errors.maxOccurrences = $_('series.mustBe1');
		}
		if (endCondition === 'date' && !recurrenceEnd) {
			errors.recurrenceEnd = $_('series.endDateRequired');
		}
		if (durationMinutes) {
			const d = parseInt(durationMinutes);
			if (isNaN(d) || d < 1) errors.durationMinutes = $_('series.durationMin');
		}
		if (maxCapacity) {
			const parsed = Number(maxCapacity);
			if (!Number.isInteger(parsed) || parsed < 1) {
				errors.maxCapacity = $_('eventForm.maxAttendeesError');
			}
		}
		if (rsvpDeadlineOffsetHours) {
			const h = parseInt(rsvpDeadlineOffsetHours);
			if (isNaN(h) || h < 1) errors.rsvpDeadlineOffsetHours = $_('series.hourMin');
		}
		if (showRetention) {
			const days = parseInt(retentionDays);
			if (isNaN(days) || days < 1 || days > 365) {
				errors.retentionDays = $_('eventForm.retentionError');
			}
		}
		return Object.keys(errors).length === 0;
	}

	async function handleSubmit() {
		if (!validate()) return;

		submitting = true;
		try {
			const body: Record<string, unknown> = {
				title: title.trim(),
				description: description.trim(),
				location: location.trim(),
				timezone,
				startDate,
				eventTime,
				recurrenceRule,
				contactRequirement,
				showHeadcount,
				showGuestList,
				retentionDays: parseInt(retentionDays)
			};
			if (durationMinutes) body.durationMinutes = parseInt(durationMinutes);
			if (endCondition === 'count' && maxOccurrences) body.maxOccurrences = parseInt(maxOccurrences);
			if (endCondition === 'date' && recurrenceEnd) body.recurrenceEnd = datetimeLocalToUTC(recurrenceEnd + 'T23:59:59', timezone);
			if (rsvpDeadlineOffsetHours) body.rsvpDeadlineOffsetHours = parseInt(rsvpDeadlineOffsetHours);
			if (maxCapacity) body.maxCapacity = parseInt(maxCapacity);

			const result = await api.post<{ data: EventSeries }>('/events/series', body);
			toast.success($_('series.createdSuccess'));
			goto(`/events/series/${result.data.id}`);
		} catch (err: unknown) {
			const apiErr = err as { message?: string };
			toast.error(apiErr.message || $_('series.createFailed'));
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>{$_('series.newPageTitle')}</title>
</svelte:head>

<AppShell>
	<div class="max-w-3xl mx-auto">
		<div class="mb-8">
			<a href="/events/series" class="text-sm text-primary hover:text-primary-hover">{$_('series.backToSeries')}</a>
			<h1 class="mt-2 text-2xl font-bold font-display text-neutral-900">{$_('series.newHeading')}</h1>
			<p class="mt-1 text-sm text-neutral-500">{$_('series.newSubtitle')}</p>
		</div>

		<form
			onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
		>
			<Card class="mb-6">
				<div class="space-y-6">
					<h2 class="text-lg font-semibold font-display text-neutral-900">{$_('eventForm.eventDetails')}</h2>

					<Input
						label={$_('series.seriesTitle')}
						name="title"
						bind:value={title}
						placeholder={$_('series.seriesTitlePlaceholder')}
						error={errors.title || ''}
						required
					/>

					<Textarea
						label={$_('eventForm.description')}
						name="description"
						bind:value={description}
						placeholder={$_('eventForm.descriptionPlaceholder')}
						rows={4}
					/>

					<Input
						label={$_('eventForm.location')}
						name="location"
						bind:value={location}
						placeholder={$_('eventForm.locationPlaceholder')}
					/>

					<Select
						label={$_('eventForm.timezone')}
						name="timezone"
						bind:value={timezone}
						options={tzOptions}
						error={errors.timezone || ''}
						required
					/>

					<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
						<div class="space-y-1">
							<label for="startDate" class="block text-sm font-medium text-neutral-700">
								{$_('series.startDate')} <span class="text-error">*</span>
							</label>
							<input
								id="startDate"
								type="date"
								bind:value={startDate}
								min={today}
								required
								class="block w-full rounded-lg border px-3 py-2 text-sm shadow-sm transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 {errors.startDate
									? 'border-error text-error focus:border-error focus:ring-error'
									: 'border-neutral-300 text-neutral-900 focus:border-primary focus:ring-primary'}"
							/>
							{#if errors.startDate}
								<p class="text-sm text-error">{errors.startDate}</p>
							{/if}
						</div>

						<div class="space-y-1">
							<label for="eventTime" class="block text-sm font-medium text-neutral-700">
								{$_('series.eventTime')} <span class="text-error">*</span>
							</label>
							<input
								id="eventTime"
								type="time"
								bind:value={eventTime}
								required
								class="block w-full rounded-lg border px-3 py-2 text-sm shadow-sm transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 {errors.eventTime
									? 'border-error text-error focus:border-error focus:ring-error'
									: 'border-neutral-300 text-neutral-900 focus:border-primary focus:ring-primary'}"
							/>
							{#if errors.eventTime}
								<p class="text-sm text-error">{errors.eventTime}</p>
							{/if}
						</div>

						<Input
							label={$_('series.duration')}
							name="durationMinutes"
							type="number"
							bind:value={durationMinutes}
							placeholder={$_('series.durationPlaceholder')}
							error={errors.durationMinutes || ''}
						/>
					</div>
				</div>
			</Card>

			<Card class="mb-6">
				<div class="space-y-6">
					<h2 class="text-lg font-semibold font-display text-neutral-900">{$_('series.recurrence')}</h2>

					<Select
						label={$_('series.repeat')}
						name="recurrenceRule"
						bind:value={recurrenceRule}
						options={recurrenceOptions}
						required
					/>

					<fieldset>
						<legend class="text-sm font-medium text-neutral-700 mb-3">{$_('series.endCondition')}</legend>
						<div class="space-y-3">
							<label class="flex items-center gap-3 cursor-pointer">
								<input
									type="radio"
									name="endCondition"
									value="none"
									bind:group={endCondition}
									class="text-primary focus:ring-primary/40"
								/>
								<span class="text-sm text-neutral-700">{$_('series.endNone')}</span>
							</label>
							<label class="flex items-center gap-3 cursor-pointer">
								<input
									type="radio"
									name="endCondition"
									value="count"
									bind:group={endCondition}
									class="text-primary focus:ring-primary/40"
								/>
								<span class="text-sm text-neutral-700">{$_('series.endCount')}</span>
							</label>
							{#if endCondition === 'count'}
								<div class="ml-7">
									<Input
										name="maxOccurrences"
										type="number"
										bind:value={maxOccurrences}
										placeholder={$_('series.occurrencesPlaceholder')}
										error={errors.maxOccurrences || ''}
									/>
								</div>
							{/if}
							<label class="flex items-center gap-3 cursor-pointer">
								<input
									type="radio"
									name="endCondition"
									value="date"
									bind:group={endCondition}
									class="text-primary focus:ring-primary/40"
								/>
								<span class="text-sm text-neutral-700">{$_('series.endDate')}</span>
							</label>
							{#if endCondition === 'date'}
								<div class="ml-7 space-y-1">
									<input
										type="date"
										bind:value={recurrenceEnd}
										min={startDate || today}
										class="block w-full rounded-lg border px-3 py-2 text-sm shadow-sm transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 {errors.recurrenceEnd
											? 'border-error text-error focus:border-error focus:ring-error'
											: 'border-neutral-300 text-neutral-900 focus:border-primary focus:ring-primary'}"
									/>
									{#if errors.recurrenceEnd}
										<p class="text-sm text-error">{errors.recurrenceEnd}</p>
									{/if}
								</div>
							{/if}
						</div>
					</fieldset>
				</div>
			</Card>

			<Card class="mb-6">
				<div class="space-y-6">
					<h2 class="text-lg font-semibold font-display text-neutral-900">{$_('series.rsvpSettings')}</h2>

					<Select
						label={$_('eventForm.contactReq')}
						name="contactRequirement"
						bind:value={contactRequirement}
						options={filteredContactOptions}
					/>

					<fieldset class="pt-2">
						<legend class="text-sm font-medium text-neutral-700 mb-3">{$_('eventForm.guestVisibility')}</legend>
						<p class="text-xs text-neutral-400 mb-3">{$_('eventForm.guestVisibilityHelp')}</p>
						<div class="space-y-2">
							<label class="flex items-center gap-3 cursor-pointer">
								<input
									type="checkbox"
									bind:checked={showHeadcount}
									class="rounded border-neutral-300 text-primary focus:ring-primary/40"
								/>
								<span class="text-sm text-neutral-700">{$_('eventForm.showCount')}</span>
							</label>
							<label class="flex items-center gap-3 cursor-pointer">
								<input
									type="checkbox"
									bind:checked={showGuestList}
									class="rounded border-neutral-300 text-primary focus:ring-primary/40"
								/>
								<span class="text-sm text-neutral-700">{$_('eventForm.showNames')}</span>
							</label>
						</div>
					</fieldset>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<Input
							label={$_('series.deadlineOffset')}
							name="rsvpDeadlineOffsetHours"
							type="number"
							bind:value={rsvpDeadlineOffsetHours}
							placeholder={$_('series.deadlineOffsetPlaceholder')}
							helper={$_('series.deadlineOffsetHelp')}
							error={errors.rsvpDeadlineOffsetHours || ''}
						/>
						<Input
							label={$_('eventForm.maxAttendees')}
							name="maxCapacity"
							type="number"
							bind:value={maxCapacity}
							placeholder={$_('eventForm.maxAttendeesPlaceholder')}
							helper={$_('series.maxAttendeesHelpSeries')}
							error={errors.maxCapacity || ''}
						/>
					</div>

					<div class="pt-2">
						{#if showRetention}
							<Input
								label={$_('eventForm.retentionDaysLabel')}
								name="retentionDays"
								type="number"
								bind:value={retentionDays}
								helper={$_('series.retentionHelpSeries')}
								error={errors.retentionDays || ''}
							/>
						{:else}
							<p class="text-xs text-neutral-400">
								{$_('series.retentionDefaultSeries')}
								<button
									type="button"
									class="text-primary hover:text-primary-hover underline underline-offset-2"
									onclick={() => (showRetention = true)}
								>
									{$_('eventForm.specifyRetention')}
								</button>
							</p>
						{/if}
					</div>
				</div>
			</Card>

			<div class="flex items-center justify-end gap-3">
				<Button variant="outline" href="/events/series">{$_('series.cancel')}</Button>
				<Button type="submit" loading={submitting}>{$_('series.createSeries')}</Button>
			</div>
		</form>
	</div>
</AppShell>
