#!/bin/bash
# Mock ezctl for testing ksk8s task runner
# Usage: mock_ezctl.sh <cluster> <step>

CLUSTER="$1"
STEP="$2"

echo "[MOCK] ezctl setup $CLUSTER $STEP started"
echo "[MOCK] Cluster: $CLUSTER"
echo "[MOCK] Step: $STEP"

for i in {1..30}; do
    echo "[MOCK] Progress line $i/30 - simulating ansible output"
    sleep 0.2
done

echo "[MOCK] ezctl setup $CLUSTER $STEP completed successfully"
exit 0
