package foundation.icon.btp.restrictions;

import foundation.icon.btp.irc2.IRC2Basic;
import com.iconloop.testsvc.Account;
import com.iconloop.testsvc.Score;
import com.iconloop.testsvc.ServiceManager;
import com.iconloop.testsvc.TestBase;
import org.junit.jupiter.api.*;

import java.math.BigInteger;

import static java.math.BigInteger.TEN;
import static org.junit.jupiter.api.Assertions.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class RestrictionsTest extends TestBase {

    final static String tokenName = "ETH";
    final static String symbol = "ETH";
    final static int decimals = 18;
    private static final BigInteger initialSupply = BigInteger.valueOf(2000);
    private static final BigInteger totalSupply = initialSupply.multiply(TEN.pow(decimals));
    private static final ServiceManager sm = getServiceManager();
    private static Account[] owners;
    private static Score token;
    private static Score restrictons;

    @BeforeAll
    public static void setUp() throws Exception {
        owners = new Account[4];
        for (int i = 0; i < owners.length; i++) {
            owners[i] = sm.createAccount(100);
        }

        token = sm.deploy(owners[0], IRC2Basic.class, tokenName, symbol, decimals, initialSupply);
        restrictons = sm.deploy(owners[0], Restrictions.class);
        BigInteger balance = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        assertEquals(totalSupply, balance);
    }

    @Order(1)
    @Test
    void validateBeforeRegister() {
        restrictons.invoke(owners[0], "validateRestriction", tokenName, owners[2].getAddress().toString(), owners[3].getAddress().toString(), BigInteger.valueOf(5));
    }

    @Order(1)
    @Test
    void registerTokenLimit() {
        restrictons.invoke(owners[0], "registerTokenLimit", tokenName, symbol, token.getAddress(), BigInteger.valueOf(10));
    }

    @Order(2)
    @Test
    void addBlackListUser() {
        restrictons.invoke(owners[0], "addBlacklistedUser", owners[1].getAddress().toString());
        restrictons.invoke(owners[0], "addBlacklistedUser", owners[2].getAddress().toString());
        restrictons.invoke(owners[0], "addBlacklistedUser", "0x0000000000000000000000000000000000000001");
        Boolean _isUserBlackListed = (Boolean) restrictons.call("isUserBlackListed", "0x0000000000000000000000000000000000000001");
        assertTrue(_isUserBlackListed);
    }

    @Order(3)
    @Test
    void removeBlacklistedUser() {
        Boolean _isUserBlackListed = (Boolean) restrictons.call("isUserBlackListed", owners[2].getAddress().toString());
        assertTrue(_isUserBlackListed);
        restrictons.invoke(owners[0], "removeBlacklistedUser", owners[2].getAddress().toString());
        _isUserBlackListed = (Boolean) restrictons.call("isUserBlackListed", owners[2].getAddress().toString());
        assertFalse(_isUserBlackListed);
    }

    @Order(4)
    @Test
    void isUserBlackListed() {
        Boolean _isUserBlackListed = (Boolean) restrictons.call("isUserBlackListed", owners[1].getAddress().toString());
        assertTrue(_isUserBlackListed);
    }

    @Order(5)
    @Test
    void validateRestriction() {
        restrictons.invoke(owners[0], "validateRestriction", tokenName, owners[2].getAddress().toString(), owners[3].getAddress().toString(), BigInteger.valueOf(5));
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                restrictons.invoke(owners[0], "validateRestriction", tokenName, "0x0000000000000000000000000000000000000001", owners[1].getAddress().toString(), BigInteger.valueOf(5))
        );
        assertTrue(thrown.getMessage().contains("_from user is Blacklisted"));
        thrown = assertThrows(AssertionError.class, () ->
                restrictons.invoke(owners[0], "validateRestriction", tokenName, "0x0000000000000000000000000000000000000002", owners[1].getAddress().toString(), BigInteger.valueOf(5))
        );
        assertTrue(thrown.getMessage().contains("_to user is Blacklisted"));
        thrown = assertThrows(AssertionError.class, () ->
                restrictons.invoke(owners[0], "validateRestriction", tokenName, "0x0000000000000000000000000000000000000002", owners[2].getAddress().toString(), BigInteger.valueOf(11))
        );
        assertTrue(thrown.getMessage().contains("Transfer amount exceeds the transaction limit"));
    }

    @Order(5)
    @Test
    void validateRestrictionb() {
        restrictons.invoke(owners[0], "addBlacklistedUser", "0xebcbd4a934a68510e21ba25b2a827138248a63e5\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000");
        restrictons.invoke(owners[0], "validateRestriction", tokenName, owners[2].getAddress().toString(), owners[3].getAddress().toString(), BigInteger.valueOf(5));
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                restrictons.invoke(owners[0], "validateRestriction", tokenName, "0x0000000000000000000000000000000000000001", owners[1].getAddress().toString(), BigInteger.valueOf(5))
        );
        assertTrue(thrown.getMessage().contains("_from user is Blacklisted"));
        thrown = assertThrows(AssertionError.class, () ->
                restrictons.invoke(owners[0], "validateRestriction", tokenName, "0x0000000000000000000000000000000000000002", owners[1].getAddress().toString(), BigInteger.valueOf(5))
        );
        assertTrue(thrown.getMessage().contains("_to user is Blacklisted"));
        thrown = assertThrows(AssertionError.class, () ->
                restrictons.invoke(owners[0], "validateRestriction", tokenName, "0x0000000000000000000000000000000000000002", owners[2].getAddress().toString(), BigInteger.valueOf(11))
        );
        assertTrue(thrown.getMessage().contains("Transfer amount exceeds the transaction limit"));
    }
}
