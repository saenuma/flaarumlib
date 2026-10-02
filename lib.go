// this package 'flaarumlib' is the golang library for communicating with the flaarum server.
package flaarumlib

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/pkg/errors"
)

type Client struct {
	ProjName       string
	KeyFromFlaarum string
}

func haveAllowedChars(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII { // unicode.MaxASCII is 127
			return false
		}
	}
	return true
}

func NewClient(projName, keyStr string) (Client, error) {
	if len(keyStr) != 20 && !haveAllowedChars(keyStr) {
		return Client{}, errors.New(fmt.Sprintf("keyStr '%s' is either not length 20 or contains non ASCII characters", keyStr))
	}
	return Client{projName, keyStr}, nil
}

func (cl *Client) Ping() error {
	deadline := time.Now().Add(2000 * time.Millisecond)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	defer cancelCtx()

	return cl.innerPing(ctx)
}

func (cl *Client) innerPing(ctx context.Context) error {
	urlValues := url.Values{}
	urlValues.Add("tinything", cl.KeyFromFlaarum)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, DEFAULT_ADDR+"is-flaarum",
		strings.NewReader(urlValues.Encode()))
	if err != nil {
		return retError(10, err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return retError(10, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return retError(10, err.Error())
	}

	if resp.StatusCode == 200 {
		if string(body) == "yeah-flaarum" {
			return nil
		} else {
			return retError(10, "Unexpected Error in confirming that the server is a flaarum store.")
		}
	} else {
		return retError(10, string(body))
	}

}

// Converts a time.Time to the date format expected in flaarum
func RightDateFormat(d time.Time) string {
	return d.Format(DATE_FORMAT)
}

// Converts a time.Time to the datetime format expected in flaarum
func RightDateTimeFormat(d time.Time) string {
	return d.Format(DATETIME_FORMAT)
}
