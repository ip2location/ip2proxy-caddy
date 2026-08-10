<?php
$prefix = 'HTTP_X_IP2PROXY_';
$fields = [
	'IS_PROXY',
	'PROXY_TYPE',
	'COUNTRY_CODE',
	'COUNTRY_NAME',
	'REGION',
	'CITY',
	'ISP',
	'DOMAIN',
	'USAGE_TYPE',
	'ASN',
	'AS',
	'LAST_SEEN',
	'THREAT',
	'PROVIDER',
	'FRAUD_SCORE',
];

foreach ($fields as $val) {
	$key = $prefix . $val;
	if (isset($_SERVER[$key])) {
		echo $key . ': ' . $_SERVER[$key] . "\n<br>";
	}
}
?>
