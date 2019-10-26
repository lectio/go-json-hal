package hal

import (
	"encoding/json"
	"log"
)

//
// Filters
//

type FilterOperator struct {
	Operator string        `json:"operator"`
	Values   []interface{} `json:"values"`
}

type Filter map[string]FilterOperator

type FilterList []Filter

type Filters struct {
	FilterList
}

func NewFilters() *Filters {
	return &Filters{make([]Filter, 0, 1)}
}

func (f *Filters) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.FilterList)
}

func (f *Filters) String() string {
	if len(f.FilterList) == 0 {
		return ""
	}
	if buf, err := json.Marshal(f); err != nil {
		log.Fatal(err)
	} else {
		return string(buf)
	}
	return ""
}

func (f *Filters) Filter(name string, operator string, values ...interface{}) *Filters {
	filter := Filter{}
	filter[name] = FilterOperator{
		Operator: operator,
		Values:   values,
	}

	f.FilterList = append(f.FilterList, filter)
	return f
}
