package constants

import (
	"net/http"
	dto "vm/internal/dtos"
	api "vm/internal/gen"
)

const (
	InvalidJSONFormatErrorCode  = "INVALID_JSON"
	SQLRecordNotFoundErrorCode  = "RECORD_NOT_FOUND"
	LoadStatusConflictErrorCode = "STATUS_CONFLICT"
	ValidationErrorCode         = "VALIDATION_ERROR"
	UnauthorizedErrorCode       = "UNAUTHORIZED"
	AuthorizationErrorCode      = "FORBIDDEN"
	InternalServerErrorCode     = "INTERNAL_ERROR"
)

var ErrorCodeToStatus = map[string]int{
	InvalidJSONFormatErrorCode:  http.StatusBadRequest,
	SQLRecordNotFoundErrorCode:  http.StatusNotFound,
	LoadStatusConflictErrorCode: http.StatusConflict,
	ValidationErrorCode:         http.StatusUnprocessableEntity,
	UnauthorizedErrorCode:       http.StatusUnauthorized,
	AuthorizationErrorCode:      http.StatusForbidden,
}

var responseRegistry = map[OperationType]map[int]func(api.ErrorResponse) any{
	VMPowerOff: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.VMPowerOffBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMPowerOffInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMPowerOffNotFound)(&e) },
	},
	VMPowerOn: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.VMPowerOnBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMPowerOnInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMPowerOnNotFound)(&e) },
	},
	VMDelete: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.VMDeleteBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMDeleteInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMDeleteNotFound)(&e) },
	},
	VMDeploy: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.HCIDeployVMBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.HCIDeployVMInternalServerError)(&e) },
	},
	VMReset: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.VMPowerResetBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMPowerResetInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMPowerResetNotFound)(&e) },
	},
	VMRefresh: {
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMRefreshInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMRefreshNotFound)(&e) },
	},
	VMRestartGuestOS: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.VMRestartGuestOSBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMRestartGuestOSInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMRestartGuestOSNotFound)(&e) },
	},
	VMShutdownGuestOS: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.VMShutdownGuestOSBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.VMShutdownGuestOSInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.VMShutdownGuestOSNotFound)(&e) },
	},
	VMReconfigure: {
		http.StatusBadRequest:          func(e api.ErrorResponse) any { return (*api.EditVMBadRequest)(&e) },
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.EditVMInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.EditVMNotFound)(&e) },
	},
	VMMachine: {
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.GetVirtualMachineRequestInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.GetVirtualMachineRequestNotFound)(&e) },
	},
	VMMachineList: {
		http.StatusInternalServerError: func(e api.ErrorResponse) any { return (*api.GetVirtualMachineRequestListInternalServerError)(&e) },
		http.StatusNotFound:            func(e api.ErrorResponse) any { return (*api.GetVirtualMachineRequestListNotFound)(&e) },
	},
}

func MapServiceError(err dto.ApiResponseError, Operation OperationType) any {
	statusCode, ok := ErrorCodeToStatus[err.ErrorCode]
	if !ok {
		statusCode = http.StatusInternalServerError
		err.ErrorCode = InternalServerErrorCode
	}

	errRes := api.ErrorResponse{
		// DebugId:        uuid.New().String(),
		ErrorCode:      err.ErrorCode,
		HttpStatusCode: statusCode,
		Message:        err.Message,
	}

	if opMap, ok := responseRegistry[Operation]; ok {
		if constructor, ok := opMap[statusCode]; ok {
			return constructor(errRes)
		}
	}
	// fallback
	return fallbackErrorResponse(Operation, errRes)
}

func fallbackErrorResponse(op OperationType, errRes api.ErrorResponse) any {
	switch op {
	case VMPowerOff:
		return (*api.VMPowerOffInternalServerError)(&errRes)
	case VMPowerOn:
		return (*api.VMPowerOnInternalServerError)(&errRes)
	case VMDelete:
		return (*api.VMDeleteInternalServerError)(&errRes)
	case VMDeploy:
		return (*api.HCIDeployVMInternalServerError)(&errRes)
	case VMReset:
		return (*api.VMPowerResetInternalServerError)(&errRes)
	case VMRefresh:
		return (*api.VMRefreshInternalServerError)(&errRes)
	case VMRestartGuestOS:
		return (*api.VMRestartGuestOSInternalServerError)(&errRes)
	case VMShutdownGuestOS:
		return (*api.VMShutdownGuestOSInternalServerError)(&errRes)
	case VMReconfigure:
		return (*api.EditVMInternalServerError)(&errRes)
	case VMMachine:
		return (*api.GetVirtualMachineRequestInternalServerError)(&errRes)
	case VMMachineList:
		return (*api.GetVirtualMachineRequestListInternalServerError)(&errRes)
	default:
		// Generic fallback if operation is unknown
		return (*api.ErrorResponse)(&errRes)
	}
}
