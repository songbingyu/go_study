#!/usr/bin/expect

set timeout 20 
spawn ssh root@115.28.216.4
expect "*password:"
send   "010d64ee\n"
expect  "#"
interact 
