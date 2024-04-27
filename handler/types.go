package handler

type EstateRequest struct {
	Width  uint16 `json:"width" validate:"required,min=1,max=50000"`
	Length uint16 `json:"length" validate:"required,min=1,max=50000"`
}
