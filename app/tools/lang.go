package tools

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func AdminLang(key string, lang string) string {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	err := filepath.Walk(filepath.Join(path.Base(), "lang/admin"), func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".json" && !info.IsDir() {
            _, err := bundle.LoadMessageFile(path)
            if err != nil {
                facades.Log().Info("failed to load translation file %s: %v", path, err)
            }
        }
        return nil
    })

    if err != nil {
        facades.Log().Fatal("failed to walk through the directory: %v", err)
    }
	localizer := i18n.NewLocalizer(bundle, "zh")
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key})

}

func Lang(key string, lang string) string {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	err := filepath.Walk(filepath.Join(path.Base(), "lang"), func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".json" && !info.IsDir() {
            _, err := bundle.LoadMessageFile(path)
            if err != nil {
                facades.Log().Info("failed to load translation file %s: %v", path, err)
            }
        }
        return nil
    })

    if err != nil {
        facades.Log().Fatal("failed to walk through the directory: %v", err)
    }
	localizer := i18n.NewLocalizer(bundle, "zh")
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: key})

}