package service

// Shell interface defines common methods for different shell types
type Shell interface {
	GetShellType() string
	GetCurrentDirectory(url string, password string) (string, error)
	ListFiles(url string, password string) ([]string, error)
	ExecCommand(url string, password string, command string) (string, error)
}

// PHPShell implements Shell interface for PHP
type PHPShell struct{}

// JavaShell implements Shell interface for Java
type JavaShell struct{}

// CSharpShell implements Shell interface for C#
type CSharpShell struct{}

// AspShell implements Shell interface for ASP
type AspShell struct{}
