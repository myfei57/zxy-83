package plan

func (c *Cycle) Recover() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.pos.Restore(); err != nil {
		return err
	}
	if err := c.brush.RefreshMap(); err != nil {
		return err
	}
	c.stage = StageIdle
	return nil
}

func (c *Cycle) EmergencyStop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	trainID := c.trainID
	if err := c.water.Stop(); err != nil {
		return err
	}
	if err := c.conv.Stop(); err != nil {
		return err
	}
	_ = c.pos.ClearHead()
	_ = c.pos.ClearPosition()
	c.stage = StageIdle
	c.trainID = ""
	_, _ = c.audit.Add("emergency_stop", trainID, "emergency stop completed")
	return nil
}
