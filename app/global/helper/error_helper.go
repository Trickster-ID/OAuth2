package helper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"oauth2/app/global/model"
	"runtime"
	"strings"
)

var DefaultStatusText = map[int]string{
	http.StatusInternalServerError: "Terjadi Kesalahan, Silahkan Coba lagi Nanti",
	http.StatusNotFound:            "Data tidak Ditemukan",
	http.StatusBadRequest:          "Ada kesalahan pada request data, silahkan dicek kembali",
}

func WriteLog(err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		output := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
		}
		outputForPrint := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
			Line:       fmt.Sprintf("%d", line),
			Filename:   file,
			Function:   funcName,
		}

		output.SystemMessage = err.Error()
		if message == nil {
			output.Message = DefaultStatusText[errorCode]
			if output.Message == "" {
				output.Message = http.StatusText(errorCode)
				outputForPrint.Message = http.StatusText(errorCode)
			}
		} else {
			output.Message = message
			outputForPrint.Message = message
		}
		if errorCode == http.StatusInternalServerError {
			output.Line = fmt.Sprintf("%d", line)
			output.Filename = file
			output.Function = funcName
		}

		logForPrint := map[string]interface{}{}
		_ = DecodeMapType(outputForPrint, &logForPrint)

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(logForPrint).Error(err)
		return output
	}

	return nil
}
