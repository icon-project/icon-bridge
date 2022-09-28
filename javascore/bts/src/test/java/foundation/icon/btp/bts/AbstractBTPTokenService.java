package foundation.icon.btp.bts;


import com.iconloop.score.test.Account;
import com.iconloop.score.test.Score;
import com.iconloop.score.test.ServiceManager;
import com.iconloop.score.test.TestBase;
import foundation.icon.btp.lib.BTPAddress;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.function.Executable;
import java.math.BigInteger;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.mockito.MockedStatic;
import org.mockito.MockedStatic.Verification;
import org.mockito.Mockito;
import score.Address;
import score.Context;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.spy;

public class AbstractBTPTokenService extends TestBase {
    public static final ServiceManager sm = getServiceManager();
    public Account owner;
    public Account nonOwner;
    public Score score;
    public Account irc2 = Account.newScoreAccount(100);
    public Account wrappedIRC2 = Account.newScoreAccount(101);
    public BTPTokenService scoreSpy;
    public Account bmcMock = Account.newScoreAccount(1);
    public String TEST_TOKEN = "Test Token";
    public String ICON = "ICON";
    public String PARA = "PARA";
    public String ETH_ADDR = "0x71C7656EC7ab88b098defB751B7401B5f6d8976F";
    public String SERVICE = "bts";
    public String METACHAIN = "metachain";
    public static final BigInteger UINT_CAP = pow(BigInteger.TWO,256).subtract(BigInteger.ONE);

    BTPAddress btpAddress1 = new BTPAddress(METACHAIN,"0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
    BTPAddress btpAddress2 = new BTPAddress("icon","0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
    BTPAddress btpAddress3 = new BTPAddress("harmony","0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
    String[] links = new String[] {btpAddress1.toString(), btpAddress2.toString(), btpAddress3.toString()};

    static MockedStatic<Context> contextMock;

    @BeforeAll
    protected static void init() {
        contextMock = Mockito.mockStatic(Context.class, Mockito.CALLS_REAL_METHODS);
    }

    @BeforeEach
    void setup() throws Exception {

        owner = sm.createAccount(100);
        nonOwner = sm.createAccount(100);
        BTPAddress btpAddress = new BTPAddress("icon", "0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");

        byte[] bytes = "OptimizedJar".getBytes();

        contextMock.when(() -> Context.call(eq(String.class),eq(bmcMock.getAddress()),
                eq("getBtpAddress"), any())).thenReturn(btpAddress.toString());

        score = sm.deploy(owner, BTPTokenService.class,
                bmcMock.getAddress(), ICON, 18, BigInteger.ZERO, BigInteger.ZERO, bytes);

        BTPTokenService instance = (BTPTokenService) score.getInstance();
        scoreSpy = spy(instance);
        score.setInstance(scoreSpy);
    }

    public void blacklistMocks() {
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), any(), any(), any(), any());

        contextMock.when(sendMessage).thenReturn(null);
        getLinksMock();

    }

    public static BigInteger pow(BigInteger num, int exponent) {
        BigInteger result = BigInteger.ONE;
        for (int i = 0; i < exponent; i++) {
            result = result.multiply(num);
        }
        return result;
    }

    public void getLinksMock() {
        Verification returnLinks = () -> Context.call(eq(String[].class),
                eq(bmcMock.getAddress()), eq("getLinks"), any());
        contextMock.when(returnLinks).thenReturn(links);
    }

    protected void icxLimitBTPMessage() {
        getLinksMock();
        for (String link: links) {
            contextMock.when(() -> Context.call(eq(bmcMock.getAddress()),
                    eq("sendMessage"), eq(link), any(),
                    eq(BigInteger.ONE), any())).thenReturn(null);
        }
    }

    protected void printSn() {
        System.out.println("SN is: "+score.call("getSn"));
    }

    protected void blackListUser(String net, String[] addr, BigInteger sn) {
        blacklistMocks();
        score.invoke(owner, "addBlacklistAddress", net, addr);
        byte[] _msg = blacklistSuccessfulResponse();
        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, sn, _msg);
    }

    protected void removeBlackListedUser(String net, String[] addr, BigInteger sn) {
        blacklistMocks();
        score.invoke(owner, "removeBlacklistAddress", net, addr);
        byte[] _msg = blacklistRemoveSuccessfulResponse();
        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, sn, _msg);
    }

    protected byte[] blacklistSuccessfulResponse() {
        BlacklistResponse response = new BlacklistResponse();
        response.setCode(BlacklistResponse.RC_OK);
        BTSMessage message = new BTSMessage();
        message.setData(response.toBytes());
        message.setServiceType(BTSMessage.BLACKLIST_MESSAGE);
        return message.toBytes();
    }

    protected byte[] blacklistUnsuccessfulResponse() {
        BlacklistResponse response = new BlacklistResponse();
        response.setCode(BlacklistResponse.RC_ERR);
        BTSMessage message = new BTSMessage();
        message.setData(response.toBytes());
        message.setServiceType(BTSMessage.BLACKLIST_MESSAGE);
        return message.toBytes();
    }

    protected byte[] blacklistRemoveSuccessfulResponse() {
        BlacklistResponse response = new BlacklistResponse();
        response.setCode(BlacklistResponse.RC_OK);
        BTSMessage message = new BTSMessage();
        message.setData(response.toBytes());
        message.setServiceType(BTSMessage.BLACKLIST_MESSAGE);
        return message.toBytes();
    }

    protected byte[] blacklistRemoveUnuccessfulResponse() {
        BlacklistResponse response = new BlacklistResponse();
        response.setCode(BlacklistResponse.RC_ERR);
        BTSMessage message = new BTSMessage();
        message.setData(response.toBytes());
        message.setServiceType(BTSMessage.BLACKLIST_MESSAGE);
        return message.toBytes();
    }

    protected void tokenLimitBTPMessage() {
        getLinksMock();
        for (String link: links) {
            BTPAddress addr = BTPAddress.valueOf(link);
            contextMock.when(() -> Context.call(eq(bmcMock.getAddress()),
                    eq("sendMessage"), eq(addr.net()), any(),
                    any(), any())).thenReturn(null);
        }
    }

    public String generateBTPAddress(String net, String addr) {
        BTPAddress btpAddress = new BTPAddress(net, addr);
        return btpAddress.toString();
    }

    public Verification sendICX() {
        Verification sendICX = () -> Context.getValue();
        return sendICX;
    }

    public Verification sendIcxToUser() {
        Verification sendICXToUser = () -> Context.transfer(any(), any());
        return sendICXToUser;
    }

    public void sendBTPMessageMock() {
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), any(), any(), any(), any());
        contextMock.when(sendMessage).thenReturn(null);
    }

    public void register(String name, Address addr) {
        score.invoke(owner, "register",
                name, name, 18, BigInteger.valueOf(10), BigInteger.ONE, addr);
    }

    public void register() {
        score.invoke(owner, "register",
                TEST_TOKEN, "TTK", 18, BigInteger.ZERO, BigInteger.ONE, irc2.getAddress());
    }

    public void registerWrapped() {

        Verification deployWrappedToken = () -> Context.deploy(any(), eq(PARA),
                eq(PARA),eq(18));
        contextMock.when(deployWrappedToken).thenReturn(wrappedIRC2.getAddress());

        score.invoke(owner, "register",PARA, PARA, 18, BigInteger.ZERO, BigInteger.TWO,
                Address.fromString("cx0000000000000000000000000000000000000000"));
    }

    public void deposit(BigInteger value) {
        contextMock.when(() -> Context.call(eq(BigInteger.class),eq(irc2.getAddress()), eq("balanceOf"), eq(owner.getAddress()))).thenReturn(BigInteger.valueOf(100));
        score.invoke(irc2,"tokenFallback", owner.getAddress(), value, "0".getBytes());
    }

    public void expectErrorMessage(Executable contractCall, String errorMessage) {
        AssertionError e = Assertions.assertThrows(AssertionError.class, contractCall);
        assertEquals(errorMessage, e.getMessage());
    }

    public void expectErrorMessageIn(Executable contractCall, String errorMessage) {
        AssertionError e = Assertions.assertThrows(AssertionError.class, contractCall);
        assert e.getMessage().contains(errorMessage);
    }
}