package validation

var ValidationMessages = map[string]FieldError{
	"required": {
		Message: "%s je obavezno polje",
		Reason:  "Ovo polje je obavezno i nedostaje u zahtjevu",
	},
	"email": {
		Message: "%s nije ispravan email",
		Reason:  "Unesena vrijednost ne odgovara standardnom formatu emaila",
	},
	"numeric": {
		Message: "%s mora biti broj",
		Reason:  "Dozvoljeni su samo numerički znakovi za ovo polje",
	},
	"datetime": {
		Message: "%s nije ispravan datum",
		Reason:  "Format datuma mora biti YYYY-MM-DD",
	},
	"min": {
		Message: "%s je prekratak",
		Reason:  "Ovo polje ne zadovoljava minimalne zahtjeve za duljinu",
	},
	"len": {
		Message: "%s ima neispravnu duljinu",
		Reason:  "Ovo polje ne zadovoljava zahtjeve za točnu duljinu prema pravilu validacije",
	},
	"oib": {
		Message: "%s mora sadržavati točno 11 znamenki",
		Reason:  "OIB zahtijeva točno 11 numeričkih znamenki",
	},
	"telephone": {
		Message: "%s nije ispravan broj telefona",
		Reason:  "Unesena vrijednost ne odgovara standardnom formatu broja telefona",
	},
}
