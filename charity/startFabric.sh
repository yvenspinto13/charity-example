#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

# launch network; create channel and join peer to channel
cd ../network
./start.sh

# Now launch the CLI container in order to install, instantiate chaincode
# and put some initial donations 
docker-compose -f ./docker-compose.yml up -d charitycli

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.charity.com/users/Admin@org1.charity.com/msp" charitycli peer chaincode install -n charity -v 1.0 -p github.com/charity
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.charity.com/users/Admin@org1.charity.com/msp" charitycli peer chaincode instantiate -o orderer.charity.com:7050 -C charitychannel -n charity -v 1.0 -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.charity.com/users/Admin@org1.charity.com/msp" charitycli peer chaincode invoke -o orderer.charity.com:7050 -C charitychannel -n charity -c '{"function":"initLedger","Args":[""]}'

printf "\nTotal setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
printf "Run 'node enrollAdmin.js' to enroll admin \n\n"
printf "Run 'node registerUser' to register a user \n\n"
printf "Run 'node invoke.js' to donate, spend \n"
printf "The 'node query.js' to check balance and all donations or particular donation \n\n"
