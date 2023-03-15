package foundation.icon.btp.token;

import java.math.BigInteger;
import java.util.Arrays;

import foundation.icon.btp.lib.BSH;
import foundation.icon.score.util.Logger;
import foundation.icon.btp.lib.BMCScoreInterface;

import score.annotation.External;
import score.Address;
import score.Context;
import score.VarDB;

public class Escrow implements BSH {
    private static final Logger logger = Logger.getLogger(Escrow.class);
    public static final String SERVICE = "i2a";
    public static final int MESSAGE_LENGTH = 29;

    private final String to;
    private final Address bmc;
    private final BigInteger asaId;
    private final Address tokenAddress;

    private String lastReceivedErrorMessage = "BTP Error Message";

    private final VarDB<BigInteger> sn = Context.newVarDB("serviceNumber", BigInteger.class);

    public Escrow (String _to, Address _bmc, BigInteger _asaId, Address _tokenAddress) {
        bmc = _bmc;
        to = _to;
        asaId = _asaId;
        tokenAddress = _tokenAddress;

        if (sn.get() == null) {
            sn.set(BigInteger.ZERO);
        }
    }

    @External(readonly=true)
    public String getLastReceivedErrorMessage() {
        return lastReceivedErrorMessage;
    }

    @External()
    public void handleBTPMessage(String _from, String _svc, BigInteger _sn, byte[] _msg) {
        Context.require(Context.getCaller().equals(bmc), "Only BMC");
        Context.require(_msg.length == MESSAGE_LENGTH, "Invalid message length");
        
        byte[] amountBytes = Arrays.copyOfRange(_msg, 0, 8); 
        BigInteger amount = new BigInteger(1, amountBytes);
        
        byte[] isContractBytes = Arrays.copyOfRange(_msg, 8, 9); 
        boolean isContract = isContractBytes[0] != 0;

        byte[] dstBytes = Arrays.copyOfRange(_msg, 9, _msg.length); 
        String dstString = byteArrayToHex(dstBytes);
        String formattedDstString = isContract ? "cx" + dstString : "hx" + dstString;
        Address dst = Address.fromString(formattedDstString);

        Context.require(amount.compareTo(BigInteger.ZERO) >= 0);

        Context.call(tokenAddress, "transfer", dst, amount);

        increaseSn();
    }

    @External()
    public void handleBTPError(String _src, String _svc, BigInteger _sn, long _code, String _msg) {
        this.lastReceivedErrorMessage = _msg;
    }

    @External()
    public void handleFeeGathering(String _fa, String _svc) {
    }

    @External
    public void tokenFallback(Address _from, BigInteger _value, byte[] _algoPubKey) {
        Context.require( _value.compareTo(BigInteger.ZERO) >= 0, "value should be positive");
        Context.require(_value.compareTo(Conversion.maxUint64) < 1, "Amount too big");
        
        byte[] assets = Conversion.bigIntToByteArray(asaId);
        byte[] amountBytes = Conversion.bigIntToByteArray(_value);
        byte[] message = new byte[2 + assets.length + _algoPubKey.length + amountBytes.length + _algoPubKey.length];

        int offset = 0;
        
        // set assets count
        message[0] = 1;
        offset += 1;
        System.arraycopy(assets, 0, message, offset, assets.length);
        offset += assets.length;
        // set accounts count
        message[offset] = 1;
        offset += 1;
        System.arraycopy(_algoPubKey, 0, message, offset, _algoPubKey.length);
        offset += _algoPubKey.length;
        System.arraycopy(amountBytes, 0, message, offset, amountBytes.length);
        offset += amountBytes.length;
        System.arraycopy(_algoPubKey, 0, message, offset, _algoPubKey.length);

        BMCScoreInterface bmc = new BMCScoreInterface(this.bmc);
        BigInteger sn = increaseSn();
        
        bmc.sendMessage(to, SERVICE, sn, message);
    }

    private String byteToHex(byte num) {
        char[] hexDigits = new char[2];
        hexDigits[0] = Character.forDigit((num >> 4) & 0xF, 16);
        hexDigits[1] = Character.forDigit((num & 0xF), 16);
        return new String(hexDigits);
    }

    private String byteArrayToHex(byte[] byteArray) {
        StringBuffer hexStringBuffer = new StringBuffer();
        for (int i = 0; i < byteArray.length; i++) {
            hexStringBuffer.append(byteToHex(byteArray[i]));
        }
        return hexStringBuffer.toString();
    }

    private BigInteger increaseSn() {
        BigInteger newSn = sn.get().add(BigInteger.ONE);
        sn.set(newSn);
        return newSn;
    }
}