package service

// Shell interface defines common methods for different shell types
type Shell interface {
	GetShellType() string
	FreshSession(id int, url, password string) (string, error)
	BaseInfo(id int, url, password string) (string, error)

	ExecCommand(id int, command, url, password string) (string, error)
	ExecCode(id int, code, url, password string) (string, error)
	ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error)

	FileZip(id int, srcPath, toPath, url, password string) (string, error)
	FileUnZip(id int, srcPath, toPath, url, password string) (string, error)

	FileList(id int, path, url, password string) (string, error)
	FileShow(id int, path, url, password string) (string, error)
}

// PHPShell implements Shell interface for PHP
type PHPShell struct{}

// JavaShell implements Shell interface for Java
type JavaShell struct{}

// CSharpShell implements Shell interface for C#
type CSharpShell struct{}

// AspShell implements Shell interface for ASP
type AspShell struct{}
