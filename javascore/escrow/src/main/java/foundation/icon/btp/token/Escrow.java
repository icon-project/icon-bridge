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

    private final String to;
    private final Address bmc;
    private final BigInteger asaId;

    private byte[] lastReceivedMessage = "BTP Message".getBytes();
    private String lastReceivedErrorMessage = "BTP Error Message";

    private final VarDB<BigInteger> sn = Context.newVarDB("serviceNumber", BigInteger.class);

    public Escrow (String _to, Address _bmc, BigInteger _asaId) {
        bmc = _bmc;
        to = _to;
        asaId = _asaId;

        if (sn.get() == null) {
            sn.set(BigInteger.ZERO);
        }
    }

    @External(readonly=true)
    public byte[] getLastReceivedMessage() {
        return lastReceivedMessage;
    }

    @External(readonly=true)
    public String getLastReceivedErrorMessage() {
        return lastReceivedErrorMessage;
    }

    @External()
    public void handleBTPMessage(String _from, String _svc, BigInteger _sn, byte[] _msg) {
        this.lastReceivedMessage = _msg;
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

    private BigInteger increaseSn() {
        BigInteger newSn = sn.get().add(BigInteger.ONE);
        sn.set(newSn);
        return newSn;
    }
}