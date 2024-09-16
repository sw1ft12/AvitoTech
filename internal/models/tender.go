package models

import "time"

type Tender struct {
	Id              string    `json:"id,omitempty" db:"id"`
	Name            string    `json:"name,omitempty" db:"name"`
	Description     string    `json:"description,omitempty" db:"description"`
	ServiceType     string    `json:"serviceType,omitempty" db:"type"`
	Status          string    `json:"status,omitempty" db:"status"`
	CreatorUsername string    `json:"creatorUsername,omitempty" db:"created_by"`
	Version         int       `json:"version,omitempty" db:"version"`
	OrganizationId  string    `json:"organizationId,omitempty" db:"organization_id"`
	CreatedAt       time.Time `json:"createdAt,omitempty" db:"created_at"`
}

type CreateTenderDto struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	Version         string `json:"version"`
	CreatorUsername string `json:"creatorUsername"`
	OrganizationId  int    `json:"organizationId"`
}
