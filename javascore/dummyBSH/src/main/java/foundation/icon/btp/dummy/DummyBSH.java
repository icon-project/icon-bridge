package foundation.icon.btp.dummy;

import java.math.BigInteger;

import score.annotation.External;

import foundation.icon.btp.lib.BSH;

public class DummyBSH implements BSH {
    private byte[] lastReceivedMessage = "BTP Message".getBytes();
    private String lastReceivedErrorMessage = "BTP Error Message";

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
    }

    @External()
    public void handleBTPError(String _src, String _svc, BigInteger _sn, long _code, String _msg) {
        this.lastReceivedErrorMessage = _msg;
    }

    @External()
    public void handleFeeGathering(String _fa, String _svc) {
    }
}