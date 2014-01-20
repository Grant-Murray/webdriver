package plog

import (
  "fmt"
  //"github.com/sourcegraph/go-selenium"
  "github.com/Grant-Murray/webdriver"
  "io/ioutil"
  "os"
  "strings"
  "testing"
  "time"
)

var (
  EmailAddr string
  Token     string
)

func Test_Verify_Setup(t *testing.T) {
  if err := webdriver.InitializeRemote(); err != nil {
    t.Fatalf("Cannot connect to selenium server: %s", err)
  }
}

type vCase struct {
  email string
  tok   string
  msg   string
}

func doVerifyCase(t *testing.T, cur vCase) {
  url := fmt.Sprintf("https://plog.org:8004/#/verify/%s/token/%s", cur.email, cur.tok)
  err := webdriver.Drv.Get(url)
  if err != nil {
    t.Fatalf("Failed to load %s: %s\n", url, err)
  }
  webdriver.WaitFor(5*time.Second, webdriver.ElementToVanish, "div[class='selenium-flag']")

  actualMsg, err := webdriver.FetchText("p[name='Message']")
  if err != nil {
    t.Fatalf("Failed to fetch the main message: %s", err)
  }

  if actualMsg != cur.msg {
    t.Fatalf("Expected message \"%s\", got \"%s\"", cur.msg, actualMsg)
  }
}

// SlurpEmail slurps the email file for EmailAddr and returns it as a string, the email file is deleted as part of the process.
func SlurpEmail(EmailAddr string, t *testing.T) string {
  mailFile := "/tmp/mailbot.boxes/" + strings.ToLower(EmailAddr[:strings.Index(EmailAddr, `@`)])
  var emailB []byte
  var err error

  for i := 0; i < 5; i++ {
    emailB, err = ioutil.ReadFile(mailFile)
    if err == nil {
      break
    } else {

      if i == 4 {
        t.Fatalf("Unable to read the verification email in %s got err %s", mailFile, err)
      } else {
        t.Logf("Waiting for email to show up in %s (%s)", mailFile, err)
      }
    }
  }

  err = os.Remove(mailFile)
  if err != nil {
    t.Fatalf("Error attempting to remove %s: %s", mailFile, err)
  }

  return string(emailB)
}

func VerifyEmailAddressFor(expectedAddr string, t *testing.T) {

  email := SlurpEmail(expectedAddr, t)

  // parse out the link to get email address and token
  s1 := email[strings.Index(email, "verify/")+7:] // s1 starts at email address
  EmailAddr = s1[:strings.Index(s1, "/")]
  s2 := email[strings.Index(email, "token/")+6:] // s2 starts at the token value
  Token = s2[:strings.Index(s2, "\n")]

  if EmailAddr != expectedAddr {
    t.Fatalf("Failed to read the email address from the file, expected <%s>, got <%s>", expectedAddr, EmailAddr)
  }

  doVerifyCase(t, vCase{EmailAddr, Token, "Success! Next step: Login and enjoy"})

}

func Test_VerifyEmail_BadAddress(t *testing.T) {

  cases := []vCase{
    {"WTF-Email", "6ba7b814-9dad-11d1-80b4-00c04fd430c8", "Failed! Server says: Verification failed; Not a valid email address"},
    {"matches@pattern.io", "toktok", "Failed! Server says: Verification failed; Not a valid token"},
    {"matches@pattern.io", "6ba7b814-9dad-11d1-80b4-00c04fd430c8", "Failed! Server says: Verification failed"},
  }

  for c := 0; c < len(cases); c++ {
    cur := cases[c]
    t.Logf("Case %d: EmailAddr=%s Token=%s", c, cur.email, cur.tok)
    doVerifyCase(t, cur)
  }
}

func Test_Verify_Success(t *testing.T) {

  VerifyEmailAddressFor(`jplain@mailbot.net`, t)
}

func Test_Verify_Breakdown(t *testing.T) {
  webdriver.Drv.Quit()
}
