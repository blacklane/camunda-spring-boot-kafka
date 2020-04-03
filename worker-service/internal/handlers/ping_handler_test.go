package handlers

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("ping", "/ping", nil)

	h := PingHandler()
	h.ServeHTTP(w, r)

	wantStatusCode := 200
	gotStatusCode := w.Result().StatusCode

	wantBody := "ok\n"
	bodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	gotBody := string(bodyBytes)

	if gotStatusCode != wantStatusCode {
		t.Errorf("got http status: %d, want: %d", gotStatusCode, wantStatusCode)
	}

	if gotBody != wantBody {
		t.Errorf("got response body: %s, want: %s", gotBody, wantBody)
	}
}
