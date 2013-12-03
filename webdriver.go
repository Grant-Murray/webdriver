// Package webdriver provides utility functions that aid in writing selenium webdriver tests
package webdriver

import (
	"fmt"
	"github.com/sourcegraph/go-selenium"
	"io/ioutil"
	"time"
)

// globals
var (
	Drv             selenium.WebDriver
	WaitForTimedOut bool
)

const RemoteURL = "http://localhost:4444/wd/hub"
const InitialWait time.Duration = 500 * time.Millisecond

// InitializeRemote establishes the connection to the remote selenium server
func InitializeRemote() (err error) {
	// GLM(self) run /opt/selenium/start-server.sh to start the server
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	if Drv, err = selenium.NewRemote(caps, RemoteURL); err != nil {
		err = fmt.Errorf("Failure calling selenium.NewRemote for %s: %s\n", RemoteURL, err)
	}
	return err
}

// ScreenshotToFile takes a screenshot and writes it to filename
func ScreenshotToFile(filename string) (err error) {

	screenshot, err := Drv.Screenshot()
	if err != nil {
		err = fmt.Errorf("Error during ScreenshotToFile using filename %s\n  Error:%s\n", filename, err)
		return err
	} else {
		ioutil.WriteFile(filename, screenshot, 0644)
	}
	return nil
}

// ElementToVanish is a WaitFor function. As long as the element is present
// waiting continues, once the element cannot be found waiting stops
func ElementToVanish(sel []interface{}) bool {
	_, err := Drv.FindElement(selenium.ByCSSSelector, sel[0].(string))
	if err.Error() == "no such element" {
		return true
	}
	return false
}

// ElementToAppear is a WaitFor function. As long as the element is absent
// waiting continues, once the element is found waiting stops
func ElementToAppear(sel []interface{}) bool {
	_, err := Drv.FindElement(selenium.ByCSSSelector, sel[0].(string))
	if err == nil {
		return true
	}
	return false
}

// DocumentIsReady can be used as a WaitFor isReady parameter
func DocumentIsReady(unused []interface{}) bool {
	result, err := Drv.ExecuteScript("return document.readyState", nil)
	if err != nil && result == "complete" {
		return true
	}
	return false
}

// UrlIsCurrent can be used as a WaitFor isReady parameter
func UrlIsCurrent(urls []interface{}) bool {
	cur, err := Drv.CurrentURL()
	if err != nil {
		return false
	}

	for i := 0; i < len(urls); i++ {
		if cur == urls[i].(string) {
			return true
		}
	}
	return false
}

// WaitFor sleeps until isReady() returns true unless it waits as long as timeoutAfter then it sets WaitForTimedOut to true and returns
func WaitFor(timeoutAfter time.Duration, isReady func([]interface{}) bool, args ...interface{}) {
	time.Sleep(InitialWait)

	if InitialWait >= timeoutAfter {
		WaitForTimedOut = true
		return
	}

	const ITERATIONS = 10
	var sleepDuration time.Duration = (timeoutAfter - InitialWait) / ITERATIONS

	for i := 0; i <= ITERATIONS; i++ {
		if isReady(args) {
			WaitForTimedOut = false
			return
		}
		time.Sleep(sleepDuration)
	}

	WaitForTimedOut = true
}

// FindNamedElements returns a map of elements with a memeber for each name in names
func FindNamedElements(names []string) (elements map[string]selenium.WebElement, err error) {

	elements = make(map[string]selenium.WebElement, 100)

	for _, n := range names {
		sel := fmt.Sprintf("[name=\"%s\"]", n)
		elements[n], err = Drv.FindElement(selenium.ByCSSSelector, sel)
		if err != nil {
			err = fmt.Errorf("Error finding element %s: %s", sel, err)
			return elements, err
		}
	}
	return elements, nil
}

// FetchText returns the msg text in an element ByCSSSelector sel
func FetchText(sel string) (msg string, err error) {
	var e selenium.WebElement

	e, err = Drv.FindElement(selenium.ByCSSSelector, sel)
	if err != nil {
		err = fmt.Errorf("Failed to find element %s (%s)\n", sel, err)
		return "", err
	}
	msg, err = e.Text()
	if err != nil {
		err = fmt.Errorf("Failed to retrieve %s text: %s", sel, err)
		return "", err
	}
	return msg, nil
}
