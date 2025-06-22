error_reporting(1);
header('Content-Type: text/html; charset=UTF-8');
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