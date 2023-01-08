const MockBTSPeriphery = artifacts.require("MockBTSPeriphery");
const BTSPeriphery = artifacts.require("BTSPeriphery");
const BTSCore = artifacts.require("MockBTSCore");
const BMC = artifacts.require("MockBMC");
const Holder = artifacts.require("Holder");
const NotPayable = artifacts.require("NotPayable");
const NonRefundable = artifacts.require("NonRefundable");
const Refundable = artifacts.require("Refundable");
const EncodeMsg = artifacts.require("EncodeMessage");
const {assert, AssertionError} = require("chai");
const truffleAssert = require("truffle-assertions");
const rlp = require("rlp");

let toHex = (buf) => {
    buf = buf.toString("hex");
    if (buf.substring(0, 2) == "0x") return buf;
    return "0x" + buf.toString("hex");
};

contract("PRA BTSCore Query and Management", (accounts) => {
    let bts_core, bts_periphery;
    let _native = "PARA";
    let service = "bts";
    let _net = "1234.iconee";
    let _bmcICON = "btp://1234.iconee/0x1234567812345678";
    let REPONSE_HANDLE_SERVICE = 2;
    let RC_OK = 0;
    let _fee = 10;
    let _fixed_fee = 500000;
    before(async () => {
        bmc = await BMC.new("1234.pra");
        bts_core = await BTSCore.new();
        bts_periphery = await BTSPeriphery.new();
        encode_msg = await EncodeMsg.new();
        await bts_core.initialize(_native, _fee, _fixed_fee);
        await bts_periphery.initialize(bmc.address, bts_core.address);
        await bmc.addService(service, bts_periphery.address);
        await bmc.addLink(_bmcICON);
    });

    it(`Scenario 1: Contract's owner to register a new coin`, async () => {
        let _name = "ICON";
        await bts_core.updateBTSPeriphery(bts_periphery.address);
        await bts_core.register(_name, "", 18, 100, 0, "0x0000000000000000000000000000000000000000");
        output = await bts_core.coinNames();
        assert(output[0] === _native && output[1] === "ICON");
    });

    it("Scenario 2: Non-ownership role client registers a new coin", async () => {
        let _name = "TRON";
        await truffleAssert.reverts(
            bts_core.register.call(_name, "", 18, 100, 0, "0x0000000000000000000000000000000000000000", {from: accounts[1]}),
            "Unauthorized"
        );
    });

    it("Scenario 3: Contract’s owner registers an existed coin", async () => {
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_core.register.call(_name, "", 18, 100, 0, "0x0000000000000000000000000000000000000000"),
            "ExistCoin"
        );
    });

    it("Scenario 4: Contract’s owner to update BTSPeriphery contract", async () => {
        await bts_core.updateBTSPeriphery(bts_periphery.address);
    });

    it("Scenario 5: Non-ownership role client updates BTSPeriphery contract", async () => {
        await truffleAssert.reverts(
            bts_core.updateBTSPeriphery.call(accounts[2], {
                from: accounts[1],
            }),
            "Unauthorized"
        );
    });

    it("Scenario 6: Contract’s owner updates BTSPeriphery while this contract has pending requests", async () => {
        let _to = "btp://1234.iconee/0x12345678";
        await bts_core.transferNativeCoin(_to, {
            from: accounts[0],
            value: 100000000,
        });
        await truffleAssert.reverts(
            bts_core.updateBTSPeriphery.call(accounts[2]),
            "HasPendingRequest"
        );
        //  Clear pending request
        let _msg = await encode_msg.encodeResponseMsg(
            REPONSE_HANDLE_SERVICE,
            RC_OK,
            ""
        );
        await bmc.receiveResponse(_net, service, 1, _msg);
    });

    it(`Scenario 9: Contract's owner updates fee ratio`, async () => {
        let new_fee = 20;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_core.setFeeRatio(_name, _feeNumerator, new_fee);

        let fees = await bts_core.feeRatio(_name);
        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fee
        );
    });

    it("Scenario 10: Non-ownership role client updates fee ratio", async () => {
        let old_fee = 20;
        let new_fee = 50;
        let _feeNumerator = 100;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_core.setFeeRatio.call(_name, _feeNumerator, new_fee, {from: accounts[1]}),
            "Unauthorized"
        );

        let fees = await bts_core.feeRatio(_name);
        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === old_fee
        );
    });

    it("Scenario 11: Fee Numerator is set higher than Fee Denominator", async () => {
        let old_feeNumerator = 100;
        let new_feeNumerator = 15000;
        let new_fee = 15000;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_core.setFeeRatio.call(_name, new_feeNumerator, new_fee),
            "InvalidSetting"
        );

        let fees = await bts_core.feeRatio(_name);

        assert(
            web3.utils.BN(fees._feeNumerator).toNumber() === old_feeNumerator
        );
    });

    it("Scenario 12: Contract owner updates fixed fee", async () => {
        let new_fixed_fee = 1000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_core.setFeeRatio(_name, _feeNumerator, new_fixed_fee);
        let fees = await bts_core.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            new_fixed_fee
        );
    });

    it("Scenario 13: Non-ownership role client updates fixed fee", async () => {
        let old_fixed_fee = 1000000;
        let new_fixed_fee = 2000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_core.setFeeRatio.call(_name, _feeNumerator, new_fixed_fee, {from: accounts[1]}),
            "Unauthorized"
        );
        let fees = await bts_core.feeRatio(_name);
        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            old_fixed_fee
        );
    });


    it("Scenario 15: Query a valid supporting coin", async () => {
        let _name1 = "wBTC";
        let _name2 = "Ethereum";
        await bts_core.register(_name1, "", 18, 100, 0, "0xDf1930A268e204c24fAA25E7E72D26166551F933");
        await bts_core.register(_name2, "", 18, 100, 0, "0xDf1930A268e204c24fAA25E7E72D26166551F933");

        let _query = "ICON";
        let result = await bts_core.coinId(_query);
        assert(
            web3.utils.toChecksumAddress(result) !==
            web3.utils.toChecksumAddress(
                "0x96EdA576D1Bd2016C8beb15bC7e741D7B3552D45"
            )
        );
    });

    it("Scenario 16: Query an invalid supporting coin", async () => {
        let _query = "EOS";
        let result = await bts_core.coinId(_query);
        assert(web3.utils.toChecksumAddress(result) ===
                    web3.utils.toChecksumAddress(
                        "0x0000000000000000000000000000000000000000"
                    ));
    });

    it("Scenario 17: Non-Owner tries to add a new Owner", async () => {
        let oldList = await bts_core.getOwners();
        try {
            await bts_core.addOwner(accounts[1], {from: accounts[2]});
        } catch (err) {
            assert(err, "exited with an error (status 0)");
        }
        let newList = await bts_core.getOwners();
        assert(
            oldList.length === 1 &&
            oldList[0] === accounts[0] &&
            newList.length === 1 &&
            newList[0] === accounts[0]
        );
    });

    it("Scenario 18: Current Owner adds a new Owner", async () => {
        let oldList = await bts_core.getOwners();
        await bts_core.addOwner(accounts[1]);
        let newList = await bts_core.getOwners();
        assert(
            oldList.length === 1 &&
            oldList[0] === accounts[0] &&
            newList.length === 2 &&
            newList[0] === accounts[0] &&
            newList[1] === accounts[1]
        );
    });

    it("Scenario 19: After adding a new Owner, owner registers a new coin", async () => {
        let _name3 = "TRON";
        await bts_core.register(_name3, "", 18, 100, 0, "0xDf1930A268e204c24fAA25E7E72D26166551F933");
        output = await bts_core.coinNames();
        console.log(output);
        assert(
            output[0] === _native &&
            output[1] === "ICON" &&
            output[2] === "wBTC" &&
            output[3] === "Ethereum" &&
            output[4] === "TRON"
        );
    });

    it("Scenario 20: New Owner registers a new coin", async () => {
        let _name3 = "BINANCE";
        await bts_core.register(_name3, "", 18, 100, 0, "0xDf1930A268e204c24fAA25E7E72D26166551F933", {from: accounts[1]});
        output = await bts_core.coinNames();
        assert(
            output[0] === _native &&
            output[1] === "ICON" &&
            output[2] === "wBTC" &&
            output[3] === "Ethereum" &&
            output[4] === "TRON" &&
            output[5] === "BINANCE"
        );
    });

    it("Scenario 21: New owner updates BTSPeriphery contract", async () => {
        let newBTSPeriphery = await BTSPeriphery.new();
        await bts_core.updateBTSPeriphery(newBTSPeriphery.address, {
            from: accounts[1],
        });
    });

    it("Scenario 22: Old owner updates BTSPeriphery contract", async () => {
        let newBTSPeriphery = await BTSPeriphery.new();
        await bts_core.updateBTSPeriphery(newBTSPeriphery.address, {
            from: accounts[0],
        });
    });

    it("Scenario 25: New owner updates new fee ratio", async () => {
        let new_fee = 30;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_core.setFeeRatio(_name, _feeNumerator, new_fee, {from: accounts[1]});
        let fees = await bts_core.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fee
        );
    });

    it("Scenario 26: Old owner updates new fee ratio - After adding new Owner", async () => {
        let new_fee = 40;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_core.setFeeRatio(_name, _feeNumerator, new_fee);
        let fees = await bts_core.feeRatio(_name);
        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fee
        );
    });

    it("Scenario 27: New owner updates new fixed fee", async () => {
        let new_fixed_fee = 3000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_core.setFeeRatio(_name, _feeNumerator, new_fixed_fee, {from: accounts[1]});
        let fees = await bts_core.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            new_fixed_fee
        );
    });

    it("Scenario 28: Old owner updates new fixed fee - After adding new Owner", async () => {
        let new_fixed_fee = 4000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_core.setFeeRatio(_name, _feeNumerator, new_fixed_fee);
        let fees = await bts_core.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            new_fixed_fee
        );
    });

    it("Scenario 29: Non-Owner tries to remove an Owner", async () => {
        let oldList = await bts_core.getOwners();
        await truffleAssert.reverts(
            bts_core.removeOwner.call(accounts[0], {from: accounts[2]}),
            "Unauthorized"
        );
        let newList = await bts_core.getOwners();
        assert(
            oldList.length === 2 &&
            oldList[0] === accounts[0] &&
            oldList[1] === accounts[1] &&
            newList.length === 2 &&
            newList[0] === accounts[0] &&
            newList[1] === accounts[1]
        );
    });

    it("Scenario 30: Current Owner removes another Owner", async () => {
        let oldList = await bts_core.getOwners();
        await bts_core.removeOwner(accounts[0], {from: accounts[1]});
        let newList = await bts_core.getOwners();
        assert(
            oldList.length === 2 &&
            oldList[0] === accounts[0] &&
            oldList[1] === accounts[1] &&
            newList.length === 1 &&
            newList[0] === accounts[1]
        );
    });

    it("Scenario 31: The last Owner removes him/herself", async () => {
        let oldList = await bts_core.getOwners();
        await truffleAssert.reverts(
            bts_core.removeOwner.call(accounts[1], {from: accounts[1]}),
            "CannotRemoveMinOwner"
        );
        let newList = await bts_core.getOwners();
        assert(
            oldList.length === 1 &&
            oldList[0] === accounts[1] &&
            newList.length === 1 &&
            newList[0] === accounts[1]
        );
    });

    it("Scenario 32: Removed Owner tries to register a new coin", async () => {
        let _name3 = "KYBER";
        await truffleAssert.reverts(
            bts_core.register.call(_name3, "", 18, 100, 0, "0x0000000000000000000000000000000000000000"),
            "Unauthorized"
        );
        output = await bts_core.coinNames();
        assert(
            output[0] === _native &&
            output[1] === "ICON" &&
            output[2] === "wBTC" &&
            output[3] === "Ethereum" &&
            output[4] === "TRON" &&
            output[5] === "BINANCE"
        );
    });

    it("Scenario 33: Removed Owner tries to update BTSPeriphery contract", async () => {
        await truffleAssert.reverts(
            bts_core.updateBTSPeriphery.call(accounts[3], {
                from: accounts[0],
            }),
            "Unauthorized"
        );
    });

    it("Scenario 35: Removed Owner tries to update new fee ratio", async () => {
        let new_fee = 30;
        let _feeNumerator = 100;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_core.setFeeRatio.call(_name, _feeNumerator, new_fee, {from: accounts[0]}),
            "Unauthorized"
        );
    });
});

contract("As a user, I want to send PRA to ICON blockchain", (accounts) => {
    let bts_periphery, bts_core, bmc, nonrefundable, refundable;
    let service = "bts";
    let _bmcICON = "btp://1234.iconee/0x1234567812345678";
    let _net = "1234.iconee";
    let _to = "btp://1234.iconee/0x12345678";
    let RC_OK = 0;
    let RC_ERR = 1;
    let _native = "PARA";
    let deposit = 1000000000000;
    let _fee = 10;
    let _fixed_fee = 500000;
    let REPONSE_HANDLE_SERVICE = 2;
    let _uri = "https://github.com/icon-project/icon-bridge";

    before(async () => {
        bts_periphery = await BTSPeriphery.new();
        bts_core = await BTSCore.new();
        bmc = await BMC.new("1234.pra");
        encode_msg = await EncodeMsg.new();
        await bts_core.initialize(_native, _fee, _fixed_fee);
        await bts_periphery.initialize(bmc.address, bts_core.address);
        await bts_core.updateBTSPeriphery(bts_periphery.address);
        nonrefundable = await NonRefundable.new();
        refundable = await Refundable.new();
        await bmc.addService(service, bts_periphery.address);
        await bmc.addLink(_bmcICON);
    });

    it("Scenario 1: Transferring native coins to an invalid BTP Address format", async () => {
        let invalid_destination = "1234.iconee/0x12345678";
        let amount = 600000;
        await truffleAssert.reverts(
            bts_core.transferNativeCoin.call(invalid_destination, {
                from: accounts[0],
                value: amount,
            }),
            "revert"
        );
        bts_coin_balance = await bts_core.balanceOf(
            bts_core.address,
            _native
        );
        account_balance = await bts_core.balanceOf(accounts[0], _native);
        assert(
            web3.utils.BN(bts_coin_balance._usableBalance).toNumber() === 0 &&
            web3.utils.BN(account_balance._lockedBalance).toNumber() === 0
        );
    });

    it("Scenario 2: Transferring zero coin", async () => {
        await truffleAssert.reverts(
            bts_core.transferNativeCoin.call(_to, {
                from: accounts[0],
                value: 0,
            }),
            "revert"
        );
    });

    it("Scenario 3: msg.value less than fixed_fee", async () => {
        //  fixed_fee = 500000;
        let amount = 100000;
        await truffleAssert.reverts(
            bts_core.transferNativeCoin.call(_to, {
                from: accounts[0],
                value: amount,
            }),
            "revert"
        );
    });

    // it("Scenario 4: Transferring to an invalid network/not supported network", async () => {
    //     let invalid_destination = "btp://1234.eos/0x12345678";
    //     let amount = 600000;
    //     await truffleAssert.reverts(
    //         bts_core.transferNativeCoin.call(invalid_destination, {
    //             from: accounts[1],
    //             value: amount,
    //         }),
    //         "LimitExceed"
    //     );
    // });

    it("Scenario 5: Account client transfers a valid native coin to a side chain", async () => {
        let amount = 600000;
        let account_balanceBefore = await bts_core.balanceOf(
            accounts[0],
            _native
        );
        let tx = await bts_core.transferNativeCoin(_to, {
            from: accounts[0],
            value: amount,
        });
        let account_balanceAfter = await bts_core.balanceOf(
            accounts[0],
            _native
        );
        let bts_coin_balance = await bts_core.balanceOf(
            bts_core.address,
            _native
        );
        let chargedFee = Math.floor(amount / 1000) + _fixed_fee;

        const transferEvents = await bts_periphery.getPastEvents(
            "TransferStart",
            {
                fromBlock: tx.receipt.blockNumber,
                toBlock: "latest",
            }
        );
        let event = transferEvents[0].returnValues;
        assert.equal(event._from, accounts[0]);
        assert.equal(event._to, _to);
        assert.equal(event._sn, 1);
        assert.equal(event._assetDetails.length, 1);
        assert.equal(event._assetDetails[0].coinName, "PARA");
        assert.equal(event._assetDetails[0].value, amount - chargedFee);
        assert.equal(event._assetDetails[0].fee, chargedFee);

        const linkStatus = await bmc.getStatus(_bmcICON);
        const bmcBtpAddress = await bmc.getBmcBtpAddress();

        const messageEvents = await bmc.getPastEvents("Message", {
            fromBlock: tx.receipt.blockNumber,
            toBlock: "latest",
        });
        event = messageEvents[0].returnValues;
        assert.equal(event._next, _bmcICON);
        assert.equal(event._seq, linkStatus.txSeq);

        const bmcMsg = rlp.decode(event._msg);

        assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
        assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
        assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
        assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 1);

        const ServiceMsg = rlp.decode(bmcMsg[4]);
        assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);

        const coinTransferMsg = rlp.decode(ServiceMsg[1]);
        assert.equal(
            web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
            accounts[0]
        );
        assert.equal(
            web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
            _to.split("/").slice(-1)[0]
        );
        assert.equal(
            web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
            _native
        );
        assert.equal(
            web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
            amount - chargedFee
        );

        assert(
            web3.utils.BN(bts_coin_balance._userBalance).toNumber() ===
            amount &&
            web3.utils
                .BN(account_balanceBefore._lockedBalance)
                .toNumber() === 0 &&
            web3.utils
                .BN(account_balanceAfter._lockedBalance)
                .toNumber() === amount
        );
    });

    it("Scenario 6: BTSPeriphery receives a successful response of a recent request", async () => {
        let amount = 600000;
        let account_balanceBefore = await bts_core.balanceOf(
            accounts[0],
            _native
        );
        let _msg = await encode_msg.encodeResponseMsg(
            REPONSE_HANDLE_SERVICE,
            RC_OK,
            ""
        );
        let tx = await bmc.receiveResponse(_net, service, 1, _msg);
        let account_balanceAfter = await bts_core.balanceOf(
            accounts[0],
            _native
        );
        let fees = await bts_core.getAccumulatedFees();

        const transferEvents = await bts_periphery.getPastEvents(
            "TransferEnd",
            {
                fromBlock: tx.receipt.blockNumber,
                toBlock: "latest",
            }
        );
        let event = transferEvents[0].returnValues;

        assert.equal(event._from, accounts[0]);
        assert.equal(event._sn, 1);
        assert.equal(event._code, 0);
        assert.equal(event._response, "");

        assert(
            fees[0].coinName === _native &&
            Number(fees[0].value) ===
            Math.floor(amount / 1000) + _fixed_fee &&
            web3.utils
                .BN(account_balanceBefore._lockedBalance)
                .toNumber() === amount &&
            web3.utils
                .BN(account_balanceAfter._lockedBalance)
                .toNumber() === 0
        );
    });

    // it("Scenario 7: BTSPeriphery receives an error response of a recent request", async () => {
    //     let amount = 600000;
    //     let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
    //     let account_balanceBefore = await bts_core.balanceOf(
    //         accounts[0],
    //         _native
    //     );
    //     let bts_coin_balance_before = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let _msg = await encode_msg.encodeResponseMsg(
    //         REPONSE_HANDLE_SERVICE,
    //         RC_ERR,
    //         ""
    //     );
    //     let tx = await bmc.receiveResponse(_net, service, 2, _msg);
    //     let account_balanceAfter = await bts_core.balanceOf(
    //         accounts[0],
    //         _native
    //     );
    //     let bts_coin_balance_after = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //
    //     const transferEvents = await bts_periphery.getPastEvents(
    //         "TransferEnd",
    //         {
    //             fromBlock: tx.receipt.blockNumber,
    //             toBlock: "latest",
    //         }
    //     );
    //     let event = transferEvents[0].returnValues;
    //
    //     assert.equal(event._from, accounts[0]);
    //     assert.equal(event._sn, 2);
    //     assert.equal(event._code, 1);
    //     assert.equal(event._response, "");
    //
    //     //  Unable to check balance of accounts[0] since this account has also paid gas fee
    //     //  It would be easier to check if this is a contract
    //     //  Requestor will be receive an amount of refund as
    //     //  refund = amount - chargeAmt
    //     assert(
    //         web3.utils.BN(account_balanceBefore._lockedBalance).toNumber() ===
    //         amount &&
    //         web3.utils
    //             .BN(account_balanceAfter._lockedBalance)
    //             .toNumber() === 0 &&
    //         web3.utils
    //             .BN(account_balanceAfter._refundableBalance)
    //             .toNumber() === 0 &&
    //         web3.utils
    //             .BN(bts_coin_balance_before._usableBalance)
    //             .toNumber() ===
    //         2 * amount &&
    //         web3.utils
    //             .BN(bts_coin_balance_after._usableBalance)
    //             .toNumber() ===
    //         amount + chargedFee
    //     );
    // });
    //
    // it("Scenario 8: Non-refundable contract transfers a valid native coin to a side chain", async () => {
    //     let amount = 600000;
    //     await nonrefundable.deposit({from: accounts[2], value: deposit});
    //     let contract_balanceBefore = await bts_core.balanceOf(
    //         nonrefundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_before = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let tx = await nonrefundable.transfer(bts_core.address, _to, amount);
    //     let contract_balanceAfter = await bts_core.balanceOf(
    //         nonrefundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_after = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
    //
    //     const transferEvents = await bts_periphery.getPastEvents(
    //         "TransferStart",
    //         {
    //             fromBlock: tx.receipt.blockNumber,
    //             toBlock: "latest",
    //         }
    //     );
    //     let event = transferEvents[0].returnValues;
    //     assert.equal(event._from, nonrefundable.address);
    //     assert.equal(event._to, _to);
    //     assert.equal(event._sn, 3);
    //     assert.equal(event._assetDetails.length, 1);
    //     assert.equal(event._assetDetails[0].coinName, "PARA");
    //     assert.equal(event._assetDetails[0].value, amount - chargedFee);
    //     assert.equal(event._assetDetails[0].fee, chargedFee);
    //
    //     const linkStatus = await bmc.getStatus(_bmcICON);
    //     const bmcBtpAddress = await bmc.getBmcBtpAddress();
    //
    //     const messageEvents = await bmc.getPastEvents("Message", {
    //         fromBlock: tx.receipt.blockNumber,
    //         toBlock: "latest",
    //     });
    //     event = messageEvents[0].returnValues;
    //     assert.equal(event._next, _bmcICON);
    //     assert.equal(event._seq, linkStatus.txSeq);
    //
    //     const bmcMsg = rlp.decode(event._msg);
    //
    //     assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
    //     assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
    //     assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
    //     assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 3);
    //
    //     const ServiceMsg = rlp.decode(bmcMsg[4]);
    //     assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
    //
    //     const coinTransferMsg = rlp.decode(ServiceMsg[1]);
    //     assert.equal(
    //         web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
    //         nonrefundable.address
    //     );
    //     assert.equal(
    //         web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
    //         _to.split("/").slice(-1)[0]
    //     );
    //     assert.equal(
    //         web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
    //         _native
    //     );
    //     assert.equal(
    //         web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
    //         amount - chargedFee
    //     );
    //
    //     assert(
    //         web3.utils.BN(contract_balanceBefore._usableBalance).toNumber() ===
    //         web3.utils.BN(contract_balanceAfter._usableBalance).toNumber() +
    //         amount &&
    //         web3.utils
    //             .BN(contract_balanceBefore._lockedBalance)
    //             .toNumber() === 0 &&
    //         web3.utils
    //             .BN(contract_balanceAfter._lockedBalance)
    //             .toNumber() === amount &&
    //         web3.utils
    //             .BN(bts_coin_balance_before._usableBalance)
    //             .toNumber() ===
    //         amount + chargedFee &&
    //         web3.utils
    //             .BN(bts_coin_balance_after._usableBalance)
    //             .toNumber() ===
    //         2 * amount + chargedFee
    //     );
    // });
    //
    // it(`Scenario 9: BTSPeriphery receives an error response of a recent request and fails to refund coins back to Non-refundable contract`, async () => {
    //     let amount = 600000;
    //     let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
    //     let contract_balanceBefore = await bts_core.balanceOf(
    //         nonrefundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_before = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let _msg = await encode_msg.encodeResponseMsg(
    //         REPONSE_HANDLE_SERVICE,
    //         RC_ERR,
    //         ""
    //     );
    //     let tx = await bmc.receiveResponse(_net, service, 3, _msg);
    //     let contract_balanceAfter = await bts_core.balanceOf(
    //         nonrefundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_after = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //
    //     const transferEvents = await bts_periphery.getPastEvents(
    //         "TransferEnd",
    //         {
    //             fromBlock: tx.receipt.blockNumber,
    //             toBlock: "latest",
    //         }
    //     );
    //     let event = transferEvents[0].returnValues;
    //
    //     assert.equal(event._from, nonrefundable.address);
    //     assert.equal(event._sn, 3);
    //     assert.equal(event._code, 1);
    //     assert.equal(event._response, "");
    //
    //     assert.equal(
    //         web3.utils.BN(contract_balanceBefore._lockedBalance).toNumber(),
    //         amount
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceAfter._lockedBalance).toNumber(),
    //         0
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceBefore._usableBalance).toNumber(),
    //         web3.utils.BN(contract_balanceAfter._usableBalance).toNumber()
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceAfter._refundableBalance).toNumber(),
    //         amount - chargedFee
    //     );
    //     assert.equal(
    //         web3.utils.BN(bts_coin_balance_before._usableBalance).toNumber(),
    //         web3.utils.BN(bts_coin_balance_after._usableBalance).toNumber()
    //     );
    // });
    //
    // it("Scenario 10: Refundable contract transfers a valid native coin to a side chain", async () => {
    //     let amount = 600000;
    //     await refundable.deposit({from: accounts[2], value: deposit});
    //     let contract_balanceBefore = await bts_core.balanceOf(
    //         refundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_before = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let tx = await refundable.transfer(bts_core.address, _to, amount);
    //     let contract_balanceAfter = await bts_core.balanceOf(
    //         refundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_after = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
    //
    //     const transferEvents = await bts_periphery.getPastEvents(
    //         "TransferStart",
    //         {
    //             fromBlock: tx.receipt.blockNumber,
    //             toBlock: "latest",
    //         }
    //     );
    //     let event = transferEvents[0].returnValues;
    //     assert.equal(event._from, refundable.address);
    //     assert.equal(event._to, _to);
    //     assert.equal(event._sn, 4);
    //     assert.equal(event._assetDetails.length, 1);
    //     assert.equal(event._assetDetails[0].coinName, "PARA");
    //     assert.equal(event._assetDetails[0].value, amount - chargedFee);
    //     assert.equal(event._assetDetails[0].fee, chargedFee);
    //
    //     const linkStatus = await bmc.getStatus(_bmcICON);
    //     const bmcBtpAddress = await bmc.getBmcBtpAddress();
    //
    //     const messageEvents = await bmc.getPastEvents("Message", {
    //         fromBlock: tx.receipt.blockNumber,
    //         toBlock: "latest",
    //     });
    //     event = messageEvents[0].returnValues;
    //     assert.equal(event._next, _bmcICON);
    //     assert.equal(event._seq, linkStatus.txSeq);
    //
    //     const bmcMsg = rlp.decode(event._msg);
    //
    //     assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
    //     assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
    //     assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
    //     assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 4);
    //
    //     const ServiceMsg = rlp.decode(bmcMsg[4]);
    //     assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
    //
    //     const coinTransferMsg = rlp.decode(ServiceMsg[1]);
    //     assert.equal(
    //         web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
    //         refundable.address
    //     );
    //     assert.equal(
    //         web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
    //         _to.split("/").slice(-1)[0]
    //     );
    //     assert.equal(
    //         web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
    //         _native
    //     );
    //     assert.equal(
    //         web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
    //         amount - chargedFee
    //     );
    //
    //     assert.equal(
    //         web3.utils.BN(contract_balanceBefore._usableBalance).toNumber(),
    //         web3.utils.BN(contract_balanceAfter._usableBalance).toNumber() +
    //         amount
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceBefore._lockedBalance).toNumber(),
    //         0
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceAfter._lockedBalance).toNumber(),
    //         amount
    //     );
    //     assert.equal(
    //         web3.utils.BN(bts_coin_balance_after._usableBalance).toNumber(),
    //         web3.utils.BN(bts_coin_balance_before._usableBalance).toNumber() +
    //         amount
    //     );
    // });
    //
    // it("Scenario 11: BTSPeriphery receives an error response of a recent request", async () => {
    //     let amount = 600000;
    //     let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
    //     let contract_balanceBefore = await bts_core.balanceOf(
    //         refundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_before = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //     let _msg = await encode_msg.encodeResponseMsg(
    //         REPONSE_HANDLE_SERVICE,
    //         RC_ERR,
    //         ""
    //     );
    //     let tx = await bmc.receiveResponse(_net, service, 4, _msg);
    //     let contract_balanceAfter = await bts_core.balanceOf(
    //         refundable.address,
    //         _native
    //     );
    //     let bts_coin_balance_after = await bts_core.balanceOf(
    //         bts_core.address,
    //         _native
    //     );
    //
    //     const transferEvents = await bts_periphery.getPastEvents(
    //         "TransferEnd",
    //         {
    //             fromBlock: tx.receipt.blockNumber,
    //             toBlock: "latest",
    //         }
    //     );
    //     let event = transferEvents[0].returnValues;
    //
    //     assert.equal(event._from, refundable.address);
    //     assert.equal(event._sn, 4);
    //     assert.equal(event._code, 1);
    //     assert.equal(event._response, "");
    //
    //     assert.equal(
    //         web3.utils.BN(contract_balanceBefore._lockedBalance).toNumber(),
    //         amount
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceAfter._lockedBalance).toNumber(),
    //         0
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceAfter._usableBalance).toNumber(),
    //         web3.utils.BN(contract_balanceBefore._usableBalance).toNumber() +
    //         amount -
    //         chargedFee
    //     );
    //     assert.equal(
    //         web3.utils.BN(contract_balanceAfter._refundableBalance).toNumber(),
    //         0
    //     );
    //     assert.equal(
    //         web3.utils.BN(bts_coin_balance_before._usableBalance).toNumber(),
    //         web3.utils.BN(bts_coin_balance_after._usableBalance).toNumber() +
    //         amount -
    //         chargedFee
    //     );
    // });
});
//
// contract("As a user, I want to send ERC1155_ICX to ICON blockchain", (accounts) => {
//     let bts_periphery, bts_core, bmc, holder;
//     let service = "Coin/WrappedCoin";
//     let _uri = "https://github.com/icon-project/icon-bridge";
//     let _native = "PARA";
//     let _fee = 10;
//     let _fixed_fee = 500000;
//     let _name = "ICON";
//     let _bmcICON = "btp://1234.iconee/0x1234567812345678";
//     let _net = "1234.iconee";
//     let _from = "0x12345678";
//     let _value = 999999999999999;
//     let REPONSE_HANDLE_SERVICE = 2;
//     let RC_OK = 0;
//     let RC_ERR = 1;
//
//     before(async () => {
//         bts_periphery = await BTSPeriphery.new();
//         bts_core = await BTSCore.new();
//         bmc = await BMC.new("1234.pra");
//         encode_msg = await EncodeMsg.new();
//         await bts_periphery.initialize(
//             bmc.address,
//             bts_core.address,
//             service
//         );
//         await bts_core.initialize(_native, _fee, _fixed_fee);
//         await bts_core.updateBTSPeriphery(bts_periphery.address);
//         holder = await Holder.new();
//         await bmc.addService(service, bts_periphery.address);
//         await bmc.addVerifier(_net, accounts[1]);
//         await bmc.addLink(_bmcICON);
//         await holder.addBSHContract(
//             bts_periphery.address,
//             bts_core.address
//         );
//         await bts_core.register(_name, "", 18);
//         let _msg = await encode_msg.encodeTransferMsgWithAddress(
//             _from,
//             holder.address,
//             _name,
//             _value
//         );
//         await bmc.receiveRequest(_bmcICON, "", service, 0, _msg);
//         id = await bts_core.coinId(_name, "", 18);
//     });
//
//     it("Scenario 1: User has not yet set approval for token being transferred out by Operator", async () => {
//         let _to = "btp://1234.iconee/0x12345678";
//         let _value = 600000;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_name, _value, _to),
//             "ERC1155: caller is not owner nor approved"
//         );
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//     });
//
//     it(`Scenario 2: User has set approval, but user's balance has insufficient amount`, async () => {
//         let _to = "btp://1234.iconee/0x12345678";
//         let _value = 9999999999999999n;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_name, _value, _to),
//             "ERC1155: insufficient balance for transfer"
//         );
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//     });
//
//     it("Scenario 3: User requests to transfer an invalid Token", async () => {
//         let _to = "btp://1234.iconee/0x12345678";
//         let _value = 9999999999999999n;
//         let _token = "EOS";
//         await holder.setApprove(bts_core.address);
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_token, _value, _to),
//             "UnregisterCoin"
//         );
//     });
//
//     it("Scenario 4: User transfers Tokens to an invalid BTP Address format", async () => {
//         let _to = "1234.iconee/0x12345678";
//         let amount = 600000;
//         let contract_balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_name, amount, _to),
//             "VM Exception while processing transaction: revert"
//         );
//         let contract_balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let bts_core_balance = await bts_core.balanceOf(
//             bts_core.address,
//             _name
//         );
//
//         assert.equal(
//             web3.utils.BN(contract_balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(contract_balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(contract_balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(contract_balanceBefore._usableBalance).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(bts_core_balance._usableBalance).toNumber(),
//             0
//         );
//     });
//
//     it("Scenario 5: User requests to transfer zero Token", async () => {
//         let _to = "1234.iconee/0x12345678";
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_name, 0, _to),
//             "revert"
//         );
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber()
//         );
//     });
//
//     it("Scenario 6: Transferring amount is less than fixed fee", async () => {
//         let _to = "1234.iconee/0x12345678";
//         let _name = "ICON";
//         let amount = 100000;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_name, amount, _to),
//             "revert"
//         );
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber()
//         );
//     });
//
//     it("Scenario 7: User requests to transfer to an invalid network/Not Supported Network", async () => {
//         let _to = "btp://1234.eos/0x12345678";
//         let _name = "ICON";
//         let amount = 600000;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         await truffleAssert.reverts(
//             holder.callTransfer.call(_name, amount, _to),
//             "BMCRevertNotExistsBMV"
//         );
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let bts_core_balance = await bts_core.balanceOf(
//             bts_core.address,
//             _name
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(bts_core_balance._usableBalance).toNumber(),
//             0
//         );
//     });
//
//     it("Scenario 8: User sends a valid transferring request", async () => {
//         let _to = "btp://1234.iconee/0x12345678";
//         let amount = 600000;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         let tx = await holder.callTransfer(_name, amount, _to);
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let bts_core_balance = await bts_core.balanceOf(
//             bts_core.address,
//             _name
//         );
//         let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferStart",
//             {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//         );
//         let event = transferEvents[0].returnValues;
//         assert.equal(event._from, holder.address);
//         assert.equal(event._to, _to);
//         assert.equal(event._sn, 1);
//         assert.equal(event._assetDetails.length, 1);
//         assert.equal(event._assetDetails[0].coinName, _name);
//         assert.equal(event._assetDetails[0].value, amount - chargedFee);
//         assert.equal(event._assetDetails[0].fee, chargedFee);
//
//         const linkStatus = await bmc.getStatus(_bmcICON);
//         const bmcBtpAddress = await bmc.getBmcBtpAddress();
//
//         const messageEvents = await bmc.getPastEvents("Message", {
//             fromBlock: tx.receipt.blockNumber,
//             toBlock: "latest",
//         });
//         event = messageEvents[0].returnValues;
//         assert.equal(event._next, _bmcICON);
//         assert.equal(event._seq, linkStatus.txSeq);
//
//         const bmcMsg = rlp.decode(event._msg);
//
//         assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
//         assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
//         assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
//         assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 1);
//
//         const ServiceMsg = rlp.decode(bmcMsg[4]);
//         assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
//
//         const coinTransferMsg = rlp.decode(ServiceMsg[1]);
//         assert.equal(
//             web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
//             holder.address
//         );
//         assert.equal(
//             web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
//             _to.split("/").slice(-1)[0]
//         );
//         assert.equal(
//             web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
//             _name
//         );
//         assert.equal(
//             web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
//             amount - chargedFee
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             amount
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber() - amount
//         );
//         assert.equal(
//             web3.utils.BN(bts_core_balance._usableBalance).toNumber(),
//             amount
//         );
//     });
//
//     it("Scenario 9: BTSPeriphery receives a successful response of a recent request", async () => {
//         let amount = 600000;
//         let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
//         let contract_balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let _msg = await encode_msg.encodeResponseMsg(
//             REPONSE_HANDLE_SERVICE,
//             RC_OK,
//             ""
//         );
//         let tx = await bmc.receiveResponse(_net, service, 1, _msg);
//         let contract_balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let fees = await bts_core.getAccumulatedFees();
//         let bts_core_balance = await bts_core.balanceOf(
//             bts_core.address,
//             _name
//         );
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferEnd",
//             {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//         );
//         let event = transferEvents[0].returnValues;
//
//         assert.equal(event._from, holder.address);
//         assert.equal(event._sn, 1);
//         assert.equal(event._code, 0);
//         assert.equal(event._response, "");
//
//         assert.equal(
//             web3.utils.BN(contract_balanceBefore._lockedBalance).toNumber(),
//             amount
//         );
//         assert.equal(
//             web3.utils.BN(contract_balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(contract_balanceBefore._usableBalance).toNumber(),
//             web3.utils.BN(contract_balanceAfter._usableBalance).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(bts_core_balance._usableBalance).toNumber(),
//             chargedFee
//         );
//         assert.equal(fees[1].coinName, _name);
//         assert.equal(Number(fees[1].value), chargedFee);
//     });
//
//     it("Scenario 8: User sends a valid transferring request", async () => {
//         let _to = "btp://1234.iconee/0x12345678";
//         let amount = 100000000000000;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let bts_core_balance_before = await bts_core.balanceOf(
//             bts_core.address,
//             _name
//         );
//         await holder.setApprove(bts_core.address);
//         let tx = await holder.callTransfer(_name, amount, _to);
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let bts_core_balance_after = await bts_core.balanceOf(
//             bts_core.address,
//             _name
//         );
//         let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferStart",
//             {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//         );
//         let event = transferEvents[0].returnValues;
//         assert.equal(event._from, holder.address);
//         assert.equal(event._to, _to);
//         assert.equal(event._sn, 2);
//         assert.equal(event._assetDetails.length, 1);
//         assert.equal(event._assetDetails[0].coinName, _name);
//         assert.equal(event._assetDetails[0].value, amount - chargedFee);
//         assert.equal(event._assetDetails[0].fee, chargedFee);
//
//         const linkStatus = await bmc.getStatus(_bmcICON);
//         const bmcBtpAddress = await bmc.getBmcBtpAddress();
//
//         const messageEvents = await bmc.getPastEvents("Message", {
//             fromBlock: tx.receipt.blockNumber,
//             toBlock: "latest",
//         });
//         event = messageEvents[0].returnValues;
//         assert.equal(event._next, _bmcICON);
//         assert.equal(event._seq, linkStatus.txSeq);
//
//         const bmcMsg = rlp.decode(event._msg);
//
//         assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
//         assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
//         assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
//         assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 2);
//
//         const ServiceMsg = rlp.decode(bmcMsg[4]);
//         assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
//
//         const coinTransferMsg = rlp.decode(ServiceMsg[1]);
//         assert.equal(
//             web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
//             holder.address
//         );
//         assert.equal(
//             web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
//             _to.split("/").slice(-1)[0]
//         );
//         assert.equal(
//             web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
//             _name
//         );
//         assert.equal(
//             web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
//             amount - chargedFee
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             amount
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber() - amount
//         );
//         assert.equal(
//             web3.utils.BN(bts_core_balance_after._usableBalance).toNumber(),
//             web3.utils
//                 .BN(bts_core_balance_before._usableBalance)
//                 .toNumber() + amount
//         );
//     });
//
//     it("Scenario 10: BTSPeriphery receives an error response of a recent request", async () => {
//         let amount = 100000000000000;
//         let chargedFee = Math.floor(amount / 1000) + _fixed_fee;
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//         let _msg = await encode_msg.encodeResponseMsg(
//             REPONSE_HANDLE_SERVICE,
//             RC_ERR,
//             ""
//         );
//         let tx = await bmc.receiveResponse(_net, service, 2, _msg);
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             _name
//         );
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferEnd",
//             {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//         );
//         let event = transferEvents[0].returnValues;
//
//         assert.equal(event._from, holder.address);
//         assert.equal(event._sn, 2);
//         assert.equal(event._code, 1);
//         assert.equal(event._response, "");
//
//         assert.equal(
//             web3.utils.BN(balanceBefore._lockedBalance).toNumber(),
//             amount
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._lockedBalance).toNumber(),
//             0
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._usableBalance).toNumber(),
//             web3.utils.BN(balanceBefore._usableBalance).toNumber() +
//             amount -
//             chargedFee
//         );
//         assert.equal(
//             web3.utils.BN(balanceAfter._refundableBalance).toNumber(),
//             0
//         );
//     });
// });
//
contract("As a user, I want to receive PRA from ICON blockchain", (accounts) => {
    let bmc, bts_periphery, bts_core, notpayable, refundable;
    let service = "bts";
    let _bmcICON = "btp://1234.iconee/0x1234567812345678";
    let _net = "1234.iconee";
    let _to = "btp://1234.iconee/0x12345678";
    let _native = "PARA";
    let _fee = 10;
    let _fixed_fee = 500000;
    let RC_ERR = 1;
    let RC_OK = 0;
    let _uri = "https://github.com/icon-project/icon-bridge";

    before(async () => {
        bts_periphery = await BTSPeriphery.new();
        bts_core = await BTSCore.new();
        bmc = await BMC.new("1234.pra");
        encode_msg = await EncodeMsg.new();
        await bts_core.initialize(_native, _fee, _fixed_fee);
        await bts_periphery.initialize(bmc.address, bts_core.address);
        await bts_core.updateBTSPeriphery(bts_periphery.address);
        notpayable = await NotPayable.new();
        refundable = await Refundable.new();
        await bmc.addService(service, bts_periphery.address);
        await bmc.addLink(_bmcICON);

        await bts_core.transferNativeCoin(_to, {
            from: accounts[0],
            value: 100000000,
        });
        btpAddr = await bmc.bmcAddress();
    });

    it("Scenario 1: Receiving address is invalid", async () => {
        let _from = "0x12345678";
        let _value = 1000;
        let _address = "0x1234567890123456789";
        let _eventMsg = await encode_msg.encodeResponseBMCMessage(
            btpAddr,
            _bmcICON,
            service,
            10,
            RC_ERR,
            "InvalidAddress"
        );
        let _msg = await encode_msg.encodeTransferMsgWithStringAddress(
            _from,
            _address,
            _native,
            _value
        );
        let output = await bmc.receiveRequest(
            _bmcICON,
            "",
            service,
            10,
            _msg
        );
        assert(
            output.logs[0].args._next === _bmcICON &&
            output.logs[0].args._msg === _eventMsg
        );
    });

    it("Scenario 2: BTSCore has insufficient funds to transfer", async () => {
        let _from = "0x12345678";
        let _value = 1000000000;
        let balanceBefore = await bmc.getBalance(accounts[1]);
        let _eventMsg = await encode_msg.encodeResponseBMCMessage(
            btpAddr,
            _bmcICON,
            service,
            10,
            RC_ERR,
            "TransferFailed"
        );
        let _msg = await encode_msg.encodeTransferMsgWithAddress(
            _from,
            accounts[1],
            _native,
            _value
        );
        let output = await bmc.receiveRequest(
            _bmcICON,
            "",
            service,
            10,
            _msg
        );
        let balanceAfter = await bmc.getBalance(accounts[1]);

        assert.equal(
            web3.utils.BN(balanceAfter).toString(),
            web3.utils.BN(balanceBefore).toString()
        );
        assert.equal(output.logs[0].args._next, _bmcICON);
        assert.equal(output.logs[0].args._msg, _eventMsg);
    });

    it(`Scenario 3: BTSCore tries to transfer PARA coins to a non-payable contract, but it fails`, async () => {
        let _from = "0x12345678";
        let _value = 1000;
        let balanceBefore = await bmc.getBalance(notpayable.address);
        let _eventMsg = await encode_msg.encodeResponseBMCMessage(
            btpAddr,
            _bmcICON,
            service,
            10,
            RC_ERR,
            "TransferFailed"
        );
        let _msg = await encode_msg.encodeTransferMsgWithAddress(
            _from,
            notpayable.address,
            _native,
            _value
        );
        let output = await bmc.receiveRequest(
            _bmcICON,
            "",
            service,
            10,
            _msg
        );
        let balanceAfter = await bmc.getBalance(notpayable.address);

        assert.equal(
            web3.utils.BN(balanceAfter).toNumber(),
            web3.utils.BN(balanceBefore).toNumber()
        );
        assert.equal(output.logs[0].args._next, _bmcICON);
        assert.equal(output.logs[0].args._msg, _eventMsg);
    });

    it("Scenario 4: BTSPeriphery receives a request of transferring coins", async () => {
        let _from = "0x12345678";
        let _value = 12345;
        let balanceBefore = await bmc.getBalance(accounts[1]);
        let _eventMsg = await encode_msg.encodeResponseBMCMessage(
            btpAddr,
            _bmcICON,
            service,
            10,
            RC_OK,
            ""
        );
        let _msg = await encode_msg.encodeTransferMsgWithAddress(
            _from,
            accounts[1],
            _native,
            _value
        );
        let output = await bmc.receiveRequest(
            _bmcICON,
            "",
            service,
            10,
            _msg
        );
        let balanceAfter = await bmc.getBalance(accounts[1]);

        assert.equal(
            web3.utils.BN(balanceAfter).toString(),
            web3.utils
                .BN(balanceBefore)
                .add(new web3.utils.BN(_value))
                .toString()
        );
        assert.equal(output.logs[0].args._next, _bmcICON);
        assert.equal(output.logs[0].args._msg, _eventMsg);
    });

    it(`Scenario 5: BTSPeriphery receives a request of transferring coins`, async () => {
        let _from = "0x12345678";
        let _value = 23456;
        let balanceBefore = await bmc.getBalance(refundable.address);
        let _eventMsg = await encode_msg.encodeResponseBMCMessage(
            btpAddr,
            _bmcICON,
            service,
            10,
            RC_OK,
            ""
        );
        let _msg = await encode_msg.encodeTransferMsgWithStringAddress(
            _from,
            refundable.address,
            _native,
            _value
        );
        let output = await bmc.receiveRequest(
            _bmcICON,
            "",
            service,
            10,
            _msg
        );
        let balanceAfter = await bmc.getBalance(refundable.address);

        assert.equal(
            web3.utils.BN(balanceAfter).toNumber(),
            web3.utils.BN(balanceBefore).toNumber() + _value
        );
        assert.equal(output.logs[0].args._next, _bmcICON);
        assert.equal(output.logs[0].args._msg, _eventMsg);
    });
});
//
// contract("As a user, I want to receive ERC1155_ICX from ICON blockchain", (accounts) => {
//     let bmc, bts_periphery, bts_core, holder, notpayable;
//     let service = "Coin/WrappedCoin";
//     let _uri = "https://github.com/icon-project/icon-bridge";
//     let _native = "PARA";
//     let _fee = 10;
//     let _fixed_fee = 500000;
//     let _name = "ICON";
//     let _bmcICON = "btp://1234.iconee/0x1234567812345678";
//     let _net = "1234.iconee";
//     let _from = "0x12345678";
//     let RC_ERR = 1;
//     let RC_OK = 0;
//
//     before(async () => {
//         bts_periphery = await BTSPeriphery.new();
//         bts_core = await BTSCore.new();
//         bmc = await BMC.new("1234.pra");
//         encode_msg = await EncodeMsg.new();
//         await bts_periphery.initialize(
//             bmc.address,
//             bts_core.address,
//             service
//         );
//         await bts_core.initialize(_native, _fee, _fixed_fee);
//         await bts_core.updateBTSPeriphery(bts_periphery.address);
//         holder = await Holder.new();
//         notpayable = await NotPayable.new();
//         await bmc.addService(service, bts_periphery.address);
//         await bmc.addVerifier(_net, accounts[1]);
//         await bmc.addLink(_bmcICON);
//         await holder.addBSHContract(
//             bts_periphery.address,
//             bts_core.address
//         );
//         await bts_core.register(_name, "", 18);
//         id = await bts_core.coinId(_name, "", 18);
//         btpAddr = await bmc.bmcAddress();
//     });
//
//     it("Scenario 1: Receiving address is invalid", async () => {
//         let _value = 1000;
//         let _address = "0x1234567890123456789";
//         let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//             btpAddr,
//             _bmcICON,
//             service,
//             10,
//             RC_ERR,
//             "InvalidAddress"
//         );
//         let _msg = await encode_msg.encodeTransferMsgWithStringAddress(
//             _from,
//             _address,
//             _name,
//             _value
//         );
//         let output = await bmc.receiveRequest(
//             _bmcICON,
//             "",
//             service,
//             10,
//             _msg
//         );
//
//         assert.equal(output.logs[0].args._next, _bmcICON);
//         assert.equal(output.logs[0].args._msg, _eventMsg);
//     });
//
//     it(`Scenario 2: Receiving contract does not implement ERC1155Holder / Receiver`, async () => {
//         let _value = 1000;
//         let balanceBefore = await bts_core.balanceOf(
//             notpayable.address,
//             id
//         );
//         let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//             btpAddr,
//             _bmcICON,
//             service,
//             10,
//             RC_ERR,
//             "TransferFailed"
//         );
//         let _msg = await encode_msg.encodeTransferMsgWithAddress(
//             _from,
//             notpayable.address,
//             _name,
//             _value
//         );
//         let output = await bmc.receiveRequest(
//             _bmcICON,
//             "",
//             service,
//             10,
//             _msg
//         );
//         let balanceAfter = await bts_core.balanceOf(notpayable.address, id);
//
//         assert.equal(
//             web3.utils.BN(balanceAfter).toNumber(),
//             web3.utils.BN(balanceBefore).toNumber()
//         );
//         assert.equal(output.logs[0].args._next, _bmcICON);
//         assert.equal(output.logs[0].args._msg, _eventMsg);
//     });
//
//     it("Scenario 3: BTSPeriphery receives a request of invalid token", async () => {
//         let _value = 3000;
//         let _tokenName = "Ethereum";
//         let invalid_coin_id = await bts_core.coinId(_tokenName);
//         let balanceBefore = await bts_core.balanceOf(
//             holder.address,
//             invalid_coin_id
//         );
//         let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//             btpAddr,
//             _bmcICON,
//             service,
//             10,
//             RC_ERR,
//             "UnregisteredCoin"
//         );
//         let _msg = await encode_msg.encodeTransferMsgWithAddress(
//             _from,
//             holder.address,
//             _tokenName,
//             _value
//         );
//         let output = await bmc.receiveRequest(
//             _bmcICON,
//             "",
//             service,
//             10,
//             _msg
//         );
//         let balanceAfter = await bts_core.balanceOf(
//             holder.address,
//             invalid_coin_id
//         );
//
//         assert.equal(
//             web3.utils.BN(balanceAfter).toNumber(),
//             web3.utils.BN(balanceBefore).toNumber()
//         );
//         assert.equal(output.logs[0].args._next, _bmcICON);
//         assert.equal(output.logs[0].args._msg, _eventMsg);
//     });
//
//     it("Scenario 4: Receiver is a ERC1155Holder contract", async () => {
//         let _value = 2500;
//         let balanceBefore = await bts_core.balanceOf(holder.address, id);
//         let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//             btpAddr,
//             _bmcICON,
//             service,
//             10,
//             RC_OK,
//             ""
//         );
//         let _msg = await encode_msg.encodeTransferMsgWithAddress(
//             _from,
//             holder.address,
//             _name,
//             _value
//         );
//         let output = await bmc.receiveRequest(
//             _bmcICON,
//             "",
//             service,
//             10,
//             _msg
//         );
//         let balanceAfter = await bts_core.balanceOf(holder.address, id);
//
//         assert.equal(
//             web3.utils.BN(balanceAfter).toNumber(),
//             web3.utils.BN(balanceBefore).toNumber() + _value
//         );
//         assert.equal(output.logs[0].args._next, _bmcICON);
//         assert.equal(output.logs[0].args._msg, _eventMsg);
//     });
//
//     it("Scenario 5: Receiver is an account client", async () => {
//         let _value = 5500;
//         let balanceBefore = await bts_core.balanceOf(accounts[1], id);
//         let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//             btpAddr,
//             _bmcICON,
//             service,
//             10,
//             RC_OK,
//             ""
//         );
//         let _msg = await encode_msg.encodeTransferMsgWithAddress(
//             _from,
//             accounts[1],
//             _name,
//             _value
//         );
//         let output = await bmc.receiveRequest(
//             _bmcICON,
//             "",
//             service,
//             10,
//             _msg
//         );
//         let balanceAfter = await bts_core.balanceOf(accounts[1], id);
//
//         assert.equal(
//             web3.utils.BN(balanceAfter).toNumber(),
//             web3.utils.BN(balanceBefore).toNumber() + _value
//         );
//         assert.equal(output.logs[0].args._next, _bmcICON);
//         assert.equal(output.logs[0].args._msg, _eventMsg);
//     });
// });
//
// contract("BSHs handle Gather Fee Service Requests", (accounts) => {
//     let bts_periphery, bts_core, bmc, holder;
//     let service = "Coin/WrappedCoin";
//     let _uri = "https://github.com/icon-project/icon-bridge";
//     let _native = "PARA";
//     let _fee = 10;
//     let _fixed_fee = 500000;
//     let _name1 = "ICON";
//     let _name2 = "BINANCE";
//     let _name3 = "ETHEREUM";
//     let _name4 = "TRON";
//     let _net1 = "1234.iconee";
//     let _net2 = "1234.binance";
//     let _from1 = "0x12345678";
//     let _from2 = "0x12345678";
//     let _value1 = 999999999999999;
//     let _value2 = 999999999999999;
//     let _to1 = "btp://1234.iconee/0x12345678";
//     let _to2 = "btp://1234.binance/0x12345678";
//     let _txAmt = 1000000;
//     let _txAmt1 = 100000000;
//     let _txAmt2 = 5000000;
//     let RC_OK = 0;
//     let RC_ERR = 1;
//     let REPONSE_HANDLE_SERVICE = 2;
//     let _bmcICON = "btp://1234.iconee/0x1234567812345678";
//     let _sn0 = 1;
//     let _sn1 = 2;
//     let _sn2 = 3;
//
//     before(async () => {
//         bts_periphery = await MockBTSPeriphery.new();
//         bts_core = await BTSCore.new();
//         bmc = await BMC.new("1234.pra");
//         encode_msg = await EncodeMsg.new();
//         await bts_periphery.initialize(bmc.address, bts_core.address, service);
//         await bts_core.initialize(_native, _fee, _fixed_fee);
//         await bts_core.updateBTSPeriphery(bts_periphery.address);
//         holder = await Holder.new();
//         btpAddr = await bmc.bmcAddress();
//         await bmc.addService(service, bts_periphery.address);
//         await bmc.addVerifier(_net1, accounts[1]);
//         await bmc.addVerifier(_net2, accounts[2]);
//         await bmc.addLink(_bmcICON);
//         await holder.addBSHContract(bts_periphery.address, bts_core.address);
//         await bts_core.register(_name1, "", 18);
//         await bts_core.register(_name2, "", 18);
//         await bts_core.register(_name3, "", 18);
//         await bts_core.register(_name4, "", 18);
//         let _msg1 = await encode_msg.encodeTransferMsgWithAddress(
//             _from1,
//             holder.address,
//             _name1,
//             _value1
//         );
//         await bmc.receiveRequest(_bmcICON, "", service, _sn0, _msg1);
//         let _msg2 = await encode_msg.encodeTransferMsgWithAddress(
//             _from2,
//             holder.address,
//             _name2,
//             _value2
//         );
//         await bmc.receiveRequest(_bmcICON, "", service, _sn1, _msg2);
//         await bts_core.transferNativeCoin(_to1, {
//             from: accounts[0],
//             value: _txAmt,
//         });
//         let _responseMsg = await encode_msg.encodeResponseMsg(
//             REPONSE_HANDLE_SERVICE,
//             RC_OK,
//             ""
//         );
//         await bmc.receiveResponse(_net1, service, _sn0, _responseMsg);
//         await holder.setApprove(bts_core.address);
//         await holder.callTransfer(_name1, _txAmt1, _to1);
//         await bmc.receiveResponse(_net1, service, _sn1, _responseMsg);
//         await holder.callTransfer(_name2, _txAmt2, _to2);
//         await bmc.receiveResponse(_net1, service, _sn2, _responseMsg);
//     });
//
//     it(`Scenario 1: Query 'Aggregation Fee'`, async () => {
//         let aggregationFee = await bts_core.getAccumulatedFees();
//
//         assert.equal(aggregationFee.length, 5);
//         assert.equal(aggregationFee[0].coinName, "PARA");
//         assert.equal(aggregationFee[1].coinName, "ICON");
//         assert.equal(aggregationFee[2].coinName, "BINANCE");
//         assert.equal(aggregationFee[3].coinName, "ETHEREUM");
//         assert.equal(aggregationFee[4].coinName, "TRON");
//
//         assert.equal(
//             Number(aggregationFee[0].value),
//             Math.floor(_txAmt / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(aggregationFee[1].value),
//             Math.floor(_txAmt1 / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(aggregationFee[2].value),
//             Math.floor(_txAmt2 / 1000) + _fixed_fee
//         );
//         assert.equal(Number(aggregationFee[3].value), 0);
//         assert.equal(Number(aggregationFee[4].value), 0);
//     });
//
//     it("Scenario 2: Receiving a FeeGathering request not from BMCService", async () => {
//         let _sn3 = 3;
//         let FA1Before = await bts_periphery.getAggregationFeeOf(_native); //  state Aggregation Fee of each type of Coins
//         let FA2Before = await bts_periphery.getAggregationFeeOf(_name1);
//         let FA3Before = await bts_periphery.getAggregationFeeOf(_name2);
//         await truffleAssert.reverts(
//             bts_periphery.handleFeeGathering.call(_to1, service, {
//                 from: accounts[1],
//             }),
//             "Unauthorized"
//         );
//         let FA1After = await bts_periphery.getAggregationFeeOf(_native);
//         let FA2After = await bts_periphery.getAggregationFeeOf(_name1);
//         let FA3After = await bts_periphery.getAggregationFeeOf(_name2);
//         let fees = await bts_periphery.getFees(_sn3); //  get pending Aggregation Fee list
//
//         assert.equal(
//             web3.utils.BN(FA1Before).toNumber(),
//             web3.utils.BN(FA1After).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(FA2Before).toNumber(),
//             web3.utils.BN(FA2After).toNumber()
//         );
//         assert.equal(
//             web3.utils.BN(FA3Before).toNumber(),
//             web3.utils.BN(FA3After).toNumber()
//         );
//         assert.equal(fees.amounts.length, 0);
//     });
//
//     //  Before:
//     //      + state Aggregation Fee of each type of Coins are set
//     //      + pendingAggregation Fee list is empty
//     //  After:
//     //      + all states of Aggregation Fee are push into pendingAggregation Fee list
//     //      + state Aggregation Fee of each type of Coins are reset
//     it("Scenario 3: Handle GatherFee request from BMCService contract", async () => {
//         let _sn3 = 4;
//         let FA1Before = await bts_periphery.getAggregationFeeOf(_native); //  state Aggregation Fee of each type of Coins
//         let FA2Before = await bts_periphery.getAggregationFeeOf(_name1);
//         let FA3Before = await bts_periphery.getAggregationFeeOf(_name2);
//         let _bmcService = await encode_msg.encodeBMCService(_to1, [service]);
//         let output = await bmc.receiveRequest(
//             _bmcICON,
//             "",
//             "bmc",
//             100,
//             _bmcService
//         );
//         let FA1After = await bts_periphery.getAggregationFeeOf(_native);
//         let FA2After = await bts_periphery.getAggregationFeeOf(_name1);
//         let FA3After = await bts_periphery.getAggregationFeeOf(_name2);
//         let fees = await bts_periphery.getFees(_sn3); //  get pending Aggregation Fee list
//         let list = [];
//         for (let i = 0; i < fees.amounts.length; i++) {
//             list[i] = [fees.coinNames[i], fees.amounts[i]];
//         }
//         let _eventMsg = await encode_msg.encodeTransferFeesBMCMessage(
//             btpAddr,
//             _bmcICON,
//             _to1,
//             service,
//             _sn3,
//             bts_core.address,
//             list
//         );
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferStart",
//             {
//                 fromBlock: output.receipt.blockNumber,
//                 toBlock: "latest",
//             }
//         );
//         let event = transferEvents[0].returnValues;
//         assert.equal(event._from, bts_core.address);
//         assert.equal(event._to, _to1);
//         assert.equal(event._sn, _sn3);
//         assert.equal(event._assetDetails.length, 3);
//         assert.equal(event._assetDetails[0].coinName, _native);
//         assert.equal(event._assetDetails[0].value, fees.amounts[0]);
//         assert.equal(event._assetDetails[0].fee, 0);
//         assert.equal(event._assetDetails[1].coinName, _name1);
//         assert.equal(event._assetDetails[1].value, fees.amounts[1]);
//         assert.equal(event._assetDetails[1].fee, 0);
//         assert.equal(event._assetDetails[2].coinName, _name2);
//         assert.equal(event._assetDetails[2].value, fees.amounts[2]);
//         assert.equal(event._assetDetails[2].fee, 0);
//
//         assert.equal(
//             web3.utils.BN(FA1Before).toNumber(),
//             Math.floor(_txAmt / 1000) + _fixed_fee
//         );
//         assert.equal(
//             web3.utils.BN(FA2Before).toNumber(),
//             Math.floor(_txAmt1 / 1000) + _fixed_fee
//         );
//         assert.equal(
//             web3.utils.BN(FA3Before).toNumber(),
//             Math.floor(_txAmt2 / 1000) + _fixed_fee
//         );
//
//         assert.equal(web3.utils.BN(FA1After).toNumber(), 0);
//         assert.equal(web3.utils.BN(FA2After).toNumber(), 0);
//         assert.equal(web3.utils.BN(FA3After).toNumber(), 0);
//
//         assert.equal(fees.coinNames[0], _native);
//         assert.equal(fees.coinNames[1], _name1);
//         assert.equal(fees.coinNames[2], _name2);
//
//         assert.equal(
//             Number(fees.amounts[0]),
//             Math.floor(_txAmt / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(fees.amounts[1]),
//             Math.floor(_txAmt1 / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(fees.amounts[2]),
//             Math.floor(_txAmt2 / 1000) + _fixed_fee
//         );
//
//         assert.equal(output.logs[0].args._next, _bmcICON);
//         assert.equal(output.logs[0].args._msg, _eventMsg);
//     });
//
//     it("Scenario 4: Receiving a successful response", async () => {
//         let _sn3 = 4;
//         let feesBefore = await bts_periphery.getFees(_sn3);
//         let _responseMsg = await encode_msg.encodeResponseMsg(
//             REPONSE_HANDLE_SERVICE,
//             RC_OK,
//             ""
//         );
//         let tx = await bmc.receiveResponse(_net1, service, _sn3, _responseMsg);
//         let feesAfter = await bts_periphery.getFees(_sn3);
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferEnd",
//             {
//                 fromBlock: tx.receipt.blockNumber,
//                 toBlock: "latest",
//             }
//         );
//         let event = transferEvents[0].returnValues;
//
//         assert.equal(event._from, bts_core.address);
//         assert.equal(event._sn, _sn3);
//         assert.equal(event._code, 0);
//         assert.equal(event._response, "");
//
//         assert.equal(feesBefore.amounts.length, 3);
//         assert.equal(feesBefore.coinNames[0], _native);
//         assert.equal(feesBefore.coinNames[1], _name1);
//         assert.equal(feesBefore.coinNames[2], _name2);
//         assert.equal(
//             Number(feesBefore.amounts[0]),
//             Math.floor(_txAmt / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(feesBefore.amounts[1]),
//             Math.floor(_txAmt1 / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(feesBefore.amounts[2]),
//             Math.floor(_txAmt2 / 1000) + _fixed_fee
//         );
//         assert.equal(feesAfter.amounts.length, 0);
//     });
//
//     it("Scenario 5: Receiving an error response", async () => {
//         let _sn4 = 5;
//         let _sn5 = 6;
//         let _sn6 = 7;
//         let _amt1 = 2000000;
//         let _amt2 = 6000000;
//         await holder.callTransfer(_name1, _amt1, _to1);
//         let _responseMsg = await encode_msg.encodeResponseMsg(
//             REPONSE_HANDLE_SERVICE,
//             RC_OK,
//             ""
//         );
//         await bmc.receiveResponse(_net1, service, _sn4, _responseMsg);
//         await holder.callTransfer(_name2, _amt2, _to2);
//         await bmc.receiveResponse(_net2, service, _sn5, _responseMsg);
//         let _bmcService = await encode_msg.encodeBMCService(_to1, [service]);
//         await bmc.receiveRequest(_bmcICON, "", "bmc", 100, _bmcService);
//
//         let FA1Before = await bts_periphery.getAggregationFeeOf(_name1);
//         let FA2Before = await bts_periphery.getAggregationFeeOf(_name2);
//         let feesBefore = await bts_periphery.getFees(_sn6);
//         let _errMsg = await encode_msg.encodeResponseMsg(
//             REPONSE_HANDLE_SERVICE,
//             RC_ERR,
//             ""
//         );
//         let tx = await bmc.receiveResponse(_net1, service, _sn6, _errMsg);
//         let FA1After = await bts_periphery.getAggregationFeeOf(_name1);
//         let FA2After = await bts_periphery.getAggregationFeeOf(_name2);
//         let feesAfter = await bts_periphery.getFees(_sn6);
//
//         const transferEvents = await bts_periphery.getPastEvents(
//             "TransferEnd",
//             {
//                 fromBlock: tx.receipt.blockNumber,
//                 toBlock: "latest",
//             }
//         );
//         let event = transferEvents[0].returnValues;
//
//         assert.equal(event._from, bts_core.address);
//         assert.equal(event._sn, _sn6);
//         assert.equal(event._code, 1);
//         assert.equal(event._response, "");
//
//         assert.equal(feesBefore.amounts.length, 2);
//         assert.equal(feesBefore.coinNames[0], _name1);
//         assert.equal(feesBefore.coinNames[1], _name2);
//         assert.equal(
//             Number(feesBefore.amounts[0]),
//             Math.floor(_amt1 / 1000) + _fixed_fee
//         );
//         assert.equal(
//             Number(feesBefore.amounts[1]),
//             Math.floor(_amt2 / 1000) + _fixed_fee
//         );
//
//         assert.equal(web3.utils.BN(FA1Before).toNumber(), 0);
//         assert.equal(web3.utils.BN(FA2Before).toNumber(), 0);
//         assert.equal(feesAfter.amounts.length, 0);
//         assert.equal(
//             web3.utils.BN(FA1After).toNumber(),
//             Math.floor(_amt1 / 1000) + _fixed_fee
//         );
//         assert.equal(
//             web3.utils.BN(FA2After).toNumber(),
//             Math.floor(_amt2 / 1000) + _fixed_fee
//         );
//     });
// });
//
// contract("As a user, I want to receive multiple Coins/Tokens from ICON blockchain", (accounts) => {
//         let bts_periphery, bts_core, bmc, holder, refundable;
//         let service = "Coin/WrappedCoin";
//         let _uri = "https://github.com/icon-project/icon-bridge";
//         let _native = "PARA";
//         let _fee = 10;
//         let _fixed_fee = 500000;
//         let _name1 = "ICON";
//         let _name2 = "BINANCE";
//         let _name3 = "ETHEREUM";
//         let _name4 = "TRON";
//         let _net1 = "1234.iconee";
//         let _bmcICON = "btp://1234.iconee/0x1234567812345678";
//         let RC_OK = 0;
//         let RC_ERR = 1;
//         let _from1 = "0x12345678";
//         let _to = "btp://1234.iconee/0x12345678";
//
//         before(async () => {
//             bts_periphery = await BTSPeriphery.new();
//             bts_core = await BTSCore.new();
//             bmc = await BMC.new("1234.pra");
//             encode_msg = await EncodeMsg.new();
//             await bts_periphery.initialize(
//                 bmc.address,
//                 bts_core.address,
//                 service
//             );
//             await bts_core.initialize(_native, _fee, _fixed_fee);
//             await bts_core.updateBTSPeriphery(bts_periphery.address);
//             holder = await Holder.new();
//             refundable = await Refundable.new();
//             btpAddr = await bmc.bmcAddress();
//             await bmc.addService(service, bts_periphery.address);
//             await bmc.addVerifier(_net1, accounts[1]);
//             await bmc.addLink(_bmcICON);
//             await holder.addBSHContract(
//                 bts_periphery.address,
//                 bts_core.address
//             );
//             await bts_core.register(_name1, "", 18);
//             await bts_core.register(_name2, "", 18);
//             await bts_core.register(_name3, "", 18);
//             await bts_core.register(_name4, "", 18);
//             await bts_core.transferNativeCoin(_to, {
//                 from: accounts[0],
//                 value: 10000000,
//             });
//         });
//
//         it("Scenario 1: Receiving address is invalid", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 40000;
//             let _address = "0x1234567890123456789";
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_ERR,
//                 "InvalidAddress"
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithStringAddress(
//                 _from1,
//                 _address,
//                 [
//                     [_native, _value1],
//                     [_name1, _value2],
//                     [_name2, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//
//         it("Scenario 2: BSHPerphery receives a request of invalid token", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 40000;
//             let _invalid_token = "EOS";
//             let balance1Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name1
//             );
//             let balance2Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name2
//             );
//             let balance3Before = await bts_core.balanceOf(
//                 holder.address,
//                 _invalid_token
//             );
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_ERR,
//                 "UnregisteredCoin"
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithAddress(
//                 _from1,
//                 holder.address,
//                 [
//                     [_name1, _value1],
//                     [_name2, _value2],
//                     [_invalid_token, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//             let balance1After = await bts_core.balanceOf(
//                 holder.address,
//                 _name1
//             );
//             let balance2After = await bts_core.balanceOf(
//                 holder.address,
//                 _name2
//             );
//             let balance3After = await bts_core.balanceOf(
//                 holder.address,
//                 _invalid_token
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3Before._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3After._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//
//         it("Scenario 3: One of requests is failed in TransferBatch", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 20000000;
//             let balance1Before = await bts_core.balanceOf(
//                 accounts[1],
//                 _name1
//             );
//             let balance2Before = await bts_core.balanceOf(
//                 accounts[1],
//                 _name2
//             );
//             let balance3Before = await bts_core.balanceOf(
//                 accounts[1],
//                 _native
//             );
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_ERR,
//                 "TransferFailed"
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithAddress(
//                 _from1,
//                 accounts[1],
//                 [
//                     [_name1, _value1],
//                     [_name2, _value2],
//                     [_native, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//             let balance1After = await bts_core.balanceOf(
//                 accounts[1],
//                 _name1
//             );
//             let balance2After = await bts_core.balanceOf(
//                 accounts[1],
//                 _name2
//             );
//             let balance3After = await bts_core.balanceOf(
//                 accounts[1],
//                 _native
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2Before._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2After._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//
//         it("Scenario 4: One of requests is failed in TransferBatch", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 40000;
//             let balance1Before = await bts_core.balanceOf(
//                 refundable.address,
//                 _native
//             );
//             let balance2Before = await bts_core.balanceOf(
//                 refundable.address,
//                 _name1
//             );
//             let balance3Before = await bts_core.balanceOf(
//                 refundable.address,
//                 _name2
//             );
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_ERR,
//                 "TransferFailed"
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithAddress(
//                 _from1,
//                 refundable.address,
//                 [
//                     [_native, _value1],
//                     [_name1, _value2],
//                     [_name2, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//             let balance1After = await bts_core.balanceOf(
//                 refundable.address,
//                 _native
//             );
//             let balance2After = await bts_core.balanceOf(
//                 refundable.address,
//                 _name1
//             );
//             let balance3After = await bts_core.balanceOf(
//                 refundable.address,
//                 _name2
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3Before._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3After._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//
//         it("Scenario 5: One of requests is failed in TransferBatch", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 40000;
//             let balance1Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name1
//             );
//             let balance2Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name2
//             );
//             let balance3Before = await bts_core.balanceOf(
//                 holder.address,
//                 _native
//             );
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_ERR,
//                 "TransferFailed"
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithAddress(
//                 _from1,
//                 holder.address,
//                 [
//                     [_name1, _value1],
//                     [_name2, _value2],
//                     [_native, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//             let balance1After = await bts_core.balanceOf(
//                 holder.address,
//                 _name1
//             );
//             let balance2After = await bts_core.balanceOf(
//                 holder.address,
//                 _name2
//             );
//             let balance3After = await bts_core.balanceOf(
//                 holder.address,
//                 _native
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3Before._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2After._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3After._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//
//         it("Scenario 6: Receiving a successful TransferBatch request", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 40000;
//             let balance1Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name1
//             );
//             let balance2Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name2
//             );
//             let balance3Before = await bts_core.balanceOf(
//                 holder.address,
//                 _name3
//             );
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_OK,
//                 ""
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithAddress(
//                 _from1,
//                 holder.address,
//                 [
//                     [_name1, _value1],
//                     [_name2, _value2],
//                     [_name3, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//             let balance1After = await bts_core.balanceOf(
//                 holder.address,
//                 _name1
//             );
//             let balance2After = await bts_core.balanceOf(
//                 holder.address,
//                 _name2
//             );
//             let balance3After = await bts_core.balanceOf(
//                 holder.address,
//                 _name3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance2Before._usableBalance).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balance3Before._usableBalance).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1After._usableBalance).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balance2After._usableBalance).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balance3After._usableBalance).toNumber(),
//                 _value3
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//
//         it("Scenario 7: Receiving a successful TransferBatch request", async () => {
//             let _value1 = 1000;
//             let _value2 = 10000;
//             let _value3 = 40000;
//             let balance1Before = await bts_core.balanceOf(
//                 accounts[1],
//                 _native
//             );
//             let balance2Before = await bts_core.balanceOf(
//                 accounts[1],
//                 _name2
//             );
//             let balance3Before = await bts_core.balanceOf(
//                 accounts[1],
//                 _name3
//             );
//             let _eventMsg = await encode_msg.encodeResponseBMCMessage(
//                 btpAddr,
//                 _bmcICON,
//                 service,
//                 10,
//                 RC_OK,
//                 ""
//             );
//             let _msg = await encode_msg.encodeBatchTransferMsgWithAddress(
//                 _from1,
//                 accounts[1],
//                 [
//                     [_native, _value1],
//                     [_name2, _value2],
//                     [_name3, _value3],
//                 ]
//             );
//             let output = await bmc.receiveRequest(
//                 _bmcICON,
//                 "",
//                 service,
//                 10,
//                 _msg
//             );
//             let balance1After = await bts_core.balanceOf(
//                 accounts[1],
//                 _native
//             );
//             let balance2After = await bts_core.balanceOf(
//                 accounts[1],
//                 _name2
//             );
//             let balance3After = await bts_core.balanceOf(
//                 accounts[1],
//                 _name3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balance1After._usableBalance).toString(),
//                 web3.utils
//                     .BN(balance1Before._usableBalance)
//                     .add(new web3.utils.BN(_value1))
//                     .toString()
//             );
//             assert.equal(
//                 web3.utils.BN(balance2After._usableBalance).toNumber(),
//                 web3.utils.BN(balance2Before._usableBalance).toNumber() +
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balance3After._usableBalance).toNumber(),
//                 web3.utils.BN(balance3Before._usableBalance).toNumber() +
//                 _value3
//             );
//
//             assert.equal(output.logs[0].args._next, _bmcICON);
//             assert.equal(output.logs[0].args._msg, _eventMsg);
//         });
//     }
// );
//
// contract("As a user, I want to send multiple coins/tokens to ICON blockchain", (accounts) => {
//         let bts_periphery, bts_core, bmc, holder;
//         let service = "Coin/WrappedCoin";
//         let _uri = "https://github.com/icon-project/icon-bridge";
//         let _native = "PARA";
//         let _fee = 10;
//         let _fixed_fee = 500000;
//         let _net = "1234.iconee";
//         let _from = "0x12345678";
//         let _value = 999999999999999;
//         let REPONSE_HANDLE_SERVICE = 2;
//         let RC_OK = 0;
//         let RC_ERR = 1;
//         let _bmcICON = "btp://1234.iconee/0x1234567812345678";
//         let _coin1 = "ICON";
//         let _coin2 = "TRON";
//         let _coin3 = "BINANCE";
//         let initAmt = 1000000000000000;
//
//         before(async () => {
//             bts_periphery = await BTSPeriphery.new();
//             bts_core = await BTSCore.new();
//             bmc = await BMC.new("1234.pra");
//             encode_msg = await EncodeMsg.new();
//             await bts_periphery.initialize(
//                 bmc.address,
//                 bts_core.address,
//                 service
//             );
//             await bts_core.initialize(_native, _fee, _fixed_fee);
//             await bts_core.updateBTSPeriphery(bts_periphery.address);
//             holder = await Holder.new();
//             await bmc.addService(service, bts_periphery.address);
//             await bmc.addVerifier(_net, accounts[1]);
//             await bmc.addLink(_bmcICON);
//             await holder.addBSHContract(
//                 bts_periphery.address,
//                 bts_core.address
//             );
//             await bts_core.register(_coin1, "", 18);
//             await bts_core.register(_coin2, "", 18);
//             await bts_core.register(_coin3, "", 18);
//             await bts_core.transferNativeCoin("btp://1234.iconee/0x12345678", {
//                 from: accounts[0],
//                 value: initAmt,
//             });
//             await holder.deposit({from: accounts[1], value: 100000000000000});
//             let _msg1 = await encode_msg.encodeTransferMsgWithAddress(
//                 _from,
//                 holder.address,
//                 _coin1,
//                 _value
//             );
//             await bmc.receiveRequest(_bmcICON, "", service, 0, _msg1);
//             let _msg2 = await encode_msg.encodeTransferMsgWithAddress(
//                 _from,
//                 holder.address,
//                 _coin2,
//                 _value
//             );
//             await bmc.receiveRequest(_bmcICON, "", service, 1, _msg2);
//             let _msg3 = await encode_msg.encodeTransferMsgWithAddress(
//                 _from,
//                 holder.address,
//                 _coin3,
//                 _value
//             );
//             await bmc.receiveRequest(_bmcICON, "", service, 2, _msg3);
//
//             _msg1 = await encode_msg.encodeTransferMsgWithAddress(
//                 _from,
//                 accounts[1],
//                 _coin1,
//                 _value
//             );
//             await bmc.receiveRequest(_bmcICON, "", service, 0, _msg1);
//             _msg2 = await encode_msg.encodeTransferMsgWithAddress(
//                 _from,
//                 accounts[1],
//                 _coin2,
//                 _value
//             );
//             await bmc.receiveRequest(_bmcICON, "", service, 1, _msg2);
//             _msg3 = await encode_msg.encodeTransferMsgWithAddress(
//                 _from,
//                 accounts[1],
//                 _coin3,
//                 _value
//             );
//             await bmc.receiveRequest(_bmcICON, "", service, 2, _msg3);
//         });
//
//         it("Scenario 1: User has not yet set approval for token being transferred out by Operator", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _values = [600000, 700000];
//             let _native_amt = 800000;
//             let _query = [_native, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it(`Scenario 2: User has set approval, but user's balance has insufficient amount`, async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _values = [600000, 9999999999999999n];
//             let _native_amt = 700000;
//             let _query = [_native, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 3: User requests to transfer an invalid Token", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let invalid_token = "EOS";
//             let _coins = [_coin1, invalid_token];
//             let _values = [600000, 700000];
//             let _native_amt = 800000;
//             let _query = [_native, _coin1, invalid_token];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 4: User transfers Tokens to an invalid BTP Address format", async () => {
//             let _to = "1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _values = [600000, 700000];
//             let _query = [_native, _coin1, _coin2];
//             let _native_amt = 800000;
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 5: User requests to transfer zero Token", async () => {
//             let _to = "1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _values = [600000, 0];
//             let _query = [_native, _coin1, _coin2];
//             let _native_amt = 800000;
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 6: Transffering amount is less than fixed fee", async () => {
//             let _to = "1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _values = [600000, 300000];
//             let _query = [_native, _coin1, _coin2];
//             let _native_amt = 800000;
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 7: User requests to transfer to an invalid network/Not Supported Network", async () => {
//             let _to = "btp://1234.eos/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _values = [600000, 700000];
//             let _query = [_native, _coin1, _coin2];
//             let _native_amt = 1000;
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     _native_amt
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 8: Account client sends an invalid request of transferBatch", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_native, _native, _native];
//             let _values = [600000, 600000, 600000];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 accounts[2],
//                 _coins
//             );
//             await truffleAssert.reverts(
//                 bts_core.transferBatch.call(_coins, _values, _to, {
//                     from: accounts[2],
//                     value: 600000,
//                 }),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 accounts[2],
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 initAmt
//             );
//         });
//
//         it("Scenario 9: Contract client sends an invalid request of transferBatch", async () => {
//             let _to = "btp://1234.eos/0x12345678";
//             let _coins = [_native, _coin1, _coin2];
//             let _values = [600000, 700000];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             await holder.setApprove(bts_core.address);
//             await truffleAssert.reverts(
//                 holder.callTransferBatch.call(
//                     bts_core.address,
//                     _coins,
//                     _values,
//                     _to,
//                     0
//                 ),
//                 "revert"
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 0
//             );
//         });
//
//         it("Scenario 10: Contract client sends a valid transferBatch request", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _values = [_value2, _value3];
//             let _query = [_native, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             let tx = await holder.callTransferBatch(
//                 bts_core.address,
//                 _coins,
//                 _values,
//                 _to,
//                 _value1
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferStart",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//             assert.equal(event._from, holder.address);
//             assert.equal(event._to, _to);
//             assert.equal(event._sn, 2);
//             assert.equal(event._assetDetails.length, 3);
//             assert.equal(event._assetDetails[0].coinName, _coin1);
//             assert.equal(event._assetDetails[0].value, _value2 - chargedFee2);
//             assert.equal(event._assetDetails[0].fee, chargedFee2);
//             assert.equal(event._assetDetails[1].coinName, _coin2);
//             assert.equal(event._assetDetails[1].value, _value3 - chargedFee3);
//             assert.equal(event._assetDetails[1].fee, chargedFee3);
//             assert.equal(event._assetDetails[2].coinName, _native);
//             assert.equal(event._assetDetails[2].value, _value1 - chargedFee1);
//             assert.equal(event._assetDetails[2].fee, chargedFee1);
//
//             const linkStatus = await bmc.getStatus(_bmcICON);
//             const bmcBtpAddress = await bmc.getBmcBtpAddress();
//
//             const messageEvents = await bmc.getPastEvents("Message", {
//                 fromBlock: tx.receipt.blockNumber,
//                 toBlock: "latest",
//             });
//             event = messageEvents[0].returnValues;
//             assert.equal(event._next, _bmcICON);
//             assert.equal(event._seq, linkStatus.txSeq);
//
//             const bmcMsg = rlp.decode(event._msg);
//
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
//             assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 2);
//
//             const ServiceMsg = rlp.decode(bmcMsg[4]);
//             assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
//
//             const coinTransferMsg = rlp.decode(ServiceMsg[1]);
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
//                 holder.address
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
//                 _to.split("/").slice(-1)[0]
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
//                 _coin1
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
//                 _value2 - chargedFee2
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][1][0])),
//                 _coin2
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][1][1])),
//                 _value3 - chargedFee3
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][2][0])),
//                 _native
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][2][1])),
//                 _value1 - chargedFee1
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber() -
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() -
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() -
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt + _value1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 _value3
//             );
//         });
//
//         it("Scenario 11: BTSPeriphery receives a successful response of a recent request", async () => {
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _coins = [_native, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let _responseMsg = await encode_msg.encodeResponseMsg(
//                 REPONSE_HANDLE_SERVICE,
//                 RC_OK,
//                 ""
//             );
//             let tx = await bmc.receiveResponse(_net, service, 2, _responseMsg);
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let fees = await bts_core.getAccumulatedFees();
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferEnd",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//
//             assert.equal(event._from, holder.address);
//             assert.equal(event._sn, 2);
//             assert.equal(event._code, 0);
//             assert.equal(event._response, "");
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber()
//             );
//
//             assert.equal(fees[0].coinName, _native);
//             assert.equal(Number(fees[0].value), chargedFee1);
//             assert.equal(fees[1].coinName, _coin1);
//             assert.equal(Number(fees[1].value), chargedFee2);
//             assert.equal(fees[2].coinName, _coin2);
//             assert.equal(Number(fees[2].value), chargedFee3);
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt + _value1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 chargedFee3
//             );
//         });
//
//         it("Scenario 12: Account client sends a valid transferBatch request", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_coin3, _coin1, _coin2];
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _values = [_value1, _value2, _value3];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 accounts[1],
//                 _coins
//             );
//             await bts_core.setApprovalForAll(bts_core.address, true, {
//                 from: accounts[1],
//             });
//             let tx = await bts_core.transferBatch(_coins, _values, _to, {
//                 from: accounts[1],
//             });
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 accounts[1],
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferStart",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//             assert.equal(event._from, accounts[1]);
//             assert.equal(event._to, _to);
//             assert.equal(event._sn, 3);
//             assert.equal(event._assetDetails.length, 3);
//             assert.equal(event._assetDetails[0].coinName, _coin3);
//             assert.equal(event._assetDetails[0].value, _value1 - chargedFee1);
//             assert.equal(event._assetDetails[0].fee, chargedFee1);
//             assert.equal(event._assetDetails[1].coinName, _coin1);
//             assert.equal(event._assetDetails[1].value, _value2 - chargedFee2);
//             assert.equal(event._assetDetails[1].fee, chargedFee2);
//             assert.equal(event._assetDetails[2].coinName, _coin2);
//             assert.equal(event._assetDetails[2].value, _value3 - chargedFee3);
//             assert.equal(event._assetDetails[2].fee, chargedFee3);
//
//             const linkStatus = await bmc.getStatus(_bmcICON);
//             const bmcBtpAddress = await bmc.getBmcBtpAddress();
//
//             const messageEvents = await bmc.getPastEvents("Message", {
//                 fromBlock: tx.receipt.blockNumber,
//                 toBlock: "latest",
//             });
//             event = messageEvents[0].returnValues;
//             assert.equal(event._next, _bmcICON);
//             assert.equal(event._seq, linkStatus.txSeq);
//
//             const bmcMsg = rlp.decode(event._msg);
//
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
//             assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 3);
//
//             const ServiceMsg = rlp.decode(bmcMsg[4]);
//             assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
//
//             const coinTransferMsg = rlp.decode(ServiceMsg[1]);
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
//                 accounts[1]
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
//                 _to.split("/").slice(-1)[0]
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
//                 _coin3
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
//                 _value1 - chargedFee1
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][1][0])),
//                 _coin1
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][1][1])),
//                 _value2 - chargedFee2
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][2][0])),
//                 _coin2
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][2][1])),
//                 _value3 - chargedFee3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber() -
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() -
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() -
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 _value2 + chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 _value3 + chargedFee3
//             );
//         });
//
//         it("Scenario 13: BTSPeriphery receives an error response of a recent request", async () => {
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _coins = [_coin3, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 accounts[1],
//                 _coins
//             );
//             let _responseMsg = await encode_msg.encodeResponseMsg(
//                 REPONSE_HANDLE_SERVICE,
//                 RC_ERR,
//                 ""
//             );
//             let tx = await bmc.receiveResponse(_net, service, 3, _responseMsg);
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 accounts[1],
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferEnd",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//
//             assert.equal(event._from, accounts[1]);
//             assert.equal(event._sn, 3);
//             assert.equal(event._code, 1);
//             assert.equal(event._response, "");
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber() +
//                 _value1 -
//                 chargedFee1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() +
//                 _value2 -
//                 chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() +
//                 _value3 -
//                 chargedFee3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 chargedFee1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 2 * chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 2 * chargedFee3
//             );
//         });
//
//         it("Scenario 14: Contract client sends a valid transferBatch request", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_coin3, _coin1, _coin2];
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _values = [_value1, _value2, _value3];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             await holder.setApprove(bts_core.address);
//             let tx = await holder.callTransferBatch(
//                 bts_core.address,
//                 _coins,
//                 _values,
//                 _to,
//                 0
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferStart",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//             assert.equal(event._from, holder.address);
//             assert.equal(event._to, _to);
//             assert.equal(event._sn, 4);
//             assert.equal(event._assetDetails.length, 3);
//             assert.equal(event._assetDetails[0].coinName, _coin3);
//             assert.equal(event._assetDetails[0].value, _value1 - chargedFee1);
//             assert.equal(event._assetDetails[0].fee, chargedFee1);
//             assert.equal(event._assetDetails[1].coinName, _coin1);
//             assert.equal(event._assetDetails[1].value, _value2 - chargedFee2);
//             assert.equal(event._assetDetails[1].fee, chargedFee2);
//             assert.equal(event._assetDetails[2].coinName, _coin2);
//             assert.equal(event._assetDetails[2].value, _value3 - chargedFee3);
//             assert.equal(event._assetDetails[2].fee, chargedFee3);
//
//             const linkStatus = await bmc.getStatus(_bmcICON);
//             const bmcBtpAddress = await bmc.getBmcBtpAddress();
//
//             const messageEvents = await bmc.getPastEvents("Message", {
//                 fromBlock: tx.receipt.blockNumber,
//                 toBlock: "latest",
//             });
//             event = messageEvents[0].returnValues;
//             assert.equal(event._next, _bmcICON);
//             assert.equal(event._seq, linkStatus.txSeq);
//
//             const bmcMsg = rlp.decode(event._msg);
//
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
//             assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 4);
//
//             const ServiceMsg = rlp.decode(bmcMsg[4]);
//             assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
//
//             const coinTransferMsg = rlp.decode(ServiceMsg[1]);
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
//                 holder.address
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
//                 _to.split("/").slice(-1)[0]
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
//                 _coin3
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
//                 _value1 - chargedFee1
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][1][0])),
//                 _coin1
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][1][1])),
//                 _value2 - chargedFee2
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][2][0])),
//                 _coin2
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][2][1])),
//                 _value3 - chargedFee3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber() -
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() -
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() -
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 _value1 + chargedFee1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 _value2 + 2 * chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 _value3 + 2 * chargedFee3
//             );
//         });
//
//         it("Scenario 15: BTSPeriphery receives an error response of a recent request", async () => {
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _coins = [_coin3, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let _responseMsg = await encode_msg.encodeResponseMsg(
//                 REPONSE_HANDLE_SERVICE,
//                 RC_ERR,
//                 ""
//             );
//             let tx = await bmc.receiveResponse(_net, service, 4, _responseMsg);
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferEnd",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//
//             assert.equal(event._from, holder.address);
//             assert.equal(event._sn, 4);
//             assert.equal(event._code, 1);
//             assert.equal(event._response, "");
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber() +
//                 _value1 -
//                 chargedFee1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() +
//                 _value2 -
//                 chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() +
//                 _value3 -
//                 chargedFee3
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 2 * chargedFee1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 3 * chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 3 * chargedFee3
//             );
//         });
//
//         //  This test is replicated from Scenario 10
//         it("Scenario 16: Contract client sends a valid transferBatch request", async () => {
//             let _to = "btp://1234.iconee/0x12345678";
//             let _coins = [_coin1, _coin2];
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _values = [_value2, _value3];
//             let _query = [_native, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             await holder.setApprove(bts_core.address);
//             let tx = await holder.callTransferBatch(
//                 bts_core.address,
//                 _coins,
//                 _values,
//                 _to,
//                 _value1
//             );
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _query
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _query
//             );
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferStart",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//             assert.equal(event._from, holder.address);
//             assert.equal(event._to, _to);
//             assert.equal(event._sn, 5);
//             assert.equal(event._assetDetails.length, 3);
//             assert.equal(event._assetDetails[0].coinName, _coin1);
//             assert.equal(event._assetDetails[0].value, _value2 - chargedFee2);
//             assert.equal(event._assetDetails[0].fee, chargedFee2);
//             assert.equal(event._assetDetails[1].coinName, _coin2);
//             assert.equal(event._assetDetails[1].value, _value3 - chargedFee3);
//             assert.equal(event._assetDetails[1].fee, chargedFee3);
//             assert.equal(event._assetDetails[2].coinName, _native);
//             assert.equal(event._assetDetails[2].value, _value1 - chargedFee1);
//             assert.equal(event._assetDetails[2].fee, chargedFee1);
//
//             const linkStatus = await bmc.getStatus(_bmcICON);
//             const bmcBtpAddress = await bmc.getBmcBtpAddress();
//
//             const messageEvents = await bmc.getPastEvents("Message", {
//                 fromBlock: tx.receipt.blockNumber,
//                 toBlock: "latest",
//             });
//             event = messageEvents[0].returnValues;
//             assert.equal(event._next, _bmcICON);
//             assert.equal(event._seq, linkStatus.txSeq);
//
//             const bmcMsg = rlp.decode(event._msg);
//
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[0])), bmcBtpAddress);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[1])), _bmcICON);
//             assert.equal(web3.utils.hexToUtf8(toHex(bmcMsg[2])), service);
//             assert.equal(web3.utils.hexToNumber(toHex(bmcMsg[3])), 5);
//
//             const ServiceMsg = rlp.decode(bmcMsg[4]);
//             assert.equal(web3.utils.hexToUtf8(toHex(ServiceMsg[0])), 0);
//
//             const coinTransferMsg = rlp.decode(ServiceMsg[1]);
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[0])),
//                 holder.address
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[1])),
//                 _to.split("/").slice(-1)[0]
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][0][0])),
//                 _coin1
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][0][1])),
//                 _value2 - chargedFee2
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][1][0])),
//                 _coin2
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][1][1])),
//                 _value3 - chargedFee3
//             );
//             assert.equal(
//                 web3.utils.hexToUtf8(toHex(coinTransferMsg[2][2][0])),
//                 _native
//             );
//             assert.equal(
//                 web3.utils.hexToNumber(toHex(coinTransferMsg[2][2][1])),
//                 _value1 - chargedFee1
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber() -
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() -
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() -
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt + 2 * _value1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 _value2 + 3 * chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 _value3 + 3 * chargedFee3
//             );
//         });
//
//         it("Scenario 17: BTSPeriphery receives an error response of a recent request", async () => {
//             let _value1 = 600000;
//             let _value2 = 700000;
//             let _value3 = 800000;
//             let _coins = [_native, _coin1, _coin2];
//             let balanceBefore = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let _responseMsg = await encode_msg.encodeResponseMsg(
//                 REPONSE_HANDLE_SERVICE,
//                 RC_ERR,
//                 ""
//             );
//             let tx = await bmc.receiveResponse(_net, service, 5, _responseMsg);
//             let balanceAfter = await bts_core.getBalanceOfBatch(
//                 holder.address,
//                 _coins
//             );
//             let bts_core_balance = await bts_core.getBalanceOfBatch(
//                 bts_core.address,
//                 _coins
//             );
//             let chargedFee1 = Math.floor(_value1 / 1000) + _fixed_fee;
//             let chargedFee2 = Math.floor(_value2 / 1000) + _fixed_fee;
//             let chargedFee3 = Math.floor(_value3 / 1000) + _fixed_fee;
//
//             const transferEvents = await bts_periphery.getPastEvents(
//                 "TransferEnd",
//                 {fromBlock: tx.receipt.blockNumber, toBlock: "latest"}
//             );
//             let event = transferEvents[0].returnValues;
//
//             assert.equal(event._from, holder.address);
//             assert.equal(event._sn, 5);
//             assert.equal(event._code, 1);
//             assert.equal(event._response, "");
//
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[0]).toNumber(),
//                 _value1
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[1]).toNumber(),
//                 _value2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._lockedBalances[2]).toNumber(),
//                 _value3
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[1]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._lockedBalances[2]).toNumber(),
//                 0
//             );
//
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[0]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[0]).toNumber()
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[1]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[1]).toNumber() +
//                 _value2 -
//                 chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._usableBalances[2]).toNumber(),
//                 web3.utils.BN(balanceBefore._usableBalances[2]).toNumber() +
//                 _value3 -
//                 chargedFee3
//             );
//             assert.equal(
//                 web3.utils.BN(balanceBefore._refundableBalances[0]).toNumber(),
//                 0
//             );
//             assert.equal(
//                 web3.utils.BN(balanceAfter._refundableBalances[0]).toNumber(),
//                 _value1 - chargedFee1
//             );
//
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[0]).toNumber(),
//                 initAmt + 2 * _value1
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[1]).toNumber(),
//                 4 * chargedFee2
//             );
//             assert.equal(
//                 web3.utils.BN(bts_core_balance._usableBalances[2]).toNumber(),
//                 4 * chargedFee3
//             );
//         });
//     }
// );
