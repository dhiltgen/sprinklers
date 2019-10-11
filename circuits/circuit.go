package circuits

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (c *Circuit) init() {
	c.pin = newPin(c.GPIONumber)

	// On init, set to output
	c.pin.Output()

	// And clear it
	//c.pin.Low()
	c.pin.High()
	c.State = false
}

func (c *Circuit) reportMetric() {
	// TODO - add log levels and make this verbose maybe?
	//log.Printf("Updating metric %d %s as %f\n", c.GPIONumber, c.Name, c.accumulation)
	accumulation.With(
		prometheus.Labels{
			"circuit": fmt.Sprintf("%d", c.GPIONumber),
			"name":    c.Name,
		}).Set(c.accumulation)
}

// Update updates the settings for the circuit given the new input
// For now, all we care about is State
func (c *Circuit) Update(new *Circuit) error {
	if new.State != c.State {
		c.State = new.State
		if c.State {
			log.Printf("Setting circuit \"%s\" ON", c.Name)
			c.started = time.Now()
			if new.TimeRemaining != "" {
				timeLeft, err := time.ParseDuration(new.TimeRemaining)
				if err != nil {
					return err
				}
				c.TimeRemaining = timeLeft.String()
				log.Printf("circuit will turn off in %s\n", timeLeft.String())
				ctx, cancel := context.WithCancel(context.Background())
				c.cancel = cancel
				timer := time.NewTimer(timeLeft)
				tick := time.NewTicker(AccumulationUpdateInterval)

				go func() {
					defer c.reportMetric()
					defer tick.Stop()
					for {
						select {
						case <-tick.C:
							prior := c.started
							c.started = time.Now()
							c.accumulation = c.accumulation + float64((c.started.Sub(prior)))/float64(time.Hour)*c.WaterConsumption
							c.reportMetric()
						case <-timer.C:
							log.Printf("Setting circuit \"%s\" OFF", c.Name)
							cancel()
							c.cancel = nil
							c.TimeRemaining = ""
							c.State = false
							c.pin.High()

							prior := c.started
							c.accumulation = c.accumulation + float64((c.started.Sub(prior)))/float64(time.Hour)*c.WaterConsumption
							return
						case <-ctx.Done():
							return
						}
					}
				}()
			}
			c.pin.Low()
		} else {
			log.Printf("Setting circuit \"%s\" OFF", c.Name)
			if c.cancel != nil {
				cancel := c.cancel
				c.TimeRemaining = ""
				c.cancel = nil
				cancel()
			}
			c.State = false
			c.pin.High()

			prior := c.started
			c.accumulation = c.accumulation + float64((c.started.Sub(prior)))/float64(time.Hour)*c.WaterConsumption
			c.reportMetric()
		}
	}
	return nil
}
