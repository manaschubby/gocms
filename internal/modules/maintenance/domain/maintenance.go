package domain

import (
	"time"

	"github.com/google/uuid"
)

type MaintenanceStatus string
type EscalationLevel string

const (
	StatusPending    MaintenanceStatus = "Pending"
	StatusAssigned   MaintenanceStatus = "Assigned"
	StatusInProgress MaintenanceStatus = "In Progress"
	StatusResolved   MaintenanceStatus = "Resolved"
	StatusRejected   MaintenanceStatus = "Rejected"
)

const (
	EscalationNone   EscalationLevel = "None"
	EscalationLevel1 EscalationLevel = "Level1" // 24h - Subcategory Supervisor
	EscalationLevel2 EscalationLevel = "Level2" // 48h - Category Manager
	EscalationLevel3 EscalationLevel = "Level3" // 72h - Dean Administration
)

// Category (e.g. Maintenance, IT)
type MaintenanceCategory struct {
	Id           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	ManagerEmail string    `json:"managerEmail" db:"manager_email"`
	Active       bool      `json:"active" db:"active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// Subcategory (e.g. Carpentry, Plumbing)
type MaintenanceSubcategory struct {
	Id              uuid.UUID `json:"id" db:"id"`
	CategoryId      uuid.UUID `json:"categoryId" db:"category_id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	SupervisorEmail string    `json:"supervisorEmail" db:"supervisor_email"`
	Active          bool      `json:"active" db:"active"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

// Detail / Child Category (specific issue types)
type MaintenanceDetail struct {
	Id            uuid.UUID `json:"id" db:"id"`
	SubcategoryId uuid.UUID `json:"subcategoryId" db:"subcategory_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Active        bool      `json:"active" db:"active"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// Worker / Field Personnel
type MaintenanceWorker struct {
	Id            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	UserEmail     string    `json:"userEmail" db:"user_email"`
	Phone         string    `json:"phone" db:"phone"`
	SubcategoryId uuid.UUID `json:"subcategoryId" db:"subcategory_id"`
	Active        bool      `json:"active" db:"active"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// Request
type MaintenanceRequest struct {
	Id               uuid.UUID         `json:"id" db:"id"`
	RequesterEmail   string            `json:"requesterEmail" db:"requester_email"`
	RequesterName    string            `json:"requesterName" db:"requester_name"`
	Location         string            `json:"location" db:"location"`
	CategoryId       uuid.UUID         `json:"categoryId" db:"category_id"`
	SubcategoryId    uuid.UUID         `json:"subcategoryId" db:"subcategory_id"`
	DetailId         *uuid.UUID        `json:"detailId" db:"detail_id"`
	Description      string            `json:"description" db:"description"`
	Status           MaintenanceStatus `json:"status" db:"status"`
	EscalationLevel  EscalationLevel   `json:"escalationLevel" db:"escalation_level"`
	LastEscalatedAt  *time.Time        `json:"lastEscalatedAt" db:"last_escalated_at"`
	AssignedWorkerId *uuid.UUID        `json:"assignedWorkerId" db:"assigned_worker_id"`
	AssignedAt       *time.Time        `json:"assignedAt" db:"assigned_at"`
	ResolvedAt       *time.Time        `json:"resolvedAt" db:"resolved_at"`
	ResolutionNotes  string            `json:"resolutionNotes" db:"resolution_notes"`
	CreatedAt        time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time         `json:"updatedAt" db:"updated_at"`
}

// Status Log (audit trail)
type MaintenanceStatusLog struct {
	Id             uuid.UUID          `json:"id" db:"id"`
	RequestId      uuid.UUID          `json:"requestId" db:"request_id"`
	UserEmail      string             `json:"userEmail" db:"user_email"`
	Action         string             `json:"action" db:"action"`
	PreviousStatus *MaintenanceStatus `json:"previousStatus" db:"previous_status"`
	NewStatus      *MaintenanceStatus `json:"newStatus" db:"new_status"`
	Comments       string             `json:"comments" db:"comments"`
	Timestamp      time.Time          `json:"timestamp" db:"timestamp"`
}

// Escalation Log
type MaintenanceEscalationLog struct {
	Id              uuid.UUID       `json:"id" db:"id"`
	RequestId       uuid.UUID       `json:"requestId" db:"request_id"`
	EscalationLevel EscalationLevel `json:"escalationLevel" db:"escalation_level"`
	NotifiedEmails  []string        `json:"notifiedEmails" db:"notified_emails"`
	Timestamp       time.Time       `json:"timestamp" db:"timestamp"`
}

// Config (singleton)
type MaintenanceConfig struct {
	DeanEmail string    `json:"deanEmail" db:"dean_email"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}