error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');
function getgbkStr($str){
    $s0 = iconv('gbk','utf-8//IGNORE',$s1);
    $s1 = iconv('utf-8','gbk//IGNORE',$str);
    if($s1 == $str){
        return $s1;
    }else{
        return iconv('utf-8','gbk//IGNORE',$str);
    }
}
function main($path)
{
	$path=getgbkStr($path);
	$dirs = explode(',',$path);
	foreach ($dirs as $dir) {
		if(is_file($dir)){
			unlink($dir);
		}
		else{
			delFile($dir);
		}
    }
	echo "Success";
}
function delFile($dir)
{
	$files = scandir($dir);
		foreach ($files as $file) {
			if($file!= "."&&$file!= ".."){
				$cpath="$dir/$file";
				if(is_file($cpath)){
					unlink($cpath);
				}
			else{
				delFile($cpath);
				rmdir($cpath);
			}
		}
	}
	rmdir($dir);
}

