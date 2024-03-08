#!/bin/sh

cd ../lib

make clean
make all
cp ./librgbmatrix.* ../go

cd ../go
