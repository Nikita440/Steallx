package browser

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/moond4rk/hackbrowserdata/browser/chromium"
	"github.com/moond4rk/hackbrowserdata/browser/firefox"
	"github.com/moond4rk/hackbrowserdata/browserdata"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
	"github.com/moond4rk/hackbrowserdata/utils/typeutil"
)

type Browser interface {
	// Name is browser's name
	Name() string
	// BrowsingData returns all browsing data in the browser.
	BrowsingData(isFullExport bool) (*browserdata.BrowserData, error)
}

// PickBrowsers returns a list of browsers that match the name and profile.
func PickBrowsers(name, profile string) ([]Browser, error) {
	var browsers []Browser
	clist := pickChromium(name, profile)
	for _, b := range clist {
		if b != nil {
			browsers = append(browsers, b)
		}
	}
	flist := pickFirefox(name, profile)
	for _, b := range flist {
		if b != nil {
			browsers = append(browsers, b)
		}
	}
	return browsers, nil
}

func pickChromium(name, profile string) []Browser {
	var browsers []Browser
	name = strings.ToLower(name)
	if name == "all" {
		for _, v := range chromiumList {
			if !fileutil.IsDirExists(filepath.Clean(v.profilePath)) {

				continue
			}
			multiChromium, err := chromium.New(v.name, v.storage, v.profilePath, v.dataTypes)
			if err != nil {

				continue
			}
			for _, b := range multiChromium {

				browsers = append(browsers, b)
			}
		}
	}
	if c, ok := chromiumList[name]; ok {
		if profile == "" {
			profile = c.profilePath
		}
		if !fileutil.IsDirExists(filepath.Clean(profile)) {

		}
		chromiumList, err := chromium.New(c.name, c.storage, profile, c.dataTypes)
		if err != nil {

		}
		for _, b := range chromiumList {

			browsers = append(browsers, b)
		}
	}
	return browsers
}

func pickFirefox(name, profile string) []Browser {
	var browsers []Browser
	name = strings.ToLower(name)
	if name == "all" || name == "firefox" {
		for _, v := range firefoxList {
			if profile == "" {
				profile = v.profilePath
			} else {
				profile = fileutil.ParentDir(profile)
			}

			if !fileutil.IsDirExists(filepath.Clean(profile)) {

				continue
			}

			if multiFirefox, err := firefox.New(profile, v.dataTypes); err == nil {
				for _, b := range multiFirefox {

					browsers = append(browsers, b)
				}
			} else {

			}
		}

		return browsers
	}

	return nil
}

func ListBrowsers() []string {
	var l []string
	l = append(l, typeutil.Keys(chromiumList)...)
	l = append(l, typeutil.Keys(firefoxList)...)
	sort.Strings(l)
	return l
}

func Names() string {
	return strings.Join(ListBrowsers(), "|")
}
