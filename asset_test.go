package ams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetAsset(t *testing.T) {
	expected := testAsset("sample-id", "Sample")
	m := http.NewServeMux()
	m.HandleFunc(fmt.Sprintf("/Assets('%v')", expected.ID), func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, http.MethodGet)
		testAMSHeader(t, r.Header, false)

		if err := json.NewEncoder(w).Encode(expected); err != nil {
			t.Fatal(err)
		}
	})
	s := httptest.NewServer(m)
	defer s.Close()

	client, err := NewClient(s.URL, testTokenSource())
	if err != nil {
		t.Fatal(err)
	}

	actual, err := client.GetAsset(context.TODO(), expected.ID)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected asset. expected: %#v, actual: %#v", expected, actual)
	}
}