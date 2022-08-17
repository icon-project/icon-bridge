// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.0;
pragma abicoder v2;
import "@openzeppelin/contracts-upgradeable/utils/math/SafeMathUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./interfaces/IBTSPeriphery.sol";
import "./interfaces/IBTSCore.sol";
import "./libraries/String.sol";
import "./libraries/Types.sol";
import "./ERC20Tradable.sol";
import "./interfaces/IERC20Tradable.sol";

/**
   @title BTSCore contract
   @dev This contract is used to handle coin transferring service
   Note: The coin of following contract can be:
   Native Coin : The native coin of this chain
   Wrapped Native Coin : A tokenized ERC20 version of another native coin like ICX
*/
contract BTSCore is Initializable, IBTSCore, ReentrancyGuardUpgradeable {
    using SafeMathUpgradeable for uint256;
    using String for string;
    event SetOwnership(address indexed promoter, address indexed newOwner);
    event RemoveOwnership(address indexed remover, address indexed formerOwner);

    struct Coin {
        address addr;
        uint256 feeNumerator;
        uint256 fixedFee;
        uint256 coinType;
    }

    modifier onlyOwner() {
        require(owners[msg.sender] == true, "Unauthorized");
        _;
    }

    modifier onlyBTSPeriphery() {
        require(msg.sender == address(btsPeriphery), "Unauthorized");
        _;
    }

    uint256 private constant FEE_DENOMINATOR = 10**4;
    uint256 private constant RC_OK = 0;
    uint256 private constant RC_ERR = 1;

    uint256 private constant NATIVE_COIN_TYPE = 0;
    uint256 private constant NATIVE_WRAPPED_COIN_TYPE = 1;
    uint256 private constant NON_NATIVE_TOKEN_TYPE = 2;

    IBTSPeriphery internal btsPeriphery;

    address[] private listOfOwners;
    uint256[] private chargedAmounts; //   a list of amounts have been charged so far (use this when Fee Gathering occurs)
    string[] internal coinsName; // a string array stores names of supported coins
    string[] private chargedCoins; //   a list of coins' names have been charged so far (use this when Fee Gathering occurs)
    string internal nativeCoinName;

    mapping(address => bool) internal owners;
    mapping(string => uint256) internal aggregationFee; // storing Aggregation Fee in state mapping variable.
    mapping(address => mapping(string => Types.Balance)) internal balances;
    mapping(string => address) internal coins; //  a list of all supported coins
    mapping(string => Coin) internal coinDetails;

    function initialize(
        string calldata _nativeCoinName,
        uint256 _feeNumerator,
        uint256 _fixedFee
    ) public initializer {
        owners[msg.sender] = true;
        listOfOwners.push(msg.sender);
        emit SetOwnership(address(0), msg.sender);
        nativeCoinName = _nativeCoinName;
        coins[_nativeCoinName] = address(0);
        coinsName.push(_nativeCoinName);
        coinDetails[_nativeCoinName] = Coin(
            address(0),
            _feeNumerator,
            _fixedFee,
            NATIVE_COIN_TYPE
        );
    }

    /**
        @notice Get name of nativecoin
        @dev caller can be any
        @return Name of nativecoin
    */
    function getNativeCoinName() external override view returns (string memory) {
        return nativeCoinName;
    }

    /**
       @notice Adding another Onwer.
       @dev Caller must be an Onwer of BTP network
       @param _owner    Address of a new Onwer.
   */
    function addOwner(address _owner) external override onlyOwner {
        require(owners[_owner] == false, "ExistedOwner");
        owners[_owner] = true;
        listOfOwners.push(_owner);
        emit SetOwnership(msg.sender, _owner);
    }

    /**
       @notice Removing an existing Owner.
       @dev Caller must be an Owner of BTP network
       @dev If only one Owner left, unable to remove the last Owner
       @param _owner    Address of an Owner to be removed.
   */
    function removeOwner(address _owner) external override onlyOwner {
        require(listOfOwners.length > 1, "CannotRemoveMinOwner");
        require(owners[_owner] == true, "NotanOwner");
        delete owners[_owner];
        _remove(_owner);
        emit RemoveOwnership(msg.sender, _owner);
    }

    function _remove(address _addr) internal {
        for (uint256 i = 0; i < listOfOwners.length; i++)
            if (listOfOwners[i] == _addr) {
                listOfOwners[i] = listOfOwners[listOfOwners.length - 1];
                listOfOwners.pop();
                break;
            }
    }

    /**
       @notice Checking whether one specific address has Owner role.
       @dev Caller can be ANY
       @param _owner    Address needs to verify.
    */
    function isOwner(address _owner) external view override returns (bool) {
        return owners[_owner];
    }

    /**
       @notice Get a list of current Owners
       @dev Caller can be ANY
       @return      An array of addresses of current Owners
    */
    function getOwners() external view override returns (address[] memory) {
        return listOfOwners;
    }

    /**
        @notice update BTS Periphery address.
        @dev Caller must be an Owner of this contract
        _btsPeriphery Must be different with the existing one.
        @param _btsPeriphery    BTSPeriphery contract address.
    */
    function updateBTSPeriphery(address _btsPeriphery)
        external
        override
        onlyOwner
    {
        require(_btsPeriphery != address(0), "InvalidSetting");
        if (address(btsPeriphery) != address(0)) {
            require(
                btsPeriphery.hasPendingRequest() == false,
                "HasPendingRequest"
            );
        }
        btsPeriphery = IBTSPeriphery(_btsPeriphery);
    }

    /**
        @notice set fee ratio.
        @dev Caller must be an Owner of this contract
        The transfer fee is calculated by feeNumerator/FEE_DEMONINATOR. 
        The feeNumetator should be less than FEE_DEMONINATOR
        _feeNumerator if it is set to `10`, which means the default fee ratio is 0.1%.
        @param _feeNumerator    the fee numerator
    */
    function setFeeRatio(
        string calldata _name,
        uint256 _feeNumerator,
        uint256 _fixedFee
    ) external override onlyOwner {
        require(_feeNumerator <= FEE_DENOMINATOR, "InvalidSetting");
        require(_name.compareTo(nativeCoinName) || coins[_name] != address(0), "TokenNotExists");
        require(_fixedFee >= 0 && _feeNumerator >= 0, "LessThan0");
        coinDetails[_name].feeNumerator = _feeNumerator;
        coinDetails[_name].fixedFee = _fixedFee;
    }

    /**
        @notice Registers a wrapped coin and id number of a supporting coin.
        @dev Caller must be an Owner of this contract
        _name Must be different with the native coin name.
        _symbol symbol name for wrapped coin.
        _decimals decimal number
        @param _name    Coin name. 
    */
    function register(
        string calldata _name,
        string calldata _symbol,
        uint8 _decimals,
        uint256 _feeNumerator,
        uint256 _fixedFee,
        address _addr
    ) external override onlyOwner {
        require(!_name.compareTo(nativeCoinName), "ExistNativeCoin");
        require(coins[_name] == address(0), "ExistCoin");
        require(_feeNumerator <= FEE_DENOMINATOR, "InvalidFeeSetting");
        require(_fixedFee >= 0 && _feeNumerator >= 0, "LessThan0");
        if (_addr == address(0)) {
            address deployedERC20 = address(
                new ERC20Tradable(_name, _symbol, _decimals)
            );
            coins[_name] = deployedERC20;
            coinsName.push(_name);
            coinDetails[_name] = Coin(
                deployedERC20,
                _feeNumerator,
                _fixedFee,
                NATIVE_WRAPPED_COIN_TYPE
            );
        } else {
            coins[_name] = _addr;
            coinsName.push(_name);
            coinDetails[_name] = Coin(
                _addr,
                _feeNumerator,
                _fixedFee,
                NON_NATIVE_TOKEN_TYPE
            );
        }
        string[] memory tokenArr = new string[](1);
        tokenArr[0] = _name;
        uint[] memory valArr = new uint[](1);
        valArr[0] = type(uint256).max;
        btsPeriphery.setTokenLimit(tokenArr, valArr);
    }

    /**
       @notice Return all supported coins names
       @dev 
       @return _names   An array of strings.
    */
    function coinNames()
        external
        view
        override
        returns (string[] memory _names)
    {
        return coinsName;
    }

    /**
       @notice  Return an _id number of Coin whose name is the same with given _coinName.
       @dev     Return nullempty if not found.
       @return  _coinId     An ID number of _coinName.
    */
    function coinId(string calldata _coinName)
        external
        view
        override
        returns (address)
    {
        return coins[_coinName];
    }

    /**
       @notice  Check Validity of a _coinName
       @dev     Call by BTSPeriphery contract to validate a requested _coinName
       @return  _valid     true of false
    */
    function isValidCoin(string calldata _coinName)
        external
        view
        override
        returns (bool _valid)
    {
        return (coins[_coinName] != address(0) ||
            _coinName.compareTo(nativeCoinName));
    }

    /**
        @notice Get fee numerator and fixed fee
        @dev caller can be any
        @param _coinName Coin name
        @return _feeNumerator Fee numerator for given coin
        @return _fixedFee Fixed fee for given coin
    */
    function feeRatio(string calldata _coinName)
        external
        override
        view
        returns (uint _feeNumerator, uint _fixedFee)
    {
        Coin memory coin = coinDetails[_coinName];
        _feeNumerator = coin.feeNumerator;
        _fixedFee = coin.fixedFee;
    }

    /**
        @notice Return a usable/locked/refundable balance of an account based on coinName.
        @return _usableBalance the balance that users are holding.
        @return _lockedBalance when users transfer the coin, 
                it will be locked until getting the Service Message Response.
        @return _refundableBalance refundable balance is the balance that will be refunded to users.
    */
    function balanceOf(address _owner, string memory _coinName)
        external
        view
        override
        returns (
            uint256 _usableBalance,
            uint256 _lockedBalance,
            uint256 _refundableBalance,
            uint256 _userBalance
        )
    {
        if (_coinName.compareTo(nativeCoinName)) {
            return (
                0,
                balances[_owner][_coinName].lockedBalance,
                balances[_owner][_coinName].refundableBalance,
                address(_owner).balance
            );
        }
        address _erc20Address = coins[_coinName];
        IERC20 ierc20 = IERC20(_erc20Address);
        _userBalance = _erc20Address != address(0)
            ? ierc20.balanceOf(_owner)
            : 0;
        uint allowance = _erc20Address != address(0)
            ? ierc20.allowance(_owner, address(this))
            : 0;
        _usableBalance = allowance > _userBalance
            ? _userBalance
            : allowance;
        return (
            _usableBalance,
            balances[_owner][_coinName].lockedBalance,
            balances[_owner][_coinName].refundableBalance,
            _userBalance
        );
    }

    /**
        @notice Return a list Balance of an account.
        @dev The order of request's coinNames must be the same with the order of return balance
        Return 0 if not found.
        @return _usableBalances         An array of Usable Balances
        @return _lockedBalances         An array of Locked Balances
        @return _refundableBalances     An array of Refundable Balances
    */
    function balanceOfBatch(address _owner, string[] calldata _coinNames)
        external
        view
        override
        returns (
            uint256[] memory _usableBalances,
            uint256[] memory _lockedBalances,
            uint256[] memory _refundableBalances,
            uint256[] memory _userBalances
        )
    {
        _usableBalances = new uint256[](_coinNames.length);
        _lockedBalances = new uint256[](_coinNames.length);
        _refundableBalances = new uint256[](_coinNames.length);
        _userBalances = new uint256[](_coinNames.length);
        for (uint256 i = 0; i < _coinNames.length; i++) {
            (
                _usableBalances[i],
                _lockedBalances[i],
                _refundableBalances[i],
                _userBalances[i]
            ) = this.balanceOf(_owner, _coinNames[i]);
        }
        return (_usableBalances, _lockedBalances, _refundableBalances, _userBalances);
    }

    /**
        @notice Return a list accumulated Fees.
        @dev only return the asset that has Asset's value greater than 0
        @return _accumulatedFees An array of Asset
    */
    function getAccumulatedFees()
        external
        view
        override
        returns (Types.Asset[] memory _accumulatedFees)
    {
        _accumulatedFees = new Types.Asset[](coinsName.length);
        for (uint256 i = 0; i < coinsName.length; i++) {
            _accumulatedFees[i] = (
                Types.Asset(coinsName[i], aggregationFee[coinsName[i]])
            );
        }
        return _accumulatedFees;
    }

    /**
       @notice Allow users to deposit `msg.value` native coin into a BTSCore contract.
       @dev MUST specify msg.value
       @param _to  An address that a user expects to receive an amount of tokens.
    */
    function transferNativeCoin(string calldata _to) external payable override {

        btsPeriphery.checkTransferRestrictions(
            nativeCoinName,
            msg.sender,
            msg.value
        );
        //  Aggregation Fee will be charged on BSH Contract
        //  A new charging fee has been proposed. `fixedFee` is introduced
        //  _chargeAmt = fixedFee + msg.value * feeNumerator / FEE_DENOMINATOR
        //  Thus, it's likely that _chargeAmt is always greater than 0
        //  require(_chargeAmt > 0) can be omitted
        //  If msg.value less than _chargeAmt, it likely fails when calculating
        //  _amount = _value - _chargeAmt
        uint256 _chargeAmt = msg
            .value
            .mul(coinDetails[nativeCoinName].feeNumerator)
            .div(FEE_DENOMINATOR)
            .add(coinDetails[nativeCoinName].fixedFee);

        //  @dev msg.value is an amount request to transfer (include fee)
        //  Later on, it will be calculated a true amount that should be received at a destination
        _sendServiceMessage(
            msg.sender,
            _to,
            coinsName[0],
            msg.value,
            _chargeAmt
        );
    }

    /**
       @notice Allow users to deposit an amount of wrapped native coin `_coinName` from the `msg.sender` address into the BTSCore contract.
       @dev Caller must set to approve that the wrapped tokens can be transferred out of the `msg.sender` account by BTSCore contract.
       It MUST revert if the balance of the holder for token `_coinName` is lower than the `_value` sent.
       @param _coinName    A given name of a wrapped coin 
       @param _value       An amount request to transfer from a Requester (include fee)
       @param _to          Target BTP address.
    */
    function transfer(
        string calldata _coinName,
        uint256 _value,
        string calldata _to
    ) external override {
        require(!_coinName.compareTo(nativeCoinName), "InvalidWrappedCoin");
        address _erc20Address = coins[_coinName];
        require(_erc20Address != address(0), "UnregisterCoin");

        btsPeriphery.checkTransferRestrictions(
            _coinName,
            msg.sender,
            _value
        );

        //  _chargeAmt = fixedFee + msg.value * feeNumerator / FEE_DENOMINATOR
        //  Thus, it's likely that _chargeAmt is always greater than 0
        //  require(_chargeAmt > 0) can be omitted
        //  If _value less than _chargeAmt, it likely fails when calculating
        //  _amount = _value - _chargeAmt
        uint256 _chargeAmt = _value
            .mul(coinDetails[_coinName].feeNumerator)
            .div(FEE_DENOMINATOR)
            .add(coinDetails[_coinName].fixedFee);

        //  Transfer and Lock Token processes:
        //  BTSCore contract calls safeTransferFrom() to transfer the Token from Caller's account (msg.sender)
        //  Before that, Caller must approve (setApproveForAll) to accept
        //  token being transfer out by an Operator
        //  If this requirement is failed, a transaction is reverted.
        //  After transferring token, BTSCore contract updates Caller's locked balance
        //  as a record of pending transfer transaction
        //  When a transaction is completed without any error on another chain,
        //  Locked Token amount (bind to an address of caller) will be reset/subtract,
        //  then emit a successful TransferEnd event as a notification
        //  Otherwise, the locked amount will also be updated
        //  but BTSCore contract will issue a refund to Caller before emitting an error TransferEnd event
        IERC20Tradable(_erc20Address).transferFrom(
            msg.sender,
            address(this),
            _value
        );
        //  @dev _value is an amount request to transfer (include fee)
        //  Later on, it will be calculated a true amount that should be received at a destination
        _sendServiceMessage(msg.sender, _to, _coinName, _value, _chargeAmt);
    }

    /**
       @notice This private function handles overlapping procedure before sending a service message to BTSPeriphery
       @param _from             An address of a Requester
       @param _to               BTP address of of Receiver on another chain
       @param _coinName         A given name of a requested coin 
       @param _value            A requested amount to transfer from a Requester (include fee)
       @param _chargeAmt        An amount being charged for this request
    */
    function _sendServiceMessage(
        address _from,
        string calldata _to,
        string memory _coinName,
        uint256 _value,
        uint256 _chargeAmt
    ) private {
        //  Lock this requested _value as a record of a pending transferring transaction
        //  @dev `_value` is a requested amount to transfer, from a Requester, including charged fee
        //  The true amount to receive at a destination receiver is calculated by
        //  _amounts[0] = _value.sub(_chargeAmt);
        require(_value > _chargeAmt, "ValueGreaterThan0");
        lockBalance(_from, _coinName, _value);
        string[] memory _coins = new string[](1);
        _coins[0] = _coinName;
        uint256[] memory _amounts = new uint256[](1);
        _amounts[0] = _value.sub(_chargeAmt);
        uint256[] memory _fees = new uint256[](1);
        _fees[0] = _chargeAmt;

        //  @dev `_amounts` is a true amount to receive at a destination after deducting a charged fee
        btsPeriphery.sendServiceMessage(_from, _to, _coins, _amounts, _fees);
    }

    /**
       @notice Allow users to transfer multiple coins/wrapped coins to another chain
       @dev Caller must set to approve that the wrapped tokens can be transferred out of the `msg.sender` account by BTSCore contract.
       It MUST revert if the balance of the holder for token `_coinName` is lower than the `_value` sent.
       In case of transferring a native coin, it also checks `msg.value`
       The number of requested coins MUST be as the same as the number of requested values
       The requested coins and values MUST be matched respectively
       @param _coinNames    A list of requested transferring wrapped coins
       @param _values       A list of requested transferring values respectively with its coin name
       @param _to          Target BTP address.
    */
    function transferBatch(
        string[] calldata _coinNames,
        uint256[] memory _values,
        string calldata _to
    ) external payable override {
        require(_coinNames.length == _values.length, "InvalidRequest");
        require(_coinNames.length > 0, "Zero length arguments");
        uint256 size = msg.value != 0
            ? _coinNames.length.add(1)
            : _coinNames.length;
        string[] memory _coins = new string[](size);
        uint256[] memory _amounts = new uint256[](size);
        uint256[] memory _chargeAmts = new uint256[](size);
        Coin memory _coin;
        string memory coinName;
        uint value;

        for (uint256 i = 0; i < _coinNames.length; i++) {
            address _erc20Addresses = coins[_coinNames[i]];
            //  Does not need to check if _coinNames[i] == native_coin
            //  If _coinNames[i] is a native_coin, coins[_coinNames[i]] = 0
            require(_erc20Addresses != address(0), "UnregisterCoin");
            coinName = _coinNames[i];
            value = _values[i];
            require(value > 0,"ZeroOrLess");

            btsPeriphery.checkTransferRestrictions(
                coinName,
                msg.sender,
                value
            );

            IERC20Tradable(_erc20Addresses).transferFrom(
                msg.sender,
                address(this),
                value
            );

            _coin = coinDetails[coinName];
            //  _chargeAmt = fixedFee + msg.value * feeNumerator / FEE_DENOMINATOR
            //  Thus, it's likely that _chargeAmt is always greater than 0
            //  require(_chargeAmt > 0) can be omitted
            _coins[i] = coinName;
            _chargeAmts[i] = value
                .mul(_coin.feeNumerator)
                .div(FEE_DENOMINATOR)
                .add(_coin.fixedFee);
            _amounts[i] = value.sub(_chargeAmts[i]);

            //  Lock this requested _value as a record of a pending transferring transaction
            //  @dev Note that: _value is a requested amount to transfer from a Requester including charged fee
            //  The true amount to receive at a destination receiver is calculated by
            //  _amounts[i] = _values[i].sub(_chargeAmts[i]);
            lockBalance(msg.sender, coinName, value);
        }

        if (msg.value != 0) {
            //  _chargeAmt = fixedFee + msg.value * feeNumerator / FEE_DENOMINATOR
            //  Thus, it's likely that _chargeAmt is always greater than 0
            //  require(_chargeAmt > 0) can be omitted
            _coins[size - 1] = nativeCoinName; // push native_coin at the end of request
            _chargeAmts[size - 1] = msg
                .value
                .mul(coinDetails[nativeCoinName].feeNumerator)
                .div(FEE_DENOMINATOR)
                .add(coinDetails[nativeCoinName].fixedFee);
            _amounts[size - 1] = msg.value.sub(_chargeAmts[size - 1]);
            lockBalance(msg.sender, nativeCoinName, msg.value);
        }

        //  @dev `_amounts` is true amounts to receive at a destination after deducting charged fees
        btsPeriphery.sendServiceMessage(
            msg.sender,
            _to,
            _coins,
            _amounts,
            _chargeAmts
        );
    }

    /**
        @notice Reclaim the token's refundable balance by an owner.
        @dev Caller must be an owner of coin
        The amount to claim must be smaller or equal than refundable balance
        @param _coinName   A given name of coin
        @param _value       An amount of re-claiming tokens
    */
    function reclaim(string calldata _coinName, uint256 _value)
        external
        override
        nonReentrant
    {
        require(
            balances[msg.sender][_coinName].refundableBalance >= _value,
            "Imbalance"
        );

        balances[msg.sender][_coinName].refundableBalance = balances[
            msg.sender
        ][_coinName].refundableBalance.sub(_value);

        this.refund(msg.sender, _coinName, _value);
    }

    //  Solidity does not allow using try_catch with interal/private function
    //  Thus, this function would be set as 'external`
    //  But, it has restriction. It should be called by this contract only
    //  In addition, there are only two functions calling this refund()
    //  + handleRequestService(): this function only called by BTSPeriphery
    //  + reclaim(): this function can be called by ANY
    //  In case of reentrancy attacks, the chance happenning on BTSPeriphery
    //  since it requires a request from BMC which requires verification fron BMV
    //  reclaim() has higher chance to have reentrancy attacks.
    //  So, it must be prevented by adding 'nonReentrant'
    function refund(
        address _to,
        string calldata _coinName,
        uint256 _value
    ) external {
        require(msg.sender == address(this), "Unauthorized");
        if (_coinName.compareTo(nativeCoinName)) {
            paymentTransfer(payable(_to), _value);
        } else {
            IERC20(coins[_coinName]).transfer(_to, _value);
        }
    }

    function paymentTransfer(address payable _to, uint256 _amount) private {
        (bool sent, ) = _to.call{ value: _amount }("");
        require(sent, "PaymentFailed");
    }

    /**
        @notice mint the wrapped coin.
        @dev Caller must be an BTSPeriphery contract
        Invalid _coinName will have an _id = 0. However, _id = 0 is also dedicated to Native Coin
        Thus, BTSPeriphery will check a validity of a requested _coinName before calling
        for the _coinName indicates with id = 0, it should send the Native Coin (Example: PRA) to user account
        @param _to    the account receive the minted coin
        @param _coinName    coin name
        @param _value    the minted amount   
    */
    function mint(
        address _to,
        string calldata _coinName,
        uint256 _value
    ) external override onlyBTSPeriphery {
        if (_coinName.compareTo(nativeCoinName)) {
            paymentTransfer(payable(_to), _value);
        } else if (
            coinDetails[_coinName].coinType == NATIVE_WRAPPED_COIN_TYPE
        ) {
            IERC20Tradable(coins[_coinName]).mint(_to, _value);
        } else if (coinDetails[_coinName].coinType == NON_NATIVE_TOKEN_TYPE) {
            IERC20(coins[_coinName]).transfer(_to, _value);
        }
    }

    /**
        @notice Handle a response of a requested service
        @dev Caller must be an BTSPeriphery contract
        @param _requester   An address of originator of a requested service
        @param _coinName    A name of requested coin
        @param _value       An amount to receive on a destination chain
        @param _fee         An amount of charged fee
    */
    function handleResponseService(
        address _requester,
        string calldata _coinName,
        uint256 _value,
        uint256 _fee,
        uint256 _rspCode
    ) external override onlyBTSPeriphery {
        //  Fee Gathering and Transfer Coin Request use the same method
        //  and both have the same response
        //  In case of Fee Gathering's response, `_requester` is this contract's address
        //  Thus, check that first
        //  -- If `_requester` is this contract's address, then check whethere response's code is RC_ERR
        //  In case of RC_ERR, adding back charged fees to `aggregationFee` state variable
        //  In case of RC_OK, ignore and return
        //  -- Otherwise, handle service's response as normal
        if (_requester == address(this)) {
            if (_rspCode == RC_ERR) {
                aggregationFee[_coinName] = aggregationFee[_coinName].add(
                    _value
                );
            }
            return;
        }
        uint256 _amount = _value.add(_fee);
        balances[_requester][_coinName].lockedBalance = balances[_requester][
            _coinName
        ].lockedBalance.sub(_amount);

        //  A new implementation has been proposed to prevent spam attacks
        //  In receiving error response, BTSCore refunds `_value`, not including `_fee`, back to Requestor
        if (_rspCode == RC_ERR) {
            try this.refund(_requester, _coinName, _value) {} catch {
                balances[_requester][_coinName].refundableBalance = balances[
                    _requester
                ][_coinName].refundableBalance.add(_value);
            }
        } else if (_rspCode == RC_OK) {
            address _erc20Address = coins[_coinName];
            if (
                !_coinName.compareTo(nativeCoinName) &&
                coinDetails[_coinName].coinType == NATIVE_WRAPPED_COIN_TYPE
            ) {
                IERC20Tradable(_erc20Address).burn(address(this), _value);
            }
        }
        aggregationFee[_coinName] = aggregationFee[_coinName].add(_fee);
    }

    /**
        @notice Handle a request of Fee Gathering
            Usage: Copy all charged fees to an array
        @dev Caller must be an BTSPeriphery contract
    */
    function transferFees(string calldata _fa)
        external
        override
        onlyBTSPeriphery
    {
        //  @dev Due to uncertainty in identifying a size of returning memory array
        //  and Solidity does not allow to use 'push' with memory array (only storage)
        //  thus, must use 'temp' storage state variable
        for (uint256 i = 0; i < coinsName.length; i++) {
            if (aggregationFee[coinsName[i]] != 0) {
                chargedCoins.push(coinsName[i]);
                chargedAmounts.push(aggregationFee[coinsName[i]]);
                delete aggregationFee[coinsName[i]];
            }
        }
        btsPeriphery.sendServiceMessage(
            address(this),
            _fa,
            chargedCoins,
            chargedAmounts,
            new uint256[](chargedCoins.length) //  chargedFees is an array of 0 since this is a fee gathering request
        );
        delete chargedCoins;
        delete chargedAmounts;
    }

    function lockBalance(
        address _to,
        string memory _coinName,
        uint256 _value
    ) private {
        balances[_to][_coinName].lockedBalance = balances[_to][_coinName]
            .lockedBalance
            .add(_value);
    }
}