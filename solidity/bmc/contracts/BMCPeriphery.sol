// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.5.0 <0.8.0;
pragma experimental ABIEncoderV2;

import "./interfaces/IBSH.sol";
import "./interfaces/IBMCPeriphery.sol";
import "./interfaces/IBMCManagement.sol";
import "./libraries/ParseAddress.sol";
import "./libraries/RLPDecodeStruct.sol";
import "./libraries/RLPEncodeStruct.sol";
import "./libraries/String.sol";
import "./libraries/Types.sol";
import "./libraries/Utils.sol";

import "@openzeppelin/contracts-upgradeable/proxy/Initializable.sol";

contract BMCPeriphery is IBMCPeriphery, Initializable {
    using String for string;
    using ParseAddress for address;
    using RLPDecodeStruct for bytes;
    using RLPEncodeStruct for Types.BMCMessage;
    using RLPEncodeStruct for Types.Response;
    using Utils for uint256;

    uint256 internal constant UNKNOWN_ERR = 0;
    uint256 internal constant BMC_ERR = 10;
    uint256 internal constant BSH_ERR = 40;

    string private bmcBtpAddress; // a network address, i.e. btp://1234.pra/0xabcd
    address private bmcManagement;
    bytes[] internal msgs;

    function initialize(string memory _network, address _bmcManagementAddr)
        public
        initializer
    {
        bmcBtpAddress = string("btp://").concat(_network).concat("/").concat(
            address(this).toString()
        );
        bmcManagement = _bmcManagementAddr;
    }

    event Message(
        string _next, //  an address of the next BMC (it could be a destination BMC)
        uint256 _seq, //  a sequence number of BMC (NOT sequence number of BSH)
        bytes _msg
    );

    // emit errors in BTP messages processing
    event ErrorOnBTPError(
        string _svc,
        int256 _sn,
        uint256 _code,
        string _errMsg,
        uint256 _svcErrCode,
        string _svcErrMsg
    );

    function getBmcBtpAddress() external view override returns (string memory) {
        return bmcBtpAddress;
    }

    function validateRelay(string calldata _prev)
        internal
        view
        returns (address)
    {
        address relay = address(0);
        address[] memory relays = IBMCManagement(bmcManagement).getLinkRelays(
            _prev
        );
        for (uint256 i = 0; i < relays.length; i++) {
            if (msg.sender == relays[i]) {
                relay = msg.sender;
                break;
            }
        }
        require(
            relay == msg.sender,
            "BMCRevertUnauthroized: relay not registered"
        );
        return relay;
    }

    /**
       @notice Verify and decode RelayMessage, and dispatch BTP Messages to registered BSHs
       @dev Caller must be a registered relayer.     
       @param _prev    BTP Address of the BMC generates the message
       @param _msg     base64 encoded string of serialized bytes of Relay Message refer RelayMessage structure
     */
    function handleRelayMessage(string calldata _prev, bytes calldata _msg)
        external
        override
    {
        validateRelay(_prev);

        uint256 linkRxSeq = IBMCManagement(bmcManagement).getLinkRxSeq(_prev);
        uint256 linkRxHeight = IBMCManagement(bmcManagement).getLinkRxHeight(
            _prev
        );

        uint256 rxSeq = linkRxSeq;
        uint256 rxHeight = linkRxHeight;

        Types.ReceiptProof[] memory rps = _msg.decodeReceiptProofs();

        Types.BMCMessage memory bmsg;
        Types.MessageEvent memory ev;

        for (uint256 i = 0; i < rps.length; i++) {
            if (rps[i].height < rxHeight) {
                revert("RevertInvalidRxHeight: rp.height < rxHeight");
            }
            rxHeight = rps[i].height;
            for (uint256 j = 0; j < rps[i].events.length; j++) {
                rxSeq++;
                ev = rps[i].events[j];
                if (ev.seq > rxSeq) {
                    revert("RevertInvalidRxSeq: ev.seq > expected rxSeq");
                }
                if (ev.seq < rxSeq) {
                    revert("RevertInvalidRxSeq: ev.seq < expected rxSeq");
                }
                if (!ev.nextBmc.compareTo(bmcBtpAddress)) {
                    continue;
                }
                try this.decodeBTPMessage(ev.message) returns (
                    Types.BMCMessage memory _decoded
                ) {
                    bmsg = _decoded;
                } catch {
                    continue;
                }
                if (bmsg.dst.compareTo(bmcBtpAddress)) {
                    handleMessage(_prev, bmsg);
                    continue;
                }
                (string memory _net, ) = bmsg.dst.splitBTPAddress();
                try IBMCManagement(bmcManagement).resolveRoute(_net) returns (
                    string memory _nextLink,
                    string memory
                ) {
                    _sendMessage(_nextLink, ev.message);
                } catch Error(string memory _error) {
                    _sendError(_prev, bmsg, BMC_ERR, _error);
                }
            }
        }

        IBMCManagement(bmcManagement).updateLinkRxSeq(_prev, rxSeq - linkRxSeq);
        IBMCManagement(bmcManagement).updateRelayStats(
            msg.sender,
            0,
            rxSeq - linkRxSeq
        );
        IBMCManagement(bmcManagement).updateLinkRxHeight(
            _prev,
            rxHeight - linkRxHeight
        );
    }

    //  @dev Despite this function was set as external, it should be called internally
    //  since Solidity does not allow using try_catch with internal function
    //  this solution can solve the issue
    function decodeBTPMessage(bytes memory _rlp)
        external
        pure
        returns (Types.BMCMessage memory)
    {
        return _rlp.decodeBMCMessage();
    }

    function handleMessage(string calldata _prev, Types.BMCMessage memory _msg)
        internal
    {
        address _bshAddr;
        if (_msg.svc.compareTo("bmc")) {
            Types.BMCService memory _sm;
            try this.tryDecodeBMCService(_msg.message) returns (
                Types.BMCService memory res
            ) {
                _sm = res;
            } catch {
                _sendError(_prev, _msg, BMC_ERR, "BMCRevertParseFailure");
                return;
            }

            if (_sm.serviceType.compareTo("FeeGathering")) {
                Types.GatherFeeMessage memory _gatherFee;
                try this.tryDecodeGatherFeeMessage(_sm.payload) returns (
                    Types.GatherFeeMessage memory res
                ) {
                    _gatherFee = res;
                } catch {
                    _sendError(_prev, _msg, BMC_ERR, "BMCRevertParseFailure");
                    return;
                }

                for (uint256 i = 0; i < _gatherFee.svcs.length; i++) {
                    _bshAddr = IBMCManagement(bmcManagement)
                        .getBshServiceByName(_gatherFee.svcs[i]);
                    //  If 'svc' not found, ignore
                    if (_bshAddr != address(0)) {
                        try
                            IBSH(_bshAddr).handleFeeGathering(
                                _gatherFee.fa,
                                _gatherFee.svcs[i]
                            )
                        {} catch {
                            //  If BSH contract throws a revert error, ignore and continue
                        }
                    }
                }
            } else if (_sm.serviceType.compareTo("Link")) {
                string memory _to = _sm.payload.decodePropagateMessage();
                Types.Link memory link = IBMCManagement(bmcManagement).getLink(
                    _prev
                );
                bool check;
                if (link.isConnected) {
                    for (uint256 i = 0; i < link.reachable.length; i++)
                        if (_to.compareTo(link.reachable[i])) {
                            check = true;
                            break;
                        }
                    if (!check) {
                        string[] memory _links = new string[](1);
                        _links[0] = _to;
                        IBMCManagement(bmcManagement).updateLinkReachable(
                            _prev,
                            _links
                        );
                    }
                }
            } else if (_sm.serviceType.compareTo("Unlink")) {
                string memory _to = _sm.payload.decodePropagateMessage();
                Types.Link memory link = IBMCManagement(bmcManagement).getLink(
                    _prev
                );
                if (link.isConnected) {
                    for (uint256 i = 0; i < link.reachable.length; i++) {
                        if (_to.compareTo(link.reachable[i]))
                            IBMCManagement(bmcManagement).deleteLinkReachable(
                                _prev,
                                i
                            );
                    }
                }
            } else if (_sm.serviceType.compareTo("Init")) {
                string[] memory _links = _sm.payload.decodeInitMessage();
                IBMCManagement(bmcManagement).updateLinkReachable(
                    _prev,
                    _links
                );
            } else if (_sm.serviceType.compareTo("Sack")) {
                // skip this case since it has been removed from internal services
            } else revert("BMCRevert: not exists internal handler");
        } else {
            _bshAddr = IBMCManagement(bmcManagement).getBshServiceByName(
                _msg.svc
            );
            if (_bshAddr == address(0)) {
                _sendError(_prev, _msg, BMC_ERR, "BMCRevertNotExistsBSH");
                return;
            }

            if (_msg.sn >= 0) {
                (string memory _net, ) = _msg.src.splitBTPAddress();
                try
                    IBSH(_bshAddr).handleBTPMessage(
                        _net,
                        _msg.svc,
                        uint256(_msg.sn),
                        _msg.message
                    )
                {} catch Error(string memory _error) {
                    /**
                     * @dev Uncomment revert to debug errors
                     */
                    //revert(_error);
                    _sendError(_prev, _msg, BSH_ERR, _error);
                }
            } else {
                Types.Response memory _errMsg = _msg.message.decodeResponse();
                try
                    IBSH(_bshAddr).handleBTPError(
                        _msg.src,
                        _msg.svc,
                        uint256(_msg.sn * -1),
                        _errMsg.code,
                        _errMsg.message
                    )
                {} catch Error(string memory _error) {
                    emit ErrorOnBTPError(
                        _msg.svc,
                        _msg.sn * -1,
                        _errMsg.code,
                        _errMsg.message,
                        BSH_ERR,
                        _error
                    );
                } catch (bytes memory _error) {
                    emit ErrorOnBTPError(
                        _msg.svc,
                        _msg.sn * -1,
                        _errMsg.code,
                        _errMsg.message,
                        UNKNOWN_ERR,
                        string(_error)
                    );
                }
            }
        }
    }

    //  @dev Solidity does not allow using try_catch with internal function
    //  Thus, work-around solution is the followings
    //  If there is any error throwing, BMC contract can catch it, then reply back a RC_ERR Response
    function tryDecodeBMCService(bytes calldata _msg)
        external
        pure
        returns (Types.BMCService memory)
    {
        return _msg.decodeBMCService();
    }

    function tryDecodeGatherFeeMessage(bytes calldata _msg)
        external
        pure
        returns (Types.GatherFeeMessage memory)
    {
        return _msg.decodeGatherFeeMessage();
    }

    function _sendMessage(string memory _to, bytes memory _serializedMsg)
        internal
    {
        IBMCManagement(bmcManagement).updateLinkTxSeq(_to);
        emit Message(
            _to,
            IBMCManagement(bmcManagement).getLinkTxSeq(_to),
            _serializedMsg
        );
    }

    function _sendError(
        string calldata _prev,
        Types.BMCMessage memory _message,
        uint256 _errCode,
        string memory _errMsg
    ) internal {
        if (_message.sn > 0) {
            bytes memory _serializedMsg = Types
                .BMCMessage(
                    bmcBtpAddress,
                    _message.src,
                    _message.svc,
                    _message.sn * -1,
                    Types.Response(_errCode, _errMsg).encodeResponse()
                )
                .encodeBMCMessage();
            _sendMessage(_prev, _serializedMsg);
        }
    }

    /**
       @notice Send the message to a specific network.
       @dev Caller must be an registered BSH.
       @param _to      Network Address of destination network
       @param _svc     Name of the service
       @param _sn      Serial number of the message, it should be positive
       @param _msg     Serialized bytes of Service Message
    */
    function sendMessage(
        string memory _to,
        string memory _svc,
        uint256 _sn,
        bytes memory _msg
    ) external override {
        require(
            msg.sender == bmcManagement ||
                IBMCManagement(bmcManagement).getBshServiceByName(_svc) ==
                msg.sender,
            "BMCRevertUnauthorized"
        );
        require(_sn >= 0, "BMCRevertInvalidSN");
        //  In case BSH sends a REQUEST_COIN_TRANSFER,
        //  but '_to' is a network which is not supported by BMC
        //  revert() therein
        (string memory _nextLink, string memory _dst) = IBMCManagement(
            bmcManagement
        ).resolveRoute(_to);
        bytes memory _rlp = Types
            .BMCMessage(bmcBtpAddress, _dst, _svc, int256(_sn), _msg)
            .encodeBMCMessage();
        _sendMessage(_nextLink, _rlp);
    }

    /*
       @notice Get status of BMC.
       @param _link        BTP Address of the connected BMC.
       @return tx_seq       Next sequence number of the next sending message.
       @return rx_seq       Next sequence number of the message to receive.
       @return verifier     VerifierStatus Object contains status information of the BMV.
    */
    function getStatus(string calldata _link)
        public
        view
        override
        returns (Types.LinkStats memory _linkStats)
    {
        Types.Link memory link = IBMCManagement(bmcManagement).getLink(_link);
        require(link.isConnected == true, "BMCRevertNotExistsLink");
        Types.RelayStats[] memory _relays = IBMCManagement(bmcManagement)
            .getRelayStatusByLink(_link);
        uint256 _rotateTerm = link.maxAggregation.getRotateTerm(
            link.blockIntervalSrc.getScale(link.blockIntervalDst)
        );
        return
            Types.LinkStats(
                link.rxSeq,
                link.txSeq,
                Types.VerifierStats(0, 0, 0, ""), //dummy
                _relays,
                link.relayIdx,
                link.rotateHeight,
                _rotateTerm,
                link.delayLimit,
                link.maxAggregation,
                link.rxHeightSrc,
                link.rxHeight,
                link.blockIntervalSrc,
                link.blockIntervalDst,
                block.number
            );
    }
}
