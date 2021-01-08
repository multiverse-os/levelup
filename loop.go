package levelup

import "time"

func (self Database) WriteLoop() {
	var (
		updating bool
		updated  time.Time
		err      error
	)

	for {
		select {
		case errc := <-self.quit:
			// Chain indexer terminating, report no failure and abort
			errc <- nil
			return

		case <-self.update:
			// Section headers completed (or rolled back), update the index
			self.Access.Lock()
			if time.Since(updated) > 8*time.Second {
				updated = time.Now()
			}
			self.Access.Unlock()
			if err != nil {
				select {
				case <-self.ctx.Done():
					<-self.quit <- nil
					return
				default:
				}
			}
			self.Access.Lock()

			// If processing succeeded and no reorgs occurred, mark the section completed
			if err == nil {
				if updating {
					updating = false
				}
			} else {
				// If processing failed, don't retry until further notification
			}
			// If there are still further sections to process, reschedule
			time.AfterFunc(self.throttle, func() {
				select {
				case self.update <- struct{}{}:
				default:
				}
			})
		}
		self.Access.Unlock()
	}
}
