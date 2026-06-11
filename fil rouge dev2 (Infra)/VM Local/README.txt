# Y-Plaza - Infrastructure Locale

## Description

Cette maquette représente l'infrastructure locale du projet Y-Plaza.

Elle comprend :

- Windows Server
- Active Directory
- DNS
- DHCP
- GPO
- pfSense (Firewall)
- VPN
- Sauvegarde automatique via Robocopy

---

## Machines virtuelles

### Windows Server

Nom :

SRV-YPLAZA-DC01

Adresse IP :

192.168.1.10

Compte :

Administrateur

Mot de passe :

Password123!

---

### Windows Client

Nom :

CLIENT-01

Adresse IP :

Attribuée automatiquement par DHCP

Compte :

Utilisateur / Administrateur

Mot de passe :

Password123!

---

## Domaine Active Directory

Nom :

yplaza.local

---

## Services installés

- Active Directory
- DNS
- DHCP
- GPO
- Partages réseau
- Robocopy automatisé
- Planificateur de tâches

---

## Sauvegarde

La synchronisation est automatique grâce à Robocopy.

Serveur :

C:\Partage

↓

Client :

C:\Partages

---

## Firewall

- pfSense
- Windows Defender Firewall

---

## Démonstration

1. Se connecter au serveur
2. Créer un fichier dans :

C:\Partage

3. Attendre l'exécution automatique
4. Vérifier sa présence sur le client