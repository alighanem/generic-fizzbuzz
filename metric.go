package main

import "sync"

type Counter struct {
	sync.RWMutex
	Metrics map[string]*Metric
}

type Metric struct {
	Converter *ValueConverter `json:"converter"`
	Hits      int             `json:"hits"`
}

func (c *Counter) Inc(converter *ValueConverter) {
	c.Lock()
	m, ok := c.Metrics[converter.Key()]
	if !ok {
		m = &Metric{
			Converter: converter,
		}
		c.Metrics[converter.Key()] = m
	}
	m.Hits++
	c.Unlock()
}

func (c *Counter) GetMaxHits() *Metric {
	c.RLock()
	defer c.RUnlock()

	if len(c.Metrics) == 0 {
		return nil
	}

	var metric *Metric
	var max int
	for _, m := range c.Metrics {
		if m.Hits < max {
			continue
		}
		max = m.Hits
		metric = m
	}

	return metric
}
