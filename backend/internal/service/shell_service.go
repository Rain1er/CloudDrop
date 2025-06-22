package service

// Shell interface defines common methods for different shell types
type Shell interface {
	GetShellType() string
	FreshSession(id int, url string, password string) (string, error)
	BaseInfo(id int, url string, password string) (string, error)
	ExecCommand(id int, command string, url string, password string) (string, error)
	FileList(id int, path string, url string, password string) (string, error)
	FileShow(id int, path string, url string, password string) (string, error)
}

// PHPShell implements Shell interface for PHP
type PHPShell struct{}

// JavaShell implements Shell interface for Java
type JavaShell struct{}

// CSharpShell implements Shell interface for C#
type CSharpShell struct{}

// AspShell implements Shell interface for ASP
type AspShell struct{}
