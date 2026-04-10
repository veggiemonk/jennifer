module github.com/veggiemonk/jennifer

go 1.26.2

retract v1.26.2 // Dict rendering missing trailing commas in multi-entry map literals

tool (
	github.com/dave/rebecca/cmd/becca
	github.com/veggiemonk/jennifer/gennames
)

require (
	github.com/dave/gopackages v0.0.0-20250212082220-e38641800008 // indirect
	github.com/dave/jennifer v1.7.1 // indirect
	github.com/dave/kerr v0.0.0-20230520060319-776b4fd985da // indirect
	github.com/dave/rebecca v0.9.2 // indirect
)
