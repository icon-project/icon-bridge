// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.0;
pragma abicoder v2;
import "../BTSPeriphery.sol";
import "../BTSCore.sol";

contract MockBTSPeriphery is BTSPeriphery {
    // using String for string;

    // function getFees(uint256 _sn)
    //     external
    //     view
    //     returns (Types.PendingTransferCoin memory)
    // {
    //     return requests[_sn];
    // }

    // function getAggregationFeeOf(string calldata _coinName)
    //     external
    //     view
    //     returns (uint256 _fee)
    // {
    //     Types.Asset[] memory _fees = btsCore.getAccumulatedFees();
    //     for (uint256 i = 0; i < _fees.length; i++) {
    //         if (_coinName.compareTo(_fees[i].coinName)) return _fees[i].value;
    //     }
    // }
}
