package i18n

type Language struct {
	Months     [12]string
	Dictionary map[string]string
}

var Languages = map[string]Language{
	"en": lang_en,
	"ru": lang_ru,
}

func getLang(lang string) Language {
	language, ok := Languages[lang]
	if !ok {
		language = Languages["en"]
	}
	return language
}

func GetMonthName(month int, lang string) string {
	language := getLang(lang)
	if month < 1 || month > 12 {
		return ""
	}
	return language.Months[month-1]
}

func T(key string, locale string) string {
	language := getLang(locale)
	if value, ok := language.Dictionary[key]; ok {
		return value
	}
	return "[" + key + "]"
}
