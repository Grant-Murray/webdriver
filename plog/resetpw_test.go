package plog

import (
	"code.grantmurray.com/webdriver"
	"fmt"
	"strings"
	"testing"
	"time"
)

/******** Tests Start Here *********/

func Test_ResetPW_Setup(t *testing.T) {
	if err := webdriver.InitializeRemote(); err != nil {
		t.Fatalf("Cannot connect to selenium server: %s", err)
	}
}

func RequestPasswordResetFor(inEmailAddr string, t *testing.T) {
	url := "https://plog.org:8004/#/password"
	err := webdriver.Drv.Get(url)
	if err != nil {
		t.Fatalf("Failed to load %s: %s\n", url, err)
	}
	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	// Get elements
	elements, err := webdriver.FindNamedElements([]string{"Message", "EmailAddr", "ResetPasswordButton"})
	if err != nil {
		t.Fatalf("FindNamedElements failed: %s", err)
	}

	// Fill
	elements["EmailAddr"].Clear()
	elements["EmailAddr"].SendKeys(inEmailAddr)

	// Submit
	elements["ResetPasswordButton"].Click()

	webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

	// Verify expected results
	actualMsg, err := webdriver.FetchText("p[name='Message']")
	if err != nil {
		t.Fatalf("Failed to fetch the main message: %s", err)
	}

	expectMsg := fmt.Sprintf("Check email (%s) for a reset token", inEmailAddr)

	if actualMsg != expectMsg {
		t.Fatalf("Expected message \"%s\", got \"%s\"", expectMsg, actualMsg)
	}
}

func Test_ResetPW_request_success(t *testing.T) {
	RequestPasswordResetFor(userOne.EmailAddr, t)

	email := SlurpEmail(userOne.EmailAddr, t)

	// parse out the link to get email address and token
	s1 := email[strings.Index(email, "reset/")+7:] // s1 starts at email address
	EmailAddr = s1[:strings.Index(s1, "/")]
	s2 := email[strings.Index(email, "token/")+6:] // s2 starts at the token value
	Token = s2[:strings.Index(s2, "\n")]

}

func Test_ResetPW_Breakdown(t *testing.T) {
	webdriver.Drv.Quit()
}
