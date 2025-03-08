package gosample_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	gosample "github.com/heyjun3/go-sample"
)

func TestCompatibilityCalculator(t *testing.T) {
	calculator := gosample.NewCompatibilityCalculator(
		gosample.DisengageAdult,
		gosample.IsCritical,
		gosample.IsFree,
		gosample.IsAdaptive,
	)
	c := 70000
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

	calculator.ExecMatching(managers, members)
	fmt.Println("done")

	// result := calculator.ExecMatching(managers, members)
	// fmt.Println(result)
	// for k, v := range result {
	// 	fmt.Println("key", k, "value", v)
	// }
}
