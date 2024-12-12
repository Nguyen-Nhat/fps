package errorz

import "fmt"

func OutputFileTypeNotSupported(outputFileType string) error {
	return fmt.Errorf("output file type %v is not supported", outputFileType)
}

func InputFileTypeNotSupported(inputFileType string) error {
	return fmt.Errorf("input file type %v is not supported", inputFileType)
}

func SlackSendingErr(status string) error {
	return fmt.Errorf("sending message to slack err, status: %s", status)
}
