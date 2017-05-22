<?php

require_once('workflows.php');

const SITE_URL = 'http://wooordhunt.ru';

$w = new Workflows();

$query = implode(' ', array_slice($argv, 1));

$url = SITE_URL . "/get_tips.php?abc=". urlencode($query);
$suggestions = json_decode($w->request($url), true);

if ($suggestions) {
  foreach($suggestions['tips'] as $suggest) {
    $w->result(
      time(),
      SITE_URL . '/word/' . $suggest['w'],
      $suggest['t'],
      $suggest['w'],
      'icon.png',
      'yes'
    );
  }
} else {
  $w->result(
    time(),
    SITE_URL . '/word/' . $query,
    $query,
    null,
    'icon.png',
    'yes'
  );
}

echo $w->toxml();
