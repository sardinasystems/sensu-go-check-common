package nagios

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Threshold represents range for value
//
// http://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT
type Threshold struct {
	Start    float64
	End      float64
	Inverted bool
}

// ParseThreshold parses Nagios Threshold format
func ParseThreshold(s string) (Threshold, error) {
	inverted := false
	if strings.HasPrefix(s, "@") {
		inverted = true
		s = s[1:]
	}

	parse := func(ss string) (float64, error) {
		switch ss {
		case "~":
			return math.Inf(-1), nil
		case "":
			return math.Inf(1), nil
		default:
			return strconv.ParseFloat(ss, 64)
		}
	}

	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)

		start, err := parse(parts[0])
		if err != nil {
			return Threshold{}, fmt.Errorf("Threshold parse error: %w", err)
		}

		end, err := parse(parts[1])
		if err != nil {
			return Threshold{}, fmt.Errorf("Threshold parse error: %w", err)
		}

		return Threshold{
			Start:    start,
			End:      end,
			Inverted: inverted,
		}, nil
	}

	f, err := parse(s)
	if err != nil {
		return Threshold{}, fmt.Errorf("Threshold parse error: %w", err)
	}

	return Threshold{
		Start:    0.0,
		End:      f,
		Inverted: inverted,
	}, nil
}

// Check returns true if value should be alerted
func (r *Threshold) Check(v float64) bool {
	if r.Inverted {
		return v >= r.Start && v <= r.End
	}

	return v < r.Start || v > r.End
}

// -*- Implement cobra.Value -*-

func (r *Threshold) String() string {

	f := func(v float64) string {
		switch {
		case math.IsInf(v, 1):
			return ""
		case math.IsInf(v, -1):
			return "~"
		default:
			return fmt.Sprintf("%v", v)
		}
	}

	inv := ""
	if r.Inverted {
		inv = "@"
	}

	return fmt.Sprintf("%s%s:%s", inv, f(r.Start), f(r.End))
}

func (r *Threshold) Set(s string) error {
	th, err := ParseThreshold(s)
	if err != nil {
		return err
	}

	*r = th
	return nil
}

func (r *Threshold) Type() string {
	return "threshold"
}
