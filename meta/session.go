package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	Span

	MarkCode      string
	Operator      string
	Agency        string
	Interval      time.Duration
	ElevationMask float64
	HeaderComment string
}

func (s Session) Less(session Session) bool {
	switch {
	case s.MarkCode < session.MarkCode:
		return true
	case s.MarkCode > session.MarkCode:
		return false
	case s.Interval < session.Interval:
		return true
	case s.Interval > session.Interval:
		return false
	case s.Start.Before(session.Start):
		return true
	case s.Start.After(session.Start):
		return false
	default:
		return false
	}
}

type SessionList []Session

func (s SessionList) Len() int           { return len(s) }
func (s SessionList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SessionList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s SessionList) encode() [][]string {
	data := [][]string{{
		"Mark Code",
		"Operator",
		"Agency",
		"Interval",
		"Elevation Mask",
		"Header Comment",
		"Start Date",
		"End Date",
	}}
	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.MarkCode),
			strings.TrimSpace(v.Operator),
			strings.TrimSpace(v.Agency),
			strings.TrimSpace(v.Interval.String()),
			strings.TrimSpace(strconv.FormatFloat(v.ElevationMask, 'g', -1, 64)),
			strings.TrimSpace(v.HeaderComment),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (c *SessionList) decode(data [][]string) error {
	var sessions []Session
	if len(data) > 1 {
		for _, v := range data[1:] {
			if len(v) != 8 {
				return fmt.Errorf("incorrect number of installed session fields")
			}
			var err error

			var interval time.Duration
			if interval, err = time.ParseDuration(v[3]); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, v[6]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, v[7]); err != nil {
				return err
			}

			var mask float64
			if mask, err = strconv.ParseFloat(v[4], 64); err != nil {
				return err
			}

			sessions = append(sessions, Session{
				MarkCode:      strings.TrimSpace(v[0]),
				Operator:      strings.TrimSpace(v[1]),
				Agency:        strings.TrimSpace(v[2]),
				Interval:      interval,
				ElevationMask: mask,
				HeaderComment: strings.TrimSpace(v[5]),
				Span: Span{
					Start: start,
					End:   end,
				},
			})
		}

		*c = SessionList(sessions)
	}
	return nil
}

func LoadSessions(path string) ([]Session, error) {
	var s []Session

	if err := LoadList(path, (*SessionList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(SessionList(s))

	return s, nil
}
