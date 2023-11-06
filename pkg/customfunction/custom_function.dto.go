package customFunc

// FuncResult ...
type FuncResult struct {
	Result       interface{}
	ErrorMessage string
}

// CustomFunction ...
type CustomFunction struct {
	FunctionPattern string   // format is $func.myFunction;1;{{$A}};{{$param.sellerId}}; ...
	Name            string   // E.g: myFunction
	ParamsRaw       []string // []string{"1", "{{$A}}", "{{$param.sellerId}}", ...}
	ParamsMapped    []string // []string{"1", "value colum A", "value of sellerId", ...}
}
