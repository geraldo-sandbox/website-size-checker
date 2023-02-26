package main

import (
	"encoding/json"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestVisitStruct(t *testing.T) {
	v := &Visit{}
	metaValue := reflect.ValueOf(v).Elem()

	for _, name := range []string{"Url", "BodySize", "Error"} {
		field := metaValue.FieldByName(name)
		if field == (reflect.Value{}) {
			t.Errorf("field %s not exist in struct", name)
		}
	}
}

func TestContains(t *testing.T) {
	if !contains([]string{"GET", "POST"}, "POST") { // item exists in the list
		t.Errorf("item 'POST' must be found in the slice {\"GET\", \"POST\"}")
	}

	if contains([]string{"GET", "POST"}, "PUT") { // item does not exists in the list
		t.Errorf("item 'PUT' must NOT be found in the slice {\"GET\", \"POST\"}")
	}

	if contains([]string{"GET", "POST"}, "get") { // is case sensitive
		t.Errorf("item 'get' must NOT be found in the slice {\"GET\", \"POST\"}")
	}
}

func TestVisitSync(t *testing.T) {
	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()
	respPayload := map[string]interface{}{
		"hello": "I'm a test",
	}
	b, err := json.Marshal(respPayload)
	assert.NoError(t, err, "marshall response payload failed")

	httpmock.RegisterResponder("GET", "http://123.de", httpmock.NewBytesResponder(200, b))

	visit := visitUrl(client, "http://123.de", "GET", 10, false, false, nil) // no need to use the WaitGroup here
	assert.NoError(t, visit.Error, "the request must not fail")
	assert.Equal(t, "http://123.de", visit.Url, "incorrect URL")
	assert.GreaterOrEqual(t, visit.BodySize, len(b), "BodySize must be greater than 22 bytes. The reason "+
		"for that is the http response contains more data - headers")
}

func TestVisitSyncError(t *testing.T) {
	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://123.de", httpmock.NewErrorResponder(errors.New("timeout")))

	visit := visitUrl(client, "http://123.de", "GET", 10, false, false, nil) // no need to use the WaitGroup here

	assert.Error(t, visit.Error, "the request MUST fail")
	assert.Equal(t, "http://123.de", visit.Url, "incorrect URL")
	assert.Equal(t, 0, visit.BodySize, "BodySize must be 0 bytes because of the "+
		"http error response")
}

func TestVisitAsync(t *testing.T) {
	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()
	respPayload := map[string]interface{}{
		"hello": "I'm a test",
	}
	b, err := json.Marshal(respPayload)
	assert.NoError(t, err, "marshall response payload failed")

	httpmock.RegisterResponder("GET", "http://1234.de", httpmock.NewBytesResponder(200, b))
	visitQ := make(chan Visit, 1)
	_ = visitUrl(client, "http://1234.de", method, timeout, true, true, visitQ) // no need to use the WaitGroup here

	select {
	case visit := <-visitQ:
		assert.NoError(t, visit.Error, "the request must not fail")
		assert.Equal(t, "http://1234.de", visit.Url, "incorrect URL")
		assert.GreaterOrEqual(t, visit.BodySize, len(b), "BodySize must be greater than 22 bytes. The reason "+
			"for that is the http response contains more data - headers")
	case <-time.After(30 * time.Second):
		t.Error("fail to process request async")
	}

}

func TestVisitAsyncError(t *testing.T) {
	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://1234.de", httpmock.NewErrorResponder(errors.New("timeout")))
	visitQ := make(chan Visit, 1)
	_ = visitUrl(client, "http://1234.de", method, timeout, true, true, visitQ) // no need to use the WaitGroup here

	select {
	case visit := <-visitQ:
		assert.Error(t, visit.Error, "the request MUST fail")
		assert.Equal(t, "http://1234.de", visit.Url, "incorrect URL")
		assert.Equal(t, 0, visit.BodySize, "BodySize must be 0 bytes because of the "+
			"http error response")
	case <-time.After(30 * time.Second):
		t.Error("fail to process request async")
	}
}

func TestSortingAscending(t *testing.T) {
	visits := []*Visit{
		{Url: "https://www.10.de/", BodySize: 1386124},
		{Url: "https://www.7.de/", BodySize: 997149},
		{Url: "https://www.6.de/", BodySize: 915400},
		{Url: "https://www.9.net/", BodySize: 1183526},
		{Url: "https://www.8.de/", BodySize: 1021858},
	}
	sortVisits(visits, true)

	assert.Equal(t, 915400, visits[0].BodySize, "First visit BodySize must be 915400 bytes")
	assert.Equal(t, "https://www.6.de/", visits[0].Url, "First URL visit must be 'https://www.6.de/'")

	assert.Equal(t, 1386124, visits[4].BodySize, "Last visit BodySize must be 1386124 bytes")
	assert.Equal(t, "https://www.10.de/", visits[4].Url, "Last URL visit must be 'https://www.10.de/'")
}

func TestSortingDescending(t *testing.T) {
	visits := []*Visit{
		{Url: "https://www.10.de/", BodySize: 1386124},
		{Url: "https://www.7.de/", BodySize: 997149},
		{Url: "https://www.6.de/", BodySize: 915400},
		{Url: "https://www.9.net/", BodySize: 1183526},
		{Url: "https://www.8.de/", BodySize: 1021858},
	}
	sortVisits(visits, false)

	assert.Equal(t, 915400, visits[4].BodySize, "First visit BodySize must be 915400 bytes")
	assert.Equal(t, "https://www.6.de/", visits[4].Url, "First URL visit must be 'https://www.6.de/'")

	assert.Equal(t, 1386124, visits[0].BodySize, "Last visit BodySize must be 1386124 bytes")
	assert.Equal(t, "https://www.10.de/", visits[0].Url, "Last URL visit must be 'https://www.10.de/'")
}
