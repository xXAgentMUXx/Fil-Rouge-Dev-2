# Y-Plaza : Infrastructure Locale

## Description

Cette maquette représente l'infrastructure locale du projet Y-Plaza réalisée sous VMware.

Elle comprend :

- Windows Server
- Windows Client
- Active Directory
- DNS
- DHCP
- GPO
- pfSense (Pare-feu / Routeur)
- VPN WireGuard
- Sauvegarde automatique via Robocopy

---

# Machines virtuelles

## 1. Windows Server

Nom :

SRV-YPLAZA-DC01

Adresse IP :

192.168.1.10

Compte :

Administrateur

Mot de passe :

Password123!

Lien OneDrive :

https://1drv.ms/u/c/026cf40436ea1da3/IQB2Cq3_3RirT5EsghmlxEXRAWNOcw9O1hD7L1Fd4twTnDo?e=V0CgB5

---

## 2. Windows Client

Nom :

CLIENT-01

Adresse IP :

Attribuée automatiquement par DHCP
(généralement 192.168.1.100)

Compte administrateur :

YPLAZA\Administrateur

Mot de passe :

Password123!

Compte utilisateur :

YPLAZA\nom_utilisateur

Mot de passe :

Yplaza@1234

Lien OneDrive :

https://1drv.ms/u/c/026cf40436ea1da3/IQCotbJzK0G6TrSU1MqJzmvpATgr7w_F1MIH3e3aZeng9nw?e=m8ur8h

---

## 3. pfSense

Fonction :

- Pare-feu
- Routeur
- NAT

Interfaces :

WAN : NAT

LAN : 192.168.1.1

Accès Web :

https://192.168.1.1

Compte :

admin

Mot de passe :

Password123!

Lien OneDrive :

https://1drv.ms/u/c/026cf40436ea1da3/IQDVNPnokzsBTYyLhWa9WJsUAbVVlfTL0LEkMFnn2D7WO_k?e=0xMdor

---

# Domaine Active Directory

Nom du domaine :

yplaza.local

---

# Services installés

- Active Directory
- DNS
- DHCP
- GPO
- Partages réseau
- VPN WireGuard
- Robocopy automatisé
- Planificateur de tâches
- pfSense

---

# Sauvegarde

La synchronisation est automatisée grâce à Robocopy.

Serveur :

```
C:\Partage
```

↓

Client :

```
\\CLIENT-01\BackupClient
```

Les sauvegardes sont exécutées automatiquement via le Planificateur de tâches.

---

# Sécurité

- Pare-feu pfSense
- Pare-feu Windows Defender
- Authentification Active Directory
- GPO de sécurité
- Sauvegarde automatisée
- VPN WireGuard

---

# Démonstration

1. Démarrer les trois machines virtuelles.
2. Se connecter au serveur.
3. Créer un fichier dans :

```
C:\Partage
```

4. Attendre l'exécution automatique de Robocopy.
5. Vérifier la présence du fichier sur le client.
6. Tester la connexion avec un utilisateur du domaine.
7. Vérifier l'application des GPO et des droits d'accès.

---

## Auteurs

- Urban Mathys
- Axel Macé
- Vittorio Gandossi

Bachelor 2 Informatique Ynov