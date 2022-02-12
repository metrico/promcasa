package model

type DataDatabasesMap struct {
	Value        string `json:"value"`
	Node         string `json:"node"`
	Name         string `json:"name"`
	DBname       string `json:"db_name"`
	Host         string `json:"host"`
	TableSeries  string `json:"time_series"`
	TableSamples string `json:"samples"`
	Primary      bool   `json:"primary"`
	Online       bool   `json:"online"`
}

type ConfigDatabasesMap struct {
	Value           string   `json:"value"`
	Name            string   `json:"name"`
	Node            string   `json:"node"`
	Host            string   `json:"host"`
	Primary         bool     `json:"primary"`
	Online          bool     `json:"online"`
	URL             string   `json:"url"`
	ProtectedTables []string `json:"-"`
	SkipTables      []string `json:"-"`
}

type ConfigURLNode struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Primary bool   `json:"primary"`
}
