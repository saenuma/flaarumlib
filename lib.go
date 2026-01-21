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
)

type Client struct {
	Addr     string
	ProjName string
}

func NewClient(ip, projName string) Client {
	return Client{"http://" + ip + ":22318/", projName}
}

// Used whenever you changed the default port
func NewClientCustomPort(ip, projName string, port int) Client {
	return Client{"http://" + ip + fmt.Sprintf(":%d/", port), projName}
}

func (cl *Client) Ping() error {
	deadline := time.Now().Add(2000 * time.Millisecond)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	defer cancelCtx()

	return cl.InnerPing(ctx)
}

func (cl *Client) InnerPing(ctx context.Context) error {
	urlValues := url.Values{}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cl.Addr+"is-flaarum",
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
