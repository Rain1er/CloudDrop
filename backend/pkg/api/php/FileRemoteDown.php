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

function main($path = "",$url="")
{
	$path=getgbkStr($path);
	$url=getgbkStr($url);
	$filename=substr($url, strrpos($url,'/')+1); 
	if (!file_exists($path)) {
        mkdir($path,0777);
        @chmod($path,0777);
    }
    ob_start(); 
    readfile($url);
    $img = ob_get_contents();
    ob_end_clean(); 
    file_put_contents($path.$filename, $img);
	echo "Success";
}