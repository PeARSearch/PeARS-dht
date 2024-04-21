#!/bin/bash
# Handler script to process events and queries.

# Read the payload from stdin (if any)
read PAYLOAD

# This would be local search pod
SEARCH_URL=https://www.searchfoss.nl/

RESPONSE=$(curl -s -X GET $SEARCH_URL?q=$PAYLOAD)

# Extract URLs using grep and sed

echo "$RESPONSE" | grep -oP 'https://[^"]+' |grep -v '^https://www.searchfoss.nl/' | sed 's/https:\/\///g'| sed 's/amp;//g' | sed 's|</a>||g' | sed 's|</p>||g' |  uniq | paste -sd ',' - | head -c 900
