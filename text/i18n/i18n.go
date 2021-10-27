package i18n

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// defaultLocale is the default locale for the system
var defaultLocale = "zh-CN"

// i18nLocales holds all the locales supported
var i18nLocales map[string]bool

// i18Translations holds all the messages in different locales
var i18nTranslations map[string]map[string]string

// Locale returns the actually locale of the messages
func Locale(locale string) string {
	if _, ok := i18nLocales[locale]; ok {
		return locale
	} else {
		return defaultLocale
	}
}

// Default set the default local string
func Default(locale string) {
	defaultLocale = locale
}

// Tr gets the formatted messages by the specified locale and message key
func Tr(locale string, key string, params ...interface{}) (output string) {
	if messages, ok := i18nTranslations[locale]; !ok {
		output = key
	} else {
		if msgFmt, ok := messages[key]; !ok {
			output = key
		} else {
			output = fmt.Sprintf(msgFmt, params...)
		}
	}
	return
}

// InitTranslations read from the translation files of different locales.
func InitTranslations(localeDir string) (err error) {
	i18nLocales = make(map[string]bool)
	i18nTranslations = make(map[string]map[string]string)

	localeDir, err = filepath.Abs(localeDir)
	if err != nil {
		return
	}
	err = filepath.Walk(localeDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// read files with suffix
		if !strings.HasSuffix(path, ".ini") {
			fmt.Fprintf(os.Stdout, "[i18n] ingore file %s not with the suffix .ini", path)
			return nil
		}

		// parse locale and messages
		translationFileRelativePath := strings.TrimPrefix(strings.TrimPrefix(path, localeDir), string(os.PathSeparator))
		translationItems := strings.Split(translationFileRelativePath, string(os.PathSeparator))
		if len(translationItems) != 2 {
			fmt.Fprintf(os.Stderr, "[i18n] ignore the invalid translation file path %s\n", translationFileRelativePath)
			return nil
		}

		translationLocale := translationItems[0]

		// read translation file
		translationFp, openErr := os.Open(path)
		if openErr != nil {
			fmt.Fprintf(os.Stderr, "[i18n] parse the translation file error, %s\n", openErr.Error())
			return nil
		}

		// check the translation map
		if _, exists := i18nLocales[translationLocale]; !exists {
			i18nLocales[translationLocale] = true
			i18nTranslations[translationLocale] = make(map[string]string)
		}

		translationScanner := bufio.NewScanner(translationFp)
		for translationScanner.Scan() {
			messageLine := strings.TrimSpace(translationScanner.Text())
			if messageLine == "" || strings.HasPrefix(messageLine, "#") {
				// empty or comments, ignore
				continue
			}

			messageLineItems := strings.Split(messageLine, "=")
			if len(messageLineItems) != 2 {
				// invalid, ignore
				continue
			}

			// take the message item
			messageKey := strings.TrimSpace(messageLineItems[0])
			messageValue := strings.TrimSpace(messageLineItems[1])
			i18nTranslations[translationLocale][messageKey] = messageValue
		}

		return nil
	}) // walk ends
	return
}
