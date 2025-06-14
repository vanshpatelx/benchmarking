#!/bin/bash

URL="http://localhost:8080/signup"
TOTAL_STEPS=10
INITIAL_REQUESTS=100
INCREMENT_REQUESTS=100
CONCURRENCY=50
SLEEP_BETWEEN=1


echo "Starting gradual signup benchmark test for $URL"

for (( i=1; i<=TOTAL_STEPS; i++ ))
do
    REQUESTS=$(( INITIAL_REQUESTS + (i - 1) * INCREMENT_REQUESTS ))
    echo ""
    echo "Step $i: Sending $REQUESTS signup requests with concurrency $CONCURRENCY..."

    # Use same payload for all requests
    echo "username=user_$RANDOM&password=pass123" > payload.txt

    hey -n "$REQUESTS" -c "$CONCURRENCY" \
        -m POST \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -D payload.txt \
        "$URL"

    echo "Sleeping for $SLEEP_BETWEEN seconds..."
    sleep "$SLEEP_BETWEEN"
done

echo "Benchmark test completed."
