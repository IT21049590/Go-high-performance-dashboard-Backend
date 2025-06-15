package config

var AppSettings = struct {
	ImportCSVOnStart bool
	CSVFilePath      string
}{
	ImportCSVOnStart: false,                     
	CSVFilePath:      "data/data.csv",           
}
