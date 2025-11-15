[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/3aLMp1Wb)
# Network Security Laboratory — Deliverable 3

This repository contains the minimal setup to start working on the assignment.

## How to use this repository

Welcome to the base repository for Lab 3. This repository provides:

- A minimal Go package definition.
- A lightweight source skeleton for the requested service.
- A `Makefile` with some helpers.

After you understand every piece of code provided in this repository, you should "make this
repository yours": that means updating this `README.md`.

Even if you use the provided helpers, explain here how your software is built and executed.

If you decide to use a different library than the example, you're encouraged to do so;
just remember to mention it in the `README.md` and provide links to its documentation.

## Helpers

- `Makefile`: it includes the build instruction for the main program. If you add packages
  to the solution or rename the existing one, update the Makefile accordingly.

## Helpers for certificates

This deliverable **requires** the usage of a self-signed certificate. In order to generate it and
make it usable in your environment, the `Makefile` includes some rules that generate and install
the needed certificates in the appropriate paths of your OS.

The code in `cmd/mydb/main.go` assumes the certificates are in the locations created by
the provided `Makefile`. Adapt the code if you make any changes.

**DISCLAIMER**

The certificate generation and installation commands assume a
**Debian GNU/Linux–based distribution**. If you use a different OS, adapt the commands
accordingly, since some commands may differ.
