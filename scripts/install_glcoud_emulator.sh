#!/bin/bash

unameOut="$(uname -s)"

if [ "$TMPDIR" == "" ];then
  TMPDIR="/tmp/"
fi

outFile="$TMPDIR/google-cloud-sdk-x86_64.tar.gz"

# -O file
case "${unameOut}" in
    Linux*)     wget -O $outFile "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-163.0.0-linux-x86_64.tar.gz";;
    Darwin*)    wget -O $outFile "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-163.0.0-darwin-x86_64.tar.gz";;
    *)          echo "not supported operation system"; exit 1;;
esac

tarOutFolder="$TMPDIR/google-cloud-sdk-x86_64"
if [ -d $tarOutFolder ];then
  rm -rf $tarOutFolder
fi
mkdir $tarOutFolder

tar -xvzf $outFile --directory $tarOutFolder
cd $tarOutFolder
sh $tarOutFolder/google-cloud-sdk/install.sh --quiet --path-update true --rc-path $HOME/.gcloud_profile
# source $(find $HOME -type f -name '*profile' -maxdepth 1)