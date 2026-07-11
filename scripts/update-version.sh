#!/usr/bin/env sh
set -eu

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <version-or-tag>" >&2
  exit 1
fi

version="$1"
version="${version#refs/tags/}"
version="${version#v}"

case "$version" in
  ''|*[!0-9.]*|*..*|.*|*.)
    echo "Invalid semantic version: $1" >&2
    exit 1
    ;;
esac

perl -0pi -e 's/var AppVersion = "[^"]+"/var AppVersion = "'"$version"'"/' config/config.go
perl -0pi -e 's/\*\*SSH Tunnel Manager \(sshtm\) v[^*]+\*\*/**SSH Tunnel Manager (sshtm) v'"$version"'**/' README.md
perl -0pi -e 's/sshtm_[0-9]+\.[0-9]+\.[0-9]+-1_amd64\.deb/sshtm_'"$version"'-1_amd64.deb/g' README.md
perl -0pi -e 's/^pkgver=.*/pkgver='"$version"'/m' packaging/arch/PKGBUILD
if [ -f packaging/arch/.SRCINFO ]; then
  perl -0pi -e 's/(pkgver = )[0-9]+\.[0-9]+\.[0-9]+/${1}'"$version"'/g; s/(sshtm-)[0-9]+\.[0-9]+\.[0-9]+(\.tar\.gz)/${1}'"$version"'${2}/g; s/(tags\/v)[0-9]+\.[0-9]+\.[0-9]+/${1}'"$version"'/g' packaging/arch/.SRCINFO
fi
perl -0pi -e 's/^sshtm \([0-9]+\.[0-9]+\.[0-9]+-1\)/sshtm ('"$version"'-1)/m' packaging/debian/changelog
