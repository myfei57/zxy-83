package main

import (
	"time"

	"trainwash/internal/audit"
	"trainwash/internal/brush"
	"trainwash/internal/chem"
	"trainwash/internal/console"
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

func BuildSystem(cfg Config) (*console.System, error) {
	st := store.New(cfg.DataDir)
	layout := ns.NewStationLayout()
	limits := ns.DefaultLimits()
	tracker := pos.NewTracker(st)
	entryService := entry.NewService(st)
	waterSystem := water.NewSystem(st, limits)
	brushSet := brush.NewSet(st, tracker, waterSystem)
	entryService.AttachPublisher(brushSet)
	chemService := chem.NewService(st, waterSystem, limits)
	roofService := roof.NewService(tracker)
	dryService := dry.NewService(waterSystem)
	convService := conv.NewService(st, waterSystem, layout, limits)
	recorder := audit.NewRecorder(st)
	cycle := plan.NewCycle(entryService, tracker, brushSet, roofService, waterSystem, chemService, dryService, convService, recorder)
	brushSet.AttachReleaser(entryService)
	waterSystem.AttachDrainGate(chemService)
	_ = cycle.Recover()
	system := console.NewSystem(cfg.DataDir, st, tracker, entryService, brushSet, waterSystem, chemService, roofService, dryService, convService, cycle, recorder, layout, limits, time.Now)
	return system, nil
}
