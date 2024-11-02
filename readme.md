# MediaTools

MediaTools est une collection d'outils pour manipuler des médias (vidéos uniquement).

## IMPORTANT

A la création d'une issue, merci d'attribuer l'issue a la personne qui s'occupera de la résoudre. Sinon aucune notification ne sera envoyée à la personne concernée (et donc personne ne sera prévenu).

## Release

Vous pouvez télécharger la dernière version de MediaTools sur la page des [releases](https://github.com/Developpeur-du-dimanche/MediaTools/releases)

## Installation


### Prérequis

Avoir GO installé sur votre machine. Pour l'installer, suivez les instructions sur le site officiel: https://golang.org/doc/install

Installation des dépendances:
```bash
go get
```

### Compilation

Pour compiler le projet, exécutez la commande suivante:
```bash
go build -o ./main.exe ./cmd/mediatools/main.go
```

Cela peut etre long, car le projet doit compiler les dépendances pour l'interface graphique. Il faut donc patienter.

### Exécution

Pour exécuter le programme, exécutez la commande suivante:
```bash
./main.exe
```

### Air (optionnel, pour le développement principalement)

Air permet de recharger automatiquement le programme lorsqu'un fichier est modifié. Cela permet de gagner du temps lors du développement.

#### Installation

Pour installer Air, exécutez la commande suivante:
```bash
go install github.com/air-verse/air@latest
```

Pour exécuter le programme avec Air, exécutez la commande suivante:
```bash
air
```
