package customFunc

const funcReUploadFile = "reUploadFile"

// reUploadFile ...
func reUploadFile(url string) FuncResult {
	return FuncResult{url + "abc", ""}
}
