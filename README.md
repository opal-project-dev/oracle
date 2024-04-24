# Opal Deployment Guide

The Opal deployment is not that combersome as the Platypus project, however still requires some elbow grease.

## 1. Deploy the mock chainlink oracle (do not use in prod)

1. Generate the Ethereum key pair (you can use Metamask for that) for the oracle wallet.
2. Deploy the [smart contract](https://github.com/opal-project-dev/oracle/blob/main/internal/service/solidity/contracts.sol) through Remix (just copy-paste the code). Save the address.
3. Transfer the ownership of the smart contract to the oracle wallet.
4. Build the oracle service via Docker.
5. Update the config [here](https://github.com/opal-project-dev/developer-edition/blob/main/config/config.yaml). Change the `internal_address` to the address of the smart contract and provide the `wallet` keys.

The oracle service fetches Bitcoin prices from Coingecko and pushes them to the oracle smart contract. The opal project then uses this feed to calculate the interest rates.

Don't forget to fund the oracle wallet.

## 2. Deploy the smart contracts

Here are the steps to deploy the smart contracts:

1. Locate the smart contracts package [here](https://github.com/opal-project-dev/fluidity/tree/main/packages/contracts).
2. Run `npm install`.
3. Update the `hardhat.config.js` with the new piccadilly RPC URL.
4. Create `.env` file referencing the `.env.example`.
5. Update the config file [here](https://github.com/opal-project-dev/fluidity/blob/main/packages/contracts/mainnetDeployment/deploymentParams.mainnet.js). `CHAINLINK_AUTUSD_PROXY` is the mock oracle from the previous step. `GENERAL_SAFE` and `OPL_SAFE` are the addresses that initially receive the OPL tokens. `OPL_SAFE` won't be able to transfer the tokens for 1 year. The `DEPLOYER` is the address of the deployer.
6. To deploy the contracts run `npx hardhat run mainnetDeployment/mainnetDeployment.js --network piccadilly`. Save the addresses.

## 3. Deploy the front end

You will need to update the configs for the `dev-frontend` and `lib-ethers` packages.

1. Locate the front end package [here](https://github.com/opal-project-dev/fluidity/tree/main/packages/dev-frontend).
2. Update all the necessary addresses and network chainid as in the [example](https://github.com/opal-project-dev/fluidity/commit/6f2964271ae38f62f2d7a26e673092c06929d167).
3. Build the front end via `yarn build`.

## Disclaimer

GLHF!
