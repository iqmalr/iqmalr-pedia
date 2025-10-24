package response

import "time"

type ApplicationResponse struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
