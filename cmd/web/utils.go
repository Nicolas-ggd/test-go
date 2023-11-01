package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)

	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(res http.ResponseWriter, status int) {
	http.Error(res, http.StatusText(status), status)
}

func (app *application) notFound(res http.ResponseWriter) {
	app.clientError(res, http.StatusNotFound)
}

func (app *application) decodePostForm(res *http.Request, decode any) error {
	err := res.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(decode, res.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}
