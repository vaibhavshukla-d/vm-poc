package constants

type Category string
type SubCategory string
type ExtraKey string

type OperationType string
type RequestStatus string
type VMDeployStatus string

const (
 General         Category = "General"
 IO              Category = "IO"
 Internal        Category = "Internal"
 MySql           Category = "MySql"
 Redis           Category = "Redis"
 Validation      Category = "Validation"
 RequestResponse Category = "RequestResponse"
)

const (
 // General
 Startup         SubCategory = "Startup"
 ExternalService SubCategory = "ExternalService"

 // SQL Server
 Migration SubCategory = "Migration"
 Select    SubCategory = "Select"
 Rollback  SubCategory = "Rollback"
 Update    SubCategory = "Update"
 Delete    SubCategory = "Delete"
 Insert    SubCategory = "Insert"

 // Internal
 Api                 SubCategory = "Api"
 HashPassword        SubCategory = "HashPassword"
 DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"
 FailedToCreateUser  SubCategory = "FailedToCreateUser"

 // Validation
 MobileValidation   SubCategory = "MobileValidation"
 PasswordValidation SubCategory = "PasswordValidation"

 // IO
 RemoveFile SubCategory = "RemoveFile"
)

const (
 AppName      ExtraKey = "AppName"
 LoggerName   ExtraKey = "Logger"
 ClientIp     ExtraKey = "ClientIp"
 HostIp       ExtraKey = "HostIp"
 Method       ExtraKey = "Method"
 StatusCode   ExtraKey = "StatusCode"
 BodySize     ExtraKey = "BodySize"
 Path         ExtraKey = "Path"
 Latency      ExtraKey = "Latency"
 RequestBody  ExtraKey = "RequestBody"
 ResponseBody ExtraKey = "ResponseBody"
 ErrorMessage ExtraKey = "ErrorMessage"
)

const (
 VMDeploy          OperationType = "VMDeploy"
 VMPowerOff        OperationType = "VMPowerOff"
 VMDelete          OperationType = "VMDelete"
 VMPowerOn         OperationType = "VMPowerOn"
 VMReset           OperationType = "VMReset"
 VMRefresh         OperationType = "VMRefresh"
 VMRestartGuestOS  OperationType = "VMRestartGuestOS"
 VMShutdownGuestOS OperationType = "VMShutdownGuestOS"

 StatusNew     RequestStatus = "New"
 StatusPending RequestStatus = "Pending"
 StatusDone    RequestStatus = "Done"

 VMINIT   VMDeployStatus = "Init"
 VMCLOSE  VMDeployStatus = "Close"
)

const VMRequestBasePath = "/virtualization/v1beta1/virtual-machines-request/"

 