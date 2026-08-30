import { register, init, locale, waitLocale, isLoading } from 'svelte-i18n';

export const SUPPORTED_LOCALES = ['en', 'fr'] as const;
export type SupportedLocale = (typeof SUPPORTED_LOCALES)[number];

const FALLBACK_LOCALE: SupportedLocale = 'en';

register('en', () => import('./locales/en.json'));
register('fr', () => import('./locales/fr.json'));

// Initialise with the fallback; the real instance locale is applied at runtime
// from /api/v1/config via setLocale() once the app config has loaded.
init({
	fallbackLocale: FALLBACK_LOCALE,
	initialLocale: FALLBACK_LOCALE
});

/** Apply the instance locale, falling back to English for unknown values. */
export async function setLocale(value: string | null | undefined): Promise<void> {
	const next = (SUPPORTED_LOCALES as readonly string[]).includes(value ?? '')
		? (value as SupportedLocale)
		: FALLBACK_LOCALE;
	locale.set(next);
	await waitLocale(next);
}

export { locale, waitLocale, isLoading };
