# Y-Plaza - Infrastructure Cloud

## Description

Cette version correspond au déploiement Cloud du projet.

Elle permet de démontrer :

- Hébergement des machines virtuelles
- Accès sécurisé
- Services Windows Server
- Active Directory
- Sauvegarde

---

## Machines

### Windows Server

Nom :

SRV-YPLAZA-DC01

Compte :

adminazure

Mot de passe :

Password123!

---

### Windows Client

Nom :

CLIENT-01

Compte :

adminazure / Administrateur

Mot de passe :

Password123!

---

## Domaine

yplaza.local

---

## Services

- Active Directory
- DNS
- DHCP
- GPO
- Firewall
- VPN
- Sauvegarde

---

## Architecture

Cloud

↓

Firewall

↓

Windows Server

↓

Windows Client

---

## Démonstration

- Connexion au serveur
- Vérification des utilisateurs AD
- Vérification des GPO
- Test des partages
- Test de la sauvegarde automatique