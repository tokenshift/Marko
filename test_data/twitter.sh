#!/bin/sh

if [ -z "$CONSUMER_KEY" ]; then
	>&2 echo '$CONSUMER_KEY is required'
	exit 1
fi

if [ -z "$CONSUMER_SECRET" ]; then
	>&2 echo '$CONSUMER_SECRET is required'
	exit 1
fi

if [ -z "$ACCESS_TOKEN" ]; then
	>&2 echo '$ACCESS_TOKEN is required'
	exit 1
fi

if [ -z "$ACCESS_TOKEN_SECRET" ]; then
	>&2 echo '$ACCESS_TOKEN_SECRET is required'
	exit 1
fi

#CREDS="$CONSUMER_KEY:$CONSUMER_SECRET"
#ENCODED_CREDS=`echo "$CREDS" | base64`

#ACCESS_TOKEN=`curl -s "https://api.twitter.com/oauth2/token" \
#	-X POST -u "$CONSUMER_KEY:$CONSUMER_SECRET" \
#	-d "grant_type=client_credentials" \
#	| jq -r ".access_token"`

#curl -X POST "https://stream.twitter.com/1.1/statuses/sample.json" \
#	-H "Authorization: Bearer $ACCESS_TOKEN"

#curl -s "$API_URL/1.1/users/suggestions.json" \
#	-H "Authorization: Bearer $ACCESS_TOKEN" --verbose
#	| jq -r ".[].slug" \
#	| while read slug; do
#		echo $slug
#	curl -s "$API_URL/1.1/users/suggestions/$slug.json" \
#		-H "Authorization: Bearer $ACCESS_TOKEN" \
#		| jq "."
#done
