package response

// Message message model
type Message struct {
	Message string `json:"message,omitempty" example:"success"`
}

// NewSuccessMessage new message success model
func NewSuccessMessage() *Message {
	return &Message{
		Message: "success",
	}
}
