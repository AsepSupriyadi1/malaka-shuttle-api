package dto

// RouteRequest DTO untuk membuat route baru
type RouteRequest struct {
	OriginCity      string `json:"origin" binding:"required,min=2,max=100"`
	DestinationCity string `json:"destination" binding:"required,min=2,max=100"`
}

// RouteUpdateRequest DTO untuk update route
type RouteUpdateRequest struct {
	OriginCity      string `json:"origin" binding:"required,min=2,max=100"`
	DestinationCity string `json:"destination" binding:"required,min=2,max=100"`
}

// RouteResponse DTO untuk response route
type RouteResponse struct {
	ID              uint   `json:"id"`
	OriginCity      string `json:"origin"`
	DestinationCity string `json:"destination"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// RouteSearchRequest DTO untuk search route
type RouteSearchRequest struct {
	OriginCity      *string `json:"origin,omitempty" form:"origin"`
	DestinationCity *string `json:"destination,omitempty" form:"destination"`
	Page            int     `json:"page" form:"page" binding:"min=1"`
	Limit           int     `json:"limit" form:"limit" binding:"min=1,max=100"`
}

// Note: For pagination responses, we use utils.PaginationResponse directly
// to avoid nested data structure. This provides clean response with "results" field
// instead of "data" -> "data" nesting.
