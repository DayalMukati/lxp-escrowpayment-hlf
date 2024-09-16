#!/bin/bash

# Initialize score
score=0
source ./scripts/setOrgPeerContext.sh 1
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

# Query the escrow status
echo "Querying the escrow status..."
QUERY_OUTPUT=$(peer chaincode query -C mychannel -n escrowpayment -c '{"Args":["QueryEscrowStatus","escrow1"]}' 2>&1)

echo "Query output: $QUERY_OUTPUT"

# Check if the escrow account exists
if [[ $QUERY_OUTPUT == *"Error"* ]]; then
    echo "Escrow account not found. It was not created."
    # Output the final score
    echo "Final Score: $score/50"
    exit 0
else
    echo "Escrow account found."
fi

# Check if the escrow was created (indicating escrow creation)
if [[ $QUERY_OUTPUT == *"escrow1"* ]]; then
    echo "Escrow creation successful."
    score=$((score + 10))
fi

# Check if the payer has approved the escrow (indicating payer approval)
if [[ $QUERY_OUTPUT == *"payerApproved\":true"* ]]; then
    echo "Payer approval successful."
    score=$((score + 10))
fi

# Check if the payee has confirmed the escrow (indicating payee confirmation)
if [[ $QUERY_OUTPUT == *"payeeConfirmed\":true"* ]]; then
    echo "Payee confirmation successful."
    score=$((score + 10))
fi

# Check if the escrow funds were released (indicating fund release)
if [[ $QUERY_OUTPUT == *"released\":true"* ]]; then
    echo "Escrow funds released successfully."
    score=$((score + 10))
fi

# Check if all conditions are met (full escrow lifecycle)
EXPECTED_OUTPUT='{"escrowID":"escrow1","payerID":"payer1","payeeID":"payee1","amount":1000,"payerApproved":true,"payeeConfirmed":true,"released":true}'
if [[ $QUERY_OUTPUT == "$EXPECTED_OUTPUT" ]]; then
    echo "Escrow query fully successful."
    score=$((score + 10))
else
    echo "Escrow query incomplete."
fi

# Final score output
echo "Final Score: $score/50"

# Exit with success
exit 0
