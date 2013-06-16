<?php

require_once('db.php');

$customer = '8518c43fbd828995ac39a1bf4c7651da4f70cc5d';
$secret = 'b833aba7d21febbcd695beff608cfd8d7fb8b568';

function get($url) {

	global $customer,$secret;
	$context = stream_context_create(array(
	    'http' => array(
		'header'  => "Authorization: Basic " . base64_encode("$customer:$secret")
	    )
	));
	$data = file_get_contents($url, false, $context);

	return json_decode($data,true);
}

$meals = get("https://api.crew.dreamhack.se/1/food/list");
foreach($meals as $meal) {
	if ( $fid = db()->fetchOne("SELECT fid FROM user_food WHERE `when`='%s'",$meal['datetime']) ) {
		db()->query("UPDATE user_food SET name='%s' WHERE fid=%d",$meal['name'],$fid);
	} else {
		$fid = db()->insert(array(
			'when' => $meal['datetime'],
			'name' => $meal['name']
		),'user_food');
	}

	if ( !$exists = db()->fetchAllOne("SELECT uid FROM user_food_selected WHERE fid=%d",$fid) )
		$exists = array();

	$diff = array_diff($meal['registered'],$exists);

	foreach($diff as $key => $line) {
		db()->query("INSERT INTO user_food_selected (uid,fid) VALUES (%d,%d)",$line,$fid);
	}

	echo $fid." - {$meal['datetime']} - {$meal['name']}\n";

}



?>
