const BTSCoreV1 = artifacts.require("BTSCoreV1");
const BTSCoreV2 = artifacts.require("BTSCoreV2");
const BTSPeriphery = artifacts.require("BTSPeriphery");
const BMC = artifacts.require("MockBMC");
const {assert} = require("chai");
const truffleAssert = require("truffle-assertions");
const {deployProxy, upgradeProxy} = require("@openzeppelin/truffle-upgrades");

contract("BTSCore Unit Tests - After Upgrading Contract", (accounts) => {
    let bts_coreV1, bts_coreV2;
    let _native = "PARA";
    let _uri = "https://github.com/icon-project/icon-bridge";
    let _fee = 10;
    let _fixed_fee = 500000;
    before(async () => {
        bts_coreV1 = await deployProxy(BTSCoreV1, [_native, _fee, _fixed_fee]);
        bts_coreV2 = await upgradeProxy(bts_coreV1.address, BTSCoreV2);

        bts_periphery = await BTSPeriphery.new();

        bmc = await BMC.new("1234.pra");
        await bts_periphery.initialize(bmc.address, bts_coreV2.address);
    });


    it(`Scenario 1: Should allow contract's owner to register a new coin`, async () => {
        let _name = "ICON";
        await bts_coreV2.updateBTSPeriphery(bts_periphery.address);

        await bts_coreV2.register(_name, "ICX", 18, 100, 0, "0x0000000000000000000000000000000000000000");
        output = await bts_coreV2.coinNames();
        assert(output[0] === _native && output[1] === "ICON");
    });

    it("Scenario 2: Should revert when an arbitrary client tries to register a new coin", async () => {
        let _name = "TRON";
        await truffleAssert.reverts(
            bts_coreV2.register.call(_name, "", 18, 100, 0, "0x0000000000000000000000000000000000000000", {from: accounts[1]}),
            "Unauthorized"
        );
    });

    it("Scenario 3: Should revert when contract owner registers an existed coin", async () => {
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_coreV2.register.call(_name, "", 18, 100, 0, "0x0000000000000000000000000000000000000000"),
            "ExistCoin"
        );
    });


    it("Scenario 4: Should revert when arbitrary client updates BTSPeriphery contract", async () => {
        await truffleAssert.reverts(
            bts_coreV2.updateBTSPeriphery.call(accounts[2], {
                from: accounts[1],
            }),
            "Unauthorized"
        );
    });

    it("Scenario 5: Should succeed when a client, which owns a refundable, tries to reclaim", async () => {
        let _coin = "ICON";
        let _value = 10000;
        let _id = await bts_coreV2.coinId(_coin);
        await bts_coreV2.mintMock(bts_coreV2.address, _id, _value);
        await bts_coreV2.setRefundableBalance(accounts[5], _coin, _value);
        let balanceBefore = await bts_coreV2.balanceOf(accounts[5], _coin);
        await bts_coreV2.reclaim(_coin, _value, {from: accounts[5]});
        let balanceAfter = await bts_coreV2.balanceOf(accounts[5], _coin);
        assert(
            web3.utils.BN(balanceAfter._userBalance).toNumber() ===
            web3.utils.BN(balanceBefore._userBalance).toNumber() + _value
        );
    });

    it("Scenario 6: Should not allow any clients (even a contract owner) to call a refund()", async () => {
        //  This function should be called only by itself (BTSCore contract)
        let _coin = "ICON";
        let _value = 10000;
        await truffleAssert.reverts(
            bts_coreV2.refund.call(accounts[0], _coin, _value),
            "Unauthorized"
        );
    });

    it("Scenario 7: Should allow contract owner to update fee ratio", async () => {
        let new_fee = 20;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fee);

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fee
        );
    });

    it("Scenario 8: Should revert when arbitrary client updates fee ratio", async () => {
        let old_fee = 20;
        let new_fee = 50;
        let _feeNumerator = 100;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_coreV2.setFeeRatio.call(_name, _feeNumerator, new_fee, {from: accounts[1]}),
            "Unauthorized"
        );

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === old_fee
        );
    });

    it("Scenario 9: Should revert when Fee Numerator is higher than Fee Denominator", async () => {
        let old_feeNumerator = 100;
        let new_feeNumerator = 15000;
        let new_fee = 15000;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_coreV2.setFeeRatio.call(_name, new_feeNumerator, new_fee),
            "InvalidSetting"
        );

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._feeNumerator).toNumber() === old_feeNumerator
        );
    });

    it("Scenario 10: Should allow contract owner to update fixed fee", async () => {
        let new_fixed_fee = 1000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fixed_fee);

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fixed_fee
        );
    });

    it("Scenario 11: Should revert when arbitrary client updates fixed fee", async () => {
        let old_fixed_fee = 1000000;
        let new_fixed_fee = 2000000;

        let _feeNumerator = 100;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_coreV2.setFeeRatio.call(_name, _feeNumerator, new_fixed_fee, {from: accounts[1]}),
            "Unauthorized"
        );

        let fees = await bts_coreV2.feeRatio(_name);


        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            old_fixed_fee
        );
    });

    it("Scenario 12: Should allow to set fee ratio greater than 0", async () => {
        let new_fixed_fee = 10;
        let _feeNumerator = 100;
        let _name = "ICON";

        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fixed_fee);
        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            new_fixed_fee
        );
    });

    it("Scenario 13: Should receive an id of a given coin name when querying a valid supporting coin", async () => {
        let _name1 = "wBTC";
        let _name2 = "Ethereum";
        await bts_coreV2.register(_name1, "", 18, 100, 0, "0x0000000000000000000000000000000000000000");
        await bts_coreV2.register(_name2, "", 18, 100, 0, "0x0000000000000000000000000000000000000000");

        let _query = "ICON";
        let id = web3.utils.keccak256(_query);
        let result = await bts_coreV2.coinId(_query);
        assert(
            web3.utils.toChecksumAddress(result) !==
            web3.utils.toChecksumAddress(
                "0x0000000000000000000000000000000000000000"
            )
        );
    });

    it("Scenario 14: Should receive an id = 0 when querying an invalid supporting coin", async () => {
        let _query = "EOS";
        let result = await bts_coreV2.coinId(_query);
        assert(
            web3.utils.toChecksumAddress(result) ===
            web3.utils.toChecksumAddress(
                "0x0000000000000000000000000000000000000000"
            )
        );
    });

    it("Scenario 15: Should revert when a non-Owner tries to add a new Owner", async () => {
        let oldList = await bts_coreV2.getOwners();
        await truffleAssert.reverts(
            bts_coreV2.addOwner.call(accounts[1], {from: accounts[2]}),
            "Unauthorized"
        );
        let newList = await bts_coreV2.getOwners();
        assert(
            oldList.length === 1 &&
            oldList[0] === accounts[0] &&
            newList.length === 1 &&
            newList[0] === accounts[0]
        );
    });

    it("Scenario 16: Should allow a current Owner to add a new Owner", async () => {
        let oldList = await bts_coreV2.getOwners();
        await bts_coreV2.addOwner(accounts[1]);
        let newList = await bts_coreV2.getOwners();
        assert(
            oldList.length === 1 &&
            oldList[0] === accounts[0] &&
            newList.length === 2 &&
            newList[0] === accounts[0] &&
            newList[1] === accounts[1]
        );
    });

    it("Scenario 17: Should allow old owner to register a new coin - After adding new Owner", async () => {
        let _name3 = "TRON";
        await bts_coreV2.register(_name3, "", 18, 100, 0, "0x0000000000000000000000000000000000000000");
        output = await bts_coreV2.coinNames();
        assert(
            output[0] === _native &&
            output[1] === "ICON" &&
            output[2] === "wBTC" &&
            output[3] === "Ethereum" &&
            output[4] === "TRON"
        );
    });

    it("Scenario 18: Should allow new owner to register a new coin", async () => {
        let _name3 = "BINANCE";
        await bts_coreV2.register(_name3, "", 18, 100, 0, "0x0000000000000000000000000000000000000000", {from: accounts[1]});
        output = await bts_coreV2.coinNames();
        assert(
            output[0] === _native &&
            output[1] === "ICON" &&
            output[2] === "wBTC" &&
            output[3] === "Ethereum" &&
            output[4] === "TRON" &&
            output[5] === "BINANCE"
        );
    });

    it("Scenario 19: Should allow contract owner to update BTSPeriphery contract", async () => {
        await bts_coreV2.updateBTSPeriphery(accounts[2]);
    });

    it("Scenario 20: Should allow new owner to update BTSPeriphery contract", async () => {
        //  Must clear BTSPeriphery setting since BTSPeriphery has been set
        //  The requirement: if (BTSPeriphery.address != address(0)) must check whether BTSPeriphery has any pending requests
        //  Instead of creating a MockBTSPeriphery, just clear BTSPeriphery setting
        await bts_coreV2.clearBTSPeripherySetting();
        await bts_coreV2.updateBTSPeriphery(accounts[3], {from: accounts[1]});
    });

    it("Scenario 21: Should also allow old owner to update BTSPeriphery contract - After adding new Owner", async () => {
        //  Must clear BTSPeriphery setting since BTSPeriphery has been set
        //  The requirement: if (BTSPeriphery.address != address(0)) must check whether BTSPeriphery has any pending requests
        //  Instead of creating a MockBTSPeriphery, just clear BTSPeriphery setting
        await bts_coreV2.clearBTSPeripherySetting();
        await bts_coreV2.updateBTSPeriphery(accounts[3], {from: accounts[0]});
    });

    it("Scenario 22: Should not allow any clients (even a contract owner) to call a mint()", async () => {
        //  This function should be called only by BTSPeriphery contract
        let _coin = "ICON";
        let _value = 10000;
        await truffleAssert.reverts(
            bts_coreV2.mint.call(accounts[0], _coin, _value),
            "Unauthorized"
        );
    });

    it("Scenario 23: Should not allow any clients (even a contract owner) to call a handleResponseService()", async () => {
        //  This function should be called only by BTSPeriphery contract
        let RC_OK = 0;
        let _coin = "ICON";
        let _value = 10000;
        let _fee = 1;
        let _rspCode = RC_OK;
        await truffleAssert.reverts(
            bts_coreV2.handleResponseService.call(
                accounts[0],
                _coin,
                _value,
                _fee,
                _rspCode
            ),
            "Unauthorized"
        );
    });


    it("Scenario 24: Should allow new owner to update new fee ratio", async () => {
        let new_fee = 30;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fee, {from: accounts[1]});

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fee
        );
    });

    it("Scenario 25: Should also allow old owner to update new fee ratio - After adding new Owner", async () => {
        let new_fee = 40;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fee, {from: accounts[0]});

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() === new_fee
        );
    });

    it("Scenario 26: Should allow new owner to update new fixed fee", async () => {
        let new_fixed_fee = 3000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fixed_fee, {from: accounts[1]});

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            new_fixed_fee
        );
    });
//
    it("Scenario 27: Should also allow old owner to update new fixed fee - After adding new Owner", async () => {
        let new_fixed_fee = 4000000;
        let _feeNumerator = 100;
        let _name = "ICON";
        await bts_coreV2.setFeeRatio(_name, _feeNumerator, new_fixed_fee, {from: accounts[0]});

        let fees = await bts_coreV2.feeRatio(_name);

        assert(
            web3.utils.BN(fees._fixedFee).toNumber() ===
            new_fixed_fee
        );
    });

    it("Scenario 28: Should revert when non-Owner tries to remove an Owner", async () => {
        let oldList = await bts_coreV2.getOwners();
        await truffleAssert.reverts(
            bts_coreV2.removeOwner.call(accounts[0], {from: accounts[2]}),
            "Unauthorized"
        );
        let newList = await bts_coreV2.getOwners();
        assert(
            oldList.length === 2 &&
            oldList[0] === accounts[0] &&
            oldList[1] === accounts[1] &&
            newList.length === 2 &&
            newList[0] === accounts[0] &&
            newList[1] === accounts[1]
        );
    });

    it("Scenario 29: Should allow one current Owner to remove another Owner", async () => {
        let oldList = await bts_coreV2.getOwners();
        await bts_coreV2.removeOwner(accounts[0], {from: accounts[1]});
        let newList = await bts_coreV2.getOwners();
        assert(
            oldList.length === 2 &&
            oldList[0] === accounts[0] &&
            oldList[1] === accounts[1] &&
            newList.length === 1 &&
            newList[0] === accounts[1]
        );
    });

    it("Scenario 30: Should revert when the last Owner removes him/herself", async () => {
        let oldList = await bts_coreV2.getOwners();
        await truffleAssert.reverts(
            bts_coreV2.removeOwner.call(accounts[1], {from: accounts[1]}),
            "CannotRemoveMinOwner"
        );
        let newList = await bts_coreV2.getOwners();
        assert(
            oldList.length === 1 &&
            oldList[0] === accounts[1] &&
            newList.length === 1 &&
            newList[0] === accounts[1]
        );
    });

    it("Scenario 31: Should revert when removed Owner tries to register a new coin", async () => {
        let _name3 = "KYBER";

        await truffleAssert.reverts(
            bts_coreV2.register.call(_name3, "", 18, 100, 0, "0x0000000000000000000000000000000000000000"),
            "Unauthorized"
        );
        output = await bts_coreV2.coinNames();
        assert(
            output[0] === _native &&
            output[1] === "ICON" &&
            output[2] === "wBTC" &&
            output[3] === "Ethereum" &&
            output[4] === "TRON" &&
            output[5] === "BINANCE"
        );

    });

    it("Scenario 32: Should revert when removed Owner tries to update BTSPeriphery contract", async () => {
        //  Must clear BTSPeriphery setting since BTSPeriphery has been set
        //  The requirement: if (BTSPeriphery.address != address(0)) must check whether BTSPeriphery has any pending requests
        //  Instead of creating a MockBTSPeriphery, just clear BTSPeriphery setting
        await bts_coreV2.clearBTSPeripherySetting();
        await truffleAssert.reverts(
            bts_coreV2.updateBTSPeriphery.call(accounts[3], {
                from: accounts[0],
            }),
            "Unauthorized"
        );
    });

    it("Scenario 33: Should not allow any clients (even a contract owner) to call a transferFees()", async () => {
        //  This function should be called only by BTSPeriphery contract
        let _fa = "btp://1234.iconee/0x1234567812345678";
        await truffleAssert.reverts(
            bts_coreV2.transferFees.call(_fa),
            "Unauthorized"
        );
    });


    it("Scenario 34: Should revert when removed Owner tries to update new fee ratio", async () => {
        let new_fee = 30;
        let _feeNumerator = 100;
        let _name = "ICON";
        await truffleAssert.reverts(
            bts_coreV2.setFeeRatio.call(_name, _feeNumerator, new_fee, {from: accounts[0]}),
            "Unauthorized"
        );
    });

    it("Scenario 35: Should allow arbitrary client to query balance of an account", async () => {
        let _coin = "TRON";
        let _id = await bts_coreV2.coinId(_coin);
        let _value = 10000;
        await bts_coreV2.mintMock(accounts[3], _id, _value);
        let balance = await bts_coreV2.balanceOf(accounts[3], _coin, {
            from: accounts[3],
        });
        assert(
            web3.utils.BN(balance._usableBalance).toNumber() === 0 &&
            web3.utils.BN(balance._lockedBalance).toNumber() === 0 &&
            web3.utils.BN(balance._refundableBalance).toNumber() === 0 &&
            web3.utils.BN(balance._userBalance).toNumber() === _value
        );
    });

    it("Scenario 36: Should allow arbitrary client to query a batch of balances of an account", async () => {
        let _coin1 = "BINANCE";
        let _coin2 = "TRON";
        let _id1 = await bts_coreV2.coinId(_coin1);
        let _id2 = await bts_coreV2.coinId(_coin2);
        let _value = 10000;
        let another_value = 2000;
        await bts_coreV2.mintMock(accounts[2], _id1, another_value);
        await bts_coreV2.mintMock(accounts[2], _id2, _value);
        let balance = await bts_coreV2.balanceOfBatch(
            accounts[2],
            [_coin1, _coin2],
            {from: accounts[2]}
        );
        assert(
            web3.utils.BN(balance._usableBalances[0]).toNumber() ===
            0 &&
            web3.utils.BN(balance._lockedBalances[0]).toNumber() === 0 &&
            web3.utils.BN(balance._refundableBalances[0]).toNumber() ===
            0 &&
            web3.utils.BN(balance._userBalances[0]).toNumber() ===
            another_value &&
            web3.utils.BN(balance._usableBalances[1]).toNumber() ===
            0 &&
            web3.utils.BN(balance._lockedBalances[1]).toNumber() === 0 &&
            web3.utils.BN(balance._refundableBalances[1]).toNumber() === 0 &&
            web3.utils.BN(balance._userBalances[1]).toNumber() ===
            _value
        );
    });

    it("Scenario 37: Should allow arbitrary client to query an Accumulated Fees", async () => {
        let _coin1 = "ICON";
        let _coin2 = "TRON";
        let _value1 = 40000;
        let _value2 = 10000;
        let _native_value = 2000;
        await bts_coreV2.setAggregationFee(_coin1, _value1);
        await bts_coreV2.setAggregationFee(_coin2, _value2);
        await bts_coreV2.setAggregationFee(_native, _native_value);
        let fees = await bts_coreV2.getAccumulatedFees({from: accounts[3]});
        assert(
            fees[0].coinName === _native &&
            Number(fees[0].value) === _native_value &&
            fees[1].coinName === _coin1 &&
            Number(fees[1].value) === _value1 &&
            fees[2].coinName === "wBTC" &&
            Number(fees[2].value) === 0 &&
            fees[3].coinName === "Ethereum" &&
            Number(fees[3].value) === 0 &&
            fees[4].coinName === _coin2 &&
            Number(fees[4].value) === _value2
        );
    });

    it("Scenario 38: Should revert when a client reclaims an exceeding amount", async () => {
        let _coin = "ICON";
        let _value = 10000;
        let _exceedAmt = 20000;
        await bts_coreV2.setRefundableBalance(accounts[2], _coin, _value);
        await truffleAssert.reverts(
            bts_coreV2.reclaim.call(_coin, _exceedAmt, {from: accounts[2]}),
            "Imbalance"
        );
    });

    it("Scenario 39: Should revert when a client, which does not own a refundable, tries to reclaim", async () => {
        let _coin = "ICON";
        let _value = 10000;
        await truffleAssert.reverts(
            bts_coreV2.reclaim.call(_coin, _value, {from: accounts[3]}),
            "Imbalance"
        );
    });


});
