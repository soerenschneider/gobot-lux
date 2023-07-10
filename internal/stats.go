package internal

import (
	"errors"
	"log"
	"sync"
	"time"
)

const maxEvents = 2048

type Measurement struct {
	timestamp time.Time
	value     float32
}

type SensorStats struct {
	stats []Measurement
	m     sync.RWMutex
}

func NewSensorStats() *SensorStats {
	return &SensorStats{
		stats: make([]Measurement, 0),
	}
}

func (s *SensorStats) NewEvent(measurement float32) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.stats) < maxEvents {
		measurement := Measurement{
			timestamp: time.Now(),
			value:     measurement,
		}
		s.stats = append(s.stats, measurement)
	} else {
		log.Println("Not adding further stats to stats, slice is full")
	}
}

func (s *SensorStats) GetStatsSliceSize() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return len(s.stats)
}

type IntervalStatistics struct {
	Min   float32 `json:"min"`
	Max   float32 `json:"max"`
	Delta float32 `json:"delta"`
	Avg   float32 `json:"avg"`
}

func (s *SensorStats) GetIntervalStats(window time.Duration) (IntervalStatistics, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	idx := s.getIndexOfStatsNewerThan(time.Now().Add(-window))
	return evalInterval(s.stats, idx)
}

func evalInterval(array []Measurement, fromIndex int) (IntervalStatistics, error) {
	if len(array) == 0 {
		return IntervalStatistics{
			Min: -1,
			Max: -1,
			Avg: -1,
		}, errors.New("no measurements available")
	}
	ret := IntervalStatistics{
		Min: array[fromIndex].value,
		Max: array[fromIndex].value,
	}

	var sum = array[fromIndex].value
	for i := fromIndex + 1; i < len(array); i++ {
		val := array[i].value
		if val < ret.Min {
			ret.Min = val
		}
		if val > ret.Max {
			ret.Max = val
		}
		sum += val
	}

	ret.Avg = sum / float32(len(array)-fromIndex)
	ret.Delta = ret.Max - ret.Min
	return ret, nil
}

func (s *SensorStats) getIndexOfStatsNewerThan(timestamp time.Time) int {
	s.m.RLock()
	defer s.m.RUnlock()
	for index, event := range s.stats {
		if event.timestamp.After(timestamp) {
			return index
		}
	}

	return len(s.stats)
}

func (s *SensorStats) PurgeStatsBefore(timestamp time.Time) {
	s.m.Lock()
	defer s.m.Unlock()
	marker := s.getIndexOfStatsNewerThan(timestamp)
	s.stats = s.stats[marker:]
}
