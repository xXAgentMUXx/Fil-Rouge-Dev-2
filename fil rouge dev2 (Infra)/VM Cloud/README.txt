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

## Connexion domaine

Sur VM : CLIENT-01

Compte : YPLAZA\nom d'utilisateur

Mot de passe : 

Yplaza@1234

Admin : 

Compte : YPLAZA\adminazure

Mot de passe : 

Password123!


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

## Auteurs

Urban Mathys
Axel Macé
Vittorio Gandossi

Bachelor 2 Informatique Ynov