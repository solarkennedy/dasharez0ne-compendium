dasharez0ne.tweets.json:
	twitter-archive -noat -nort @dasharez0ne > $@

update: dasharez0ne.tweets.json
	twitter-archive -noat -nort -a dasharez0ne.tweets.json
	jq . dasharez0ne.tweets.json | jq . | grep media_url\" | sort | cut -f 4 -d '"' | uniq > archive/media_urls.txt
	cd archive && wget --mirror -i media_urls.txt
