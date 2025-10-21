package modals

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interface for setting RequestID
type UUIDModel interface {
	SetRequestID(string)
}

// Utility to set a new UUID string
func SetUUID(model UUIDModel) {
	model.SetRequestID(uuid.New().String())
}

// VMRequest model
type VMRequest struct {
	RequestID       string     `gorm:"column:request_id;primaryKey;type:char(36)"`      // UUID stored as string
	Operation       string     `gorm:"column:operation;not null;type:varchar(50)"`      // Operation type
	RequestStatus   string     `gorm:"column:request_status;not null;type:varchar(50)"` // Status of the request
	WorkspaceId     string     `gorm:"column:workspace_id;type:varchar(50);default:''"`
	DatacenterId    string     `gorm:"column:datacenter_id;type:varchar(50);default:''"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime;type:timestamp"` // Auto-set on creation
	CompletedAt     *time.Time `gorm:"column:completed_at;type:timestamp"`              // Completion time
	RequestMetadata string     `gorm:"column:request_metadata;type:text"`               // JSON stored as string
}

// VMDeployInstance model
type VMDeployInstance struct {
	RequestID      string     `gorm:"column:request_id;primaryKey;type:char(36)"` // UUID primary key
	VMID           string     `gorm:"column:vm_id;type:varchar(50)"`              // Optional VM identifier
	VMName         string     `gorm:"column:vm_name;not null;type:varchar(255)"`  // Required VM name
	VMStatus       string     `gorm:"column:vm_status;not null;type:varchar(50)"` // Optional status
	VMStateMessage string     `gorm:"column:vm_state_message;type:text"`          // Optional message
	CompletedAt    *time.Time `gorm:"column:completed_at;type:timestamp"`         //
}

// Implement UUIDModel for VMRequest
func (r *VMRequest) SetRequestID(id string) {
	r.RequestID = id
}

func (r *VMRequest) BeforeCreate(tx *gorm.DB) (err error) {
	SetUUID(r)
	return nil
}

// Implement UUIDModel for VMDeployInstance
func (r *VMDeployInstance) SetRequestID(id string) {
	r.RequestID = id
}

func (r *VMDeployInstance) BeforeCreate(tx *gorm.DB) (err error) {
	SetUUID(r)
	return nil
}
