# Web3GoTokio
This is a project to test how the Web3Go (Go-Ethereum) is working to integrate a Smart Contract into a common application.

In this case, the Smart contract used is based in **ERC-20 token**, but the process to integrate any of this kind of applications is exactly the same.

Once the Smart Contract is ready in Solidity language, the process to transform the code and integrate in any code language is the following:

```shell
solc --abi TokioToken.sol

abigen --abi=TokioToken.abi --pkg=contract --out=TokioToken.go

solc --bin TokioToken.sol

abigen --bin=TokioToken.bin --abi=TokioToken.abi --pkg=contract --out=TokioToken.go
```

At the end of this process a new file is extracted named TokioToken.go. Inside the file, there is the SmartContract code based on GoLang language and a binary section where the Smart Contract is encoded to be deployed into the Blockchain network.
