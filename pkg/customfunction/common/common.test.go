package customFunc

const (
	FuncTest      = "testFunc"
	FuncTestError = "testFuncError"
)

// TestFunc ... return (a + b)
func TestFunc(a int, b int) FuncResult {
	return FuncResult{a + b, ""}
}

// TestFuncError ... return nil, error
func TestFuncError() FuncResult {
	return FuncResult{nil, "this is testing error function"}
}
