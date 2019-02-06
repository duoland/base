package i18n

import "testing"

func TestUsingTextI18n(t *testing.T) {
	msgKeys := []string{
		"service.a.name",
		"service.a.description",
		"service.b.name",
		"service.b.description",
	}

	msgKeysWithParam := []string{
		"service.a.executor",
		"service.b.executor",
	}

	locale := "zh-CN"
	Default(locale)

	InitTranslations("./locales")
	for _, key := range msgKeys {
		t.Log(Tr(locale, key))
	}

	for _, key := range msgKeysWithParam {
		t.Log(Tr(locale, key, "John Smith"))
	}

}
