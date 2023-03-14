package foundation.icon.btp.token;

import com.iconloop.score.token.irc2.IRC2Basic;

import java.math.BigInteger;
import java.util.Arrays;

import foundation.icon.btp.lib.BSH;
import foundation.icon.btp.lib.BMCScoreInterface;
import foundation.icon.score.util.Logger;

import score.annotation.External;
import score.annotation.EventLog;
import score.Address;
import score.Context;
import score.VarDB;

public class WrappedToken extends IRC2Basic implements BSH {
    private static final Logger logger = Logger.getLogger(WrappedToken.class);
    public static final String SERVICE = "wtt";
    public static final int MESSAGE_LENGTH = 29;

    private final String to;
    private final Address bmc;
    private final BigInteger asaId;

    private String lastReceivedErrorMessage = "BTP Error Message";

    private final VarDB<BigInteger> sn = Context.newVarDB("serviceNumber", BigInteger.class);

    public WrappedToken (String _to, Address _bmc, BigInteger _asaId, String _name, String _symbol, int _decimals) {
        super(_name, _symbol, _decimals);

        bmc = _bmc;
        to = _to;
        asaId = _asaId;

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
        // TODO check _from if needed

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
        _mint(dst, amount);

        increaseSn();
    }

    @External()
    public void handleBTPError(String _src, String _svc, BigInteger _sn, long _code, String _msg) {
        this.lastReceivedErrorMessage = _msg;
    }

    @External()
    public void handleFeeGathering(String _fa, String _svc) {
    }

    @External()
    public void burn(byte[] algoPubKey, BigInteger _amount) {
        Context.require(_amount.compareTo(Conversion.maxUint64) < 1, "Amount too big");
        Context.require(_amount.compareTo(BigInteger.ZERO) > 0, "Amount should be positive");

        _burn(Context.getCaller(), _amount);

        BigInteger sn = increaseSn();
        BMCScoreInterface bmc = new BMCScoreInterface(this.bmc);

        byte[] assets = Conversion.bigIntToByteArray(asaId);
        byte[] amountBytes = Conversion.bigIntToByteArray(_amount);
        byte[] message = new byte[2 + assets.length + algoPubKey.length + amountBytes.length + algoPubKey.length];

        int offset = 0;

        // set assets count
        message[0] = 1;
        offset += 1;
        System.arraycopy(assets, 0, message, offset, assets.length);
        offset += assets.length;
        // set accounts count
        message[offset] = 1;
        offset += 1;
        System.arraycopy(algoPubKey, 0, message, offset, algoPubKey.length);
        offset += algoPubKey.length;
        System.arraycopy(amountBytes, 0, message, offset, amountBytes.length);
        offset += amountBytes.length;
        System.arraycopy(algoPubKey, 0, message, offset, algoPubKey.length);

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