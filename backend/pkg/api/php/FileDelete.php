error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');

function main($path)
{
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

