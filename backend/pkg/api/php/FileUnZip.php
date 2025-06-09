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

function main($srcPath = "",$toPath="")
{
	$srcPath=getgbkStr($srcPath);
	$toPath=getgbkStr($toPath);
	$zip = new ZipArchive();
	$openRes = $zip->open($srcPath);
	$msg="Success";
	if ($openRes === TRUE) {
	  $zip->extractTo($toPath);
	  $zip->close();
	}
	else{
		$msg="Failed";
	}
	echo $msg;
}
