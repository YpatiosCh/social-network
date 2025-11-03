package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func buildHandlerUnderTest() http.Handler {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ms := Chain().
		AllowedMethod(http.MethodGet, http.MethodPost).
		Auth().       // include if your route requires it, otherwise omit
		BindReqMeta() // optional

	return ms.Finalize(next)
}

func TestOPTIONSConcurrency_NoDataRaces(t *testing.T) {
	h := buildHandlerUnderTest()

	var wg sync.WaitGroup
	const N = 300

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodOptions, "/ws", nil)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if rr.Code != http.StatusNoContent && rr.Code != http.StatusOK {
				t.Errorf("unexpected status: %d", rr.Code)
			}
		}()
	}
	wg.Wait()
}
