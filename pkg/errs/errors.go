package errs

import (
	"errors"
	"net/http"

	"github.com/stellar-payment/sp-payment/pkg/constant"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrBrokenUserReq       = errors.New("invalid request")
	ErrInvalidCred         = errors.New("invalid user credentials")
	ErrDuplicatedResources = errors.New("item already existed")
	ErrNoAccess            = errors.New("user does not have required access privilege")
	ErrUnknown             = errors.New("internal server error")
	ErrNotFound            = errors.New("resources not found")
)

const (
	ErrCodeUndefined           constant.ErrCode = 1
	ErrCodeBadRequest          constant.ErrCode = 2
	ErrCodeNoAccess            constant.ErrCode = 3
	ErrCodeInvalidCred         constant.ErrCode = 4
	ErrCodeDuplicatedResources constant.ErrCode = 5
	ErrCodeBrokenUserReq       constant.ErrCode = 6
	ErrCodeNotFound            constant.ErrCode = 7
)

const (
	ErrStatusUnknown     = http.StatusInternalServerError
	ErrStatusClient      = http.StatusBadRequest
	ErrStatusNotLoggedIn = http.StatusUnauthorized
	ErrStatusNoAccess    = http.StatusForbidden
	ErrStatusReqBody     = http.StatusUnprocessableEntity
	ErrStatusNotFound    = http.StatusNotFound
)

var errorMap = map[error]dto.ErrorResponse{
	ErrUnknown:             ErrorResponse(ErrStatusUnknown, ErrCodeUndefined, ErrUnknown),
	ErrBadRequest:          ErrorResponse(ErrStatusClient, ErrCodeBadRequest, ErrBadRequest),
	ErrInvalidCred:         ErrorResponse(ErrStatusNotLoggedIn, ErrCodeInvalidCred, ErrInvalidCred),
	ErrNoAccess:            ErrorResponse(ErrStatusNoAccess, ErrCodeNoAccess, ErrNoAccess),
	ErrDuplicatedResources: ErrorResponse(ErrStatusClient, ErrCodeDuplicatedResources, ErrDuplicatedResources),
	ErrBrokenUserReq:       ErrorResponse(ErrStatusReqBody, ErrCodeBrokenUserReq, ErrBrokenUserReq),
	ErrNotFound:            ErrorResponse(ErrStatusNotFound, ErrCodeNotFound, ErrNotFound),
}

var errorRPCMap = map[error]error{
	ErrUnknown:             ErrorRPCResponse(codes.Internal, ErrUnknown),
	ErrBadRequest:          ErrorRPCResponse(codes.InvalidArgument, ErrBadRequest),
	ErrInvalidCred:         ErrorRPCResponse(codes.PermissionDenied, ErrInvalidCred),
	ErrNoAccess:            ErrorRPCResponse(codes.PermissionDenied, ErrNoAccess),
	ErrDuplicatedResources: ErrorRPCResponse(codes.AlreadyExists, ErrDuplicatedResources),
	ErrBrokenUserReq:       ErrorRPCResponse(codes.InvalidArgument, ErrBrokenUserReq),
	ErrNotFound:            ErrorRPCResponse(codes.NotFound, ErrNotFound),
}

func ErrorResponse(status int, code constant.ErrCode, err error) dto.ErrorResponse {
	return dto.ErrorResponse{
		Status:  status,
		Code:    code,
		Message: err.Error(),
	}
}

func ErrorRPCResponse(code codes.Code, err error) error {
	return status.Error(code, err.Error())
}

func GetErrorResp(err error) (errResponse dto.ErrorResponse) {
	errResponse, ok := errorMap[err]
	if !ok {
		errResponse = errorMap[ErrUnknown]
	}

	return
}

func GetErrorRPC(err error) (rpcerr error) {
	rpcerr, ok := errorRPCMap[err]
	if !ok {
		rpcerr = errorRPCMap[ErrUnknown]
	}

	return
}
