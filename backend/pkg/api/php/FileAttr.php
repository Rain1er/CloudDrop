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
function main($path,$attr,$time)
{
	$path=getgbkStr($path);
	$res=false;
	switch($attr){
		case "read":
		$res=@chmod($path,0555);
		break;
		case "write":
		$res=@chmod($path,0777);
		break;
		case "time":
		$res=touch($path,$time/1000);
		break;
		
	}
	if($res){
		echo "Success";
	}
	else{
		echo "Failed";
	}
}
