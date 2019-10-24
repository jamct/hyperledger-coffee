# Bürokaffeemaschine mit Hyperledger

Dies ist das Beispielprojekt zum c't-Artikel ["Vertraue niemandem"](https://ct.de/y8b9) aus c't 23/2019. Beschrieben wird, wie man den Putz- und Auffüllplan einer Bürokaffeemaschine in der Blockchain realisiert. Es dient als Einstiegsprojekt in das Open-Source-Blockchain-Framework Hyperledger.

## Einrichtung

Die detaillierte Einrichtung finden Sie im Artikel. Hier finden Sie alle Kommandozeilenbefehle zum Herauskopieren.

Download der Beispielumgebung von Hyperledger:

```
curl -sSL http://bit.ly/2ysbOFE | bash -s
```

Navigieren Sie in den Unterordner "fabric-samples/chaincode-docker-devmode". Danach starten Sie den Sie das Setup:

```
docker-compose -f docker-compose-simple.yaml up
```

## Code vorbereiten

Den Code finden Sie in der Datei "main.go" in diesem Artikel im Ordner "chaincode". Legen Sie diese Datei in den Ordner "fabric-samples/chaincode/coffee".

Um den Code auszuführen, müssen Sie in den Container springen:

```
docker exec -it chaincode bash
```

Damit befinden Sie sich im Container und können den Code builden und dann ausführen:

```
cd coffee
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./coffee
```

In einem neuen Terminalfenster in den Container "cli" springen:

```
docker exec -it cli bash
```

Chaincode initialisieren:

```
peer chaincode install -p chaincodedev/chaincode/coffee -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc
```

## Code ausführen

Beispielnutzer anlegen:

```
peer chaincode invoke -n mycc -c '{"Args":["storeUser", "Max"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["storeUser", "Klaus"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["storeUser", "Claudia"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["storeUser", "Peter"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["storeUser", "Christine"]}' -C myc
```

Maschine reinigen, auffüllen und Kaffee entnehmen:

```
peer chaincode invoke -n mycc -c '{"Args":["cleanMachine"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["refillCoffee"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["drawCoffee"]}' -C myc
```


## Vertiefung

Wer sich tiefer in die Materie einarbeiten will, kann sich Hyperledger Explorer ansehen. Damit kann man visualisieren, was in der Blockchain passiert:

- https://www.hyperledger.org/projects/explorer
- https://github.com/hyperledger/blockchain-explorer


Für macOS und Windows gibt es das Projekt Chaincoder, das aktuell noch in der Beta-Phase ist. Die Oberfläche macht es einfacher, kompliziertere Projekte zu entwickeln und zu debuggen:

- https://www.chaincoder.org/