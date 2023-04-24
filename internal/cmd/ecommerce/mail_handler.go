package ecommerce

import (
	"net/http"
)

type MailService interface {
	Send(recipient string, data any, patterns ...string) error
}

// postMailHandler is an example of posting email
func postMailHandler(w http.ResponseWriter, r *http.Request) {
	//checking query params
	emailPrm := r.URL.Query().Get("email")
	namePrm := r.URL.Query().Get("name")
	if emailPrm == "" || namePrm == "" {
		logger.Errorf("Error sending email. Email and|or Name params not specified.")
		writeResponse(w, "Error sending email. Email and|or Name params not specified.", http.StatusBadRequest)
		return
	}

	//preparing the dto
	type Person struct {
		Name string
	}
	p := Person{Name: namePrm}

	//sending email using hello.tmpl template
	err := mailService.Send(emailPrm, p, "hello.tmpl")
	if err != nil {
		logger.Errorf("Error. Can't send smtp message.", err)
		writeResponse(w, "Error. Can't send smtp message.", http.StatusInternalServerError)
		return
	}

	//success
	writeResponse(w, "Mail was sent", http.StatusOK)
}

func writeResponse(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(text))
	if err != nil {
		logger.Errorf("can't write response", err)
	}
}
