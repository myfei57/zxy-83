package plan

import (
	"errors"
	"fmt"
	"sync"

	"trainwash/internal/audit"
	"trainwash/internal/brush"
	"trainwash/internal/chem"
	"trainwash/internal/conv"
	"trainwash/internal/dry"
	"trainwash/internal/entry"
	"trainwash/internal/pos"
	"trainwash/internal/roof"
	"trainwash/internal/water"
)

var (
	ErrCycleBusy  = errors.New("plan: another wash is in progress")
	ErrWrongStage = errors.New("plan: wash is not at the expected stage")
	ErrBadConfig  = errors.New("plan: invalid wash configuration")
)

type Cycle struct {
	entry   *entry.Service
	pos     *pos.Tracker
	brush   *brush.Set
	roof    *roof.Service
	water   *water.System
	chem    *chem.Service
	dry     *dry.Service
	conv    *conv.Service
	audit   *audit.Recorder
	cfg     WashConfig
	mu      sync.Mutex
	stage   Stage
	trainID string
}

func NewCycle(entryService *entry.Service, tracker *pos.Tracker, brushSet *brush.Set, roofService *roof.Service, waterSystem *water.System, chemService *chem.Service, dryService *dry.Service, convService *conv.Service, recorder *audit.Recorder) *Cycle {
	return &Cycle{
		entry: entryService,
		pos:   tracker,
		brush: brushSet,
		roof:  roofService,
		water: waterSystem,
		chem:  chemService,
		dry:   dryService,
		conv:  convService,
		audit: recorder,
		cfg:   DefaultWashConfig(),
		stage: StageIdle,
	}
}

func (c *Cycle) StartWash(train entry.Train, position pos.Position) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stage != StageIdle {
		return ErrCycleBusy
	}
	if err := c.cfg.Validate(); err != nil {
		return err
	}
	if err := c.entry.Latch(); err != nil {
		return err
	}
	if _, err := c.entry.Recognize(train); err != nil {
		return err
	}
	if err := c.pos.Persist(position); err != nil {
		return err
	}
	if err := c.conv.Inbound(); err != nil {
		return err
	}
	if err := c.brush.LowerSide(); err != nil {
		return err
	}
	if err := c.water.Start(); err != nil {
		return err
	}
	c.dry.ResetCycle()
	if err := c.chem.Spray(c.cfg.SprayMS); err != nil {
		return err
	}
	c.stage = StageWash
	c.trainID = train.ID
	_, _ = c.audit.Add("wash_start", train.ID, "wash cycle started")
	return nil
}

func (c *Cycle) CompleteWash() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stage != StageWash {
		return ErrWrongStage
	}
	if !c.water.CanRinse() {
		return water.ErrNotFlowing
	}
	lengthMM := c.pos.Current().LengthMM
	scrubMS, err := c.brush.Scrub(lengthMM)
	if err != nil {
		return err
	}
	sweeps, err := c.brush.ScrubSweeps(lengthMM, c.cfg.ScrubGapMM)
	if err != nil {
		return err
	}
	if err := c.water.BeginRinse(); err != nil {
		return err
	}
	if _, err := c.brush.Rinse(c.cfg.RinseLevelMM); err != nil {
		return err
	}
	if err := c.water.RinseDone(); err != nil {
		return err
	}
	if err := c.dry.Start(); err != nil {
		return err
	}
	if err := c.dry.DryFor(c.cfg.DrySeconds); err != nil {
		return err
	}
	if err := c.dry.Stop(); err != nil {
		return err
	}
	if err := c.brush.Retract(); err != nil {
		return err
	}
	if err := c.water.Stop(); err != nil {
		return err
	}
	if err := c.conv.Outbound(); err != nil {
		return err
	}
	if err := c.conv.Stop(); err != nil {
		return err
	}
	if !c.water.CanDrain() {
		return water.ErrNotDrainable
	}
	if err := c.water.Drain(); err != nil {
		return err
	}
	c.stage = StageIdle
	trainID := c.trainID
	c.trainID = ""
	_, _ = c.audit.Add("wash_complete", trainID, fmt.Sprintf("wash cycle completed scrub_ms=%d sweeps=%d", scrubMS, sweeps))
	return nil
}

func (c *Cycle) Stage() Stage {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stage
}

func (c *Cycle) TrainID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.trainID
}
