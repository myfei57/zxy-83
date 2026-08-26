package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"trainwash/internal/audit"
	"trainwash/internal/brush"
	"trainwash/internal/chem"
	"trainwash/internal/conv"
	"trainwash/internal/dry"
	"trainwash/internal/entry"
	"trainwash/internal/ns"
	"trainwash/internal/plan"
	"trainwash/internal/pos"
	"trainwash/internal/roof"
	"trainwash/internal/store"
	"trainwash/internal/water"
)

type System struct {
	DataDir string
	Store   *store.FileStore
	Pos     *pos.Tracker
	Entry   *entry.Service
	Brush   *brush.Set
	Water   *water.System
	Chem    *chem.Service
	Roof    *roof.Service
	Dry     *dry.Service
	Conv    *conv.Service
	Cycle   *plan.Cycle
	Audit   *audit.Recorder
	Layout  ns.StationLayout
	Limits  ns.Limits
	Now     func() time.Time
}

func NewSystem(dataDir string, st *store.FileStore, tracker *pos.Tracker, entryService *entry.Service, brushSet *brush.Set, waterSystem *water.System, chemService *chem.Service, roofService *roof.Service, dryService *dry.Service, convService *conv.Service, cycle *plan.Cycle, recorder *audit.Recorder, layout ns.StationLayout, limits ns.Limits, now func() time.Time) *System {
	return &System{
		DataDir: dataDir,
		Store:   st,
		Pos:     tracker,
		Entry:   entryService,
		Brush:   brushSet,
		Water:   waterSystem,
		Chem:    chemService,
		Roof:    roofService,
		Dry:     dryService,
		Conv:    convService,
		Cycle:   cycle,
		Audit:   recorder,
		Layout:  layout,
		Limits:  limits,
		Now:     now,
	}
}

func NewRouter(s *System) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Get("/api/state", s.handleState)
	r.Get("/api/audit", s.handleAudit)
	r.Post("/api/brush/lower", s.handleBrushLower)
	r.Post("/api/brush/retract", s.handleBrushRetract)
	r.Post("/api/brush/raise", s.handleBrushRaise)
	r.Post("/api/roof/lower", s.handleRoofLower)
	r.Post("/api/roof/raise", s.handleRoofRaise)
	r.Post("/api/roof/brush", s.handleRoofBrush)
	r.Post("/api/chem/spray", s.handleChemSpray)
	r.Post("/api/chem/alarm", s.handleChemAlarm)
	r.Post("/api/chem/alarm/clear", s.handleChemAlarmClear)
	r.Post("/api/water/start", s.handleWaterStart)
	r.Post("/api/water/stop", s.handleWaterStop)
	r.Post("/api/water/rinse", s.handleWaterRinse)
	r.Post("/api/water/recalibrate", s.handleWaterRecalibrate)
	r.Post("/api/water/drain", s.handleWaterDrain)
	r.Post("/api/dry/start", s.handleDryStart)
	r.Post("/api/dry/stop", s.handleDryStop)
	r.Post("/api/pos/persist", s.handlePosPersist)
	r.Post("/api/pos/calibrate", s.handlePosCalibrate)
	r.Post("/api/pos/head", s.handlePosHead)
	r.Post("/api/entry/type", s.handleEntryType)
	r.Post("/api/conv/inbound", s.handleConvInbound)
	r.Post("/api/conv/move", s.handleConvMove)
	r.Post("/api/conv/stop", s.handleConvStop)
	r.Post("/api/plan/wash", s.handlePlanWash)
	r.Post("/api/plan/complete", s.handlePlanComplete)
	r.Post("/api/plan/stop", s.handlePlanStop)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))
	r.Get("/", s.handleIndex)
	r.Get("/wash", s.handlePage("wash.html"))
	r.Get("/brushes", s.handlePage("brushes.html"))
	r.Get("/water", s.handlePage("water.html"))
	r.Get("/alarms", s.handlePage("alarms.html"))
	return r
}
