# go-fil-markets
[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![CircleCI](https://circleci.com/gh/chenjianmei111/go-fil-markets.svg?style=svg)](https://circleci.com/gh/chenjianmei111/go-fil-markets)
[![codecov](https://codecov.io/gh/chenjianmei111/go-fil-markets/branch/master/graph/badge.svg)](https://codecov.io/gh/chenjianmei111/go-fil-markets)
[![GoDoc](https://godoc.org/github.com/chenjianmei111/go-fil-markets?status.svg)](https://godoc.org/github.com/chenjianmei111/go-fil-markets)

This repository contains modular implementations of the [storage and retrieval market subsystems](https://chenjianmei111.github.io/specs/#systems__filecoin_markets) of Filecoin. 
They are guided by the [v1.0 and 1.1 Filecoin specification updates](https://chenjianmei111.github.io/specs/#intro__changelog). 

Separating implementations into a blockchain component and one or more mining and market components presents an opportunity to encourage implementation diversity while reusing non-security-critical components.

## Disclaimer: Reporting issues

This repo shared the issue tracker with lotus. Please report your issues at the [lotus issue tracker](https://github.com/chenjianmei111/lotus/issues)

## Components

* **[storagemarket](./storagemarket)**: for finding, negotiating, and consummating deals to
 store data between clients and providers (storage miners).
* **[retrievalmarket](./retrievalmarket)**: for finding, negotiating, and consummating deals to
 retrieve data between clients and providers (retrieval miners).
* **[filestore](./filestore)**: a wrapper around os.File for use by pieceio, storagemarket, and retrievalmarket.
* **[pieceio](./pieceio)**: utilities that take IPLD graphs and turn them into pieces. Used by storagemarket.
* **[piecestore](./piecestore)**:  a database for storing deal-related PieceInfo and CIDInfo. 
Used by storagemarket and retrievalmarket.

Related components in other repos:
* **[go-data-transfer](https://github.com/chenjianmei111/go-data-transfer)**: for exchanging piece data between clients and miners, used by storage & retrieval market modules.

### Background reading

* The [Markets in Filecoin](https://chenjianmei111.github.io/specs/#systems__filecoin_markets) 
section of the Filecoin Specification contains the canonical spec.

### Technical Documentation
* [GoDoc for Storage Market](https://godoc.org/github.com/chenjianmei111/go-fil-markets/storagemarket) contains an architectural overview and robust API documentation
* [GoDoc for Retrieval Market](https://godoc.org/github.com/chenjianmei111/go-fil-markets/retrievalmarket) contains an architectural overview and robust API documentation

## Installation
```bash
go get "github.com/chenjianmei111/go-fil-markets/<MODULENAME>"`
```

## Usage
Documentation is in the README for each module, listed in [Components](#Components).

## Contributing
Issues and PRs are welcome! Please first read the [background reading](#background-reading) and [CONTRIBUTING](.go-fil-markets/CONTRIBUTING.md) guide, and look over the current code. PRs against master require approval of at least two maintainers. 

Day-to-day discussion takes place in the #fil-components channel of the [Filecoin project chat](https://github.com/chenjianmei111/community#chat). Usage or design questions are welcome.

## Project-level documentation
The chenjianmei111 has a [community repo](https://github.com/chenjianmei111/community) with more detail about our resources and policies, such as the [Code of Conduct](https://github.com/chenjianmei111/community/blob/master/CODE_OF_CONDUCT.md).

## License
This repository is dual-licensed under Apache 2.0 and MIT terms.

Copyright 2019. Protocol Labs, Inc.
