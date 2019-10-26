package hal

import "time"

//
// TimeEntry
//

type TimeEntry struct {
	ResourceObject
}

func NewTimeEntry() *TimeEntry {
	return &TimeEntry{
		ResourceObject{
			Type: "TimeEntry",
		},
	}
}

func (res *TimeEntry) Id() int {
	return res.GetInt("id")
}

func (res *TimeEntry) Comment() *Formattable {
	f, err := DecodeFormattable(res.GetField("comment"))
	if err != nil {
		return nil
	}
	return f
}

func (res *TimeEntry) SpentOn() *time.Time {
	if dt, err := res.GetDateTime("spentOn"); err == nil {
		return &dt
	}
	return nil
}

func (res *TimeEntry) Hours() *time.Duration {
	if dt, err := res.GetDuration("hours"); err == nil {
		return &dt
	}
	return nil
}

// Register Factories
func init() {
	resourceTypes["TimeEntry"] = func() Resource {
		return NewTimeEntry()
	}
}
