package config

var AppSettings = struct {
	ImportCSVOnStart bool
	CSVFilePath      string
	RefreshTime      int
}{
	ImportCSVOnStart: false,
	CSVFilePath:      "data/data.csv",
	RefreshTime:      6,
}
