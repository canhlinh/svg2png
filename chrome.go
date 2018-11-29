package svg2png

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	gover "github.com/mcuadros/go-version"
	"github.com/sirupsen/logrus"
)

const (
	// MinChromeVersion minimize version allowed chrome headless
	MinChromeVersion = "59"
)

var (

	// DefaultChromPaths posible chrome paths
	DefaultChromPaths = []string{
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe"}

	// DefaultChromeHeight default height value of chrome brower
	DefaultChromeHeight = 1200
	// DefaultChromeWidth default width value of chrome brower
	DefaultChromeWidth = 720
	// DefaultTimeout default timeout
	DefaultTimeout = 3 * time.Second
)

// Chrome contains information about a Google Chrome
// instance, with methods to run on it.
type Chrome struct {
	path    string
	version string
	height  int
	width   int
	timeout time.Duration
}

// NewChrome intalizes new Chrome
func NewChrome() *Chrome {
	path := getChromePath()
	if path == "" {
		logrus.Fatal("Chrome not found")
	}

	version := getChromeVersion(path)
	if gover.Compare(MinChromeVersion, version, ">") {
		logrus.Fatal("Chrome version must greater than 59")
	}

	logrus.Infof("Chrome version %s", version)

	return &Chrome{
		width:   DefaultChromeWidth,
		height:  DefaultChromeHeight,
		path:    path,
		version: version,
		timeout: DefaultTimeout,
	}
}

// SetWith sets chrome's width
func (chrome *Chrome) SetWith(w int) *Chrome {
	chrome.width = w
	return chrome
}

// SetHeight sets chrome's width
func (chrome *Chrome) SetHeight(h int) *Chrome {
	chrome.height = h
	return chrome
}

//Resolution returns chrome's resolution
func (chrome Chrome) Resolution() string {
	return fmt.Sprintf("%d,%d", chrome.width, chrome.height)
}

//SetTimeout set chrome's timeout
func (chrome *Chrome) SetTimeout(timeout time.Duration) *Chrome {
	chrome.timeout = timeout
	return chrome
}

func getChromePath() string {
	for _, path := range DefaultChromPaths {

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path
		}
	}

	return ""
}

// Version returns Chrome version
func getChromeVersion(path string) string {
	out, err := exec.Command(path, "-version").Output()
	if err != nil {
		logrus.Fatal(err)
	}

	match := regexp.MustCompile(`\d+(\.\d+)+`).FindStringSubmatch(string(out))
	if len(match) <= 0 {
		logrus.Fatal("Unable to determine Chrome version.")
	}

	return match[0]
}

func isValidDestination(destination string) bool {
	return path.Base(destination) != "." && path.Base(destination) != "/" && strings.EqualFold(path.Ext(destination), ".png")
}

// Screenshoot takes a screenshot by the given website url
func (chrome Chrome) Screenshoot(websiteURL string, destination string) error {
	if _, err := url.Parse(websiteURL); err != nil {
		return err
	}

	if !isValidDestination(destination) {
		return errors.New("Destination must be a png path")
	}

	args := []string{
		"--headless",
		"--no-sandbox",
		"--disable-crash-reporter",
		"--hide-scrollbars",
		"--default-background-color=00000000",
		"--disable-gpu",
		"--window-size=" + chrome.Resolution(),
		"--screenshot=" + destination,
		websiteURL,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), chrome.timeout)
	defer cancel()

	if err := exec.CommandContext(ctx, chrome.path, args...).Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("Takes screenshoot got timeout")
		}

		return err
	}

	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		return err
	}

	return nil
}
