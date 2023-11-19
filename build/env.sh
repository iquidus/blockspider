#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
dir="$workspace/src/github.com/iquidus"
if [ ! -L "$dir/blockspider" ]; then
    mkdir -p "$dir"
    cd "$dir"
    ln -s ../../../../../. blockspider
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$dir/blockspider"
PWD="$dir/blockspider"

# Launch the arguments with the configured environment.
exec "$@"