package customFunc

const (
	funcTest      = "testFunc"
	funcTestError = "testFuncError"
)

// testFunc ... return (a + b)
func testFunc(a int, b int) FuncResult {
	return FuncResult{a + b, ""}
}

// testFuncError ... return nil, error
func testFuncError() FuncResult {
	return FuncResult{nil, "this is testing error function"}
}
