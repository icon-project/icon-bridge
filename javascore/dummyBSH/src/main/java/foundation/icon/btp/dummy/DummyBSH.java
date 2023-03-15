package foundation.icon.btp.dummy;

import java.math.BigInteger;

import foundation.icon.btp.lib.BSH;
import foundation.icon.btp.lib.BMCScoreInterface;
import foundation.icon.score.util.Logger;

import score.annotation.External;
import score.annotation.EventLog;
import score.Address;

public class DummyBSH implements BSH {
    private static final Logger logger = Logger.getLogger(DummyBSH.class);
    public static final String SERVICE = "dbsh";

    private final String to;
    private final Address bmc;

    private byte[] lastReceivedMessage = "BTP Message".getBytes();
    private String lastReceivedErrorMessage = "BTP Error Message";

    public DummyBSH (String _to, Address _bmc) {
        bmc = _bmc;
        to = _to;
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
    }

    @External()
    public void handleBTPError(String _src, String _svc, BigInteger _sn, long _code, String _msg) {
        this.lastReceivedErrorMessage = _msg;
    }

    @External()
    public void handleFeeGathering(String _fa, String _svc) {
    }

    @External()
    public void sendServiceMessage() {
        byte[] referencesCount = new byte[1];
        referencesCount[0] = 0;

        BigInteger sn = BigInteger.valueOf(1);
        byte[] dummyMessage = "Hello Algorand".getBytes();

        byte[] message = new byte[referencesCount.length + referencesCount.length + dummyMessage.length];

        System.arraycopy(referencesCount, 0, message, 0, referencesCount.length);
        System.arraycopy(referencesCount, 0, message, referencesCount.length, referencesCount.length);
        System.arraycopy(dummyMessage, 0, message, referencesCount.length + referencesCount.length, dummyMessage.length);

        BMCScoreInterface bmc = new BMCScoreInterface(this.bmc);
        bmc.sendMessage(to, SERVICE, sn, message);
    }

}