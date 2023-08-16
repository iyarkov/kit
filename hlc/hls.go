package hlc

import (
	"time"
)

const Jan012020 = 1577854800

type Stamp struct {
	Time  int32
	Count int32
	Node  int32
}

func (a Stamp) Before(b Stamp) bool {
	if a.Time == b.Time {
		if a.Count == b.Count {
			if a.Node == b.Node {
				return false
			} else if a.Node > b.Node {
				return false
			} else {
				return true
			}
		} else if a.Count > b.Count {
			return false
		} else {
			return true
		}
	} else if a.Time > b.Time {
		return false
	} else {
		return true
	}
}

func (a Stamp) After(b Stamp) bool {
	return b.Before(a)
}

type LogicalClock interface {
	Next() Stamp
	Update(Stamp) Stamp
}

func New(node int32) LogicalClock {
	wallClock := systemClock{}
	return &hybridLogicalClock{
		wallClock: wallClock,
		ts:        wallClock.now(),
		count:     0,
		node:      node,
	}
}

type StaticClock struct {
	Now Stamp
}

func (c *StaticClock) Next() Stamp {
	return c.Now
}

func (c *StaticClock) Update(Stamp) Stamp {
	return c.Now
}

type hybridLogicalClock struct {
	wallClock wallClock
	ts        int32
	count     int32
	node      int32
}

func (c *hybridLogicalClock) Next() Stamp {
	wallTs := c.wallClock.now()
	if wallTs > c.ts {
		c.ts = wallTs
		c.count = 0
	} else {
		c.count += 1
	}
	return c.now()
}

func (c *hybridLogicalClock) Update(ts2 Stamp) Stamp {
	wallTs := c.wallClock.now()
	if wallTs > c.ts && wallTs > ts2.Time {
		c.ts = wallTs
		c.count = 0
	} else if c.ts == ts2.Time {
		if c.count < ts2.Count {
			c.count = ts2.Count + 1
		} else {
			c.count = c.count + 1
		}
	} else if c.ts > ts2.Time {
		c.count = c.count + 1
	} else {
		c.ts = ts2.Time
		c.count = ts2.Count + 1
	}
	return c.now()
}

func (c *hybridLogicalClock) now() Stamp {
	return Stamp{
		Time:  c.ts,
		Count: c.count,
		Node:  c.node,
	}
}

type wallClock interface {
	now() int32
}

type systemClock struct{}

func (s systemClock) now() int32 {
	return int32(time.Now().Unix() - Jan012020)
}
