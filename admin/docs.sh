#! /bin/bash 
###########################################
#
###########################################

# constants
baseDir=$(cd `dirname "$0"`;pwd)
# functions

# main 
[ -z "${BASH_SOURCE[0]}" -o "${BASH_SOURCE[0]}" = "$0" ] || return
cd $baseDir/..
if [ -d docs ]; then
    rm -rf docs
fi

mkdir docs
cd docs
godoc -html  .. > index.html
sed -i 's/\/src\/target\//https:\/\/github.com\/chatopera\/chatopera-go-sdk\/blob\/master\//' index.html