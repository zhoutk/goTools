package common

type OpFlag struct {
	Tables bool
	Datum bool
	Views bool
	Funcs bool
}

type DbConnFields struct {
	DbHost    string
	DbPort    int
	DbUser    string
	DbPass    string
	DbName    string
	DbCharset string
	FileAlias string
}

