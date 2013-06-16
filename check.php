<?php
$start = time();
$i = 0;

while(true) {
    $i++;
    //echo "\ncheck... ";
    if ( $i % 2 ) {
        unset($pids);
        exec("pidof crew-food-checkin", $pids);
    }

    if(empty($pids)) {
        //echo "gone";
        if ( !isset($setup) ) {
            $setup = 1;
            exec('/usr/local/bin/gpio export 17 out');
            exec('/usr/local/bin/gpio export 18 out');
            exec('PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/pi/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian/bin" nice -n -21 /home/pi/crew-food-checkin/bin/crew-food-checkin&');
            //exec("renice -10 `ps -C soundbox -o pid=`");
        }
        if ( $i % 2 ) {
            exec('echo "1" > /sys/class/gpio/gpio17/value');
            exec('echo "0" > /sys/class/gpio/gpio18/value');
        } else {
            exec('echo "0" > /sys/class/gpio/gpio17/value');
            exec('echo "1" > /sys/class/gpio/gpio18/value');
        }
    }
    usleep(500000);
    if( time()-$start > 59 )
        die();
}

?>
