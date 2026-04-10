// Package jennifer is a code generator for Go
package jennifer

//go:generate go tool gennames -output "jen/hints.go" -package "jen" -name "standardLibraryHints" -standard -novendor -path "./..."
//go:generate go tool becca -package=github.com/veggiemonk/jennifer/jen
