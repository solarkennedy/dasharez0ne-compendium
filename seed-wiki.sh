#!/usr/bin/env bash


function embiggen {
  echo "$1" | sed 's/\.jpg$/\?format=jpg\&name=large/'
}

IFS="
"
for tweet in `cat dasharez0ne.tweets.json`; do

  id=$(echo "$tweet" | jq .created_at | xargs gdate '+%s' --date )
  if [[ $? -ne 0 ]]; then
    echo "Couldn't parse this tweet?:"
    echo "$tweet"
  else
    echo $id
    original_text="$(echo "$tweet" | jq .text)"
    image="$(echo "$tweet" | jq -r .entities.media[0].media_url_https)"
    if [[ "$image" == "null" ]]; then
      echo "Tweet: $original_text has no image. Continuing."
      continue
    fi
    image=$(embiggen $image)
    url="$(echo "$tweet" | jq -r .entities.media[0].url)"
    echo "---
original_text: "$original_text"
url: "$url"
tags:
  - untagged
image: "$image"
caption: |
  Please caption me as close to verbatium as possible.
  This is yaml so indent with two spaces on each line.
---
![]($image)
" > ../dasharez0ne-compendium.wiki/macros/${id}.md
  fi

done
