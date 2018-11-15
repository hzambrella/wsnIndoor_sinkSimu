//parse param.first is required,second is unrequired
package ginIOC_param

var builtinParseMethod = map[string][2][]string{
	"bool": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseBool(%sStr);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=pr_%s",
		},
		[]string{"	pr_%s,_:=strconv.ParseBool(%sStr)", "	%s:=pr_%s"},
	},
	"int": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.Atoi(%sStr);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=pr_%s",
		},
		[]string{"	pr_%s,_:=strconv.Atoi(%sStr)", "%s:=pr_%s"},
	},
	"int8": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseInt(%sStr, 10, 8);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=int8(pr_%s)",
		},
		[]string{"	pr_%s,_:=strconv.ParseInt(%sStr, 10, 8)", "%s:=int8(pr_%s)"},
	},
	"int16": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseInt(%sStr, 10, 16);if err!=nil{",
			dealErr,
			"	} ",
			"	%s:=int16(pr_%s)",
		},
		[]string{"	pr_%s,_:=strconv.ParseInt(%sStr, 10, 16)", " %s:=int16(pr_%s)"},
	},
	"int32": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseInt(%sStr, 10, 32);if err!=nil{",
			dealErr,
			"	}",
			" 	%s:=int32(pr_%s)",
		},
		[]string{"pr_%s,_:=strconv.ParseInt(%sStr, 10, 32)", " %s:=int32(pr_%s)"},
	},
	"int64": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseInt(%sStr, 10, 64);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=pr_%s",
		},
		[]string{"	%pr_s,_:=strconv.ParseInt(%sStr, 10, 64)", " %s:=pr_%s"},
	},
	"uint": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseUint(%sStr, 10, 0);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=uint(pr_%s)",
		},
		[]string{"	pr_%s,_:=strconv.ParseUint(%sStr, 10, 0)", " %s:=uint(pr_%s)"},
	},
	"uint16": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseUint(%sStr, 10, 16);if err!=nil{",
			dealErr,
			"	} ",
			"	%s:=uint16(pr_%s)",
		},
		[]string{"	pr_%s,_:=strconv.ParseUint(%sStr, 10, 16)", " %s:=uint16(pr_%s)"},
	},
	"uint32": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseUint(%sStr, 10, 32);if err!=nil{",
			dealErr,
			"	} ",
			"	%s:=uint32(pr_%s)",
		},
		[]string{"	pr_%s,_:=strconv.ParseUint(%sStr, 10, 32) ", "%s:=uint32(pr_%s)"},
	},
	"uint64": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseUint(%sStr, 10, 64);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=pr_%s",
		},
		[]string{"	pr_%s,_:=strconv.ParseUint(%sStr, 10, 64) ", "%s:=pr_%s"},
	},
	"float32": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseFloat(%sStr, 32);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=float32(pr_%s)",
		},
		[]string{"	pr_%s,_:=strconv.ParseFloat(%sStr, 32)", " %s:=float32(pr_%s)"},
	},
	"float64": [2][]string{
		[]string{
			"	pr_%s,err:=strconv.ParseFloat(%sStr, 64);if err!=nil{",
			dealErr,
			"	}",
			"	%s:=pr_%s",
		},
		[]string{"	pr_%s,_:=strconv.ParseFloat(%sStr, 64)", " %s:=pr_%s"},
	},
	"string": [2][]string{
		[]string{
			"	if %sStr==\"\"{ //%s",
			dealErr,
			"	}",
			"	%s:=%sStr"},
		[]string{"	pr_%s:=%sStr", " %s:=pr_%s"},
	},
}
