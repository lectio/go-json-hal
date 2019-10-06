package hal

//
// Attachment
//

type Attachment struct {
	ResourceObject
}

func NewAttachment() *Attachment {
	return &Attachment{
		ResourceObject{
			Type: "Attachment",
		},
	}
}

func (res *Attachment) Description() *Formattable {
	f, err := DecodeFormattable(res.GetField("description"))
	if err != nil {
		return nil
	}
	return f
}

func (res *Attachment) Id() int {
	return res.GetInt("id")
}

func (res *Attachment) FileName() string {
	return res.GetString("fileName")
}

func (res *Attachment) FileSize() int {
	return res.GetInt("fileSize")
}

func (res *Attachment) ContentType() string {
	return res.GetString("contentType")
}

// Register Factories
func init() {
	resourceTypes["Attachment"] = func() Resource {
		return NewAttachment()
	}
}
