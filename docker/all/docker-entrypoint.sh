#/usr/bin/env sh

./backend &
node server.js &

ret=$(wait)
kill -INT $(jobs -p)
exit $ret
