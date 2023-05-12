package main

import (
	"rng/api"
	"rng/engine"
)

/** Chan for histogram */
func randStream(e *engine.Engine, h uint64, l uint64) chan float64 {
	c := make(chan float64)
	go func() {
		for {
			c <- float64(e.GetRand(h, l))
		}
	}()
	return c
}

func main() {
	rngEngine := engine.NewEngine()
	rngEngine.Run()

	/** Run HTTP API */
	rngApi := api.NewAPI(rngEngine)
	rngApi.Run()

	/** Write binary directly to stdout for statistical and analytical purposes */
	//for {
	//	binary.Write(os.Stdout, binary.BigEndian, rngEngine.RandomBytes(2048))
	//}

	/** Histogram for visual check */
	//h := 100
	//l := 0
	//n := 200000
	//hist := hist.NewHist(nil, "Example histogram", "fixed", h, false)
	//c := randStream(rngEngine, uint64(h), uint64(l))
	//i := 0
	//for {
	//	// add data point to hsitogram
	//	hist.Update(<-c)
	//	if i%n == 0 {
	//		// draw histogram
	//		fmt.Println(hist.DrawSimple())
	//		time.Sleep(time.Second)
	//	}
	//	i++
	//}
}
