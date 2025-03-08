package gosample

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	// "sync"
)

type User struct {
	ID    string
	score *Score
}

func NewUser(id string, a, b, c, d, e int) (*User, error) {
	score, err := NewScore(a, b, c, d, e)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    id,
		score: score,
	}, nil
}

type Score struct {
	a uint8
	b uint8
	c uint8
	d uint8
	e uint8
}

func NewScore(a, b, c, d, e int) (*Score, error) {
	for _, v := range []int{a, b, c, d, e} {
		if 0 <= v && v <= 20 {
			continue
		}
		return nil, fmt.Errorf("invalid score range. value: %v", v)
	}
	return &Score{
		a: uint8(a),
		b: uint8(b),
		c: uint8(c),
		d: uint8(d),
		e: uint8(e),
	}, nil
}

type Rule func(*Score, *Score) uint8

func DisengageAdult(first, second *Score) uint8 {
	if math.Abs(float64(first.c)-float64(second.c)) < 16 {
		return 2
	}
	return 0
}
func IsCritical(first, second *Score) uint8 {
	if first.a >= 5 && second.a >= 5 {
		return 0
	}
	return 1
}
func IsFree(first, second *Score) uint8 {
	if math.Abs(float64(first.d)-float64(first.e)) < 2 {
		return 4
	}
	if math.Abs(float64(second.d)-float64(second.e)) < 2 {
		return 4
	}
	if first.d > first.e && second.d > second.e {
		return 2
	}
	if first.d < first.e && second.d < second.e {
		return 2
	}
	return 0
}
func IsAdaptive(first, second *Score) uint8 {
	max := uint8(0)
	for _, v := range []uint8{first.a, first.b, first.c, first.d, first.e} {
		if max < v {
			max = v
		}
	}
	if math.Abs(float64(max)-float64(first.b)) <= 2 {
		return 1
	}
	return 0
}

type CompatibilityCalculator struct {
	rules []Rule
}

func NewCompatibilityCalculator(rules ...Rule) *CompatibilityCalculator {
	return &CompatibilityCalculator{
		rules: rules,
	}
}

func (c *CompatibilityCalculator) calcCompatibility(first, second *Score) uint8 {
	var result uint8 = 0
	for _, rule := range c.rules {
		result += rule(first, second)
	}
	return result
}

func (c *CompatibilityCalculator) MostMatchingCompatibility(managers []*User, member *User) *User {
	max := uint8(0)
	matching := make([]*User, 0, len(managers))
	for _, manager := range managers {
		result := c.calcCompatibility(manager.score, member.score)
		if max < result {
			max = result
			matching = []*User{manager}
		} else if max == result {
			matching = append(matching, manager)
		}
	}
	return matching[rand.Int63n(int64(len(matching)))]
}

func (c *CompatibilityCalculator) ExecMatching(managers, members []*User) map[string][]string {
	m := make(map[string][]string)
	ch := make(chan string, len(members))
	go func() {
		defer close(ch)
		wg := &sync.WaitGroup{}
		for _, member := range members {
			wg.Add(1)
			go func() {
				match := c.MostMatchingCompatibility(managers, member)
				ch <- match.ID + ";" + member.ID
				wg.Done()
			}()
		}
		wg.Wait()
	}()

	for v := range ch {
		str := strings.Split(v, ";")
		if v, ok := m[str[0]]; ok {
			v = append(v, str[1])
			m[str[0]] = v
		} else {
			m[str[0]] = []string{str[1]}
		}
	}
	return m
}
