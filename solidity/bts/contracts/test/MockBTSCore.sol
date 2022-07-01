// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.0;
pragma abicoder v2;
import "../BTSCore.sol";

contract MockBTSCore is BTSCore {
    function mintMock(
        address _acc,
        address _erc20Address,
        uint256 _value
    ) external {
        IERC20Tradable(_erc20Address).mint(_acc, _value);
    }

    function burnMock(
        address _acc,
        address _erc20Address,
        uint256 _value
    ) external {
        IERC20Tradable(_erc20Address).mint(_acc, _value);
    }

    function setAggregationFee(string calldata _coinName, uint256 _value)
        external
    {
        aggregationFee[_coinName] += _value;
    }

    function clearAggregationFee() external {
        for (uint256 i = 0; i < coinsName.length; i++) {
            delete aggregationFee[coinsName[i]];
        }
    }

    function clearBTSPeripherySetting() external {
        btsPeriphery = IBTSPeriphery(address(0));
    }

    function setRefundableBalance(
        address _acc,
        string calldata _coinName,
        uint256 _value
    ) external {
        balances[_acc][_coinName].refundableBalance += _value;
    }
}
