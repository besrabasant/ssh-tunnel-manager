#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(dirname -- "$script_dir")"
package_dir="$repo_root/packaging/arch"
generated_srcinfo="$(mktemp "$package_dir/.SRCINFO.generated.XXXXXX")"

cleanup() {
	rm -f "$generated_srcinfo"
}
trap cleanup EXIT

if [[ "$(id -u)" -eq 0 ]]; then
	if ! id build >/dev/null 2>&1; then
		useradd -m build
	fi
	chown -R build:build "$package_dir"
	runuser -u build -- bash -c "cd '$package_dir' && makepkg --printsrcinfo > '$generated_srcinfo'"
else
	(
		cd "$package_dir"
		makepkg --printsrcinfo > "$generated_srcinfo"
	)
fi

if ! diff -u "$package_dir/.SRCINFO" "$generated_srcinfo"; then
	echo "Error: packaging/arch/.SRCINFO is out of date. Regenerate it with makepkg --printsrcinfo." >&2
	exit 1
fi

echo ".SRCINFO is up to date."
