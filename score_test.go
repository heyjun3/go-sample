package gosample_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	gosample "github.com/heyjun3/go-sample"
)

func calculatorFactory(t *testing.T, c int) (
	*gosample.CompatibilityCalculator,
	[]*gosample.User,
	[]*gosample.User,
) {
	managers := make([]*gosample.User, 0, c)
	members := make([]*gosample.User, 0, c)
	for i := range c {
		manager, err := gosample.NewUser(
			"manager"+strconv.Itoa(i),
			rand.Intn(21),
			rand.Intn(21),
			rand.Intn(21),
			rand.Intn(21),
			rand.Intn(21),
		)
		if err != nil {
			t.Error(err)
		}
		managers = append(managers, manager)
		member, err := gosample.NewUser(
			"member"+strconv.Itoa(i),
			rand.Intn(21),
			rand.Intn(21),
			rand.Intn(21),
			rand.Intn(21),
			rand.Intn(21),
		)
		if err != nil {
			t.Error(err)
		}
		members = append(members, member)
	}
	calculator := gosample.NewCompatibilityCalculator(
		// gosample.DisengageAdult{},
		gosample.IsCritical{},
		gosample.IsFree{},
		gosample.IsAdaptive{},
	)
	return calculator, managers, members
}

func TestCompatibilityCalculator(t *testing.T) {
	calculator, managers, members := calculatorFactory(t, 50000)
	start := time.Now()
	calculator.ExecMatching(managers, members)
	fmt.Println("done exec matching.", time.Since(start).Round(time.Millisecond))
}

func TestCompatibilityCalculatorConcurrency(t *testing.T) {
	calculator, managers, members := calculatorFactory(t, 50000)
	start := time.Now()
	calculator.ExecMatchingConcurrency(managers, members)
	fmt.Println("done exec matching concurrency.", time.Since(start).Round(time.Millisecond))
}
