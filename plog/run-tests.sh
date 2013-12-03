#!/bin/bash

# start sessiond
$GOPATH/src/github.com/Grant-Murray/sessiond/run.bash

# verify sessiond is running
if [ "$(pidof sessiond)" = "" ]; then
  echo sessiond did not start
  exit 1
fi

sudo rm -f /tmp/mailbot.boxes/*
sudo rm -f /tmp/webdriver.*
sudo rm -f /tmp/sessdb.*
sudo rm -f /tmp/session.test*

cd $GOPATH/src/github.com/Grant-Murray/webdriver/session
go test register_test.go verifyemail_test.go login_test.go -v 2>&1 | grep -v '^.selenium] '

PSQL="psql --username=postgres --dbname=sessdb"
$PSQL -c "select * from session.log" > /tmp/webdriver.db.log
$PSQL -c 'select * from session.user' --expanded > /tmp/webdriver.db.user
$PSQL -c 'select * from session.session' --expanded > /tmp/webdriver.db.session

cd /tmp/mailbot.boxes
for MBOX in * ; do
  sudo mv -v $MBOX /tmp/webdriver.email.$MBOX
done

echo "You may want to cleanup: killall sessiond"
