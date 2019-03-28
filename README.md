# TVM(Fabric-sdk)
##Environment

- 1.Fabric-v1.4.0
- 2.docker
- 3.docker-compose

##Install

- Install fabric single peer
    
        - mkdir -p $GOPATH/github.com/hyperledger/fabric/singlepeer
        - cd $GOPATH/github.com/hyperledger/fabric/singlepeer
        - wget https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/linux-amd64-1.4.0/hyperledger-fabric-linux-amd64-1.4.0.tar.gz
        - tar -xvf hyperledger-fabric-linux-amd64-1.4.0.tar.gz
        - chmod -R 777 ./bin
        - cp $GOPATH/src/tvm-sdk/config/fabric-single/*.yaml .
        - cp $GOPATH/src/tvm-sdk/config/docker-compose/*.yaml .
        - mkdir channel-artifacts
        - ./bin/cryptogen generate --config=./crypto-config.yaml
        - ./bin/configtxgen -profile OneOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
        - ./bin/configtxgen -profile OneOrgsChannel -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannel

- Prepare couchdb path

        - mkdir -p /var/hyperledger/data/orderer0
        - mkdir -p /var/hyperledger/data/couchdb0/_data
        - mkdir -p /var/hyperledger/data/peer0

- Prepare project path

        - mkdir -p /trias-contract/contract/
        - mkdir -p /trias-contract/backup/
        
  
##Interfaces
- type:POST,application/json
#### 1.executeContract
- address:http://address/executeContract
- requestData:
    
        {
            "address": "QmStXVUdqAbKDTecFQkawuvSNqPs6su5KJB94Uvd9MiCny",
            "checkMD5": "7a3f59dd79140c6ce5de2d6a6ef5e352",
            "command": "{\"Args\":[\"query\",\"b\"]}",
            "contractName": "firstContract",
            "contractType": "fabric",
            "contractVersion": "v1.0",
            "vmVersion": "1.0",
            "sequence": "10",
            "timestamp": 1545117978737,
            "user": "user1",
            "signature": "",
            "operation": "query"
        }
- details

|字段|说明|
|-----|------|
|address|合约的ipfs地址|
|checkMD5|合约的md5校验码|
|command|合约执行的命令|
|contractName|合约名|
|contractType|合约类型,目前只有fabric|
|contractVersion|合约版本|
|vmVersion|虚拟机版本|
|sequence|执行序号（预留字段，是否可用作以后共识状态使用）|
|timestamp|时间戳|
|user|用户|
|signature|用户签名|
|operation|操作类型,install:安装;instantiate:初始化;query:查询;……|

#### 2.upLoadPackage
- address:http://address/upLoadPackage
- requestData:
        
        {}
- details
- responseData:

        {
            "Code":1
            "Message":"success"
            "Data":"QmTfBETxQcXe19rWDowNNLs5hBekRFecmpXaCdqh2nnpnD"
        }
- details

|字段|说明|
|-----|------|
|data|ipfs文件地址，同步数据使用|

#### 3.asyncTVM
- address:http://address/asyncTVM
- requestData:

        {
            "ipfsHash":"QmTfBETxQcXe19rWDowNNLs5hBekRFecmpXaCdqh2nnpnD"
        }
- details

|字段|说明|
|-----|------|
|ipfsHash|ipfs文件地址|

#### 4.getCurrentHash
- address:http://address/getCurrentHash
- requestData:
        
        {}
- details
- responseData
        {
            "Code":1
            "Message":"success"
            "Data":"NmE5Yzc5MTE5NWNmYzZlYTAyYmFhZTVhMTNkZWFhYWUwNWE3OWQ1OWQ1NjlhYjU4OWRjMWQ3ZmQwNzJhNjc2Zg=="
        }
- details

|字段|说明|
|-----|------|
|data|当前状态的hash| 

#### 5.getCurrentDataAddress
- address:http://address/getCurrentDataAddress
- requestData:
        
        {}
- details
- responseData
        {
            "Code":1
            "Message":"success"
            "Data":"QmZ6Wr9Jcw2rkrMHTN8zFSBhHTxRpoAX31PVLQzqGUekuQ"
        }
- details

|字段|说明|
|-----|------|
|data|ipfs文件地址| 