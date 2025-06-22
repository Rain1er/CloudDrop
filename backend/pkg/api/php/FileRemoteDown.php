error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');



function main($path = "",$url="")
{
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