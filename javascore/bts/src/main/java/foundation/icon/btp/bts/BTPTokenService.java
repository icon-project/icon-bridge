/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package foundation.icon.btp.bts;

import java.math.BigInteger;
import java.util.List;
import java.util.Map;

import foundation.icon.btp.bts.irc2.IRC2ScoreInterface;
import foundation.icon.btp.bts.irc2.IRC2SupplierScoreInterface;
import foundation.icon.btp.bts.utils.EnumerableSet;
import foundation.icon.btp.lib.BMCScoreInterface;
import foundation.icon.btp.lib.BSH;
import foundation.icon.btp.lib.BTPAddress;
import foundation.icon.btp.lib.OwnerManager;
import foundation.icon.btp.lib.OwnerManagerImpl;
import foundation.icon.score.util.ArrayUtil;
import foundation.icon.score.util.Logger;
import foundation.icon.score.util.StringUtil;
import score.Address;
import score.BranchDB;
import score.ByteArrayObjectWriter;
import score.Context;
import score.DictDB;
import score.RevertedException;
import score.UserRevertedException;
import score.VarDB;
import score.annotation.EventLog;
import score.annotation.External;
import score.annotation.Optional;
import score.annotation.Payable;
import scorex.util.ArrayList;
import scorex.util.HashMap;

public class BTPTokenService implements BTS, BTSEvents, BSH, OwnerManager {
    private static final Logger logger = Logger.getLogger(BTPTokenService.class);

    public static final String SERVICE = "bts";
    public static final BigInteger NATIVE_COIN_ID = BigInteger.ZERO;
    public static final BigInteger FEE_DENOMINATOR = BigInteger.valueOf(10000);

    public static final int NATIVE_COIN_TYPE = 0;
    public static final int NATIVE_WRAPPED_COIN_TYPE = 1;
    public static final int NON_NATIVE_TOKEN_TYPE = 2;

    public static final Address ZERO_SCORE_ADDRESS = Address.fromString("cx0000000000000000000000000000000000000000");
    public static final BigInteger UINT_CAP = new BigInteger("115792089237316195423570985008687907853269984665640564039457584007913129639935"); //BigInteger.TWO.pow(256).subtract(BigInteger.ONE);

    private final Address bmc;
    private final String net;
    private final byte[] serializedIrc2;
    private final String name;

    private final OwnerManager ownerManager = new OwnerManagerImpl("owners");

    private final EnumerableSet<String > coinNames = new EnumerableSet<>("coinNames", String.class);
    private final DictDB<String, Address> coinAddresses = Context.newDictDB("coinAddresses", Address.class);
    private final DictDB<Address, String> coinAddressName = Context.newDictDB("coinAddressNames", String.class);

    private final BranchDB<String, DictDB<Address, Balance>> balances = Context.newBranchDB("balances", Balance.class);
    private final DictDB<String, BigInteger> feeBalances = Context.newDictDB("feeBalances", BigInteger.class);
    private final DictDB<BigInteger, TransferTransaction> transactions = Context.newDictDB("transactions",
            TransferTransaction.class);
    private final DictDB<BigInteger, BlacklistTransaction> blacklistTxn  = Context.newDictDB("blacklistTransaction",
            BlacklistTransaction.class);
    private final DictDB<BigInteger, TokenLimitTransaction> tokenLimitTxn  = Context.newDictDB("tokenLimitTransaction",
            TokenLimitTransaction.class);
    private final BranchDB<String, DictDB<String, Boolean>> tokenLimitStatus = Context.newBranchDB("tokenLimitStatus",
            Boolean.class);
    private final DictDB<String, Coin> coinDb = Context.newDictDB("coins", Coin.class);
    private final DictDB<String, BigInteger> tokenLimit = Context.newDictDB("tokenLimit", BigInteger.class);
    private final VarDB<BigInteger> sn = Context.newVarDB("serviceNumber", BigInteger.class);

    VarDB<Boolean> restriction = Context.newVarDB("restriction", Boolean.class);
    private final BlacklistDB blacklistDB;

    public BTPTokenService(Address _bmc, String _name, int _decimals,
            BigInteger _feeNumerator, BigInteger _fixedFee, byte[] _serializedIrc2) {
        bmc = _bmc;
        BMCScoreInterface bmcInterface = new BMCScoreInterface(bmc);
        BTPAddress btpAddress = BTPAddress.valueOf(bmcInterface.getBtpAddress());
        net = btpAddress.net();
        name = _name;
        serializedIrc2 = _serializedIrc2;
        blacklistDB = new BlacklistDB();
        restriction.set(true);

        // set sn to zero
        sn.set(BigInteger.ZERO);
        require(_feeNumerator.compareTo(BigInteger.ZERO) >= 0 &&
                        _feeNumerator.compareTo(FEE_DENOMINATOR) < 0,
                "The feeNumerator should be less than FEE_DENOMINATOR and feeNumerator should be greater than 1");
        require(_fixedFee.compareTo(BigInteger.ZERO) >= 0, "Fixed fee cannot be less than zero");


        coinDb.set(_name, new Coin(ZERO_SCORE_ADDRESS, _name, "", _decimals,
                _feeNumerator, _fixedFee, NATIVE_COIN_TYPE));
    }

    @External(readonly = true)
    public String name() {
        return "BTP Token Service";
    }

    /**
     * To change the Coin Fee setting( Fixed fee and fee percentage)
     *
     * @param _name         name of the Coin
     * @param _feeNumerator fee numerator to calculate the fee percentage, Set Zero
     *                      to retain existing value
     * @param _fixedFee     to update the fixed fee
     */
    @External
    public void setFeeRatio(String _name, BigInteger _feeNumerator, BigInteger _fixedFee) {
        requireOwnerAccess();
        Context.require(_feeNumerator.compareTo(BigInteger.ZERO) >= 0 &&
                _feeNumerator.compareTo(FEE_DENOMINATOR) < 0,
                "The feeNumerator should be less than FEE_DENOMINATOR and feeNumerator should be greater than 1");
        Context.require(_fixedFee.compareTo(BigInteger.ZERO) >= 0, "Fixed fee cannot be less than zero");
        require( isRegistered(_name), "Not supported Coin");
        Coin _coin = coinDb.get(_name);
        if (_coin == null) {
            throw BTSException.unknown("Coin Not Registered");
        } else {
            _coin.setFeeNumerator(_feeNumerator);
            _coin.setFixedFee(_fixedFee);
        }
        coinDb.set(_name, _coin);
    }

    @External(readonly = true)
    public Map<String, BigInteger> feeRatio(String _name) {
        Coin coinDetail = coinDb.get(_name);
        if ( coinDetail != null ) {
            return Map.of(
                    "fixedFee", coinDetail.getFixedFee(),
                    "feeNumerator", coinDetail.getFeeNumerator()
            );
        }
        return Map.of(
                "fixedFee", BigInteger.ZERO,
                "feeNumerator", BigInteger.ZERO
        );
    }

    @External
    public void register(String _name, String _symbol, int _decimals, BigInteger _feeNumerator, BigInteger _fixedFee,
            @Optional Address _addr) {
        requireOwnerAccess();

        require(!isRegistered(_name), "already existed");

        coinNames.add(_name);
        if (_addr == null || _addr.equals(ZERO_SCORE_ADDRESS)) {
            Address irc2Address = Context.deploy(serializedIrc2, _name, _symbol, _decimals);
            coinAddresses.set(_name, irc2Address);
            coinAddressName.set(irc2Address, _name);
            coinDb.set(_name, new Coin(irc2Address, _name, _symbol, _decimals, _feeNumerator, _fixedFee,
                    NATIVE_WRAPPED_COIN_TYPE));
        } else {
            coinAddresses.set(_name, _addr);
            coinDb.set(_name,
                    new Coin(_addr, _name, _symbol, _decimals, _feeNumerator, _fixedFee, NON_NATIVE_TOKEN_TYPE));
            coinAddressName.set(_addr, _name);
        }
    }

    @External
    public void setTokenLimit(String[] _coinNames, BigInteger[] _tokenLimits) {
        requireOwnerAccess();
        int size = _coinNames.length;
        require(size == _tokenLimits.length, "Invalid arguments");
        require(size > 0, "Array can't be empty.");
        for (int i = 0; i < size; i++) {
            String coinName = _coinNames[i];
            BigInteger coinLimit = _tokenLimits[i];
            checkUintLimit(coinLimit);
            require((_tokenLimits[i].compareTo(BigInteger.ZERO) >= 0),
                    "Invalid value");
            tokenLimit.set(coinName, coinLimit);
            tokenLimitStatus.at(this.net).set(coinName, true);
        }

        BigInteger sn = increaseSn();
        String[] links = getLinks();
        String[] networks = new String[links.length];
        for (int i = 0; i < links.length; i++) {
            String link = links[i];
            BTPAddress linkAddr = BTPAddress.valueOf(link);
            String net = linkAddr.net();
            networks[i] = net;
            for (String name: _coinNames) {
                tokenLimitStatus.at(net).set(name, false);
            }
            TokenLimitRequest request = new TokenLimitRequest(_coinNames, _tokenLimits, net);
            sendMessage(net, BTSMessage.CHANGE_TOKEN_LIMIT, sn, request.toBytes());
        }

        // to save to tokenLimitTxn Db
        TokenLimitTransaction request = new TokenLimitTransaction(_coinNames, _tokenLimits, networks);
        tokenLimitTxn.set(sn, request);
    }

    @External(readonly = true)
    public BigInteger getTokenLimit(String _name) {
        return tokenLimit.getOrDefault(_name, UINT_CAP);
    }

    @External(readonly = true)
    public TokenLimitTransaction getTokenLimitTxn(BigInteger _sn) {
        return tokenLimitTxn.get(_sn);
    }

    @External(readonly = true)
    public BigInteger getSn() {
        return sn.get();
    }

    @External
    public void addBlacklistAddress(String _net, String[] _addresses) {
        requireOwnerAccess();

        // check for valid link
        require(isValidLink(_net) || _net.equals(net), "Invalid link");

        List<String> blacklist = new ArrayList<>();

        for (String addr: _addresses) {
            addr = lowercase(addr);
            if (! isUserBlackListed(_net, addr) && addr.length() > 0) {
                if (_net.equals(net) && !isValidIconAddress(addr)) {
                    continue;
                }
                blacklist.add(addr);
                blacklistDB.addToBlacklist(_net, addr);
            }
        }

        if (blacklist.size() == 0 || net.equals(_net))  {
            return;
        }

        int size = blacklist.size();
        String[] addresses = new String[size];
        for (int i = 0; i < size; i++) {
            addresses[i] = blacklist.get(i);
        }

        BigInteger sn = increaseSn();
        BlacklistTransaction request = new BlacklistTransaction(
                BlacklistTransaction.ADD_TO_BLACKLIST, addresses, _net);

        blacklistTxn.set(sn, request);

        sendMessage(_net, BTSMessage.BLACKLIST_MESSAGE, sn, request.toBytes());
    }

    @External
    public void removeBlacklistAddress(String _net, String[] _addresses) {
        requireOwnerAccess();

        // check for valid link
        require(isValidLink(_net) || _net.equals(net), "Invalid link");

        List<String> blacklist = new ArrayList<>();

        for (String addr: _addresses) {
            if ( isUserBlackListed(_net, addr)) {
                blacklist.add(addr);
                blacklistDB.removeFromBlacklist(_net, addr);
            }
        }

        if (net.equals(_net) || blacklist.size() == 0) {
            return;
        }

        int size = blacklist.size();
        String[] addresses = new String[size];
        for (int i = 0; i < size; i++) {
            addresses[i] = blacklist.get(i);
        }

        BigInteger sn = increaseSn();
        BlacklistTransaction request = new BlacklistTransaction(
                BlacklistTransaction.REMOVE_FROM_BLACKLIST,addresses, _net);

        blacklistTxn.set(sn, request);

        sendMessage(_net, BTSMessage.BLACKLIST_MESSAGE, sn, request.toBytes());
    }

    @External(readonly = true)
    public boolean isUserBlackListed(String _net, String _address) {
        return blacklistDB.contains(_net, _address);
    }

    @External(readonly = true)
    public List<String> getBlackListedUsers(String _net, int _start, int _end) {
        if ((_end - _start) > 100) {
            throw BTSException.unknown("Can only fetch 100 users at a time");
        }
        return blacklistDB.range(_net, _start, _end);
    }

    @External(readonly = true)
    public int blackListedUsersCount(String _net) {
        return blacklistDB.length(_net);
    }

    @External(readonly = true)
    public int getRegisteredTokensCount() {
        return coinNames.length();
    }

    @External(readonly = true)
    public boolean tokenLimitStatus(String _net, String _coinName) {
        return tokenLimitStatus.at(_net).getOrDefault(_coinName, false);
    }

    @External(readonly = true)
    public List<String> coinNames() {
        // for consistency
        List<String> names = new ArrayList<>();
        names.add(name);
        names.addAll(getCoinNamesAsList());
        return names;
    }

    @External(readonly = true)
    public Address coinId(String _coinName) {
        return this.coinAddresses.getOrDefault(_coinName, null);
    }

    /**
     * Usable => Amount transferred/approved to BTS by owner
     * Tradeable Usable =>  minimumOf(availableBalance, approvedBalance)
     *
     * Locked:
     * IRC Locked: usable <- usable-value,  locked_amount<- locked_amount + value
     * Tradeable IRC Locked: same
     *
     * Refundable:
     * IRC refund:
     *     locked <- locked-value
     *     if failure:
     *         refundable = refundable + value
     *         if canRefund():
     *             refundable = refundable - value
     *
     * Userbalance:
     * IRC          :  account_balance
     * TradeableIRC :  account_balance
     */
    @External(readonly = true)
    public Map<String, BigInteger> balanceOf(Address _owner, String _coinName) {
        Balance balance = getBalance(_coinName, _owner);
        Address _addr = coinAddresses.get(_coinName);
        if (_addr == null && !_coinName.equals(name)) {
            return balance.addUserBalance(BigInteger.ZERO);
        }
        Coin _coin = coinDb.get(_coinName);
        if (_coinName.equals(name)) {
            BigInteger icxBalance = Context.getBalance(_owner);
            return balance.addUserBalance(icxBalance);
        } else if (_coin.getCoinType() == NATIVE_WRAPPED_COIN_TYPE) {
            IRC2SupplierScoreInterface _irc2 = new IRC2SupplierScoreInterface(_coin.getAddress());
            BigInteger allowance = _irc2.allowance(_owner, Context.getAddress());
            BigInteger tokenBalance = _irc2.balanceOf(_owner);
            balance.setUsable(allowance.min(tokenBalance));
            return balance.addUserBalance(tokenBalance);
        } else {
            IRC2ScoreInterface _irc2 = new IRC2ScoreInterface(_coin.getAddress());
            BigInteger tokenBalance = _irc2.balanceOf(_owner);
            return balance.addUserBalance(tokenBalance);
        }
    }

    @External(readonly = true)
    public List<Map<String, BigInteger>> balanceOfBatch(Address _owner, String[] _coinNames) {
        List<Map<String, BigInteger>> balances = new ArrayList<>();
        for (String coinName : _coinNames) {
            balances.add(balanceOf(_owner, coinName));
        }
        return balances;
    }

    @External(readonly = true)
    public Map<String, BigInteger> getAccumulatedFees() {
        Map<String, BigInteger> fees = new HashMap<>();
        for (String coinName: coinNames()) {
            fees.put(coinName, feeBalances.getOrDefault(coinName, BigInteger.ZERO));
        }
        return fees;
    }

    // To receive IRC2 token from existing Contract
    @External
    public void tokenFallback(Address _from, BigInteger _value, byte[] _data) {
        checkUintLimit(_value);
        String _coinName = coinAddressName.get(Context.getCaller());
        if (_coinName != null && _from != Context.getAddress()) {
            Context.require(coinAddresses.get(_coinName) != null, "CoinNotExists");
            Balance _userBalance = getBalance(_coinName, _from);
            _userBalance.setUsable(_userBalance.getUsable().add(_value));
            setBalance(_coinName, _from, _userBalance);
        } else {
            throw BTSException.unknown("Token not registered");
        }
    }

    @External
    public void reclaim(String _coinName, BigInteger _value) {
        require(_value.compareTo(BigInteger.ZERO) > 0, "_value must be positive");
        checkUintLimit(_value);

        Address owner = Context.getCaller();
        Balance balance = getBalance(_coinName, owner);
        require(balance.getRefundable().add(balance.getUsable()).compareTo(_value) > -1, "invalid value");
        require(isRegistered(_coinName), "Not registered");
        balance.setRefundable(balance.getRefundable().add(balance.getUsable()));
        balance.setUsable(BigInteger.ZERO);
        balance.setRefundable(balance.getRefundable().subtract(_value));
        setBalance(_coinName, owner, balance);

        if (name.equals(_coinName)) {
            Context.transfer(owner, _value);
        } else {
            _transferBatch(Context.getAddress(), owner, List.of(_coinName), List.of(_value));
        }
    }

    @Payable
    @External
    public void transferNativeCoin(String _to) {
        BigInteger value = Context.getValue();
        checkUintLimit(value);
        BTPAddress to = BTPAddress.valueOf(_to);
        require(value != null && value.compareTo(BigInteger.ZERO) > 0, "Invalid amount");
        checkRestrictions(name, Context.getCaller().toString(), to, value);
        sendRequest(Context.getCaller(), to, List.of(name), List.of(value));
    }

    @External
    public void transfer(String _coinName, BigInteger _value, String _to) {
        require(_value != null && _value.compareTo(BigInteger.ZERO) > 0, "Invalid amount");
        checkUintLimit(_value);
        require(isRegistered(_coinName), "Not supported Token");

        Address owner = Context.getCaller();
        BTPAddress to = BTPAddress.valueOf(_to);
        checkRestrictions(_coinName, Context.getCaller().toString(), to, _value);
        // only for wrapped coins
        transferFrom(owner, Context.getAddress(), _coinName, _value);
        sendRequest(owner, to, List.of(_coinName), List.of(_value));
    }

    @Payable
    @External
    public void transferBatch(String[] _coinNames, BigInteger[] _values, String _to) {
        require(_coinNames.length == _values.length, "Invalid arguments");
        List<String> coinNameList = new ArrayList<>();
        List<BigInteger> values = new ArrayList<>();
        int len = _coinNames.length;
        require(len > 0, "Zero length arguments");
        Address owner = Context.getCaller();
        BTPAddress to = BTPAddress.valueOf(_to);
        
        for (int i = 0; i < len; i++) {
            String coinName = _coinNames[i];
            BigInteger value = _values[i];
            require(!name.equals(coinName) && this.coinNames.contains(coinName), "Not supported Token");
            checkUintLimit(value);
            coinNameList.add(coinName);
            values.add(_values[i]);
            checkRestrictions(coinName, owner.toString(), to, value);
        }

        transferFromBatch(owner, Context.getAddress(), _coinNames, _values);

        BigInteger value = Context.getValue();
        if (value != null && value.compareTo(BigInteger.ZERO) > 0) {
            coinNameList.add(name);
            values.add(value);
        }
        sendRequest(owner, to, coinNameList, values);
    }
    
    @EventLog(indexed = 1)
    public void TransferStart(Address _from, String _to, BigInteger _sn, byte[] _assetDetails) {
    }

    @EventLog(indexed = 1)
    public void TransferEnd(Address _from, BigInteger _sn, BigInteger _code, byte[] _msg) {
    }

    @EventLog(indexed = 2)
    protected void TransferReceived(String _from, Address _to, BigInteger _sn, byte[] _assetDetails) {
    }

    @EventLog(indexed = 1)
    public void UnknownResponse(String _from, BigInteger _sn) {
    }

    @EventLog(indexed = 1)
    public void AddedToBlacklist(BigInteger sn, byte[] bytes) { }

    @EventLog(indexed = 1)
    public void RemovedFromBlacklist(BigInteger sn, byte[] bytes) { }

    @EventLog(indexed = 1)
    public void TokenLimitSet(BigInteger sn, byte[] bytes) { }

    @External(readonly = true)
    public TransferTransaction getTransaction(BigInteger _sn) {
        return transactions.get(_sn);
    }

    private void sendRequest(Address owner, BTPAddress to, List<String> coinNames, List<BigInteger> amounts) {
        logger.println("sendRequest", "begin");

        int len = coinNames.size();
        AssetTransferDetail[] assetTransferDetails = new AssetTransferDetail[len];
        Asset[] assets = new Asset[len];
        for (int i = 0; i < len; i++) {
            String coinName = coinNames.get(i);
            BigInteger amount = amounts.get(i);
            AssetTransferDetail assetTransferDetail = newAssetTransferDetail(coinName, amount, owner);
            lock(coinName, owner, amount);
            assetTransferDetails[i] = assetTransferDetail;
            assets[i] = new Asset(assetTransferDetail);
        }

        TransferRequest request = new TransferRequest();
        request.setFrom(owner.toString());
        request.setTo(to.account());
        request.setAssets(assets);

        TransferTransaction transaction = new TransferTransaction();
        transaction.setFrom(owner.toString());
        transaction.setTo(to.toString());
        transaction.setAssets(assetTransferDetails);

        BigInteger sn = increaseSn();
        transactions.set(sn, transaction);

        sendMessage(to.net(), BTSMessage.REQUEST_COIN_TRANSFER, sn, request.toBytes());
        TransferStart(owner, to.toString(), sn, encode(assetTransferDetails));
        logger.println("sendRequest", "end");
    }

    static byte[] encode(AssetTransferDetail[] assetTransferDetails) {
        ByteArrayObjectWriter writer = Context.newByteArrayObjectWriter("RLPn");
        writer.beginList(assetTransferDetails.length);
        for (AssetTransferDetail v : assetTransferDetails) {
            writer.write(v);
        }
        writer.end();
        return writer.toByteArray();
    }

    static byte[] encode(Asset[] assets) {
        ByteArrayObjectWriter writer = Context.newByteArrayObjectWriter("RLPn");
        writer.beginList(assets.length);
        for (Asset v : assets) {
            Asset.writeObject(writer, v);
        }
        writer.end();
        return writer.toByteArray();
    }

    private void sendMessage(String net, int serviceType, BigInteger sn, byte[] data) {
        logger.println("sendMessage", "begin");
        BTSMessage message = new BTSMessage();
        message.setServiceType(serviceType);
        message.setData(data);

        BMCScoreInterface bmc = new BMCScoreInterface(this.bmc);
        bmc.sendMessage(net, SERVICE, sn, message.toBytes());
        logger.println("sendMessage", "end");
    }

    private void responseSuccess(String net, BigInteger sn) {
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_OK);
        response.setMessage(TransferResponse.OK_MSG);
        sendMessage(net, BTSMessage.REPONSE_HANDLE_SERVICE, sn, response.toBytes());
    }

    private void responseError(String net, BigInteger sn, String message) {
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_ERR);
        response.setMessage(message);
        sendMessage(net, BTSMessage.REPONSE_HANDLE_SERVICE, sn, response.toBytes());
    }

    /**
     *
     * @param _from net
     */
    @External
    public void handleBTPMessage(String _from, String _svc, BigInteger _sn, byte[] _msg) {
        require(Context.getCaller().equals(bmc), "Only BMC");
        require(_svc.equals(SERVICE), "InvalidSvc");

        BTSMessage message = BTSMessage.fromBytes(_msg);
        int serviceType = message.getServiceType();

        if (serviceType == BTSMessage.REQUEST_COIN_TRANSFER) {
            TransferRequest request = TransferRequest.fromBytes(message.getData());
            handleRequest(request, _from, _sn);
        } else if (serviceType == BTSMessage.REPONSE_HANDLE_SERVICE) {
            TransferResponse response = TransferResponse.fromBytes(message.getData());
            handleResponse(_sn, response);
        } else if (serviceType == BTSMessage.BLACKLIST_MESSAGE) {
            BlacklistResponse response = BlacklistResponse.fromBytes(message.getData());
            handleBlacklist(_sn, response);
        } else if (serviceType == BTSMessage.CHANGE_TOKEN_LIMIT) {
            TokenLimitResponse response = TokenLimitResponse.fromBytes(message.getData());
            handleChangeTokenLimit(_from, _sn, response);
        } else if (serviceType == BTSMessage.UNKNOWN_TYPE) {
            // If receiving a RES_UNKNOWN_TYPE, ignore this message
            // or re-send another correct message
            UnknownResponse(_from, _sn);
        } else {
            // If none of those types above, BSH responds a message of RES_UNKNOWN_TYPE
            TransferResponse response = new TransferResponse();
            response.setCode(TransferResponse.RC_ERR);
            response.setMessage(TransferResponse.ERR_MSG_UNKNOWN_TYPE);
            sendMessage(_from, BTSMessage.UNKNOWN_TYPE, _sn, response.toBytes());
        }
    }

    /**
     *
     * @param _src net
     */
    @External
    public void handleBTPError(String _src, String _svc, BigInteger _sn, long _code, String _msg) {
        require(Context.getCaller().equals(bmc), "Only BMC");

        TransferTransaction tTxn = transactions.get(_sn);
        if (tTxn != null) {
            TransferResponse response = new TransferResponse();
            response.setCode(TransferResponse.RC_ERR);
            response.setMessage("BTPError [code:" + _code + ",msg:" + _msg);
            handleResponse(_sn, response);
            return;
        }

        TokenLimitTransaction tlTxn = tokenLimitTxn.get(_sn);
        if (tlTxn != null) {
            String[] coinNames = tlTxn.getCoinName();
            int size = coinNames.length;
            for (String coinName : coinNames) {
                tokenLimitStatus.at(_src).set(coinName, false);
            }
            return;
        }

        BlacklistTransaction bTxn = blacklistTxn.get(_sn);
        if (bTxn != null) {
            Integer service = bTxn.getServiceType();
            if (service == BlacklistTransaction.ADD_TO_BLACKLIST) {
                handleAddToBlacklistFailResponse(bTxn);
            } else if (service == BlacklistTransaction.REMOVE_FROM_BLACKLIST) {
                handleRemoveFromBlacklistFailResponse(bTxn);
            }
            else {
                throw BTSException.unknown("BLACKLIST HANDLE ERROR");
            }
            return;
        }
    }

    @External
    public void handleFeeGathering(String _fa, String _svc) {
        require(Context.getCaller().equals(bmc), "Only BMC");
        BTPAddress from = BTPAddress.valueOf(_fa);
        Address owner = Context.getAddress();

        List<String> coinNames = new ArrayList<>();
        List<BigInteger> feeAmounts = new ArrayList<>();
        for (String coinName : coinNames()) {
            BigInteger feeAmount = clearFee(coinName);
            if (feeAmount.compareTo(BigInteger.ZERO) > 0) {
                coinNames.add(coinName);
                feeAmounts.add(feeAmount);
            }
        }

        if (coinNames.size() > 0) {
            if (from.net().equals(net)) {
                Address fa = Address.fromString(from.account());
                int idx = coinNames.indexOf(name);
                if (idx >= 0) {
                    coinNames.remove(idx);
                    BigInteger feeAmount = feeAmounts.remove(idx);
                    Context.transfer(fa, feeAmount);
                }
                _transferBatch(owner, fa, coinNames, feeAmounts);
            } else {
                sendRequest(owner, from, coinNames, feeAmounts);
            }
        }
    }

    public Balance getBalance(String coinName, Address owner) {
        Balance balance = balances.at(coinName).get(owner);
        if (balance == null) {
            balance = new Balance();
            balance.setUsable(BigInteger.ZERO);
            balance.setLocked(BigInteger.ZERO);
            balance.setRefundable(BigInteger.ZERO);
        }
        return balance;
    }

    private void setBalance(String coinName, Address owner, Balance balance) {
        balances.at(coinName).set(owner, balance);
    }

    private void lock(String coinName, Address owner, BigInteger value) {
        logger.println("lock", "coinName:", coinName, "owner:", owner, "value:", value);
        Balance balance = getBalance(coinName, owner);
        Coin _coin = coinDb.get(coinName);
        if (_coin.getCoinType() == NON_NATIVE_TOKEN_TYPE) {
            if (balance.getUsable().compareTo(value) >= 0) {
                balance.setUsable(balance.getUsable().subtract(value));
            } else {
                Context.revert("InSufficient Usable Balance");
            }
        }
        balance.setLocked(balance.getLocked().add(value));
        setBalance(coinName, owner, balance);
    }

    private void unlock(String coinName, Address owner, BigInteger value) {
        logger.println("unlock", "coinName:", coinName, "owner:", owner, "value:", value);
        Balance balance = getBalance(coinName, owner);
        balance.setLocked(balance.getLocked().subtract(value));
        setBalance(coinName, owner, balance);
    }

    private void refund(String coinName, Address owner, BigInteger locked, BigInteger fee) {
        logger.println("refund", "coinName:", coinName, "owner:", owner, "locked:", locked, "fee: ", fee);
        // unlock and add refundable
        Balance balance = getBalance(coinName, owner);
        BigInteger value = locked.subtract(fee);
        balance.setLocked(balance.getLocked().subtract(locked));
        try {
            if (name.equals(coinName)) {
                Context.transfer(owner, value);
            } else {
                _transferBatch(Context.getAddress(), owner, List.of(coinName), List.of(value));
            }
        } catch (Exception e) {
            if (!owner.equals(Context.getAddress())) {
                balance.setRefundable(balance.getRefundable().add(value));
            }
        }
        setBalance(coinName, owner, balance);
    }

    private void addFee(String coinName, BigInteger amount) {
        BigInteger fee = feeBalances.getOrDefault(coinName, BigInteger.ZERO);
        feeBalances.set(coinName, fee.add(amount));
    }

    private BigInteger clearFee(String coinName) {
        BigInteger fee = feeBalances.getOrDefault(coinName, BigInteger.ZERO);
        if (fee.compareTo(BigInteger.ZERO) > 0) {
            feeBalances.set(coinName, BigInteger.ZERO);
        }
        return fee;
    }

    private void handleRequest(TransferRequest request, String from, BigInteger sn) {
        logger.println("handleRequest", "begin", "sn:", sn);
        Address to;
        try {
            to = Address.fromString(request.getTo());
        } catch (Exception e) {
            throw BTSException.unknown(e.getMessage());
        }

        BigInteger nativeCoinTransferAmount = null;
        Asset[] assets = request.getAssets();
        List<String> coinNamesList = new ArrayList<>();
        List<BigInteger> amounts = new ArrayList<>();
        for (Asset asset : assets) {
            String coinName = asset.getCoinName();
            BigInteger amount = asset.getAmount();
            if (amount == null || amount.compareTo(BigInteger.ZERO) < 1) {
                throw BTSException.unknown("Amount must be positive value");
            }

            // not coin
            if (this.coinNames.contains(coinName)) {
                coinNamesList.add(coinName);
                amounts.add(amount);
            } // nativecoin
            else if (name.equals(coinName)) {
                nativeCoinTransferAmount = amount;
            } else {
                throw BTSException.unknown("Invalid Token");
            }
            checkTransferRestrictions(net, coinName, from, request.getTo(), amount);
        }

        if (nativeCoinTransferAmount != null) {
            try {
                Context.transfer(to, nativeCoinTransferAmount);
            } catch (Exception e) {
                throw BTSException.unknown("fail to transfer err:" + e.getMessage());
            }
        }

        if (coinNamesList.size() > 0) {
            mintBatch(to, ArrayUtil.toStringArray(coinNamesList), ArrayUtil.toBigIntegerArray(amounts));
        }

        logger.println("handleRequest", "responseSuccess");
        responseSuccess(from, sn);
        logger.println("handleRequest", "end");
        TransferReceived(from, to, sn, encode(assets));
    }

    private void handleResponse(BigInteger sn, TransferResponse response) {
        logger.println("handleResponse", "begin", "sn:", sn);
        TransferTransaction transaction = transactions.get(sn);
        // ignore when not exists pending request
        if (transaction != null) {
            BigInteger code = response.getCode();
            Address owner = Address.fromString(transaction.getFrom());
            AssetTransferDetail[] assets = transaction.getAssets();

            logger.println("handleResponse", "code:", code);
            if (TransferResponse.RC_OK.equals(code)) {
                List<String> coinNameList = new ArrayList<>();
                List<BigInteger> amounts = new ArrayList<>();
                for (AssetTransferDetail asset : assets) {
                    String coinName = asset.getCoinName();
                    BigInteger amount = asset.getAmount();
                    BigInteger fee = asset.getFee();
                    BigInteger locked = amount.add(fee);
                    boolean isNativeCoin = name.equals(coinName);
                    if (isRegistered(coinName)) {
                        unlock(coinName, owner, locked);
                        addFee(coinName, fee);
                        if (!isNativeCoin) {
                            coinNameList.add(coinName);
                            amounts.add(amount);
                        }
                    } else {
                        // This should not happen
                        throw BTSException.unknown("invalid transaction, invalid coinName");
                    }
                }

                if (coinNameList.size() > 0) {
                    burnBatch(ArrayUtil.toStringArray(coinNameList), ArrayUtil.toBigIntegerArray(amounts));
                }
            } else {
                for (AssetTransferDetail asset : assets) {
                    String coinName = asset.getCoinName();
                    BigInteger amount = asset.getAmount();
                    BigInteger fee = asset.getFee();
                    BigInteger locked = amount.add(fee);
                    if (isRegistered(coinName)) {
                        refund(coinName, owner, locked, fee);
                    } else {
                        // This should not happen
                        throw BTSException.unknown("invalid transaction, invalid coinName");
                    }
                }
            }

            transactions.set(sn, null);
            TransferEnd(owner, sn, code, response.getMessage() != null ? response.getMessage().getBytes() : null);
        }
        logger.println("handleResponse", "end");
    }

    private void handleBlacklist(BigInteger sn, BlacklistResponse response) {
        BlacklistTransaction txn = blacklistTxn.get(sn);
        if (txn != null) {
            Integer serviceType = txn.getServiceType();
            BigInteger code = response.getCode();
            if (serviceType == BlacklistTransaction.ADD_TO_BLACKLIST) {
                logger.println("handleAddToBlacklist", "begin", "sn:", sn);
                if (BlacklistResponse.RC_OK.equals(code)) {
                    blacklistTxn.set(sn, null);
                    AddedToBlacklist(sn, response.getMessage() != null ? response.getMessage().getBytes() : null);
                } else {
                    handleAddToBlacklistFailResponse(txn);
                }
                logger.println("handleAddToBlacklist", "end");
            } else if (serviceType == BlacklistTransaction.REMOVE_FROM_BLACKLIST) {
                logger.println("handleRemoveFromBlacklist", "begin", "sn:", sn);
                if (BlacklistResponse.RC_OK.equals(code)) {
                    blacklistTxn.set(sn, null);
                    RemovedFromBlacklist(sn, response.getMessage() != null ? response.getMessage().getBytes() : null);
                } else {
                    handleRemoveFromBlacklistFailResponse(txn);
                }
                logger.println("handleRemoveFromBlacklist", "end");
            } else {
                throw BTSException.unknown("Invalid Blacklist Txn");
            }
        }
    }

    private void handleAddToBlacklistFailResponse(BlacklistTransaction txn) {
        String[] addresses = txn.getAddress();
        String net = txn.getNet();
        for(String addr: addresses) {
            removeFromBlacklistInternal(net, addr);
        }
    }

    private void handleRemoveFromBlacklistFailResponse(BlacklistTransaction txn) {
        String[] addresses = txn.getAddress();
        String net = txn.getNet();
        for(String addr : addresses) {
            addToBlacklistInternal(net, addr);
        }
    }

    private void handleChangeTokenLimit(String net, BigInteger sn, TokenLimitResponse response) {
        logger.println("handleChangeTokenLimit", "begin", "sn:", sn);
        TokenLimitTransaction txn = tokenLimitTxn.get(sn);
        if (txn != null) {
            BigInteger code = response.getCode();
            if (BlacklistResponse.RC_OK.equals(code)) {
                String[] newNetworks = removeFromArray(txn.getNet(), net);
                txn.setNet(newNetworks);
                String[] coinNames = txn.getCoinName();
                for (String coinName : coinNames) {
                    tokenLimitStatus.at(net).set(coinName, true);
                }
                if (newNetworks.length == 0) {
                    tokenLimitTxn.set(sn, null);
                } else {
                    tokenLimitTxn.set(sn, txn);
                }
            } else {
                throw BTSException.unknown("Invalid change limit transaction");
            }
            TokenLimitSet(sn, response.getMessage() != null ? response.getMessage().getBytes() : null);

        }
        logger.println("handleChangeTokenLimit", "end");
    }

    private String[] removeFromArray(String[] arr, String element) {
        boolean inArray = isInArray(arr, element);
        if ( inArray ) {
            int size = arr.length;
            String[] newArr = new String[size - 1];
            for (int i = 0, k = 0; i < size; i++) {
                if (!arr[i].equals(element)) {
                    newArr[k] = arr[i];
                    k++;
                }
            }
            return newArr;
        }
        return arr;
    }

    private boolean isInArray(String[] arr, String element) {
        for (String s : arr) {
            if (element.equals(s)) {
                return true;
            }
        }
        return false;
    }

    private void addToBlacklistInternal(String net, String addr) {
        blacklistDB.addToBlacklist(net, addr);
    }

    private void removeFromBlacklistInternal(String net, String addr) {
        blacklistDB.removeFromBlacklist(net, addr);
    }

    private AssetTransferDetail newAssetTransferDetail(String coinName, BigInteger amount, Address owner) {
        logger.println("newAssetTransferDetail", "begin");
        Coin _coin = coinDb.get(coinName);
        if (_coin == null) {
            throw BTSException.unknown("Coin Not Registered");
        }
        BigInteger feeRatio = _coin.getFeeNumerator();
        BigInteger fixedFee = _coin.getFixedFee();
        if (owner.equals(Context.getAddress())) {
            feeRatio = BigInteger.ZERO;
        }
        BigInteger fee = amount.multiply(feeRatio).divide(FEE_DENOMINATOR).add(fixedFee);
        if (feeRatio.compareTo(BigInteger.ZERO) > 0 && fee.compareTo(BigInteger.ZERO) == 0) {
            fee = BigInteger.ONE;
        }
        BigInteger transferAmount = amount.subtract(fee);
        logger.println("newAssetTransferDetail", "amount:", amount, "fee:", fee);
        if (transferAmount.compareTo(BigInteger.ZERO) < 1) {
            throw BTSException.unknown("not enough value");
        }
        AssetTransferDetail asset = new AssetTransferDetail();
        asset.setCoinName(coinName);
        asset.setAmount(transferAmount);
        asset.setFee(fee);
        logger.println("newAssetTransferDetail", "end");
        return asset;
    }

    /* Intercall with IRC2Supplier */
    private void transferFrom(Address from, Address to, String coinName, BigInteger amount) {
        logger.println("transferFrom", from, to, coinName, amount);
        try {
            Coin _coin = coinDb.get(coinName);
            if (_coin.getCoinType() == NATIVE_WRAPPED_COIN_TYPE) {
                Address coinAddress = this.getCoinAddress(coinName);
                IRC2SupplierScoreInterface irc2 = new IRC2SupplierScoreInterface(coinAddress);
                irc2.transferFrom(from, to, amount, null);
            }
        } catch (UserRevertedException e) {
            logger.println("transferFrom", "code:", e.getCode(), "msg:", e.getMessage());
            throw BTSException.irc31Reverted("code:" + e.getCode() + "msg:" + e.getMessage());
        } catch (IllegalArgumentException | RevertedException e) {
            logger.println("transferFrom", "Exception:", e.toString());
            throw BTSException.irc31Failure("Exception:" + e);
        }
    }

    private void approve(Address from, String coinName, BigInteger amount) {
        logger.println("approve", from, coinName, amount);
        Address coinAddress = this.getCoinAddress(coinName);
        IRC2SupplierScoreInterface irc2 = new IRC2SupplierScoreInterface(coinAddress);
        try {
            irc2.approve(from, amount);
        } catch (UserRevertedException e) {
            logger.println("approve", "code:", e.getCode(), "msg:", e.getMessage());
            throw BTSException.irc31Reverted("code:" + e.getCode() + "msg:" + e.getMessage());
        } catch (IllegalArgumentException | RevertedException e) {
            logger.println("approve", "Exception:", e.toString());
            throw BTSException.irc31Failure("Exception:" + e);
        }
    }

    private void _transferBatch(Address from, Address to, List<String> coinNames, List<BigInteger> amounts) {
        logger.println("transferFromBatch", from, to, coinNames, StringUtil.toString(amounts));
        int len = coinNames.size();
        for (int i = 0; i < len; i++) {
            Coin _coin = coinDb.get(coinNames.get(i));
            if (_coin.getCoinType() == NATIVE_WRAPPED_COIN_TYPE) {
                this.approve(from, coinNames.get(i), amounts.get(i));
                this.transferFrom(from, to, coinNames.get(i), amounts.get(i));
            } else if (_coin.getCoinType() == NON_NATIVE_TOKEN_TYPE) {
                this.transfer(to, coinNames.get(i), amounts.get(i));
            }
        }
    }

    private void transferFromBatch(Address from, Address to, String[] coinNames, BigInteger[] amounts) {
        logger.println("transferFromBatch", from, to, coinNames, StringUtil.toString(amounts));
        for (int i = 0; i < coinNames.length; i++) {
            this.transferFrom(from, to, coinNames[i], amounts[i]);
        }
    }

    private void mint(Address to, String coinName, BigInteger amount) {
        logger.println("mint", to, coinName, amount);
        try {
            Coin _coin = coinDb.get(coinName);
            if (_coin.getCoinType() == NATIVE_WRAPPED_COIN_TYPE) {
                IRC2SupplierScoreInterface irc2 = new IRC2SupplierScoreInterface(_coin.getAddress());
                irc2.mint(to, amount);
            } else {
                this.transfer(to, coinName, amount);
            }
        } catch (UserRevertedException e) {
            logger.println("mint", "code:", e.getCode(), "msg:", e.getMessage());
            throw BTSException.irc31Reverted("code:" + e.getCode() + "msg:" + e.getMessage());
        } catch (IllegalArgumentException | RevertedException e) {
            logger.println("mint", "Exception:", e.toString());
            throw BTSException.irc31Failure("Exception:" + e);
        }
    }

    private void transfer(Address to, String coinName, BigInteger amount) {
        Address _coinAddr = coinId(coinName);
        if (_coinAddr == null) {
            throw BTSException.irc31Failure("Exception: CoinNotFound");
        }
        try {
            IRC2ScoreInterface _irc2 = new IRC2ScoreInterface(_coinAddr);
            _irc2.transfer(to, amount, null);
        } catch (UserRevertedException e) {
            logger.println("transfer", "code:", e.getCode(), "msg:", e.getMessage());
            throw BTSException.irc31Reverted("code:" + e.getCode() + "msg:" + e.getMessage());
        } catch (IllegalArgumentException | RevertedException e) {
            logger.println("transfer", "Exception:", e.toString());
            throw BTSException.irc31Failure("Exception:" + e);
        }
    }

    private void mintBatch(Address to, String[] coinNames, BigInteger[] amounts) {
        logger.println("mintBatch", to, StringUtil.toString(coinNames), StringUtil.toString(amounts));
        for (int i = 0; i < coinNames.length; i++) {
            this.mint(to, coinNames[i], amounts[i]);
        }
    }

    private void burn(String coinName, BigInteger amount) {
        logger.println("burn", coinName, amount);
        try {
            Coin _coin = coinDb.get(coinName);
            if (_coin.getCoinType() == NATIVE_WRAPPED_COIN_TYPE) {
                IRC2SupplierScoreInterface irc2 = new IRC2SupplierScoreInterface(_coin.getAddress());
                irc2.burn(amount);
            }
        } catch (UserRevertedException e) {
            logger.println("burn", "code:", e.getCode(), "msg:", e.getMessage());
            throw BTSException.irc31Reverted("code:" + e.getCode() + "msg:" + e.getMessage());
        } catch (IllegalArgumentException | RevertedException e) {
            logger.println("burn", "Exception:", e.toString());
            throw BTSException.irc31Failure("Exception:" + e);
        }
    }

    private void burnBatch(String[] coinNames, BigInteger[] amounts) {
        logger.println("burnBatch", StringUtil.toString(coinNames), StringUtil.toString(amounts));
        for (int i = 0; i < coinNames.length; i++) {
            this.burn(coinNames[i], amounts[i]);
        }
    }

    private Address getCoinAddress(String coinName) {
        return this.coinAddresses.get(coinName);
    }

    /* Delegate OwnerManager */
    private void requireOwnerAccess() {
        if (!ownerManager.isOwner(Context.getCaller())) {
            throw BTSException.unauthorized("require owner access");
        }
    }

    @External
    public void addOwner(Address _addr) {
        try {
            ownerManager.addOwner(_addr);
        } catch (IllegalStateException e) {
            throw BTSException.unauthorized(e.getMessage());
        } catch (IllegalArgumentException e) {
            throw BTSException.unknown(e.getMessage());
        }
    }

    @External
    public void removeOwner(Address _addr) {
        try {
            ownerManager.removeOwner(_addr);
        } catch (IllegalStateException e) {
            throw BTSException.unauthorized(e.getMessage());
        } catch (IllegalArgumentException e) {
            throw BTSException.unknown(e.getMessage());
        }
    }

    @External(readonly = true)
    public Address[] getOwners() {
        return ownerManager.getOwners();
    }

    @External(readonly = true)
    public boolean isOwner(Address _addr) {
        return ownerManager.isOwner(_addr);
    }

    @External
    public void addRestriction() {
        requireOwnerAccess();
        restriction.set(true);
    }

    @External
    public void disableRestrictions() {
        requireOwnerAccess();
        restriction.set(false);
    }

    @External(readonly = true)
    public boolean isRestrictionEnabled() {
        return restriction.get();
    }

    private void checkRestrictions(String coinName, String from, BTPAddress to, BigInteger value) {
        checkTransferRestrictions(to.net(), coinName, from, to.account(), value);
        checkTransferRestrictions(net, coinName, from, to.account(), value);
    }

    private void checkTransferRestrictions(String _net, String _tokenName, String _from, String _to, BigInteger _value) {
        if (restriction.get() != null && restriction.get()) {
            validateRestriction(_net, _tokenName, _from, _to, _value);
        }
    }

    private void validateRestriction(String _net, String _token, String _from, String _to, BigInteger _value) {
        if (isUserBlackListed(_net, _from)) {
            throw BTSException.restricted("_from user is Blacklisted");
        }
        if (isUserBlackListed(_net, _to)) {
            throw BTSException.restricted("_to user is Blacklisted");
        }
        BigInteger tokenLimit = getTokenLimit(_token);
        if (_value.compareTo(tokenLimit) > 0) {
            throw BTSException.restricted("Transfer amount exceeds the transaction limit");
        }
    }

    private boolean isRegistered(String name) {
        return coinNames.contains(name) || this.name.equals(name);
    }

    private List<String> getCoinNamesAsList() {
        return coinNames.range(0, coinNames.length());
    }

    private boolean isValidIconAddress(String str) {
        try {
            Address.fromString(str);
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    static void require(boolean condition, String message) {
        if (!condition) {
            throw BTSException.unknown(message);
        }
    }

    private BigInteger increaseSn() {
        BigInteger newSn = sn.get().add(BigInteger.ONE);
        sn.set(newSn);
        return newSn;
    }

    private String lowercase(String word) {
        return word.trim().toLowerCase();
    }

    private String[] getLinks() {
        BMCScoreInterface bmcInterface = new BMCScoreInterface(bmc);
        String[] links = bmcInterface.getLinks();
        return links;
    }

    private boolean isValidLink(String net) {
        String[] links = getLinks();
        for (String link : links) {
            BTPAddress btpAddress = BTPAddress.valueOf(link);
            if (btpAddress.net().equals(net)) {
                return true;
            }
        }
        return false;
    }

    private void checkUintLimit(BigInteger value) {
        require(UINT_CAP.compareTo(value) >= 0, "Value cannot exceed uint(256)-1");
    }

}