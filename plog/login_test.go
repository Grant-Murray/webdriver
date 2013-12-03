package plog

import (
	"fmt"
	"github.com/Grant-Murray/webdriver"
	"github.com/sourcegraph/go-selenium"
	"testing"
	"time"
)

type Login struct {
	UserIdentifier string
	ClearPassword  string
}

func ExpectOnLoginPage(t *testing.T) {

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	_, err := webdriver.Drv.FindElement(selenium.ByCSSSelector, "[name=\"loginForm\"]")
	if err != nil {
		t.Fatalf("Expected to be on the login page, but finding loginForm: %s", err)
	}

}

func GotoLogin(t *testing.T) {
	err := webdriver.Drv.Get(`https://plog.org:8004/`)
	if err != nil {
		t.Fatalf("Goto login failed: %s", err)
	}
}

func SubmitLogin(in Login, t *testing.T) {

	// Get elements
	elements, err := webdriver.FindNamedElements([]string{"UserIdentifier", "ClearPassword", "LoginButton"})
	if err != nil {
		t.Fatalf("FindNamedElements failed: %s", err)
	}

	// Fill in
	elements["UserIdentifier"].Clear()
	elements["UserIdentifier"].SendKeys(in.UserIdentifier)
	elements["ClearPassword"].Clear()
	elements["ClearPassword"].SendKeys(in.ClearPassword)

	// Submit
	elements["LoginButton"].Click()
}

func checkSessionToken(t *testing.T) bool {

	js := `
    var $injector = angular.injector(['ngCookies']);
    return $injector.invoke(function($cookies) {
      if ($cookies.SessionToken && ($cookies.SessionToken.length == 36)) {
        return "present";
      } else {
        return "absent";
      }
    });
    `

	// check for SessionToken - it should be present if we logged in successfully, absent otherwise
	result, err := webdriver.Drv.ExecuteScript(js, nil)
	if err != nil {
		t.Fatalf("Failed to ExecuteScript: %s", err)
	} else if result != "present" {
		return false
	}
	return true
}

func ExpectSessionToken(t *testing.T) {
	if !checkSessionToken(t) {
		t.Fatalf("SessionToken expected to be present but was missing")
	}
}

func ExpectNoSessionToken(t *testing.T) {
	if checkSessionToken(t) {
		t.Fatalf("SessionToken expected to be absent but was found")
	}
}

const (
	LogoutPageUrl = "https://plog.org:8004/#/album"
)

func Test_Login_Setup(t *testing.T) {
	if err := webdriver.InitializeRemote(); err != nil {
		t.Fatalf("Cannot connect to selenium server: %s", err)
	}
}

type loginCase struct {
	idTyp string
	lin   Login
	msg   string
	tokc  string
}

func Test_Login_table(t *testing.T) {

	cases := []loginCase{
		{"UserId", Login{"no such dude", "passwordpassword"}, "Authentication failed", "absent"},
		{"UserId", Login{"wrong uid", userOne.ClearPassword}, "Authentication failed", "absent"},
		{"UserId", Login{userOne.UserId, "wrong password"}, "Authentication failed", "absent"},
		{"UserId", Login{userOne.UserId, userOne.ClearPassword}, "SKIP", "present"},
		// if already logged in, automatic login will happen. So for testing should log out
	}

	for c := 0; c < len(cases); c++ {
		cur := cases[c]
		t.Logf("Case %d: UserIdentifier (%s)=%s ClearPassword=%s", c, cur.idTyp, cur.lin.UserIdentifier, cur.lin.ClearPassword)

		if len(cur.lin.ClearPassword) < 10 {
			t.Fatalf("Case is not valid since ClearPassword (%s) is too short")
		}

		GotoLogin(t)
		ExpectOnLoginPage(t)
		SubmitLogin(cur.lin, t)

		webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

		if cur.tokc == "absent" {
			// Login failure was expected
			msg, err := webdriver.FetchText("p[name='LoginMessage']")
			if err != nil {
				t.Fatalf("Failed to fetch the main message: %s", err)
			}

			if msg != "Authentication failed" {
				t.Errorf("Expected failure message, got \"%s\"", msg)
			}

			ExpectNoSessionToken(t)
		} else {
			ExpectSessionToken(t)
		}
	}
}

func Logout(t *testing.T) {

	ExpectSessionToken(t)

	err := webdriver.Drv.Get(LogoutPageUrl)
	if err != nil {
		t.Fatalf("Failed to load %s: %s\n", LogoutPageUrl, err)
	}
	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	sel := "a[name=\"Logout\"]"
	LogoutLink, err := webdriver.Drv.FindElement(selenium.ByCSSSelector, sel)
	if err != nil {
		t.Fatalf("Failed to find element %s (%s)\n", sel, err)
	}

	LogoutLink.Click()
	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	ExpectNoSessionToken(t)
}

func Test_Logout(t *testing.T) {
	// happy path logout
	t.Logf("Case: happy path logout")
	Logout(t)

	// visit while logged out; expect to get login page
	// TODO add an Image URL
	pagesNeedingLogin := []string{"/#/album", "/#/album/2013-09%20September"}

	for p := 0; p < len(pagesNeedingLogin); p++ {
		page := fmt.Sprintf("https://plog.org:8004%s", pagesNeedingLogin[p])
		t.Logf("Case: %s", pagesNeedingLogin[p])

		err := webdriver.Drv.Get(page)
		if err != nil {
			t.Fatalf("Failed to load %s: %s\n", page, err)
		}

		webdriver.WaitFor(5*time.Second, webdriver.UrlIsCurrent, "https://plog.org:8004/#/login")
		if webdriver.WaitForTimedOut {
			t.Fatalf("Failed to land at expected URL")
		}

		msg, err := webdriver.FetchText("p[name='LoginMessage']")
		if err != nil {
			t.Fatalf("Failed to fetch the main message: %s", err)
		}

		if msg != "" {
			t.Errorf("Message was not blank: %s", msg)
		}

		result, err := webdriver.Drv.ExecuteScript("return localStorage['SessionToken']", nil)
		if err != nil {
			t.Fatalf("Failed to ExecuteScript: %s", err)
		}

		if result != nil {
			t.Errorf("Found a SessionToken in localStorage: %s", result)
		}

	}

	// automatic login when SessionToken is still valid: goto login, then autologin and end up on album page
	t.Logf("Case: Automatic login that works")
	GotoLogin(t)
	ExpectOnLoginPage(t)
	SubmitLogin(Login{userOne.UserId, userOne.ClearPassword}, t)

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	ExpectSessionToken(t)
	GotoLogin(t)

	webdriver.WaitFor(5*time.Second, webdriver.UrlIsCurrent, `https://plog.org:8004/#/`)
	if webdriver.WaitForTimedOut {
		t.Fatalf("Failed to land at expected URL")
	}

	// automatic login when SessionToken is present but not valid
	t.Logf("Case: Automatic login with bad session token")
	_, err := webdriver.Drv.ExecuteScript("localStorage['SessionToken'] = '00000000-0000-0000-dead-beef00000000'", nil)
	if err != nil {
		t.Fatalf("Failed to ExecuteScript: %s", err)
	}

	GotoLogin(t)

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")
	if webdriver.WaitForTimedOut {
		t.Fatalf("Failed to end up at login page")
	}

	ExpectOnLoginPage(t)

	msg, err := webdriver.FetchText("p[name='LoginMessage']")
	if err != nil {
		t.Fatalf("Failed to fetch the login message: %s", err)
	}

	if msg != "" {
		// automatic logins don't produce "Authentication failed" messages
		t.Errorf("Expected LoginMessage to be blank, got \"%s\"", msg)
	}

	// TODO periodic re-login to extend session : Problem is how to test this without waiting 5 minutes

	// TODO login, visit page that calls loggedIn() - delete session in background - visit page that calls loggedIn() + expect it to return true - wait for background login attempt - visit page that calles loggedIn() + expect it to return false

}

func Test_Login_Breakdown(t *testing.T) {
	webdriver.Drv.Quit()
}
