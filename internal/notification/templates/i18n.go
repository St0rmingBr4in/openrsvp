package templates

import "sync/atomic"

// Msg holds every translatable string used by the email templates, plain-text
// fallbacks, and subjects. One instance per supported locale.
//
// The instance locale is global and immutable after startup (set once from
// config.DefaultLocale), so a package-level default avoids threading a locale
// argument through every Render call and email call site.
type Msg struct {
	// Shared
	FooterTagline  string // "Simple event RSVPs"
	ButtonFallback string // "If the button does not work, copy and paste this link into your browser:"
	LabelEvent     string
	LabelDate      string
	LabelLocation  string

	// RSVP status labels
	StatusAttending  string
	StatusMaybe      string
	StatusDeclined   string
	StatusPending    string
	StatusWaitlisted string

	// Magic link
	MagicHeading    string
	MagicIntro      string
	MagicButton     string
	MagicExpiryPre  string // "This link expires in "
	MagicExpiryPost string // " minutes. If you did not request this link, you can safely ignore this email."
	MagicPlainTitle string
	MagicPlainBody  string
	MagicPlainExpiry string // "This link expires in %d minutes." (one %d)
	MagicPlainIgnore string

	// RSVP confirmation
	RSVPConfHeading string
	RSVPConfIntro   string
	LabelYourRSVP   string
	ModifyButton    string
	RSVPConfPlainTitle string

	// Event reminder
	ReminderHeading      string
	ReminderIntro        string
	OrganizerMessage     string // "Message from the organizer"
	ViewInviteButton     string
	ReminderPlainTitle   string
	ReminderPlainOrgMsg  string // "Message from the organizer:"
	ReminderPlainInvite  string // "View your invitation:"

	// Retention warning
	RetentionHeading  string // text after the ⚠ glyph
	RetentionIntro    string
	LabelDeletionDate string
	RetentionBody     string
	ViewEventButton   string
	RetentionFooter1  string
	RetentionFooter2  string
	RetentionPlainTitle string
	RetentionPlainBody  string // "Your event \"%s\" is scheduled for automatic deletion on %s." (%s, %s)
	RetentionPlainNote  string
	RetentionPlainExtend string // "To extend the retention period, visit:"

	// Organizer RSVP notification
	OrgNotifyHeading   string
	OrgNotifyIntroPre  string // "Someone has responded to your event "
	OrgNotifyIntroPost string // "."
	LabelGuest         string
	LabelResponse      string
	LabelEmail         string
	LabelPhone         string
	LabelAdditional    string // "Additional Guests"
	ViewDashboardButton string
	OrgNotifyPlainTitle string
	OrgNotifyPlainView  string // "View your event dashboard:"

	// Feedback confirmation
	FeedbackHeading    string
	FeedbackIntroPre   string // "We received your "
	FeedbackIntroPost  string // " submission and appreciate you taking the time to share it with us."
	FeedbackFollowUp   string
	FeedbackOutro      string
	FeedbackPlainTitle string
	FeedbackPlainBody  string // "We received your %s submission and appreciate you taking the time to share it with us." (%s)
	FeedbackPlainOutro string

	// RSVP lookup
	LookupHeading    string
	LookupIntroPre   string // "Click the button below to view and manage your RSVP for "
	LookupIntroPost  string // "."
	ViewMyRSVPButton string
	LookupPersonal   string
	LookupPlainTitle string
	LookupPlainBody  string // "Click the link below to view and manage your RSVP for %s:" (%s)
	LookupPlainNote  string

	// Waitlist promotion
	WaitlistHeading    string
	WaitlistIntroPre   string // "Great news! A spot opened up for "
	WaitlistIntroPost  string // ". You are now attending."
	ViewYourRSVPButton string
	WaitlistPlainTitle string
	WaitlistPlainBody  string // "Great news! A spot opened up for %s. You are now attending." (%s)
	WaitlistPlainView  string // "View your RSVP:"

	// Co-host invitation
	CohostHeading   string
	CohostIntroMid  string // " has added you as a co-host for "
	CohostIntroPost string // ". You can now manage RSVPs, send messages, and help run the event."
	CohostPlainTitle string
	CohostPlainBody  string // "%s has added you as a co-host for %s." (%s, %s)
	CohostPlainView  string // "View the event dashboard:"

	// Subjects
	SubjMagic           string
	SubjRSVPConfirmation string
	SubjNewRSVP         string
	SubjInvited         string
	SubjCohost          string
	SubjWaitlist        string
	SubjNewMessage      string // "New message from"
	SubjCancelled       string
	SubjRetention       string
	SubjReminder        string
	SubjFeedbackConfirm string // full subject line
}

var enMsg = Msg{
	FooterTagline:  "Simple event RSVPs",
	ButtonFallback: "If the button does not work, copy and paste this link into your browser:",
	LabelEvent:     "Event",
	LabelDate:      "Date",
	LabelLocation:  "Location",

	StatusAttending:  "Attending",
	StatusMaybe:      "Maybe",
	StatusDeclined:   "Can't make it",
	StatusPending:    "Pending",
	StatusWaitlisted: "Waitlisted",

	MagicHeading:     "Sign in to your account",
	MagicIntro:       "Click the button below to sign in to OpenRSVP. No password needed.",
	MagicButton:      "Sign In",
	MagicExpiryPre:   "This link expires in ",
	MagicExpiryPost:  " minutes. If you did not request this link, you can safely ignore this email.",
	MagicPlainTitle:  "Sign in to OpenRSVP",
	MagicPlainBody:   "Click the link below to sign in:",
	MagicPlainExpiry: "This link expires in %d minutes.",
	MagicPlainIgnore: "If you did not request this link, you can safely ignore this email.",

	RSVPConfHeading:    "RSVP Confirmed",
	RSVPConfIntro:      "Your RSVP has been recorded. Here are the details:",
	LabelYourRSVP:      "Your RSVP",
	ModifyButton:       "Modify RSVP",
	RSVPConfPlainTitle: "RSVP Confirmed",

	ReminderHeading:     "Event Reminder",
	ReminderIntro:       "This is a friendly reminder about an upcoming event.",
	OrganizerMessage:    "Message from the organizer",
	ViewInviteButton:    "View Invitation",
	ReminderPlainTitle:  "Event Reminder",
	ReminderPlainOrgMsg: "Message from the organizer:",
	ReminderPlainInvite: "View your invitation:",

	RetentionHeading:     "Data Retention Notice",
	RetentionIntro:       "Your event data is scheduled for automatic deletion soon. This is a courtesy reminder so you can take action if needed.",
	LabelDeletionDate:    "Deletion Date",
	RetentionBody:        "After this date, all event data including attendee RSVPs, messages, and invite cards will be permanently deleted. If you would like to keep this data, please log in and extend the retention period.",
	ViewEventButton:      "View Event",
	RetentionFooter1:     "This is an automated notice from OpenRSVP.",
	RetentionFooter2:     "You received this because you are the organizer of this event.",
	RetentionPlainTitle:  "Data Retention Notice",
	RetentionPlainBody:   "Your event \"%s\" is scheduled for automatic deletion on %s.",
	RetentionPlainNote:   "After this date, all event data including attendee RSVPs, messages, and invite cards will be permanently deleted.",
	RetentionPlainExtend: "To extend the retention period, visit:",

	OrgNotifyHeading:    "New RSVP Received",
	OrgNotifyIntroPre:   "Someone has responded to your event ",
	OrgNotifyIntroPost:  ".",
	LabelGuest:          "Guest",
	LabelResponse:       "Response",
	LabelEmail:          "Email",
	LabelPhone:          "Phone",
	LabelAdditional:     "Additional Guests",
	ViewDashboardButton: "View Event Dashboard",
	OrgNotifyPlainTitle: "New RSVP Received",
	OrgNotifyPlainView:  "View your event dashboard:",

	FeedbackHeading:    "Thanks for your feedback!",
	FeedbackIntroPre:   "We received your ",
	FeedbackIntroPost:  " submission and appreciate you taking the time to share it with us.",
	FeedbackFollowUp:   "Since you opted in to follow-up contact, we may reach out to you at this email address if we have questions or updates related to your feedback.",
	FeedbackOutro:      "Your feedback helps make OpenRSVP better for everyone.",
	FeedbackPlainTitle: "Thanks for your feedback!",
	FeedbackPlainBody:  "We received your %s submission and appreciate you taking the time to share it with us.",
	FeedbackPlainOutro: "Your feedback helps make OpenRSVP better for everyone.",

	LookupHeading:    "Find Your RSVP",
	LookupIntroPre:   "Click the button below to view and manage your RSVP for ",
	LookupIntroPost:  ".",
	ViewMyRSVPButton: "View My RSVP",
	LookupPersonal:   "This link is personal — please don't share it.",
	LookupPlainTitle: "Find Your RSVP",
	LookupPlainBody:  "Click the link below to view and manage your RSVP for %s:",
	LookupPlainNote:  "This link is personal — please don't share it.",

	WaitlistHeading:    "A Spot Opened Up!",
	WaitlistIntroPre:   "Great news! A spot opened up for ",
	WaitlistIntroPost:  ". You are now attending.",
	ViewYourRSVPButton: "View Your RSVP",
	WaitlistPlainTitle: "A Spot Opened Up!",
	WaitlistPlainBody:  "Great news! A spot opened up for %s. You are now attending.",
	WaitlistPlainView:  "View your RSVP:",

	CohostHeading:    "You've Been Added as a Co-Host",
	CohostIntroMid:   " has added you as a co-host for ",
	CohostIntroPost:  ". You can now manage RSVPs, send messages, and help run the event.",
	CohostPlainTitle: "You've Been Added as a Co-Host",
	CohostPlainBody:  "%s has added you as a co-host for %s.",
	CohostPlainView:  "View the event dashboard:",

	SubjMagic:            "Sign in to OpenRSVP",
	SubjRSVPConfirmation: "RSVP Confirmation",
	SubjNewRSVP:          "New RSVP",
	SubjInvited:          "You're Invited",
	SubjCohost:           "You've been added as a co-host",
	SubjWaitlist:         "A spot opened up!",
	SubjNewMessage:       "New message from",
	SubjCancelled:        "Event Cancelled",
	SubjRetention:        "Data Retention Notice",
	SubjReminder:         "Event Reminder",
	SubjFeedbackConfirm:  "We received your feedback — OpenRSVP",
}

var frMsg = Msg{
	FooterTagline:  "Des invitations simples",
	ButtonFallback: "Si le bouton ne fonctionne pas, copiez et collez ce lien dans votre navigateur :",
	LabelEvent:     "Événement",
	LabelDate:      "Date",
	LabelLocation:  "Lieu",

	StatusAttending:  "Participe",
	StatusMaybe:      "Peut-être",
	StatusDeclined:   "Ne peut pas venir",
	StatusPending:    "En attente",
	StatusWaitlisted: "Sur liste d'attente",

	MagicHeading:     "Connectez-vous à votre compte",
	MagicIntro:       "Cliquez sur le bouton ci-dessous pour vous connecter à OpenRSVP. Aucun mot de passe requis.",
	MagicButton:      "Se connecter",
	MagicExpiryPre:   "Ce lien expire dans ",
	MagicExpiryPost:  " minutes. Si vous n'avez pas demandé ce lien, vous pouvez ignorer cet e-mail.",
	MagicPlainTitle:  "Connexion à OpenRSVP",
	MagicPlainBody:   "Cliquez sur le lien ci-dessous pour vous connecter :",
	MagicPlainExpiry: "Ce lien expire dans %d minutes.",
	MagicPlainIgnore: "Si vous n'avez pas demandé ce lien, vous pouvez ignorer cet e-mail.",

	RSVPConfHeading:    "Réponse confirmée",
	RSVPConfIntro:      "Votre réponse a été enregistrée. Voici les détails :",
	LabelYourRSVP:      "Votre réponse",
	ModifyButton:       "Modifier ma réponse",
	RSVPConfPlainTitle: "Réponse confirmée",

	ReminderHeading:     "Rappel d'événement",
	ReminderIntro:       "Ceci est un petit rappel concernant un événement à venir.",
	OrganizerMessage:    "Message de l'organisateur",
	ViewInviteButton:    "Voir l'invitation",
	ReminderPlainTitle:  "Rappel d'événement",
	ReminderPlainOrgMsg: "Message de l'organisateur :",
	ReminderPlainInvite: "Voir votre invitation :",

	RetentionHeading:     "Avis de conservation des données",
	RetentionIntro:       "Les données de votre événement sont bientôt programmées pour une suppression automatique. Ceci est un rappel de courtoisie afin que vous puissiez agir si nécessaire.",
	LabelDeletionDate:    "Date de suppression",
	RetentionBody:        "Après cette date, toutes les données de l'événement, y compris les réponses des invités, les messages et les cartes d'invitation, seront définitivement supprimées. Si vous souhaitez conserver ces données, connectez-vous et prolongez la période de conservation.",
	ViewEventButton:      "Voir l'événement",
	RetentionFooter1:     "Ceci est un avis automatique d'OpenRSVP.",
	RetentionFooter2:     "Vous recevez cet e-mail car vous êtes l'organisateur de cet événement.",
	RetentionPlainTitle:  "Avis de conservation des données",
	RetentionPlainBody:   "Votre événement « %s » est programmé pour une suppression automatique le %s.",
	RetentionPlainNote:   "Après cette date, toutes les données de l'événement, y compris les réponses des invités, les messages et les cartes d'invitation, seront définitivement supprimées.",
	RetentionPlainExtend: "Pour prolonger la période de conservation, rendez-vous sur :",

	OrgNotifyHeading:    "Nouvelle réponse reçue",
	OrgNotifyIntroPre:   "Quelqu'un a répondu à votre événement ",
	OrgNotifyIntroPost:  ".",
	LabelGuest:          "Invité",
	LabelResponse:       "Réponse",
	LabelEmail:          "E-mail",
	LabelPhone:          "Téléphone",
	LabelAdditional:     "Invités supplémentaires",
	ViewDashboardButton: "Voir le tableau de bord",
	OrgNotifyPlainTitle: "Nouvelle réponse reçue",
	OrgNotifyPlainView:  "Voir le tableau de bord de votre événement :",

	FeedbackHeading:    "Merci pour votre retour !",
	FeedbackIntroPre:   "Nous avons bien reçu votre ",
	FeedbackIntroPost:  " et vous remercions d'avoir pris le temps de le partager avec nous.",
	FeedbackFollowUp:   "Comme vous avez accepté d'être recontacté, nous pourrons vous écrire à cette adresse si nous avons des questions ou des nouvelles concernant votre retour.",
	FeedbackOutro:      "Vos retours nous aident à améliorer OpenRSVP pour tout le monde.",
	FeedbackPlainTitle: "Merci pour votre retour !",
	FeedbackPlainBody:  "Nous avons bien reçu votre %s et vous remercions d'avoir pris le temps de le partager avec nous.",
	FeedbackPlainOutro: "Vos retours nous aident à améliorer OpenRSVP pour tout le monde.",

	LookupHeading:    "Retrouver ma réponse",
	LookupIntroPre:   "Cliquez sur le bouton ci-dessous pour consulter et gérer votre réponse pour ",
	LookupIntroPost:  ".",
	ViewMyRSVPButton: "Voir ma réponse",
	LookupPersonal:   "Ce lien est personnel — merci de ne pas le partager.",
	LookupPlainTitle: "Retrouver ma réponse",
	LookupPlainBody:  "Cliquez sur le lien ci-dessous pour consulter et gérer votre réponse pour %s :",
	LookupPlainNote:  "Ce lien est personnel — merci de ne pas le partager.",

	WaitlistHeading:    "Une place s'est libérée !",
	WaitlistIntroPre:   "Bonne nouvelle ! Une place s'est libérée pour ",
	WaitlistIntroPost:  ". Vous participez désormais.",
	ViewYourRSVPButton: "Voir ma réponse",
	WaitlistPlainTitle: "Une place s'est libérée !",
	WaitlistPlainBody:  "Bonne nouvelle ! Une place s'est libérée pour %s. Vous participez désormais.",
	WaitlistPlainView:  "Voir votre réponse :",

	CohostHeading:    "Vous avez été ajouté comme co-organisateur",
	CohostIntroMid:   " vous a ajouté comme co-organisateur de ",
	CohostIntroPost:  ". Vous pouvez désormais gérer les réponses, envoyer des messages et aider à organiser l'événement.",
	CohostPlainTitle: "Vous avez été ajouté comme co-organisateur",
	CohostPlainBody:  "%s vous a ajouté comme co-organisateur de %s.",
	CohostPlainView:  "Voir le tableau de bord de l'événement :",

	SubjMagic:            "Connexion à OpenRSVP",
	SubjRSVPConfirmation: "Confirmation de réponse",
	SubjNewRSVP:          "Nouvelle réponse",
	SubjInvited:          "Vous êtes invité",
	SubjCohost:           "Vous avez été ajouté comme co-organisateur",
	SubjWaitlist:         "Une place s'est libérée !",
	SubjNewMessage:       "Nouveau message de",
	SubjCancelled:        "Événement annulé",
	SubjRetention:        "Avis de conservation des données",
	SubjReminder:         "Rappel d'événement",
	SubjFeedbackConfirm:  "Nous avons bien reçu votre retour — OpenRSVP",
}

// defaultLocale is the instance-wide email locale. Stored atomically so the
// startup writer and request-time readers do not race under -race.
var defaultLocale atomic.Value // string

// SetDefaultLocale sets the instance-wide email locale ("en" or "fr").
// Called once at startup from the loaded config.
func SetDefaultLocale(locale string) {
	if locale != "fr" {
		locale = "en"
	}
	defaultLocale.Store(locale)
}

// Locale returns the message catalog for the configured instance locale.
func Locale() Msg {
	if v, ok := defaultLocale.Load().(string); ok && v == "fr" {
		return frMsg
	}
	return enMsg
}
