# lab2324omada7

Συστημα καταχωρησης κριτικων ταινιων.

## Stack

Go + Chi για routing (μέσω go-blueprint)

React για frontend

## Οδηγιες:

Αρχικά, θα χρειαστεις μια MySQL (5.6+) βάση δεδομένων.

Στην συνέχεια θα πρέπει να δημιουργήσεις το αρχείο: ``.env`` σύμφωνα με το ``.env.example``

## Make

Παρέχεται ένα Makefile, μέσα από αυτό μπορείς να φτιάξεις και να τρέξεις το docker container.

Build docker container:
``
make docker-build
``

Run the docker container:
``
make docker-run
``

Removes old binaries:
``
make clean
``
