
### Quickstart

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
### View Documents
```
./daoctl get documents
```
### View all Documents of a Type
```
./daoctl get documents --type assignment
```
### View a Specific Document and Subgraph
```
./daoctl get document <hash>
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


## Treasury Commands

Submitting a new payment against a Redemption Request 
```bash
# in the below command, the redemption_id is 6
./daoctl --vault-file hyphanewyork.json treasury newpayment 6 "3500.00 HUSD" --network BTC --trxid b475e94c6a86dd18cce0ab7a1dfc9d0f94e20baf6c91317c14ce669da4111e1c --memo "just a memo field for any additional context"
```

Attesting to an existing payment. Treasurers have the responsibility of reviewing and validating transaction payments. To attest that a posted payment is true and accurate, use the ```attest``` command.
```bash
# in the below example, the paymentID is 4 and the requestID is 6
./daoctl --vault-file hyphanewyork.json treasury attest 4 6 "3500.00 HUSD" 
```

## Multisig Deployment Proposals
```
DEBUG=true ./daoctl propose deployment create --proposal-name testprop --commit d431c59dfd0fe284eee979965160fd326cae0e73 --developer hyphanewyork --notes "this is a test deployment proposal" --account dao.hypha --config daoctl-test.yaml --vault-file ../m.hypha.json 
```