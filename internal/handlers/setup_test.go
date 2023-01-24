package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dominic-wassef/gopher-runner/internal/channeldata"
	"github.com/dominic-wassef/gopher-runner/internal/config"
	"github.com/dominic-wassef/gopher-runner/internal/driver"
	"github.com/dominic-wassef/gopher-runner/internal/helpers"
	"github.com/dominic-wassef/gopher-runner/internal/repository/dbrepo"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
)

var testSession *scs.SessionManager

func TestMain(m *testing.M) {
	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false

	mailQueue := make(chan channeldata.MailJob, 5)

	a := config.AppConfig{
		DB:           &driver.DB{},
		Session:      testSession,
		InProduction: false,
		Domain:       "localhost",
		MailQueue:    mailQueue,
	}

	app = &a

	preferenceMap := make(map[string]string)
	app.PreferenceMap = preferenceMap

	wsClient := pusher.Client{
		AppID:  "1",
		Secret: "123abc",
		Key:    "abc123",
		Secure: false,
		Host:   "localhost:4001",
	}

	app.WsClient = wsClient

	monitorMap := make(map[int]cron.EntryID)
	app.MonitorMap = monitorMap

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	app.Scheduler = scheduler

	repo := NewTestHandlers(app)
	NewHandlers(repo, app)

	helpers.NewHelpers(app)

	helpers.SetViews("./../../views")

	os.Exit(m.Run())
}

func getCtx(req *http.Request) context.Context {
	ctx, err := testSession.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func NewTestHandlers(a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}
