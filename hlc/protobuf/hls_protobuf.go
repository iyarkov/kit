package protobuf

import "github.com/iyarkov/kit/hlc"

func Compare(a *Stamp, b *Stamp) int {
	if a.Time == b.Time {
		if a.Count == b.Count {
			if a.Node == b.Node {
				return 0
			} else if a.Node > b.Node {
				return 1
			} else {
				return -1
			}
		} else if a.Count > b.Count {
			return 1
		} else {
			return -1
		}
	} else if a.Time > b.Time {
		return 1
	} else {
		return -1
	}
}

type LogicalClock interface {
	Next() *Stamp
	Update(in *Stamp) *Stamp
}

func New(delegate hlc.LogicalClock) LogicalClock {
	return &delegatingLogicalClock{
		delegate: delegate,
	}
}

type delegatingLogicalClock struct {
	delegate hlc.LogicalClock
}

func (d *delegatingLogicalClock) Next() *Stamp {
	return toProto(d.delegate.Next())
}

func (d *delegatingLogicalClock) Update(in *Stamp) *Stamp {
	return toProto(d.delegate.Update(fromProto(in)))
}

func toProto(in hlc.Stamp) *Stamp {
	return &Stamp{
		Time:  in.Time,
		Count: in.Count,
		Node:  in.Node,
	}
}

func fromProto(in *Stamp) hlc.Stamp {
	return hlc.Stamp{
		Time:  in.Time,
		Count: in.Count,
		Node:  in.Node,
	}
}
