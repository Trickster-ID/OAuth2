package helper

import (
	"github.com/gofiber/fiber/v3"
	"oauth2/app/global/model"
	"reflect"
)

func Response(ctx fiber.Ctx, data any) error {
	totalData := 0
	if data != nil {
		v := reflect.ValueOf(data)
		t := v.Type()

		// Check if the input is a slice or array
		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
			totalData = v.Len()
		}
	}
	response := model.Api{
		StatusCode:    200,
		StatusMessage: "success",
		Data:          data,
		TotalCount:    totalData,
	}
	return ctx.Status(200).JSON(response)
}

func ResponseError(ctx fiber.Ctx, errorLog *model.ErrorLog) error {
	response := model.Api{
		StatusCode:    errorLog.StatusCode,
		StatusMessage: "error",
		Data:          nil,
		TotalCount:    0,
		ErrorLog:      errorLog,
	}
	return ctx.Status(errorLog.StatusCode).JSON(response)
}
