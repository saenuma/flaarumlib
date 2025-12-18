// this package 'flaarumlib' is the golang library for communicating with the flaarum server.
package flaarumlib

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpCl *http.Client

func init() {
	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}

	httpCl = &http.Client{Transport: tr}
}

type Client struct {
	Addr     string
	KeyStr   string
	ProjName string
}

func NewClient(ip, keyStr, projName string) Client {
	return Client{"https://" + ip + ":22318/", keyStr, projName}
}

// Used whenever you changed the default port
func NewClientCustomPort(ip, keyStr, projName string, port int) Client {
	return Client{"https://" + ip + fmt.Sprintf(":%d/", port), keyStr, projName}
}

func (cl *Client) Ping() error {
	deadline := time.Now().Add(2000 * time.Millisecond)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	defer cancelCtx()

	return cl.InnerPing(ctx)
}

func (cl *Client) InnerPing(ctx context.Context) error {
	urlValues := url.Values{}
	urlValues.Set("key-str", cl.KeyStr)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cl.Addr+"is-flaarum",
		strings.NewReader(urlValues.Encode()))
	if err != nil {
		return retError(10, err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpCl.Do(req)
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
