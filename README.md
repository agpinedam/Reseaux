## Simulateur de protocole d'information sur le routage (RIP) en Go

Ce projet met en œuvre un simulateur de protocole d'information sur le routage (RIP) en utilisant le langage de programmation Go. RIP est un protocole de routage à vecteur de distance utilisé pour échanger des informations de routage dans les réseaux IP.

### Description et résumé du projet

L'objectif principal de ce projet est de fournir un outil de simulation permettant de comprendre et d'expérimenter le fonctionnement du protocole RIP. Le simulateur se compose de plusieurs parties, chacune remplissant une fonction spécifique :

- **Client:** Le client UDP gère la communication avec le serveur à l'aide de messages RIP. Il envoie des messages RIP au serveur pour mettre à jour sa table de routage et reçoit des réponses du serveur avec des informations de routage mises à jour.

- **Serveur:** Le serveur UDP reçoit les messages RIP des clients, traite les informations reçues et met à jour sa propre table de routage en conséquence. Il envoie également des réponses aux clients lorsqu'ils reçoivent des demandes d'informations de routage.

- **Routeur:** Le routeur est responsable de la construction et de la gestion de la table de routage du serveur. Il utilise des algorithmes de routage pour calculer les meilleurs itinéraires et mettre à jour la table de routage en conséquence.

- **Table routeur:** Il représente la topologie du réseau et stocke les informations nécessaires pour déterminer la meilleure route pour envoyer les paquets de données. Elle est mise à jour dynamiquement lorsque des mises à jour de routage sont reçues du serveur.

### Exécution du projet

Pour exécuter le projet, procédez comme suit

1. **Activer l'interface réseau virtuelle:** 2.

   Avant de lancer le serveur, assurez-vous d'activer l'interface réseau virtuelle virbr10 et de lui attribuer une adresse IP. Ceci peut être fait avec les commandes suivantes :

````bash
sudo ip link set dev virbr10 up
```

````bash
sudo ip addr add 10.1.1.3/24 dev virbr10
```

**Run Server:** ``````````````````````

```` ````bash
sudo go run server/server.go
```

**Exécution des tests du serveur:** 

````bash 
sudo go run server/server.go 
````

````bash
sudo go test ./server
```

**Exécutez le client:**


````bash
go run client/client.go
```

**Exécution des tests du client:** 

````bash 
go run client/client.go 
```

````bash
go test ./client
````

