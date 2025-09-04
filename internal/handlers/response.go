// SetJSONResponse is a helper to set JSON response with a customizable JSON encoder.
func SetJSONResponse(ctx echo.Context, data interface{}, httpStatus int) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().WriteHeader(httpStatus)

	response := entity.GeneralAPIResponse{
		Data:  data,
		Error: nil,
	}

	return json.NewEncoder(ctx.Response()).Encode(response)
}
func SetJSONErrorResponse(ctx echo.Context, e error, httpStatusCode int32) error {
	errorResponse := internal_error.NewErrorResponse(e)

	// Set HTTP Status Code if the error defined a specific HTTP Status Code
	statusCode := httpStatusCode
	if errorResponse.HttpStatusCode != 0 {
		statusCode = int32(errorResponse.HttpStatusCode)
	}
	ctx.Response().WriteHeader(int(statusCode))
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := entity.GeneralAPIResponse{
		Error: errorResponse,
		Data:  nil,
	}

	return json.NewEncoder(ctx.Response()).Encode(response)
}
