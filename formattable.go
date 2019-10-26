package hal

import "errors"

//
// Formattable
//

type Formattable struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	Html   string `json:"html"`
}

func DecodeFormattable(data interface{}) (*Formattable, error) {
	if data == nil {
		return nil, nil
	}
	// convert to map
	mData, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("Expected a map")
	}

	f := &Formattable{}
	for key, val := range mData {
		s, ok := val.(string)
		switch key {
		case "format":
			f.Format = s
		case "raw":
			f.Raw = s
		case "html":
			f.Html = s
		default:
			// Ignore unknown fields
			continue
		}
		if !ok {
			return nil, errors.New("Expected a string.")
		}
	}
	return f, nil
}

func NewFormattable(format, raw, html string) *Formattable {
	return &Formattable{
		Format: format,
		Raw:    raw,
		Html:   html,
	}
}
