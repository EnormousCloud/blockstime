package timeslice

import "time"

func BlockBefore(blocks []int64, src *time.Time) *int64 {
	if src == nil {
		return nil
	}
	for blockNum, blockTime := range blocks {
		if blockTime > 0 && src.Unix() > blockTime {
			res := int64(blockNum - 1)
			return &res
		}
	}
	return nil
}

func BlockAfter(blocks []int64, src *time.Time) *int64 {
	if src == nil {
		return nil
	}
	var res int64
	for blockNum, blockTime := range blocks {
		if blockTime > 0 && src.Unix() < blockTime {
			res = int64(blockNum + 1)
		}
	}
	return &res
}

func TimeBefore(blocks []int64, src *int64) *time.Time {
	if src == nil {
		return nil
	}

	return nil
}

func TimeAfter(blocks []int64, src *int64) *time.Time {
	if src == nil {
		return nil
	}

	return nil
}
