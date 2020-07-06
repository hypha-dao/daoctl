
### Build

Install Go

Build the binary
```
git clone https://github.com/hypha-dao/daoctl/
cd daoctl
go build
./daoctl get roles
```

If you change the go code, you can run without rebuilding by running:
```
go run main.go [args]
```

## Sample daoctl.yaml
Here's a sample configuration file for the testnet. You can save as a file and use ```--config``` then the file name as a parameter.
```yaml
EosioEndpoint: https://test.telos.kitchen
AssetsAsFloat: true
DAOContract: dao.hypha
Treasury:
  TokenContract: husd.hypha
  Symbol: HUSD
  Contract: bank.hypha
  EthUSDTContract: 0xdac17f958d2ee523a2206206994597c13d831ec7
  EthUSDTAddress: 0xC20f453a4B4995CA032570f212988F4978B35dDd
  BtcAddress: 35hfgfaUouzYixTUDV6CFqMiTSZcuNtTf9
TelosDecideContract: trailservice
DAOUser: treasurermmm
HyperionEndpoint: https://testnet.telosusa.io/v2
```

## Commands
### View Roles
```
./daoctl get roles
```
#### Include proposals
```
./daoctl get roles --include-proposals
```
### View Assignments
```
./daoctl get payouts
```
### View Treasury
```
./daoctl get treasury
```
### View Treasury Redemption Requests
```
./daoctl treasury get requests
```
### View Treasury Payments (fulfilling requests)
```
./daoctl treasury get payments
```
### View Details of a Ballot
```
./daoctl get ballot d4
```

With much credit and appreciation to ```eosc``` at https://github.com/eoscanada/eosc
