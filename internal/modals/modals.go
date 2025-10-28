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
    RequestID       string     `gorm:"column:request_id;primaryKey;type:char(36)" json:"request_id"`
    Operation       string     `gorm:"column:operation;not null;type:varchar(50)" json:"operation"`
    RequestStatus   string     `gorm:"column:request_status;not null;type:varchar(50)" json:"request_status"`
    WorkspaceId     string     `gorm:"column:workspace_id;type:varchar(50);default:''" json:"workspace_id"`
    DatacenterId    string     `gorm:"column:datacenter_id;type:varchar(50);default:''" json:"datacenter_id"`
    CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime;type:timestamp" json:"created_at"`
    CompletedAt     *time.Time `gorm:"column:completed_at;type:timestamp" json:"completed_at"`
    RequestMetadata string     `gorm:"column:request_metadata;type:text" json:"request_metadata"`
}
 
// VMDeployInstance model
type VMDeployInstance struct {
    RequestID      string     `gorm:"column:request_id;primaryKey;type:char(36)" json:"request_id"`
    VMName         string     `gorm:"column:vm_name;primaryKey;type:varchar(255)" json:"vm_name"`
    VMID           string     `gorm:"column:vm_id;type:varchar(50)" json:"vm_id"`
    VMStatus       string     `gorm:"column:vm_status;not null;type:varchar(50)" json:"vm_status"`
    VMStateMessage string     `gorm:"column:vm_state_message;type:text" json:"vm_state_message"`
    CompletedAt    *time.Time `gorm:"column:completed_at;type:timestamp" json:"completed_at"`
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
// func (r *VMDeployInstance) SetRequestID(id string) {
//  r.RequestID = id
// }
 
// func (r *VMDeployInstance) BeforeCreate(tx *gorm.DB) (err error) {
//  SetUUID(r)
//  return nil
// }