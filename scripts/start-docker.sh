#!/bin/bash

KEY="dev0"
CHAINID="serv_43970-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t serv-datadir.XXXXX)

echo "create and add new keys"
./servd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Serv with moniker=$MONIKER and chain-id=$CHAINID"
./servd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./servd add-genesis-account \
"$(./servd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000aserv,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./servd gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./servd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./servd validate-genesis --home $DATA_DIR

echo "starting serv node $i in background ..."
./servd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started serv node"
tail -f /dev/null