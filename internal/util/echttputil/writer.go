package echttputil

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func WriteSuccessResponse(ec echo.Context, data interface{}) error {
	return ec.JSON(http.StatusOK, dto.BaseResponse{
		Data:   data,
		Errors: nil,
	})
}

func WriteErrorResponse(ec echo.Context, err error) error {
	errResp := errs.GetErrorResp(err)
	return ec.JSON(errResp.Status, dto.BaseResponse{
		Data:   nil,
		Errors: errResp,
	})
}

func WriteFileAttachment(ec echo.Context, path string, filename string) error {
	return ec.Attachment(path, filename)
}

func WriteFileBufferAttachment(ec echo.Context, file *bytes.Buffer, contentType string, filename string) error {
	ec.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", filename))
	return ec.Blob(http.StatusOK, contentType, file.Bytes())
}
