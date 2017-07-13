// Copyright 2015
// Author: huangjunwei@youmi.net

package blog4go

import (
	log "github.com/ccpaging/log4go"
	"sync"
	"testing"
	"time"
)

type T struct {
	A int
	B string
}

type timeFormatCacheType struct {
	now    time.Time
	format string
}

var filename string = "output_log4go.log"

func BenchmarkLogrusSingleGoroutine(b *testing.B) {
	log.Close()
	log.AddFilter("file", log.FINE, log.NewFileLogWriter(filename, false))
	defer log.Close()

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "INFO", "eddie", 18, 3.1415, true, t)
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", "eddie", 18, 3.1415, true, t)
	}
}

func BenchmarkLogrusWithTimecacheSingleGoroutine(b *testing.B) {
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Second)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				now := time.Now()
				timeCache.now = now
				timeCache.format = now.Format("[2006-01-02 15:04:05]")
			}
		}
	}()

	log.Close()
	log.AddFilter("file", log.FINE, log.NewFileLogWriter(filename, false))
	defer log.Close()

	t := T{123, "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t)
		log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t)
	}
}

func BenchmarkLogrusWithTimecacheMultiGoroutine(b *testing.B) {
	now := time.Now()

	timeCache := timeFormatCacheType{
		now:    now,
		format: now.Format("2006-01-02 15:04:05"),
	}

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Second)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				now := time.Now()
				timeCache.now = now
				timeCache.format = now.Format("[2006-01-02 15:04:05]")
			}
		}
	}()

	log.Close()
	log.AddFilter("file", log.FINE, log.NewFileLogWriter(filename, false))
	defer log.Close()

	t := T{123, "test"}

	var wg sync.WaitGroup
	var beginWg sync.WaitGroup

	f := func() {
		defer wg.Done()
		beginWg.Wait()
		for i := 0; i < b.N; i++ {
			log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "INFO", "eddie", 18, 3.1415, true, t)
			log.Printf("%s [%s] haha %s. en\\en, always %d and %f, %t, %+v\n", timeCache.format, "ERROR", "eddie", 18, 3.1415, true, t)
		}
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		beginWg.Add(1)
	}

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		go f()
		beginWg.Done()
	}

	wg.Wait()
}
