error_reporting(1);
header('Content-Type: text/html; charset=UTF-8');
set_time_limit(0);
ini_set('max_execution_time', '0');
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
    $read_buffer = 4096;
    $handle = fopen($path,'rb');
    $sum_buffer = 0;
    while(!feof($handle)) {
        echo fread($handle,$read_buffer);
		ob_flush();
		flush();
    }
    fclose($handle);
}