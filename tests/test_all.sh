#!/bin/sh
MYCC=./mycc
CC=gcc

echo ========================
echo C compiler test suite
echo ========================
for i in *.c ; do
  /usr/bin/printf "testing %.8s... " $i
  $MYCC $i > /dev/null 2>&1

  if [ ! $? -eq 0 ]; then 
    echo "failed to compile!"
  fi
  
  ./a.out
  actual=$?
  CC $i > /dev/null 2>&1
  ./a.out
  correct=$?

  if [ $actual -eq $correct ]; then
    echo "ok!"
  fi
done

# clean
rm -rf a.out mycc