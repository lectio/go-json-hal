package hal

import "time"

//
// TimeEntry Activity
//

type TimeEntriesActivity struct {
	ResourceObject
}

func NewTimeEntriesActivity() *TimeEntriesActivity {
	return &TimeEntriesActivity{
		ResourceObject{
			Type: "TimeEntriesActivity",
		},
	}
}

func (res *TimeEntriesActivity) Id() int {
	return res.GetInt("id")
}

func (res *TimeEntriesActivity) Name() string {
	return res.GetString("name")
}

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

func (res *TimeEntry) SetComment(format, raw, html string) {
	res.SetField("comment", NewFormattable(format, raw, html))
}

func (res *TimeEntry) SpentOn() *time.Time {
	if dt, err := res.GetDateTime("spentOn"); err == nil {
		return &dt
	}
	return nil
}

func (res *TimeEntry) SetSpentOn(spentOn time.Time) {
	res.SetDate("spentOn", spentOn)
}

func (res *TimeEntry) Hours() *time.Duration {
	if dt, err := res.GetDuration("hours"); err == nil {
		return &dt
	}
	return nil
}

func (res *TimeEntry) SetHours(hours time.Duration) {
	res.SetDuration("hours", hours)
}

func (res *TimeEntry) SetActivity(activity string) {
	actLink := NewLink(activity)
	res.AddLink("activity", *actLink)
}

// Register Factories
func init() {
	resourceTypes["TimeEntriesActivity"] = func() Resource {
		return NewTimeEntriesActivity()
	}
	resourceTypes["TimeEntry"] = func() Resource {
		return NewTimeEntry()
	}
}
