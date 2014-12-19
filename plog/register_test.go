package plog

import (
	"fmt"
	//"github.com/sourcegraph/go-selenium"
	"code.grantmurray.com/webdriver"
	"testing"
	"time"
)

// RegisterUser is used to register new users in tests
type RegisterUser struct {
	UserId          string
	FirstName       string
	LastName        string
	EmailAddr       string
	ClearPassword   string
	ConfirmPassword string
}

// SubmitRegistration fills in the form and presses the register button. It does not wait after the click.
func SubmitRegistration(regU RegisterUser, t *testing.T) {

	// Get and wait
	err := webdriver.Drv.Get(`https://plog.org:8004/#/register`)
	if err != nil {
		t.Fatalf("Failed to load page: %s", err)
	}
	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	// Check page title
	pageTitle := "ZM Plog"
	if title, err := webdriver.Drv.Title(); err == nil {
		if title != pageTitle {
			t.Fatalf("Title page, expected \"%s\" but actually got \"%s\"\n", pageTitle, title)
		}
	} else {
		t.Fatalf("Failed to get page title: %s", err)
	}

	// Get elements
	elements, err := webdriver.FindNamedElements([]string{"UserId", "FirstName", "LastName", "EmailAddr",
		"TzName", "ClearPassword", "ConfirmPassword", "RegisterButton"})
	if err != nil {
		t.Fatalf("FindNamedElements failed: %s", err)
	}

	// Fill
	elements["UserId"].SendKeys(regU.UserId)
	elements["FirstName"].SendKeys(regU.FirstName)
	elements["LastName"].SendKeys(regU.LastName)
	elements["EmailAddr"].SendKeys(regU.EmailAddr)
	elements["ClearPassword"].SendKeys(regU.ClearPassword)
	elements["ConfirmPassword"].SendKeys(regU.ConfirmPassword)

	// Submit
	elements["RegisterButton"].Click()

	// NOT waiting here
}

// ExpectRegistrationSuccess gets call immediately after SubmitRegistration and verifies a successful registration
func ExpectRegistrationSuccess(email string, t *testing.T) {

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	msg, err := webdriver.FetchText("p[name='Message']")
	if err != nil {
		t.Fatalf("Failed to fetch the main message: %s", err)
	}

	expectMsg := fmt.Sprintf("Registration successful. Next step: Check your %s inbox and verify your email address.", email)

	if msg != expectMsg {
		t.Fatalf("Main message did not match expected, got %s", msg)
	}
}

var userOne RegisterUser = RegisterUser{"Selenium-One", "George", "Katsiopolous", "GeorgeK@mailbot.NET", "sldkfjeowir9", "sldkfjeowir9"}

// UserProfile is used to make preofile update
type UserProfile struct {
	UserId          string
	FirstName       string
	LastName        string
	EmailAddr       string
	ClearPassword   string
	ConfirmPassword string
}

var profExpected map[string]string = map[string]string{
	"UserId":          "selenium-one",
	"FirstName":       "George",
	"LastName":        "Katsiopolous",
	"EmailAddr":       "georgek@mailbot.net",
	"ClearPassword":   "",
	"ConfirmPassword": ""}

// GotoProfile attempts to load the url, but we could end up on the login page if we are not logged in
func GotoProfile(t *testing.T) {
	// Get and wait
	err := webdriver.Drv.Get(`https://plog.org:8004/#/profile`)
	if err != nil {
		t.Fatalf("Failed to load page: %s", err)
	}
}

func ExpectOnProfilePage(t *testing.T) {

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	elist := []string{"editProfileForm", "Message", "UserId", "FirstName", "LastName", "EmailAddr",
		"ClearPassword", "ConfirmPassword", "SaveProfileButton"}

	// Get elements
	elements, err := webdriver.FindNamedElements(elist)
	if err != nil {
		t.Fatalf("FindNamedElements failed: %s", err)
	}

	var val string

	for k, v := range profExpected {
		val, err = elements[k].GetAttribute("value")
		if err != nil {
			t.Fatalf("Failed to get vlaue in %s: %s", k, err)
		}

		if val != v {
			t.Fatalf("Expected %s value \"%s\", got \"%s\"", k, v, val)
		}
	}

}

// SubmitProfileChange fills in the form and presses the save button. It does not wait after the click.
func SubmitProfileChange(profU UserProfile, t *testing.T) {

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	// Get elements
	elements, err := webdriver.FindNamedElements([]string{"UserId", "FirstName", "LastName", "EmailAddr",
		"TzName", "ClearPassword", "ConfirmPassword", "SaveProfileButton"})
	if err != nil {
		t.Fatalf("FindNamedElements failed: %s", err)
	}

	// Fill
	if profU.UserId != "" {
		elements["UserId"].Clear()
		elements["UserId"].SendKeys(profU.UserId)
	}
	if profU.FirstName != "" {
		elements["FirstName"].Clear()
		elements["FirstName"].SendKeys(profU.FirstName)
	}
	if profU.LastName != "" {
		elements["LastName"].Clear()
		elements["LastName"].SendKeys(profU.LastName)
	}
	if profU.EmailAddr != "" {
		elements["EmailAddr"].Clear()
		elements["EmailAddr"].SendKeys(profU.EmailAddr)
	}
	if profU.ClearPassword != "" {
		elements["ClearPassword"].Clear()
		elements["ClearPassword"].SendKeys(profU.ClearPassword)
	}
	if profU.ConfirmPassword != "" {
		elements["ConfirmPassword"].Clear()
		elements["ConfirmPassword"].SendKeys(profU.ConfirmPassword)
	}

	// Submit
	elements["SaveProfileButton"].Click()

	// NOT waiting here
}

func ExpectProfileChangeSuccess(t *testing.T) {

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	// Get elements
	elements, err := webdriver.FindNamedElements([]string{"Message", "UserId", "FirstName", "LastName", "EmailAddr",
		"TzName", "ClearPassword", "ConfirmPassword", "SaveProfileButton"})
	if err != nil {
		t.Fatalf("FindNamedElements failed: %s", err)
	}

	msg, err := elements["Message"].Text()
	if err != nil {
		t.Fatalf("Could not get Message")
	}
	if msg != "Save successful" {
		t.Fatalf("Main message did not reflect success, got %s", msg)
	}
}

/******** Tests Start Here *********/

func Test_Setup(t *testing.T) {
	if err := webdriver.InitializeRemote(); err != nil {
		t.Fatalf("Cannot connect to selenium server: %s", err)
	}
}

func Test_Register_Success(t *testing.T) {
	SubmitRegistration(userOne, t)
	ExpectRegistrationSuccess(userOne.EmailAddr, t)
	VerifyEmailAddressFor(`georgek@mailbot.net`, t)

	var userTwo RegisterUser = RegisterUser{"Selenium-Two", "Jane", "Plain", "jplain@mailbot.NET", "neverguess", "neverguess"}
	SubmitRegistration(userTwo, t)
	ExpectRegistrationSuccess(userTwo.EmailAddr, t)
}

func Test_Register_AlreadyRegistered(t *testing.T) {

	SubmitRegistration(userOne, t)
	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	// Case 1: non-unique email
	const emailErr = "Already associated with a user"

	msg, err := webdriver.FetchText("span[name='EmailAddrErrorMsg']")
	if err != nil {
		t.Fatalf("Failed to fetch the email error message: %s", err)
	}

	if msg != emailErr {
		t.Fatalf("Expected error \"%s\" was in fact, \"%s\"", emailErr, msg)
	}

	// Case 2: non-unique user id
	const userErr = "Not available"
	origEmailAddr := userOne.EmailAddr
	userOne.EmailAddr = "anoterh@mailbot.NET  "

	SubmitRegistration(userOne, t)
	userOne.EmailAddr = origEmailAddr
	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	msg, err = webdriver.FetchText("span[name='UserIdErrorMsg']")
	if err != nil {
		t.Fatalf("Failed to fetch the user id error message: %s", err)
	}

	if msg != userErr {
		t.Fatalf("Expected error \"%s\" was in fact, \"%s\"", userErr, msg)
	}
}

func Test_Profile_NeedsLogin(t *testing.T) {

	GotoProfile(t)
	ExpectOnLoginPage(t)
	SubmitLogin(Login{userOne.UserId, userOne.ClearPassword}, t)
	ExpectOnProfilePage(t)

}

func Test_Profile_Success(t *testing.T) {

	var profU UserProfile
	profU.UserId = `Selenium-Changed`

	GotoProfile(t)
	ExpectOnProfilePage(t)
	SubmitProfileChange(profU, t)
	profExpected["UserId"] = "Selenium-Changed"
	ExpectProfileChangeSuccess(t)
	ExpectOnProfilePage(t)
}

func Test_Profile_FirstName_change(t *testing.T) {

	var profU UserProfile
	profU.FirstName = `NewFirstName`
	profExpected["UserId"] = "selenium-changed"

	GotoProfile(t)
	ExpectOnProfilePage(t)
	SubmitProfileChange(profU, t)
	profExpected["FirstName"] = profU.FirstName
	ExpectProfileChangeSuccess(t)
	ExpectOnProfilePage(t)
}

func Test_Profile_LastName_change(t *testing.T) {

	var profU UserProfile
	profU.LastName = `NewLastName`

	GotoProfile(t)
	ExpectOnProfilePage(t)
	SubmitProfileChange(profU, t)
	profExpected["LastName"] = profU.LastName
	ExpectProfileChangeSuccess(t)
	ExpectOnProfilePage(t)
}

func Test_Profile_EmailAddr_change(t *testing.T) {

	var profU UserProfile
	profU.EmailAddr = `bigdeal@little-planet.com`

	GotoProfile(t)
	ExpectOnProfilePage(t)
	SubmitProfileChange(profU, t)
	profExpected["EmailAddr"] = profU.EmailAddr
	ExpectProfileChangeSuccess(t)
	ExpectOnProfilePage(t)
	VerifyEmailAddressFor(profU.EmailAddr, t)
}

func Test_Profile_password_change(t *testing.T) {

	var profU UserProfile
	profU.ClearPassword = `New-Password-1234`
	profU.ConfirmPassword = profU.ClearPassword

	GotoProfile(t)
	ExpectOnProfilePage(t)
	SubmitProfileChange(profU, t)
	ExpectProfileChangeSuccess(t)
	ExpectOnProfilePage(t)

	// test password change by logging in
	Logout(t)
	GotoProfile(t)
	ExpectOnLoginPage(t)
	SubmitLogin(Login{profExpected["UserId"], profU.ClearPassword}, t)
	ExpectOnProfilePage(t)

}

func Test_Profile_change_all_back(t *testing.T) {

	var profU UserProfile = UserProfile{userOne.UserId, userOne.FirstName,
		userOne.LastName, userOne.EmailAddr, userOne.ClearPassword, userOne.ConfirmPassword}

	GotoProfile(t)
	ExpectOnProfilePage(t)
	SubmitProfileChange(profU, t)

	profExpected = map[string]string{
		"UserId":          "Selenium-One",
		"FirstName":       "George",
		"LastName":        "Katsiopolous",
		"EmailAddr":       "GeorgeK@mailbot.NET",
		"ClearPassword":   "",
		"ConfirmPassword": ""}
	ExpectProfileChangeSuccess(t)
	ExpectOnProfilePage(t)
	VerifyEmailAddressFor(`georgek@mailbot.net`, t)

	// test password change by logging in
	Logout(t)
	GotoProfile(t)
	ExpectOnLoginPage(t)
	SubmitLogin(Login{userOne.UserId, userOne.ClearPassword}, t)
	profExpected["UserId"] = `selenium-one`
	profExpected["EmailAddr"] = `georgek@mailbot.net`
	ExpectOnProfilePage(t)

}

func Test_Breakdown(t *testing.T) {
	webdriver.Drv.Quit()
}
