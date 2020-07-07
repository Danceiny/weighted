package weighted

import (
	"strconv"
	"testing"
	"time"
)

func TestSW_Next(t *testing.T) {
	w := &SW{}
	w.Add("server1", 5)
	w.Add("server2", 2)
	w.Add("server3", 3)

	results := make(map[string]int)

	for i := 0; i < 100; i++ {
		s := w.Next().(string)
		results[s]++
	}

	if results["server1"] != 50 || results["server2"] != 20 || results["server3"] != 30 {
		t.Error("the algorithm is wrong")
	}

	w.Reset()
	results = make(map[string]int)

	for i := 0; i < 100; i++ {
		s := w.Next().(string)
		results[s]++
	}

	if results["server1"] != 50 || results["server2"] != 20 || results["server3"] != 30 {
		t.Error("the algorithm is wrong")
	}

	w.RemoveAll()
	w.Add("server1", 7)
	w.Add("server2", 9)
	w.Add("server3", 13)

	results = make(map[string]int)

	for i := 0; i < 29000; i++ {
		s := w.Next().(string)
		results[s]++
	}

	if results["server1"] != 7000 || results["server2"] != 9000 || results["server3"] != 13000 {
		t.Error("the algorithm is wrong")
	}
}

func TestSW_NextWithCallback(t *testing.T) {
	const forcount = 999
	// positive: weight more, cost less
	var doFuncPositive = func(input string) int {
		index, _ := strconv.Atoi(input[len(input)-1:])
		index = 10 - index
		// time.Sleep(time.Duration(index) * time.Microsecond)
		return index
	}

	// negative: weight more, cost more
	var doFuncNegative = func(input string) int {
		index, _ := strconv.Atoi(input[len(input)-1:])
		time.Sleep(time.Duration(index) * time.Microsecond)
		return index
	}

	// with no callback for dynamic adjustment
	w := &SW{}
	w.Add("server5", 100)
	w.Add("server2", 100)
	w.Add("server3", 100)
	results := make(map[string]int)
	for i := 0; i < forcount; i++ {
		so, _ := w.NextWithCallback()
		s, _ := so.(string)
		results[s]++

		_ = doFuncPositive(s)
	}
	t.Log("non-callback", w.All(), results)
	if !(results["server5"] == results["server3"] && results["server3"] == results["server2"]) {
		t.Error("the algorithm is wrong")
	}

	// with positive callback
	w2 := &SW{}
	w2.Add("server5", 100)
	w2.Add("server2", 100)
	w2.Add("server3", 100)
	results2 := make(map[string]int)
	for i := 0; i < forcount; i++ {
		so, f := w2.NextWithCallback()
		s, _ := so.(string)
		results2[s]++

		callback := doFuncPositive(s)
		f(1000 / callback)
	}
	t.Log("positive-callback", w2.All(), results2)
	if !(results2["server5"] > results2["server3"] && results2["server3"] > results2["server2"]) {
		t.Error("the algorithm is wrong")
	}

	// with negative callback
	w3 := &SW{}
	w3.Add("server5", 100)
	w3.Add("server2", 100)
	w3.Add("server3", 100)
	results3 := make(map[string]int)
	for i := 0; i < forcount; i++ {
		so, f := w3.NextWithCallback()
		s, _ := so.(string)
		results3[s]++

		callback := doFuncNegative(s)
		f(1000 / callback)
	}
	t.Log("negative-callback", w3.All(), results3)
	if !(results3["server5"] < results3["server3"] && results3["server3"] < results3["server2"]) {
		t.Error("the algorithm is wrong")
	}
}
