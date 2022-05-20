// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.0 <0.8.5;
pragma abicoder v2;

import "../interfaces/IBSH.sol";
import "../libraries/RLPDecodeStruct.sol";
import "../libraries/ParseAddress.sol";
import "../libraries/String.sol";

contract MockBSH is IBSH {
    using RLPDecodeStruct for bytes;
    using String for string;
    using ParseAddress for address;
    using ParseAddress for string;

    constructor() {}

    function handleBTPMessage(
        string calldata _from,
        string calldata _svc,
        uint256 _sn,
        bytes calldata _msg
    ) external override {
        require(_svc.compareTo("TokenBSH") == true, "Invalid Service Name");
        Types.ServiceMessage memory _sm = _msg.decodeServiceMessage();
        if (_sm.serviceType == Types.ServiceType.REQUEST_COIN_TRANSFER) {
            Types.TransferAssets memory _ta = _sm.data.decodeTransferAsset();
            string memory _statusMsg;
            uint256 _status;

            try this.handleRequest(_ta) {
                _statusMsg = "Transfer Success";
                _status = 0;
            } catch Error(string memory _err) {
                /**
                 * @dev Uncomment revert to debug errors
                 */
                //revert(_err);
                _statusMsg = _err;
                _status = 1;
            }
        }
    }

    function handleRequest(Types.TransferAssets memory transferAssets)
        external
    {
        string memory _toAddress = transferAssets.to;

        try this.checkParseAddress(_toAddress) {} catch {
            revert("Invalid Address");
        }
        Types.TokenAsset[] memory _asset = transferAssets.asset;
        for (uint256 i = 0; i < _asset.length; i++) {
            // Check if the _toAddress is invalid
            uint256 _amount = _asset[i].value;
            string memory _tokenName = _asset[i].name;
            // Check if the token is registered already
        }
    }

    function checkParseAddress(string calldata _to) external pure {
        _to._toAddress();
    }

    function handleBTPError(
        string calldata,
        string calldata,
        uint256 _sn,
        uint256,
        string calldata
    ) external pure override {
        require(_sn != 1000, "Mocking error message on handleBTPError");
        assert(_sn != 100); // mocking invalid opcode
    }

    /**
       @notice Handle Gather Fee Request from ICON.
       @dev Every BSH must implement this function
       @param _fa    BTP Address of Fee Aggregator in ICON
       @param _svc   Name of the service
   */
    function handleFeeGathering(string calldata _fa, string calldata _svc)
        external
        override
    {}
}
