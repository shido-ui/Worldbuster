package organization

type Type string

const (
	TypeCompany Type = "COMPANY"
	TypeFaction Type = "FACTION"
	TypeGuild   Type = "GUILD"
)

type Organization struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       Type   `json:"type"`
	OwnerID    string `json:"ownerId"`
	MaxMembers int    `json:"maxMembers"`
}
type Membership struct {
	OrganizationID string `json:"organizationId"`
	CharacterID    string `json:"characterId"`
	Role           string `json:"role"`
}
