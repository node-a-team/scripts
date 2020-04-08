#!/bin/bash

printf "Enter version for update: "
read version

printf "\nv=%s" $version

#solDir=$HOME/.local/share/solana/install
solDir=$(echo $(which solana) |  sed "s/\/active_release\/bin\/solana//")
binaryURL=https://github.com/solana-labs/solana/releases/download/v$version/solana-release-x86_64-unknown-linux-gnu.tar.bz2

printf "\n%s" $newDir

mkdir $solDir/releases/$version
wget -P $solDir/releases/$version $binaryURL
tar -xvf $solDir/releases/$version/solana-release-x86_64-unknown-linux-gnu.tar.bz2 -C $solDir/releases/$version
rm $solDir/releases/$version/solana-release-x86_64-unknown-linux-gnu.tar.bz2

rm $solDir/active_release &&
ln -s $solDir/releases/$version/solana-release $solDir/active_release

printf "\n\033[32mUpdate complete\033[0m: "
solana --version

