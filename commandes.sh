# Étape 1 : Prérequis et téléchargement
# 1. Installer Docker et Docker Compose (si pas déjà fait)
sudo apt-get update && sudo apt-get install docker.io docker-compose

# 2. Installer Go (nécessaire pour ton chaincode en Go)
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 3. Cloner fabric-samples (uniquement pour le déploiement)
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples

# 4. Télécharger les binaires Fabric (docker images + outils CLI)
./scripts/bootstrap.sh 2.5.5 1.5.6 0.4.23


# Étape 2 : Lancer le réseau de test
# 5. Démarrer le réseau (2 peers, 1 orderer)
./network.sh up

# 6. Créer un channel dédié aux virements
./network.sh createChannel -c payments

# 7. Copier TON chaincode personnalisé dans le dossier adéquat
cp -r /chemin/vers/ton/chaincode-go ./asset-transfer-basic/



# Étape 3 : Déployer ton chaincode
# 8. Déployer TON smart contract (⚠️ Nom = "transfer", pas "basic" !)
./network.sh deployCC -ccn transfer -ccp ../asset-transfer-basic/chaincode-go -ccl go

# 9. Vérifier que le chaincode est actif
peer lifecycle chaincode querycommitted -C payments

# Étape 4 : Interagir avec ton chaincode
# 1. Utiliser l'identité Admin de Org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# 2. Créer un virement (état "pending")
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C payments -n transfer \
  -c '{
    "Args": [
      "InitTransfer",
      "TX123",
      "CUST_FR01",
      "CUST_DE01",
      "5000",
      "EUR"
    ]
  }'


# 3. Approuver le virement (état "approved")
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C payments -n transfer \
  -c '{
    "Args": [
      "ApproveTransfer",
      "TX123",
      "APPROVER_FINANCE"
    ]
  }'


# 4. Exécuter le virement (état "executed")
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C payments -n transfer \
  -c '{
    "Args": [
      "ExecuteTransfer",
      "TX123"
    ]
  }'


# 5. Vérifier l’état final
peer chaincode query -C payments -n transfer \
  -c '{
    "Args": [
      "ReadTransfer",
      "TX123"
    ]
  }'


# Étape 5 : Déboguer les erreurs courantes
# 1. Vérifier l'identité active
peer identity list

# 2. Changer d'identité si nécessaire (ex: vers un approbateur)
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/APPROVER_FINANCE@org1.example.com/msp


# 3. Approuver pour Org2 (si nécessaire)
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

peer lifecycle chaincode approveformyorg -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C payments --channelID payments --name transfer --version 1 \
  --package-id transfer_1:abcd1234 --sequence 1
