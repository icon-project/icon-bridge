package foundation.icon.btp.bmc;

import com.iconloop.score.test.Account;
import com.iconloop.score.test.Score;
import com.iconloop.score.test.ServiceManager;
import com.iconloop.score.test.TestBase;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.function.Executable;
import org.mockito.MockedStatic;
import org.mockito.MockedStatic.Verification;
import org.mockito.Mockito;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.spy;

import score.Address;
import score.Context;

public class AbstractBTPMessageCenterTest extends TestBase {

    public static final ServiceManager sm = getServiceManager();
    static MockedStatic<Context> contextMock;
    public Account owner;
    public Account nonOwner;
    public Score score;
    public BTPMessageCenter scoreSpy;
    public String NETWORK = "0x1.icon";
    public String BTS = "bts";
    public Account BTSScore = sm.createAccount();

    public String REQUIRE_OWNER_ACCESS = "require owner access";


    @BeforeAll
    protected static void init() {
        contextMock = Mockito.mockStatic(Context.class, Mockito.CALLS_REAL_METHODS);
        long CURRENT_TIMESTAMP = System.currentTimeMillis() / 1_000L;
        sm.getBlock().increase(CURRENT_TIMESTAMP / 2);
    }

    @BeforeEach
    void setup() throws Exception {

        owner = sm.createAccount(100);
        nonOwner = sm.createAccount(100);

        score = sm.deploy(owner, BTPMessageCenter.class, NETWORK);

        BTPMessageCenter instance = (BTPMessageCenter) score.getInstance();
        scoreSpy = spy(instance);
        score.setInstance(scoreSpy);
    }

    public long getCurrentBlockHeight() {
        return sm.getBlock().getHeight();
    }

    public String getBTPAddress(String address) {
        return "btp://" + NETWORK + "/" + address;
    }

    public String getBTPAddress(Address address) {
        return "btp://" + NETWORK + "/" + address.toString();
    }

    public void expectErrorMessage(Executable contractCall, String errorMessage) {
        AssertionError e = Assertions.assertThrows(AssertionError.class, contractCall);
        assertEquals(errorMessage, e.getMessage());
    }

    public Verification mockICXSent() {
        Verification sendICX = () -> Context.getValue();
        return sendICX;
    }
    public Verification mockICXBalance(Address addr) {
        Verification icxBalance = () -> Context.getBalance(addr);
        return icxBalance;
    }

    public Verification mockBlockHeight() {
        Verification blockHeight = () -> Context.getBlockHeight();
        return blockHeight;
    }

    public void addLink(String link) {
        score.invoke(owner, "addLink", link);
    }
}
