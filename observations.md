# 🧠 Retour d'expérience : Déploiement d'un smart contract pour virements B2B

## ✅ Ce qui a fonctionné
- **Workflow métier** :  
  La séparation `pending → approved → executed` correspond parfaitement aux processus de trésorerie (ex : validation par le service finance).
- **Traçabilité immuable** :  
  Chaque appel (`InitTransfer`, `ApproveTransfer`) est horodaté et stocké sur la blockchain → **audit simplifié**.
- **Accès à une technologie enterprise sans coût initial** :  
  Hyperledger Fabric est **open-source et gratuit**, contrairement aux solutions managées (AWS Managed Blockchain, IBM Blockchain Platform).  
  → J’ai pu **tester un cas métier concret** en 3 semaines grâce aux **tutoriels officiels pré-configurés** (ex : `test-network`), sans investissement financier.  
  *À noter* : Ce POC reste **simplifié** (1 ordinateur, 2 peers en Docker), mais il montre qu’une entreprise peut **évaluer la pertinence de la blockchain sans budget** avant un déploiement client.

## ⚠️ Défis rencontrés
- **Gestion des identités et permissions** :  
  Dans le réseau de test (`test-network`), chaque peer/orderer utilise un **système de certificats (MSP)** distinct.  
  → J’ai dû comprendre comment les **rôles (admin, client)** et les **endorsement policies** influencent les appels au chaincode.  
  *Exemple concret* : Une erreur `ENDORSEMENT_POLICY_FAILURE` lors de `ApproveTransfer` car l’identité n’avait pas les droits.  
  **Solution** : Utiliser `export CORE_PEER_TLS_ROOTCERT_FILE=...` pour charger le bon certificat avant chaque commande.
- **Limitation métier** :  
  La blockchain ne gère **pas les règles métier dynamiques**. Toute règle insérée dans le code est immuable.  
  *Exemple* : Impossible de coder *"2 approbations si montant > 10k€"* directement dans le chaincode → nécessite un **middleware** (ex : API Gateway avec règles configurables).
- **Questions clés pour un déploiement client** :  
  Ce projet m’a permis d’identifier **les bonnes questions dimensionnant ce type de solution** :  
  - *Évolutivité* : "Combien de transactions par seconde sont à prévoir ? Vérifier si le nombre concorde avec la limite supportée par hyperledger"  
  - *Correction d’erreurs* : "Comment gérer un bug dans le Smart Contract ? Solutions : mise à jour via `upgradeCC` ou middleware"  
  - *Stockage* : "Quel volume de données stockerez-vous hors-chain ? La blockchain ne doit pas contenir de fichiers lourds"  


## 💡 Insights clés
| Leçon apprise | Application métier |
|---------------|---------------------|
| **La blockchain n'est pas une base de données** | À réserver aux cas où l'**immuabilité** et la **confiance décentralisée** sont critiques (ex : KYC partagé entre banques). |
| **Les "smart contracts" ne sont pas si smart** | Ils exécutent des règles prédéfinies → **pas de remplacement** pour les processus métier complexes (nécessite un orchestrateur). |
| **Valeur réelle = gouvernance** | L'apport majeur est la **traçabilité des décisions** (qui a approuvé, quand ?), pas la technologie en elle-même. |

## 🔮 Prochaines étapes
1. Explorer l'**intégration avec les middleware** (ex : API Gateway) pour connecter la blockchain aux ERP clients.
2. Participer à un **atelier réglementaire** sur l'intégration de la blockchain dans les systèmes informatique des grands groupes
3. Documenter un **framework de décision** : *"Quand utiliser la blockchain ?"* pour les clients.
