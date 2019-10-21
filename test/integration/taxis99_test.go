package integration

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
	"unsafe"

	"github.com/mobilitee-smartmob/taxis99"
)

const (
	envKey99TaxisAPIKey     = "TAXIS99_API_KEY"
	envKey99TaxisExternalID = "TAXIS99_EXTERNAL_ID"
	envKey99TaxisCompanyID  = "TAXIS99_COMPANY_ID"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var employeeFilter taxis99.Filter

// 99Taxis Client
var tx99 *taxis99.Client

var src = rand.NewSource(time.Now().UnixNano())

var logging = flag.Bool("log", false, "Define if tests should log the requests")

// Transport used to log the requests.
type transportLogger struct {
	base http.RoundTripper
}

func (t *transportLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		reqBody []byte
		err     error
	)

	if req.Body != nil {
		reqBody, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()

		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	}

	res, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)

	log.Println("Request")
	log.Println("Headers:")
	for k, v := range req.Header {
		log.Printf("%s: %s\n", k, v)
	}
	log.Printf("/%s %s %s --> Response %s %s",
		req.Method, req.URL.String(), string(reqBody), res.Status, string(resBody))

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	return res, nil
}

func TestMain(m *testing.M) {
	flag.Parse()

	if externalID := os.Getenv(envKey99TaxisExternalID); externalID != "" {
		_, err := strconv.ParseInt(externalID, 10, 64)
		if err != nil {
			fmt.Printf("Invalid external ID: %s. Some tests may not run.", externalID)
		}

		employeeFilter = taxis99.Filter{"search": externalID}
	}

	companyID := os.Getenv(envKey99TaxisCompanyID)
	apiKey := os.Getenv(envKey99TaxisAPIKey)
	if apiKey == "" {
		fmt.Println("No API Key.")
		return
	}

	hc := &http.Client{
		Transport: &taxis99.Transport{
			Key:       apiKey,
			CompanyID: companyID,
			Base:      loggingHTTPClient().Transport,
		},
	}
	tx99 = taxis99.NewClient(hc)

	os.Exit(m.Run())
}

func randString(max int, rangeBytes string) string {
	b := make([]byte, max)

	for i, cache, remain := max-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMax); idx < len(rangeBytes) {
			b[i] = rangeBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func loggingHTTPClient() *http.Client {
	if !*logging {
		return http.DefaultClient
	}

	hc := &http.Client{
		Transport: &transportLogger{http.DefaultTransport},
	}

	return hc
}
