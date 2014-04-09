#!/usr/bin/expect

set timeout 20
spawn  scp  Http  root@115.28.216.4:/data
expect "*password"
send "010d64ee\n"
interact
