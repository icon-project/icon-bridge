import * as dotenv from "dotenv";
import { TezosToolkit } from "@taquito/taquito";
import { InMemorySigner } from "@taquito/signer";
import { NETWORK } from "../config/config";

dotenv.config();

const deploy = async () => {
  type networkTypes = keyof typeof NETWORK;

  const { ORIGINATOR_PRIVATE_KEY } = process.env;
  let file = process.argv[2];
  let rpcType: networkTypes = process.argv[3]
    ?.toUpperCase()
    .slice(1) as networkTypes;

  console.info("Staring validation...");

  if (!ORIGINATOR_PRIVATE_KEY) {
    console.error("Private key missing in the environment.");
    return;
  }

  if (!file) {
    console.log("Pass filename to deploy! Read the docs.");
    return;
  }

  if (!rpcType || !Object.keys(NETWORK).includes(rpcType)) {
    console.log("Valid networks:", Object.keys(NETWORK));
    console.error("Invalid network");
    return;
  }

  file = file.toLowerCase();

  const signer = await InMemorySigner.fromSecretKey(ORIGINATOR_PRIVATE_KEY);
  const Tezos = new TezosToolkit(NETWORK[rpcType].url);
  Tezos.setProvider({ signer: signer });

  console.log(`Validation checks out... \nDeploying the ${file} contract..\n`);
  try {
    const { hash, contractAddress } = await Tezos.contract.originate({
      code: require(`../contracts/build/${file}.json`),
      init: require(`../contracts/build/${file}_storage.json`),
    });

    console.log("Successfully deployed contract");
    console.log(`>> Transaction hash: ${hash}`);
    console.log(`>> Contract address: ${contractAddress}`);
  } catch (error) {
    console.log(
      `Oops... Deployment faced an issue.\nHere's the detailed info about the error,\n${error}`
    );

    console.log(
      "Make sure of these things. \n1. Compiled contracts are inside the build folder.\n2. You have enough Tezos balance in the wallet for deployment."
    );
  }
};

deploy();
