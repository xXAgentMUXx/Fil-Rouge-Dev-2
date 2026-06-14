# Y-Plaza – Real Estate Platform

## Résumé

Y-Plaza est une application web permettant de gérer l’achat et la vente de biens immobiliers.

La plateforme permet :

- aux clients de consulter et acheter des biens immobiliers
- aux agents immobiliers de publier et gérer des propriétés
- aux administrateurs de superviser la plateforme

L’application est développée en **Go** avec une architecture backend orientée services et utilise une base de données **PostgreSQL**.

Le paiement des propriétés est réalisé via **Stripe Checkout**.

L'utilisation de connexions tierces comme **Google** ou **GitHub** est également disponible.

---

## Fonctionnalités

### Authentification

- inscription utilisateur
- connexion sécurisée
- hash des mots de passe avec **bcrypt**
- gestion des rôles utilisateurs
- connexion avec **Google** et **GitHub**

---

## Gestion des propriétés

Les agents peuvent :

- ajouter une propriété
- modifier une propriété
- supprimer une propriété
- vendre une propriété
- ajouter une image

Les utilisateurs peuvent :

- consulter les propriétés
- rechercher par ville
- filtrer par prix
- filtrer par surface
- acheter une propriété

---

## Paiement sécurisé

Les achats de propriétés sont réalisés via **Stripe**.

---

## Tableau de bord

Le dashboard permet de visualiser :

- le total des ventes
- les propriétés vendues
- les villes les plus populaires
- les propriétés les plus chères
- les derniers biens ajoutés

---

## Base de données

L’application utilise **PostgreSQL**.

Tables présentes :

- users
- properties
- payments
- sales
- agencies

---

## Commandes SQL

### Utilisateurs

Par défaut un utilisateur possède le rôle :

client

Pour modifier le rôle dans la base de données :

Passer un utilisateur en agent :


UPDATE users
SET role = 'agent'
WHERE id = 1;

Passer un utilisateur en administrateur :

UPDATE users
SET role = 'admin'
WHERE id = 1;

Voir les utilisateurs et leurs rôles :

SELECT email, role FROM users;

Ajouter une agence : 

INSERT INTO agencies (id, name, city)
VALUES ('550e8400-e29b-41d4-a716-446655440000', 'Agence Lyon Sud', 'Lyon');

## Installation avec Docker

Cloner le projet :

https://github.com/xXAgentMUXx/Fil-Rouge-Dev-2.git

Télécharger et installer docker docker : 

https://www.docker.com/products/docker-desktop/

Lancer l’application :

sur le terminale tapez :


docker compose up --build

Cela lance :

le backend Go

la base de données PostgreSQL

Note :

Si vous voulez modifier la base de données postgres sql aprés avoir lancé l'application sur un nouveau terminale tapez : 

docker exec -it postgres_real_estate psql -U postgres -d realestate

Accéder au site web :

Puis maintenant vous pouvez accéder au site web depuis :

http://localhost:8080


## Technologies utilisées

Backend :

Go

Base de données :

PostgreSQL

Paiement :

Stripe

Conteneurisation :

Docker

## Auteur

Projet réalisé dans le cadre du Bachelor Informatique – Ynov par :

Axel Macé

Vittorio Gandossi

Mathys Urban
