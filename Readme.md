# Taiko Node Starter UI

1. Clone the repository `git clone https://github.com/taikoxyz/simple-taiko-node.git`
2. Fetch simple-taiko-node submodule `git submodule update --init`
3. Copy env variable `cp simple-taiko-node/.env.sample simple-taiko-node/.env`
4. Update required variables in `.env` file `L1_ENDPOINT_HTTP=` and `L1_ENDPOINT_WS=` use [chianlist.com](https://chainlist.org/chain/17000) to gain apis.
5. Run the server `go run main.go`
