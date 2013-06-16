#!/usr/bin/php
<?php

//if ( trim($argv[1]) == '73513513') 
//	die('success');
//else
//	die('fail');

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

print_r(get("https://api.crew.dreamhack.se/1/food/list"));

?>
