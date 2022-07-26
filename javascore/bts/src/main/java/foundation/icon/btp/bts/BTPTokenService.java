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

import foundation.icon.btp.lib.*;
import foundation.icon.btp.bts.irc2.IRC2ScoreInterface;
import foundation.icon.btp.bts.irc2.IRC2SupplierScoreInterface;
import foundation.icon.btp.restrictions.RestrictionsScoreInterface;
import foundation.icon.score.util.ArrayUtil;
import foundation.icon.score.util.Logger;
import foundation.icon.score.util.StringUtil;
import java.util.Map;
import score.*;
import score.annotation.EventLog;
import score.annotation.External;
import score.annotation.Optional;
import score.annotation.Payable;
import scorex.util.ArrayList;

import java.math.BigInteger;
import java.util.List;

public class BTPTokenService implements BTS, BTSEvents, BSH, OwnerManager {
    private static final Logger logger = Logger.getLogger(BTPTokenService.class);

    public static final String SERVICE = "bts";
    public static final BigInteger NATIVE_COIN_ID = BigInteger.ZERO;
    public static final BigInteger FEE_DENOMINATOR = BigInteger.valueOf(10000);

    public static final int NATIVE_COIN_TYPE = 0;
    public static final int NATIVE_WRAPPED_COIN_TYPE = 1;
    public static final int NON_NATIVE_TOKEN_TYPE = 2;

    //
    private final Address bmc;
    private final String net;
    private final byte[] serializedIrc2;
    private final String name;
    private final VarDB<BTSProperties> properties = Context.newVarDB("properties", BTSProperties.class);

    //
    private final OwnerManager ownerManager = new OwnerManagerImpl("owners");

    //
    private final ArrayDB<String> coinNames = Context.newArrayDB("coinNames", String.class);
    private final DictDB<String, Address> coinAddresses = Context.newDictDB("coinAddresses", Address.class);

    //
    private final BranchDB<String, DictDB<Address, Balance>> balances = Context.newBranchDB("balances", Balance.class);
    private final DictDB<String, BigInteger> feeBalances = Context.newDictDB("feeBalances", BigInteger.class);
    private final DictDB<BigInteger, TransferTransaction> transactions = Context.newDictDB("transactions",
            TransferTransaction.class);
    private final DictDB<String, Coin> coinDb = Context.newDictDB("coins", Coin.class);
    private final DictDB<Address, String> coinAddressName = Context.newDictDB("coinAddressNames", String.class);

    //
    private final VarDB<Address> bsrDb = Context.newVarDB("bsr", Address.class);
    VarDB<Boolean> restriction = Context.newVarDB("restricton", Boolean.class);
    RestrictionsScoreInterface restrictionsInterface;

    public BTPTokenService(Address _bmc, String _name, byte[] _serializedIrc2) {
        bmc = _bmc;
        BMCScoreInterface bmcInterface = new BMCScoreInterface(bmc);
        BTPAddress btpAddress = BTPAddress.valueOf(bmcInterface.getBtpAddress());
        net = btpAddress.net();
        name = _name;
        serializedIrc2 = _serializedIrc2;
        coinDb.set(_name, new Coin(Address.fromString("cx0000000000000000000000000000000000000000"), _name, _name, 0,
                BigInteger.ZERO, BigInteger.ZERO, NATIVE_COIN_TYPE));
    }

    public BTSProperties getProperties() {
        return properties.getOrDefault(BTSProperties.DEFAULT);
    }

    public void setProperties(BTSProperties properties) {
        this.properties.set(properties);
    }

    private boolean isRegistered(String name) {
        int len = coinNames.size();
        for (int i = 0; i < len; i++) {
            if (coinNames.get(i).equals(name)) {
                return true;
            }
        }
        return false;
    }

    private List<String> getCoinNamesAsList() {
        List<String> coinNames = new ArrayList<>();
        int len = this.coinNames.size();
        for (int i = 0; i < len; i++) {
            coinNames.add(this.coinNames.get(i));
        }
        return coinNames;
    }

    static void require(boolean condition, String message) {
        if (!condition) {
            throw BTSException.unknown(message);
        }
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
        Context.require(_feeNumerator.compareTo(BigInteger.ONE) >= 0 &&
                _feeNumerator.compareTo(FEE_DENOMINATOR) < 0,
                "The feeNumerator should be less than FEE_DENOMINATOR and feeNumerator should be greater than 1");
        require(name.equals(_name) || isRegistered(_name), "Not supported Coin");
        Coin _coin = coinDb.get(_name);
        if (_coin == null) {
            throw BTSException.unknown("Coin Not Registered");
        } else {
            _coin.setFeeNumerator(_feeNumerator);
            _coin.setFixedFee(_fixedFee);
        }
        coinDb.set(_name, _coin);
    }

    @External
    public void register(String _name, String _symbol, int _decimals, BigInteger _feeNumerator, BigInteger _fixedFee,
            @Optional Address _addr) {
        requireOwnerAccess();

        require(!name.equals(_name) && !isRegistered(_name), "already existed");
        coinNames.add(_name);
        if (_addr == null) {
            Address irc2Address = Context.deploy(serializedIrc2, _name, _symbol, _decimals);
            coinAddresses.set(_name, irc2Address);
            coinDb.set(_name, new Coin(irc2Address, _name, _symbol, _decimals, _feeNumerator, _fixedFee,
                    NATIVE_WRAPPED_COIN_TYPE));
        } else {
            coinAddresses.set(_name, _addr);
            coinDb.set(_name,
                    new Coin(_addr, _name, _symbol, _decimals, _feeNumerator, _fixedFee, NON_NATIVE_TOKEN_TYPE));
            coinAddressName.set(_addr, _name);
        }
    }

    @External(readonly = true)
    public String[] coinNames() {
        int len = coinNames.size();
        String[] names = new String[len + 1];
        names[0] = name;
        for (int i = 0; i < len; i++) {
            names[i + 1] = coinNames.get(i);
        }
        return names;
    }

    @External(readonly = true)
    public Address coinId(String _coinName) {
        return this.coinAddresses.getOrDefault(_coinName, null);
    }

    @External(readonly = true)
    public Map<String, BigInteger> balanceOf(Address _owner, String _coinName) {
        Balance balance = getBalance(_coinName, _owner);
        Address _addr = coinAddresses.get(_coinName);
        if (_addr == null) {
            return balance.addUserBalance(BigInteger.ZERO);
        }
        Coin _coin = coinDb.get(_coinName);
        if (_coinName.equals(name)) {
            BigInteger icxBalance = Context.getBalance(_owner);
            balance.setUsable(icxBalance);
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

    // To receive IRC2 token from existing Contract
    @External
    public void tokenFallback(Address _from, BigInteger _value, byte[] _data) {
        String _coinName = coinAddressName.get(Context.getCaller());
        if (_coinName != null && _from != Context.getAddress()) {
            Context.require(coinAddresses.get(_coinName) != null, "CoinNotExists");
            Balance _userBalance = getBalance(_coinName, _from);
            _userBalance.setUsable(_userBalance.getUsable().add(_value));
            setBalance(_coinName, _from, _userBalance);
        }
    }

    @External
    public void reclaim(String _coinName, BigInteger _value) {
        require(_value.compareTo(BigInteger.ZERO) > 0, "_value must be positive");

        Address owner = Context.getCaller();
        Balance balance = getBalance(_coinName, owner);
        require(balance.getRefundable().add(balance.getUsable()).compareTo(_value) > -1, "invalid value");
        Coin _coin = coinDb.get(_coinName);
        if (_coin.getCoinType() == NON_NATIVE_TOKEN_TYPE) {
            balance.setRefundable(balance.getRefundable().add(balance.getUsable()));
            balance.setUsable(BigInteger.ZERO);
        }
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
        require(value != null && value.compareTo(BigInteger.ZERO) > 0, "Invalid amount");
        checkTransferRestrictions(name, Context.getCaller().toString(), BTPAddress.valueOf(_to).account(), value);
        sendRequest(Context.getCaller(), BTPAddress.valueOf(_to), List.of(name), List.of(value));
    }

    @External
    public void transfer(String _coinName, BigInteger _value, String _to) {
        require(_value != null && _value.compareTo(BigInteger.ZERO) > 0, "Invalid amount");
        require(!name.equals(_coinName) && isRegistered(_coinName), "Not supported Token");

        Address owner = Context.getCaller();
        checkTransferRestrictions(_coinName, owner.toString(), BTPAddress.valueOf(_to).account(), _value);
        transferFrom(owner, Context.getAddress(), _coinName, _value);
        sendRequest(owner, BTPAddress.valueOf(_to), List.of(_coinName), List.of(_value));
    }

    @Payable
    @External
    public void transferBatch(String[] _coinNames, BigInteger[] _values, String _to) {
        require(_coinNames.length == _values.length, "Invalid arguments");

        List<String> registeredCoinNames = getCoinNamesAsList();
        List<String> coinNames = new ArrayList<>();
        List<BigInteger> values = new ArrayList<>();
        int len = _coinNames.length;
        for (int i = 0; i < len; i++) {
            String coinName = _coinNames[i];
            require(!name.equals(coinName) && registeredCoinNames.contains(coinName), "Not supported Token");
            coinNames.add(coinName);
            values.add(_values[i]);
        }

        Address owner = Context.getCaller();
        transferFromBatch(owner, Context.getAddress(), _coinNames, _values);

        BigInteger value = Context.getValue();
        if (value != null && value.compareTo(BigInteger.ZERO) > 0) {
            coinNames.add(name);
            values.add(value);
        }
        sendRequest(owner, BTPAddress.valueOf(_to), coinNames, values);
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

    @External(readonly = true)
    public TransferTransaction getTransaction(BigInteger _sn) {
        return transactions.get(_sn);
    }

    private void sendRequest(Address owner, BTPAddress to, List<String> coinNames, List<BigInteger> amounts) {
        logger.println("sendRequest", "begin");
        BTSProperties properties = getProperties();

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

        BigInteger sn = properties.getSn().add(BigInteger.ONE);
        properties.setSn(sn);
        setProperties(properties);
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

    private void sendMessage(String to, int serviceType, BigInteger sn, byte[] data) {
        logger.println("sendMessage", "begin");
        BTSMessage message = new BTSMessage();
        message.setServiceType(serviceType);
        message.setData(data);

        BMCScoreInterface bmc = new BMCScoreInterface(this.bmc);
        bmc.sendMessage(to, SERVICE, sn, message.toBytes());
        logger.println("sendMessage", "end");
    }

    private void responseSuccess(String to, BigInteger sn) {
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_OK);
        response.setMessage(TransferResponse.OK_MSG);
        sendMessage(to, BTSMessage.REPONSE_HANDLE_SERVICE, sn, response.toBytes());
    }

    private void responseError(String to, BigInteger sn, String message) {
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_ERR);
        response.setMessage(message);
        sendMessage(to, BTSMessage.REPONSE_HANDLE_SERVICE, sn, response.toBytes());
    }

    @External
    public void handleBTPMessage(String _from, String _svc, BigInteger _sn, byte[] _msg) {
        require(Context.getCaller().equals(bmc), "Only BMC");

        BTSMessage message = BTSMessage.fromBytes(_msg);
        int serviceType = message.getServiceType();
        if (serviceType == BTSMessage.REQUEST_COIN_TRANSFER) {
            TransferRequest request = TransferRequest.fromBytes(message.getData());
            handleRequest(request, _from, _sn);
        } else if (serviceType == BTSMessage.REPONSE_HANDLE_SERVICE) {
            TransferResponse response = TransferResponse.fromBytes(message.getData());
            handleResponse(_sn, response);
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

    @External
    public void handleBTPError(String _src, String _svc, BigInteger _sn, long _code, String _msg) {
        require(Context.getCaller().equals(bmc), "Only BMC");
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_ERR);
        response.setMessage("BTPError [code:" + _code + ",msg:" + _msg);
        handleResponse(_sn, response);
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

    private Balance getBalance(String coinName, Address owner) {
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

    private void refund(String coinName, Address owner, BigInteger value) {
        logger.println("refund", "coinName:", coinName, "owner:", owner, "value:", value);
        // unlock and add refundable
        Balance balance = getBalance(coinName, owner);
        balance.setLocked(balance.getLocked().subtract(value));
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
        } catch (IllegalArgumentException | NullPointerException e) {
            throw BTSException.unknown(e.getMessage());
        }

        BigInteger nativeCoinTransferAmount = null;
        Asset[] assets = request.getAssets();
        List<String> coinNames = new ArrayList<>();
        List<BigInteger> amounts = new ArrayList<>();
        List<String> registeredCoinNames = getCoinNamesAsList();
        for (Asset asset : assets) {
            String coinName = asset.getCoinName();
            BigInteger amount = asset.getAmount();
            if (amount == null || amount.compareTo(BigInteger.ZERO) < 1) {
                throw BTSException.unknown("Amount must be positive value");
            }

            if (registeredCoinNames.contains(coinName)) {
                coinNames.add(coinName);
                amounts.add(amount);
            } else if (name.equals(coinName)) {
                nativeCoinTransferAmount = amount;
            } else {
                throw BTSException.unknown("Invalid Token");
            }
            checkTransferRestrictions(coinName, from, request.getTo(), amount);
        }

        if (nativeCoinTransferAmount != null) {
            try {
                Context.transfer(to, nativeCoinTransferAmount);
            } catch (Exception e) {
                throw BTSException.unknown("fail to transfer err:" + e.getMessage());
            }
        }

        if (coinNames.size() > 0) {
            mintBatch(to, ArrayUtil.toStringArray(coinNames), ArrayUtil.toBigIntegerArray(amounts));
        }

        logger.println("handleRequest", "responseSuccess");
        responseSuccess(from, sn);
        logger.println("handleRequest", "end");
        TransferReceived(from, to, sn, encode(assets));
    }

    private void handleResponse(BigInteger sn, TransferResponse response) {
        logger.println("handleResponse", "begin", "sn:", sn);
        TransferTransaction transaction = transactions.get(sn);
        List<String> registeredCoinNames = getCoinNamesAsList();
        // ignore when not exists pending request
        if (transaction != null) {
            BigInteger code = response.getCode();
            Address owner = Address.fromString(transaction.getFrom());
            AssetTransferDetail[] assets = transaction.getAssets();

            logger.println("handleResponse", "code:", code);
            if (TransferResponse.RC_OK.equals(code)) {
                List<String> coinNames = new ArrayList<>();
                List<BigInteger> amounts = new ArrayList<>();
                for (AssetTransferDetail asset : assets) {
                    String coinName = asset.getCoinName();
                    BigInteger amount = asset.getAmount();
                    BigInteger fee = asset.getFee();
                    BigInteger locked = amount.add(fee);
                    boolean isNativeCoin = name.equals(coinName);
                    if (isNativeCoin || registeredCoinNames.contains(coinName)) {
                        unlock(coinName, owner, locked);
                        addFee(coinName, fee);
                        if (!isNativeCoin) {
                            coinNames.add(coinName);
                            amounts.add(amount);
                        }
                    } else {
                        // This should not happen
                        throw BTSException.unknown("invalid transaction, invalid coinName");
                    }
                }

                if (coinNames.size() > 0) {
                    burnBatch(ArrayUtil.toStringArray(coinNames), ArrayUtil.toBigIntegerArray(amounts));
                }
            } else {
                for (AssetTransferDetail asset : assets) {
                    String coinName = asset.getCoinName();
                    BigInteger amount = asset.getAmount();
                    BigInteger fee = asset.getFee();
                    BigInteger locked = amount.add(fee);
                    boolean isNativeCoin = name.equals(coinName);
                    if (isNativeCoin || registeredCoinNames.contains(coinName)) {
                        refund(coinName, owner, locked);
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

    @External(readonly = true)
    public BigInteger feeRatio() {
        BTSProperties properties = getProperties();
        return properties.getFeeRatio();
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
            logger.println("mint", "code:", e.getCode(), "msg:", e.getMessage());
            throw BTSException.irc31Reverted("code:" + e.getCode() + "msg:" + e.getMessage());
        } catch (IllegalArgumentException | RevertedException e) {
            logger.println("mint", "Exception:", e.toString());
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
    public void addRestrictor(Address _address) {
        requireOwnerAccess();
        bsrDb.set(_address);
        restriction.set(true);
        restrictionsInterface = new RestrictionsScoreInterface(_address);
    }

    @External
    public void disableRestrictions(Address _address) {
        requireOwnerAccess();
        bsrDb.set(_address);
        restriction.set(false);
    }

    private void checkTransferRestrictions(String _tokenName, String _from, String _to, BigInteger _value) {
        if (restriction.get() != null && bsrDb.get() != null && restriction.get()) {
            // restictonsInterface.validateRestriction(_tokenName, _from, _to, _value);
            Context.call(bsrDb.get(), "validateRestriction", _tokenName, _from, _to, _value);
        }
    }

}