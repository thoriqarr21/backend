package models

import "time"

type TVMReportHistory struct {
	ID          	uint         `gorm:"primaryKey"`
 	TVMReportID 	uint       	 `gorm:"column:tvm_report_id"`
	Report      	*TVMReport   `gorm:"foreignKey:TVMReportID;references:ID" json:"report,omitempty"`
	FromStatus  	ReportStatus `gorm:"type:varchar(30)"`
	ToStatus    	ReportStatus `gorm:"type:varchar(30)"`
	ChangedBy   	uint
	ChangedByUser   User         `gorm:"foreignKey:ChangedBy" json:"changed_by,omitempty"`
	Role        	Role         `gorm:"type:varchar(20)"`
	CreatedAt   	time.Time
}

func (TVMReportHistory) TableName() string {
    return "tvm_report_histories"
}