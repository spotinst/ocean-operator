// Copyright 2021 NetApp, Inc. All Rights Reserved.

package cli

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

// StartProfiling initializes profiling for the current process.
func StartProfiling(profile, output string) error {
	switch profile {
	case "none":
		return nil
	case "cpu":
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(f)
	// Block and mutex profiles need a call to Set{Block,Mutex}ProfileRate to
	// output anything. We choose to sample all events.
	case "block":
		runtime.SetBlockProfileRate(1)
		return nil
	case "mutex":
		runtime.SetMutexProfileFraction(1)
		return nil
	default:
		if profile != "" {
			// Check the profile name is valid.
			if prof := pprof.Lookup(profile); prof == nil {
				return fmt.Errorf("unknown profile %q", profile)
			}
		}
	}
	return nil
}

// StopProfiling stops the profiler.
func StopProfiling(profile, output string) error {
	switch profile {
	case "none":
		return nil
	case "cpu":
		pprof.StopCPUProfile()
	case "heap":
		runtime.GC()
		fallthrough
	default:
		if profile != "" {
			prof := pprof.Lookup(profile)
			if prof == nil {
				return nil
			}
			f, err := os.Create(output)
			if err != nil {
				return err
			}
			_ = prof.WriteTo(f, 0)
		}
	}
	return nil
}
