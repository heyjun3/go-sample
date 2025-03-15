package gosample_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	gosample "github.com/heyjun3/go-sample"
)

func userFactory(name string) *gosample.User {
	user, _ := gosample.NewUser(
		name,
		rand.Intn(21),
		rand.Intn(21),
		rand.Intn(21),
		rand.Intn(21),
		rand.Intn(21),
	)
	return user
}

func calculatorFactory(c int) (
	*gosample.CompatibilityCalculator,
	[]*gosample.User,
	[]*gosample.User,
) {
	managers := make([]*gosample.User, 0, c)
	members := make([]*gosample.User, 0, c)
	for i := range c {
		manager := userFactory("manager" + strconv.Itoa(i))
		managers = append(managers, manager)
		member := userFactory("member" + strconv.Itoa(i))
		members = append(members, member)
	}
	calculator := gosample.NewCompatibilityCalculator(
		gosample.DisengageAdult{},
		gosample.IsCritical{},
		gosample.IsFree{},
		gosample.IsAdaptive{},
	)
	return calculator, managers, members
}

func TestCompatibilityCalculator(t *testing.T) {
	calculator, managers, members := calculatorFactory(10000)
	start := time.Now()
	calculator.ExecMatching(managers, members)
	fmt.Println("done exec matching.", time.Since(start).Round(time.Millisecond))
}

func TestCompatibilityCalculatorConcurrency(t *testing.T) {
	calculator, managers, members := calculatorFactory(10000)
	start := time.Now()
	calculator.ExecMatchingConcurrency(managers, members)
	fmt.Println("done exec matching concurrency.", time.Since(start).Round(time.Millisecond))
}

func BenchmarkCompatibilityCalculator(b *testing.B) {
	calculator, managers, members := calculatorFactory(b.N)
	b.ResetTimer()
	calculator.ExecMatching(managers, members)
}

func BenchmarkCalcCompatibility(b *testing.B) {
	calculator, managers, members := calculatorFactory(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculator.CalcCompatibility(managers[i].Score, members[i].Score)
	}
}

func BenchmarkMostMatchingCompatibility(b *testing.B) {
	calculator, managers, members := calculatorFactory(b.N)
	b.ResetTimer()
	for _, member := range members {
		calculator.MostMatchingCompatibility(managers, member)
	}
}

func BenchmarkDisengageAdult(b *testing.B) {
	rule := gosample.DisengageAdult{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		first := userFactory("member" + strconv.Itoa(i))
		second := userFactory("manager" + strconv.Itoa(i))
		rule.Fn(first.Score, second.Score)
	}
}
func BenchmarkDisengageAdultV2(b *testing.B) {
	rule := gosample.NewDisengageAdultV2()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		first := userFactory("member" + strconv.Itoa(i))
		second := userFactory("manager" + strconv.Itoa(i))
		rule.Fn(first.Score, second.Score)
	}
}
