#!/usr/bin/env bash

# usage: ./git-deploy.sh -m "commit message"

while getopts "m:" opt; do
  case $opt in
    m) MSG="$OPTARG" ;;
  esac
done

if [ -z "$MSG" ]; then
  echo "Commit message required: -m \"message\""
  exit 1
fi

set -e

git add .
git commit -m "$MSG"
git push origin main

git switch render-deploy
git cherry-pick main
git push origin render-deploy

git switch main

echo "Commit complete"
