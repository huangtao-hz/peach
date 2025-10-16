module peach

go 1.24.5

require golang.org/x/text v0.29.0

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/extrame/xls v0.0.1
	github.com/huangtao-hz/excelize v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.32
)

replace github.com/huangtao-hz/excelize => ../excelize

require (
	github.com/extrame/ole2 v0.0.0-20160812065207-d69429661ad7 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/tiendc/go-deepcopy v1.6.1 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a9 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/net v0.44.0 // indirect
)
