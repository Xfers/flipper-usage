#!/bin/bash
MYSQL_BIN=$(which mysql)
if [ "x$MYSQL_BIN" == "x" ]; then
    echo "mysql cli doesn't exist"
    exit 1
fi

if [ $# -lt 5 ]; then
    echo "usage: dump_flipper_flags.sh user password host dbname output.txt"
    exit 1
fi
mysql -u $1 -p$2 -h $3 $4 -e "select flipper_features.key from flipper_features" > $5
sed -i '1d' $5