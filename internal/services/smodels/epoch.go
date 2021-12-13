package smodels

import "time"

type EpochInfo struct {
	Epoch        uint64
	SlotsInEpoch uint64
	SPS          float64
	EndEpoch     time.Time
	Progress     uint8
}
