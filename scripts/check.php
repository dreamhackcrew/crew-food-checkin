#!/usr/bin/php
<?php

require_once('db.php');

$uid = trim($argv[1])+0;
$uid = substr($uid,0,-1);

$fid = db()->fetchOne("SELECT *  FROM `user_food` WHERE `when` < '".date('Y-m-d H:i:s')."' ORDER BY `when` DESC LIMIT 1");

//echo "Meal: $fid, uid: $uid\n";

db()->query("INSERT INTO user_food_checkin (fid,uid) VALUES (%d,%d)",$fid,$uid);

die('success');
//die('fail');

?>
