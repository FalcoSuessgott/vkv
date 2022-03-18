package ui

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mingrammer/cfmt"
)

func iface(list []string) []interface{} {
	// https://groups.google.com/d/msg/golang-nuts/yTxmAjoc_vw/FheJAz0Q5MYJ
	vals := make([]interface{}, len(list))
	for i, v := range list {
		vals[i] = v
	}

	return vals
}

// PromptYN prompts yes/no question.
func PromptYN(msg string, params ...string) bool {
	s := fmt.Sprintf(msg, iface(params)...)
	res := false
	prompt := &survey.Confirm{
		Message: s,
	}

	err := survey.AskOne(prompt, &res)
	if err != nil {
		fmt.Println("Exiting.")
		os.Exit(0)
	}

	return res
}

// FSuccessMsg prints out success message to the specified writer.
func FSuccessMsg(w io.Writer, msg string, params ...string) {
	if _, err := cfmt.Fsuccessln(w, fmt.Sprintf(msg, iface(params)...)); err != nil {
		log.Fatal(err)
	}
}

// FInfoMsg prints out info message to the specified writer.
func FInfoMsg(w io.Writer, msg string, params ...string) {
	if _, err := cfmt.Finfoln(w, fmt.Sprintf(msg, iface(params)...)); err != nil {
		log.Fatal(err)
	}
}

// FErrorMsg prints out error message to the specified writer.
func FErrorMsg(w io.Writer, err error, msg string, params ...string) {
	if _, err := cfmt.Ferrorln(w, fmt.Sprintf(msg, iface(params)...)); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		if _, err := cfmt.Ferrorln(w, fmt.Sprintf("\nError: %s.\nExiting.", err.Error())); err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(1)
}

// FWarningMsg prints out warning message to the specified writer.
func FWarningMsg(w io.Writer, err error, msg string, params ...string) {
	if _, err := cfmt.Ferrorln(w, fmt.Sprintf(msg, iface(params)...)); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		if _, err := cfmt.Ferrorln(w, fmt.Sprintf("\nError: %s.\nNot fatal - continuing.", err.Error())); err != nil {
			log.Fatal(err)
		}
	}
}
