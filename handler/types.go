package handler

type EstateRequest struct {
	Width  uint16 `json:"width" validate:"required,min=1,max=50000"`
	Length uint16 `json:"length" validate:"required,min=1,max=50000"`
}

type TreeRequest struct {
	Height int `json:"height" validate:"required,min=1,max=30"`
	X      int `json:"x" validate:"required,min=1,max=50000"`
	Y      int `json:"y" validate:"required,min=1,max=50000"`
}

type DronePlanRequest struct {
	MaxDistance *uint16 `query:"max_distance"`
}
