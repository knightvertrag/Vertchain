# Doge's Pub

The idea of a blockchain powered closed economy with power in the hands of the stakeholders. Tokens give more power to members, including voting rights on new events, menu items and infrastructure.

## Build Instructions

1. Download and install go `https://golang.org/dl/`

2. Clone the repository

```
git clone https://github.com/knightvertrag/Vertchain.git
```

3. Build the CLI binary using 

```bash
go build ./cmd/dop
```

4. Run the server using 

```bash
./dop run --datadir="path/to/database/directory"
```
Display the available commands using 

```bash
./dop help
```