package gosample

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
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
	a int8
	b int8
	c int8
	d int8
	e int8
}

func NewScore(a, b, c, d, e int) (*Score, error) {
	for _, v := range []int{a, b, c, d, e} {
		if 0 <= v && v <= 20 {
			continue
		}
		return nil, fmt.Errorf("invalid score range. value: %v", v)
	}
	return &Score{
		a: int8(a),
		b: int8(b),
		c: int8(c),
		d: int8(d),
		e: int8(e),
	}, nil
}

type Rule interface {
	Fn(*Score, *Score) int8
}

type DisengageAdult struct{}

func (d DisengageAdult) Fn(first, second *Score) int8 {
	diff := first.c - second.c
	if -16 < diff && diff < 16 {
		return 2
	}
	return 0
}

type DisengageAdultV2 struct {
	m map[string]int8
}

func NewDisengageAdultV2() *DisengageAdultV2 {
	m := make(map[string]int8)
	for i := range 20 {
		for j := range 20 {
			key := strconv.Itoa(i) + strconv.Itoa(j)
			if math.Abs(float64(i-j)) < 16 {
				m[key] = 2
			} else {
				m[key] = 0
			}
		}
	}
	return &DisengageAdultV2{
		m: m,
	}
}
func (d *DisengageAdultV2) Fn(first, second *Score) int8 {
	return d.m[strconv.Itoa(int(first.c))+strconv.Itoa(int(second.c))]
}

type IsCritical struct{}

func (i IsCritical) Fn(first, second *Score) int8 {
	if first.a >= 5 && second.a >= 5 {
		return 0
	}
	return 1
}

type IsFree struct{}

func (i IsFree) Fn(first, second *Score) int8 {
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

type IsAdaptive struct{}

func (i IsAdaptive) Fn(first, second *Score) int8 {
	max := int8(0)
	for _, v := range []int8{first.a, first.b, first.c, first.d, first.e} {
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

func (c *CompatibilityCalculator) calcCompatibility(first, second *Score, max int8) int8 {
	var result int8 = 0
	for _, rule := range c.rules {
		result += rule.Fn(first, second)
	}
	return result
}

func (c *CompatibilityCalculator) MostMatchingCompatibility(managers []*User, member *User) *User {
	max := int8(0)
	matching := make([]*User, 0)
	for _, manager := range managers {
		result := c.calcCompatibility(manager.score, member.score, 0)
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
	for _, member := range members {
		match := c.MostMatchingCompatibility(managers, member)
		if v, ok := m[match.ID]; ok {
			v = append(v, member.ID)
			m[match.ID] = v
		} else {
			m[match.ID] = []string{member.ID}
		}
	}
	return m
}

func (c *CompatibilityCalculator) ExecMatchingConcurrency(managers, members []*User) map[string][]string {
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
