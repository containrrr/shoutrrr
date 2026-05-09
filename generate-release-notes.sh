#!/bin/bash

current=$1
if [ -z "$current" ]; then
  echo "Missing argument VERSION"
  exit 1
fi

tags=($(git tag --list))

for i in "${!tags[@]}"; do
   if [[ "${tags[$i]}" = "$current" ]]; then
       previous="${tags[$i - 1]}"
       break
   fi
done

if [ -z "$previous" ]; then
  echo "Invalid tag, or could not find previous tag"
  exit 1
fi

echo -e "\e[97mListing changes from \e[96m$previous\e[97m to \e[96m$current\e[0m:\n"

changes=$(git log --pretty=format:"* %h %s" $previous...$current)

echo "## Changelog"
echo "$changes" | grep -v "chore(deps)" | grep -v "Merge " | grep -v "chore(ci)"
echo
echo "### Dependencies"
echo "$changes" | grep "chore(deps)"