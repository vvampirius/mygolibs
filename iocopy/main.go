package iocopy

import (
	"context"
	"io"
	"time"
)

type Report struct {
	Copied int64
	Error  error
	Spent  time.Duration
}

func (report *Report) Speed(duration time.Duration) float64 {
	if report.Copied == 0 || report.Spent.Nanoseconds() == 0 {
		return 0
	}
	x := (float64(report.Copied) / float64(report.Spent.Nanoseconds())) * float64(duration.Nanoseconds())
	return x
}

func Copy(dst io.Writer, src io.Reader, blockSize int64, reportInterval time.Duration, reportChan chan<- Report,
	ctx context.Context) {
	active := true
	inProgress := false
	copyChan := make(chan struct {
		int64
		error
	}, 0)
	tickChan := time.NewTicker(reportInterval)
	copied := int64(0)
	intervalStartedAt := time.Now()
	for active {
		if !inProgress {
			inProgress = true
			go func() {
				written, err := io.CopyN(dst, src, blockSize)
				inProgress = false
				copyChan <- struct {
					int64
					error
				}{int64: written, error: err}
			}()
		}
		select {
		case cn := <-copyChan:
			copied = copied + cn.int64
			if cn.error != nil {
				active = false
				reportChan <- Report{
					Copied: copied,
					Error:  cn.error,
					Spent:  time.Now().Sub(intervalStartedAt),
				}
			}
		case <-tickChan.C:
			reportChan <- Report{
				Copied: copied,
				Error:  nil,
				Spent:  time.Now().Sub(intervalStartedAt),
			}
			copied = 0
			intervalStartedAt = time.Now()
		case <-ctx.Done():
			active = false
			reportChan <- Report{
				Copied: copied,
				Error:  nil,
				Spent:  time.Now().Sub(intervalStartedAt),
			}
		}
	}
}
