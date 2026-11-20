package web

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/petaki/inertia-go"
	"github.com/petaki/satellite/internal/config"
	"github.com/petaki/satellite/internal/models"
)

type loadProbeRepository struct {
	models.ProbeRepository
}

func (r *loadProbeRepository) Exists(probe models.Probe) (bool, error) {
	return probe == "p1", nil
}

type loadSeriesRepository struct {
	models.SeriesRepository
}

func (r *loadSeriesRepository) FindDiskPaths(models.Probe) ([]string, error) {
	return []string{}, nil
}

func (r *loadSeriesRepository) FindLoad(models.Probe, models.SeriesType) (models.Series, models.Series, models.Series, error) {
	return models.Series{}, models.Series{}, models.Series{}, nil
}

func (r *loadSeriesRepository) ChunkSize(models.SeriesType) int {
	return 1
}

type loadLogRepository struct {
	models.LogRepository
}

func (r *loadLogRepository) FindLogPaths(models.Probe) ([]string, error) {
	return []string{}, nil
}

type loadAlarmRepository struct{}

func (r *loadAlarmRepository) Find(models.Probe) (*models.Alarm, error) {
	return nil, models.ErrNoRecord
}

func newSeriesTestApp(store scs.Store) *app {
	sessionManager := scs.New()
	sessionManager.Store = store

	return &app{
		appConfig: &config.Config{
			AllowedHosts:  []string{"*"},
			SeriesButtons: []models.SeriesType{"last_5_minutes", "last_7_days"},
		},
		infoLog:          log.New(io.Discard, "", 0),
		errorLog:         log.New(io.Discard, "", 0),
		sessionManager:   sessionManager,
		inertiaManager:   inertia.New("", "app.gohtml", "v"),
		probeRepository:  &loadProbeRepository{},
		seriesRepository: &loadSeriesRepository{},
		logRepository:    &loadLogRepository{},
		alarmRepository:  &loadAlarmRepository{},
	}
}

func viewLoad(t *testing.T, routes http.Handler, query string, cookies ...*http.Cookie) *http.Response {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/load?probe=p1"+query, nil)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("X-Inertia-Version", "v")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rr := httptest.NewRecorder()

	routes.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET /load?probe=p1%s: status %d, body %s", query, rr.Code, rr.Body.String())
	}

	return rr.Result()
}

func TestChartViewWithoutSeriesChoiceCreatesNoSession(t *testing.T) {
	store := memstore.New()
	routes := newSeriesTestApp(store).routes()

	for _, query := range []string{"", "&type=last_5_minutes"} {
		resp := viewLoad(t, routes, query)

		if cookies := resp.Cookies(); len(cookies) > 0 {
			t.Errorf("GET /load?probe=p1%s without a cookie set %v, want no session cookie", query, cookies)
		}
	}

	sessions, err := store.All()
	if err != nil {
		t.Fatal(err)
	}

	if len(sessions) != 0 {
		t.Errorf("store has %d sessions, want 0", len(sessions))
	}
}

func TestChartViewRemembersChosenSeriesType(t *testing.T) {
	routes := newSeriesTestApp(memstore.New()).routes()

	cookies := viewLoad(t, routes, "&type=last_7_days").Cookies()
	if len(cookies) != 1 {
		t.Fatalf("choosing a series type set cookies %v, want one session cookie", cookies)
	}

	body, _ := io.ReadAll(viewLoad(t, routes, "", cookies...).Body)

	if want := `"seriesType":"last_7_days"`; !strings.Contains(string(body), want) {
		t.Errorf("follow-up view body %s, want %s", body, want)
	}
}
