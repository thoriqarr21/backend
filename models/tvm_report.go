package models

import (
	"time"


)

type ReportStatus string

const (
	StatusPending    ReportStatus = "pending"
	StatusInProgress ReportStatus = "in_progress"
	StatusResolved   ReportStatus = "resolved"
	StatusRejected   ReportStatus = "rejected"
)

type TvmReportResponse struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type TVMReport struct {
	ID          uint           `gorm:"primarykey;column:id" json:"id"`
	TVMCode     string         `gorm:"type:varchar(50);not null;column:tvm_code" json:"tvm_code"`
	Location    string         `gorm:"type:varchar(255);not null;column:location" json:"location"`
	IssueType   string         `gorm:"type:varchar(100);not null;column:issue_type" json:"issue_type"`
	Description string         `gorm:"type:text;not null;column:description" json:"description"`
	Status      ReportStatus   `gorm:"type:varchar(20);default:'pending';not null;column:status" json:"status"`
	Priority    string         `gorm:"type:varchar(20);default:'medium';column:priority" json:"priority"`
	ImageURL    string         `gorm:"type:varchar(500);column:image_url" json:"image_url,omitempty"`
	ReportedBy  uint           `gorm:"not null;column:reported_by" json:"reported_by"`
	Reporter    User           `gorm:"foreignKey:ReportedBy" json:"reporter,omitempty"`
	ResolvedBy  *uint          `gorm:"column:resolved_by" json:"resolved_by,omitempty"`
	Resolver    *User          `gorm:"foreignKey:ResolvedBy" json:"resolver,omitempty"`
	ResolvedAt  *time.Time     `gorm:"column:resolved_at" json:"resolved_at,omitempty"`
	Notes       string         `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name
func (TVMReport) TableName() string {
	return "tvm_reports"
}

type CreateReportRequest struct {
	TVMCode     string `json:"tvm_code" binding:"required"`
	Location    string `json:"location" binding:"required"`
	IssueType   string `json:"issue_type" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ImageURL    string `json:"image_url"`
}

type UpdateReportRequest struct {
	TVMCode     string `json:"tvm_code" binding:"required"`
	Location    string `json:"location" binding:"required"`
	IssueType   string `json:"issue_type" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      ReportStatus `json:"status" binding:"omitempty,oneof=pending in_progress resolved rejected"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ResolvedBy  *uint          `json:"resolved_by,omitempty"`
	ImageURL    string `json:"image_url"`
}

type ReportFilterRequest struct {
	Status     string `form:"status"`
	Priority   string `form:"priority"`
	TVMCode    string `form:"tvm_code"`
	ReportedBy uint   `form:"reported_by"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

type UpdateReportStatusRequest struct {
	Status      ReportStatus `json:"status" binding:"omitempty,oneof=pending in_progress resolved rejected"`
}