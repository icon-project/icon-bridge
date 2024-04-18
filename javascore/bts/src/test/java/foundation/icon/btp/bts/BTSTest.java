package foundation.icon.btp.bts;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;

import com.iconloop.score.test.Account;
import foundation.icon.btp.lib.BTPAddress;
import java.math.BigInteger;
import java.util.List;
import java.util.Map;
import java.util.Random;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestMethodOrder;
import org.junit.jupiter.api.function.Executable;
import org.mockito.MockedStatic.Verification;
import score.Address;
import score.Context;

@TestMethodOrder(OrderAnnotation.class)
public class BTSTest extends AbstractBTPTokenService {

    @Order(1)
    @Test
    public void name() {
        assertEquals("BTP Token Service", score.call("name"));
        assertEquals(BigInteger.ZERO, score.call(("getSn")));
    }

    @Order(2)
    @Test
    public void addToLockListTest() {

        assertEquals(BigInteger.ZERO, score.call(("getSn")));

        // non-owner tries
        String[] addr = new String[]{"Hell "};
        Executable call = () -> score.invoke(nonOwner, "addBlacklistAddress", "hello world ", addr);
        expectErrorMessage(call, "require owner access");

        blacklistMocks();
        call = () -> score.invoke(owner, "addBlacklistAddress", "networking", addr);
        expectErrorMessage(call, "Invalid link");

        String[] addr1 = new String[]{"metachain address 2"};
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr1);

        byte[] _msg = blacklistSuccessfulResponse();

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.ONE, _msg);
        verify(scoreSpy).AddedToBlacklist(eq(BigInteger.ONE), any());

        String[] addr2 = new String[]{"   metachain address 2       "};
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr2);

        List<String> actual = (List<String>) score.call("getBlackListedUsers", METACHAIN, 0, 10);
        List<String> expected = List.of("metachain address 2");
        assertEquals(expected, actual);

        String[] addr3 = new String[]{"  metachain address 1  "};
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr3);
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.TWO, _msg);
        verify(scoreSpy).AddedToBlacklist(eq(BigInteger.TWO), any());

        expected = List.of("metachain address 2", "metachain address 1");
        actual = (List<String>) score.call("getBlackListedUsers", METACHAIN, 0, 10);
        assertEquals(expected, actual);

        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 1 "));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " yu belong with me "));

        assertEquals(2, score.call("blackListedUsersCount", METACHAIN));

        assertEquals(BigInteger.valueOf(2), score.call("getSn"));

        score.invoke(owner, "addBlacklistAddress", "icon", new String[]{" invalid icon address "});
        assertEquals(0, score.call("blackListedUsersCount", "icon"));

        score.invoke(owner, "addBlacklistAddress", "icon",
                new String[]{" cx42bd7394a8272fdb8683a41b92921247c34c522a  "});
        assertEquals(1, score.call("blackListedUsersCount", "icon"));
        assertEquals(true, score.call("isUserBlackListed", "icon", " cx42bd7394a8272fdb8683a41b92921247c34c522a "));

        // unsuccessful blacklist response

        String[] addr4 = new String[]{"  metachain address 3  ", "metachain address 4", "metachain address 5"};
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr4);

        // temporarily added to blacklist
        assertEquals(5, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 3"));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 4   "));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, "metachain address 5 "));

        // interscore call, so sn increases by one
        assertEquals(BigInteger.valueOf(3), score.call("getSn"));

        _msg = blacklistUnsuccessfulResponse();
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(3), _msg);

        // remove users from blacklist
        assertEquals(2, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 3"));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 4   "));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, "metachain address 5 "));

        expected = List.of("metachain address 2", "metachain address 1");
        actual = (List<String>) score.call("getBlackListedUsers", METACHAIN, 0, 10);
        assertEquals(expected, actual);

        // check handleBTPError for blacklisted users
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr4);

        // temporarily added to blacklist
        assertEquals(5, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 3"));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 4   "));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, "metachain address 5 "));

        // interscore call, so sn increases by one
        assertEquals(BigInteger.valueOf(4), score.call("getSn"));

        score.invoke(bmcMock, "handleBTPError", METACHAIN, SERVICE, BigInteger.valueOf(4), Long.valueOf(1), "message");
        assertEquals(2, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 3"));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 4   "));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, "metachain address 5 "));

        expected = List.of("metachain address 2", "metachain address 1");
        actual = (List<String>) score.call("getBlackListedUsers", METACHAIN, 0, 10);
        assertEquals(expected, actual);
    }

    @Order(3)
    @Test
    public void removeFromLocklist() {
        assertEquals(BigInteger.ZERO, score.call(("getSn")));

        blacklistMocks();

        String[] addr1 = new String[]{"metachain address 2"};
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr1);
        String[] addr2 = new String[]{"  metachain address 1  "};
        score.invoke(owner, "addBlacklistAddress", METACHAIN, addr2);
        byte[] _msg = blacklistSuccessfulResponse();

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.ONE, _msg);
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.TWO, _msg);

        Executable call = () -> score.invoke(owner, "removeBlacklistAddress", "icx", new String[]{" hell "});
        expectErrorMessage(call, "Invalid link");

        // non-owner tries
        call = () -> score.invoke(nonOwner, "removeBlacklistAddress", "harmony", new String[]{"hell"});
        expectErrorMessage(call, "require owner access");

        // try to remove non blacklisted
        score.invoke(owner, "removeBlacklistAddress", "harmony", new String[]{"hell"});

        // remove legit user
        score.invoke(owner, "removeBlacklistAddress", METACHAIN, addr1);

        _msg = blacklistRemoveSuccessfulResponse();

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(3), _msg);
        verify(scoreSpy).RemovedFromBlacklist(eq(BigInteger.valueOf(3)), any());
        List<String> expected = List.of("metachain address 1");
        List<String> actual = (List<String>) score.call("getBlackListedUsers", METACHAIN, 0, 1);
        assertEquals(expected, actual);

        assertEquals(1, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 2 "));

        assertEquals(BigInteger.valueOf(3), score.call(("getSn")));

        // try to remove from locklist failed response
        _msg = blacklistRemoveUnuccessfulResponse();

        score.invoke(owner, "removeBlacklistAddress", METACHAIN, new String[]{"metachain address 1"});
        // remove from blacklist temporarily
        assertEquals(0, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 1 "));

        assertEquals(BigInteger.valueOf(4), score.call(("getSn")));

        // get unsuccessful response
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(4), _msg);
        assertEquals(1, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 1 "));

        // handle BTP Error for remove from blacklist
        score.invoke(owner, "removeBlacklistAddress", METACHAIN, new String[]{"metachain address 1"});
        assertEquals(0, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(false, score.call("isUserBlackListed", METACHAIN, " metachain address 1 "));
        assertEquals(BigInteger.valueOf(5), score.call(("getSn")));
        score.invoke(bmcMock, "handleBTPError", METACHAIN, SERVICE, BigInteger.valueOf(5), Long.valueOf(1), "message");
        assertEquals(1, score.call("blackListedUsersCount", METACHAIN));
        assertEquals(true, score.call("isUserBlackListed", METACHAIN, " metachain address 1 "));
    }


    @Test
    @Order(4)
    @DisplayName("register wrapped coin")
    public void registerWrappedToken() {
        Address ZERO_SCORE = Address.fromString("cx0000000000000000000000000000000000000000");

        Verification deployWrappedToken = () -> Context.deploy(any(), eq(PARA), eq(PARA), eq(18));
        contextMock.when(deployWrappedToken).thenReturn(wrappedIRC2.getAddress());

        Executable call = () -> score.invoke(nonOwner, "register", PARA, PARA, 18, BigInteger.ZERO, BigInteger.ZERO,
                ZERO_SCORE);
        expectErrorMessage(call, "require owner access");

        score.invoke(owner, "register", PARA, PARA, 18, BigInteger.ZERO, BigInteger.ZERO,
                ZERO_SCORE);

        assertEquals(1, score.call("getRegisteredTokensCount"));

        List<String> registered = List.of(ICON, PARA);
        assertEquals(registered, score.call("coinNames"));
        assertEquals(wrappedIRC2.getAddress(), score.call("coinId", PARA));
        assertEquals(UINT_CAP, score.call("getTokenLimit", PARA));

        call = () -> score.invoke(owner, "register", "COIN", "COIN", 18, BigInteger.ONE.negate(), BigInteger.ZERO,
                ZERO_SCORE);
        expectErrorMessageIn(call, "The feeNumerator should be less than FEE_DENOMINATOR "
                + "and feeNumerator should be greater than 1");

        call = () -> score.invoke(owner, "register", "COIN", "COIN", 18, BigInteger.valueOf(100000L), BigInteger.ZERO,
                ZERO_SCORE);
        expectErrorMessageIn(call, "The feeNumerator should be less than FEE_DENOMINATOR "
                + "and feeNumerator should be greater than 1");

        call = () -> score.invoke(owner, "register", "COIN", "COIN", 18, BigInteger.ZERO, BigInteger.TEN.negate(),
                ZERO_SCORE);
        expectErrorMessageIn(call, "Fixed fee cannot be less than zero");
        
        call = () -> score.invoke(owner, "register", "META", "META", 18, BigInteger.ZERO, BigInteger.ZERO,
                wrappedIRC2.getAddress());
        expectErrorMessage(call, "coin with that address already registered");
    }

    @Test
    @Order(5)
    public void registerIRC2Token() {
        // non owner tries to call (Scenario 2)
        Executable call = () -> score.invoke(nonOwner, "register", TEST_TOKEN, "TTK", 18, ICX.divide(BigInteger.TEN),
                ICX, irc2.getAddress());
        expectErrorMessage(call, "require owner access");

        // owner registers new coin (Scenario 1)
        score.invoke(owner, "register", TEST_TOKEN, "TTK", 18, BigInteger.TEN, ICX, irc2.getAddress());

        // owner registers coin that exists (Scenario 3)
        call = () -> score.invoke(owner, "register", TEST_TOKEN, "TTK", 18, BigInteger.TEN, ICX,
                irc2.getAddress());
        expectErrorMessage(call, "already existed");

        List<String> expected = List.of(ICON, TEST_TOKEN);
        assertEquals(expected, score.call("coinNames"));

        assertEquals(UINT_CAP, score.call("getTokenLimit", TEST_TOKEN));

        assertEquals(1, score.call("getRegisteredTokensCount"));
    }

    @Test
    @Order(6)
    public void operationsOnUnRegisteredToken() {
        Account IRC2Token = Account.newScoreAccount(10);
        Executable call = () -> score.invoke(IRC2Token, "tokenFallback", owner.getAddress(), BigInteger.valueOf(10),
                "0".getBytes());
        expectErrorMessage(call, "Token not registered");
    }

    @Test
    @Order(7)
    @Disabled
    public void operationOnRegisteredTokens() {
        register();

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(irc2.getAddress()), eq("balanceOf"),
                eq(owner.getAddress()))).thenReturn(BigInteger.valueOf(100));
        score.invoke(irc2, "tokenFallback", owner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        Map<String, BigInteger> balance = (Map<String, BigInteger>) score.call("balanceOf", owner.getAddress(),
                TEST_TOKEN);
        assertEquals(balance.get("usable"), BigInteger.valueOf(10));

        Executable call = () -> score.invoke(irc2, "tokenFallback", owner.getAddress(), BigInteger.ONE.negate(),
                "0".getBytes());
        expectErrorMessage(call, "value should be positive");
    }

    @Test
    @Order(8)
    @Disabled
    public void reclaimDepositedTokens() {
        // for IRC2 tokens
        register();

        // ABC is a unregistered token
        Executable call = () -> score.invoke(owner, "reclaim", "ABC", BigInteger.valueOf(0));
        expectErrorMessage(call, "_value must be positive");
        call = () -> score.invoke(owner, "reclaim", "ABC", BigInteger.valueOf(-10));
        expectErrorMessage(call, "_value must be positive");
        call = () -> score.invoke(owner, "reclaim", "ABC", BigInteger.valueOf(10));
        expectErrorMessage(call, "invalid value");

        deposit(BigInteger.valueOf(10));
        // locked = 10
        contextMock.when(() -> Context.call(eq(irc2.getAddress()), eq("transfer"), eq(owner.getAddress()),
                eq(BigInteger.valueOf(5)), any())).thenReturn(null);
        score.invoke(owner, "reclaim", TEST_TOKEN, BigInteger.valueOf(5));

        Map<String, BigInteger> balance = (Map<String, BigInteger>) score.call("balanceOf", owner.getAddress(),
                TEST_TOKEN);
        assertEquals(BigInteger.valueOf(0), balance.get("refundable"));
        assertEquals(BigInteger.valueOf(0), balance.get("locked"));
        assertEquals(BigInteger.valueOf(5), balance.get("usable"));

        // still 5 left to reclaim
        contextMock.when(() -> Context.call(eq(irc2.getAddress()), eq("transfer"), eq(owner.getAddress()),
                eq(BigInteger.valueOf(2)), any())).thenReturn(null);
        score.invoke(owner, "reclaim", TEST_TOKEN, BigInteger.valueOf(2));

        balance = (Map<String, BigInteger>) score.call("balanceOf", owner.getAddress(), TEST_TOKEN);
        assertEquals(BigInteger.valueOf(0), balance.get("refundable"));
        assertEquals(BigInteger.valueOf(0), balance.get("locked"));
        assertEquals(BigInteger.valueOf(3), balance.get("usable"));

        // wrapped IRC2 tradable is not directly transferred

    }

    @Test
    @Order(9)
    @Disabled
    public void transferNativeCoin() {
        sendBTPMessageMock();
        String btpAddress = generateBTPAddress(METACHAIN, owner.getAddress().toString());
        Account user = sm.createAccount(1000);
        // general condition
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));
        score.invoke(user, "transferNativeCoin", btpAddress);

        // blacklist reciever
        blacklistMocks();
        score.invoke(owner, "addRestriction");
        String addr = generateBTPAddress(METACHAIN, user.getAddress().toString());
        String[] finalAdddr = new String[]{user.getAddress().toString()};
        blackListUser(METACHAIN, finalAdddr, BigInteger.TWO);

//        score.invoke(owner, "addBlacklistAddress", METACHAIN, user.getAddress().toString());

        Executable call = () -> score.invoke(owner, "transferNativeCoin", addr);
        expectErrorMessage(call, "_to user is Blacklisted");

        // blacklist sender
        String newAddr = generateBTPAddress(METACHAIN, owner.getAddress().toString());
        call = () -> score.invoke(user, "transferNativeCoin", newAddr);
        expectErrorMessage(call, "_from user is Blacklisted");

        // remove blacklisted user
        removeBlackListedUser(METACHAIN, finalAdddr, BigInteger.valueOf(3));

        score.invoke(owner, "setTokenLimit", new String[]{ICON}, new BigInteger[]{BigInteger.valueOf(90)});

        call = () -> score.invoke(owner, "transferNativeCoin", addr);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        score.invoke(owner, "setTokenLimit", new String[]{ICON}, new BigInteger[]{BigInteger.valueOf(100)});

        score.invoke(owner, "transferNativeCoin", addr);

        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(99));
        score.invoke(owner, "transferNativeCoin", addr);
    }

    @Test
    @Order(9)
    @Disabled
    public void transferNativeCoinResponse() {
        sendBTPMessageMock();
        String btpAddress = generateBTPAddress(METACHAIN, owner.getAddress().toString());
        Account user = sm.createAccount();

        score.invoke(owner, "setFeeRatio", ICON, BigInteger.ZERO, BigInteger.TEN);

        // handleBTPError
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));
        score.invoke(user, "transferNativeCoin", btpAddress);

        BigInteger sn = (BigInteger) score.call("getSn");

        assertNotEquals(score.call("getTransaction", sn), null);
        // do not mock Context.transfer
        // raises exception
        // balance set to refundable
        score.invoke(bmcMock, "handleBTPError", METACHAIN, SERVICE, sn, Long.valueOf(1), "");
        Map<String, BigInteger> balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), ICON);
        // 10 as fee, 90 is refundable
        assertEquals(BigInteger.valueOf(90), balance.get("refundable"));

        verify(scoreSpy).TransferEnd(eq(user.getAddress()), eq(sn), eq(BigInteger.ONE), any());

        assertEquals(score.call("getTransaction", sn), null);

        // handleBTPError
        score.invoke(user, "transferNativeCoin", btpAddress);

        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), ICON);
        assertEquals(BigInteger.valueOf(100), balance.get("locked"));

        sn = (BigInteger) score.call("getSn");
        // mock context.transfer any amount
        contextMock.when(() -> Context.transfer(user.getAddress(), BigInteger.valueOf(90)))
                .then(invocationOnMock -> null);
        score.invoke(bmcMock, "handleBTPError", METACHAIN, SERVICE, sn, Long.valueOf(1), "");

        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), ICON);
        assertEquals(BigInteger.valueOf(90), balance.get("refundable"));
        assertEquals(BigInteger.ZERO, balance.get("locked"));

        // try to reclaim 100 first
        Executable call = () -> score.invoke(user, "reclaim", ICON, BigInteger.valueOf(100));
        expectErrorMessageIn(call, "invalid value");

        // reclaim the refundable 90
        // 90 transfer has been mocked already
        score.invoke(user, "reclaim", ICON, BigInteger.valueOf(90));
        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), ICON);
        assertEquals(BigInteger.ZERO, balance.get("refundable"));
        assertEquals(BigInteger.ZERO, balance.get("locked"));

        // handle BTPMessage
        score.invoke(user, "transferNativeCoin", btpAddress);
        sn = (BigInteger) score.call("getSn");

        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
        message.setData("message".getBytes());

        // successful response
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, sn, message.toBytes());
        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), ICON);
        assertEquals(BigInteger.ZERO, balance.get("refundable"));
        assertEquals(BigInteger.ZERO, balance.get("locked"));
    }

    @Test
    @Order(10)
    @Disabled
    public void transfer() {
        String btpaddr = generateBTPAddress("harmony", ETH_ADDR);

        Executable call = () -> score.invoke(nonOwner, "transfer", "Token1", BigInteger.valueOf(-1), btpaddr);
        expectErrorMessage(call, "Invalid amount");

        call = () -> score.invoke(nonOwner, "transfer", "Token1", BigInteger.ZERO, btpaddr);
        expectErrorMessage(call, "Invalid amount");

        call = () -> score.invoke(nonOwner, "transfer", "Token1", BigInteger.valueOf(10), btpaddr);
        expectErrorMessage(call, "Not supported Token");

        register();

        // mock message to bmc
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("harmony"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(irc2, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(nonOwner, "transfer", TEST_TOKEN, BigInteger.valueOf(10), btpaddr);

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(irc2.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));
        Map<String, BigInteger> balance = (Map<String, BigInteger>) score.call("balanceOf", nonOwner.getAddress(),
                TEST_TOKEN);
        // currently, still in locked state
        // will be updated once response comes from relayer via bmc
        assertEquals(BigInteger.valueOf(10), balance.get("locked"));

        // should not transfer native-coin
        call = () -> score.invoke(nonOwner, "transfer", ICON, BigInteger.TEN, btpaddr);
        expectErrorMessage(call, "Only for IRC2 Token");
    }

    @Test
    @Order(10)
    @Disabled
    public void transferResponseBack() {
        sendBTPMessageMock();
        String btpAddress = generateBTPAddress(METACHAIN, owner.getAddress().toString());
        Account user = sm.createAccount();

        register();

        // mock message to bmc
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq(METACHAIN),
                eq(SERVICE), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(irc2, "tokenFallback", user.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(user, "transfer", TEST_TOKEN, BigInteger.valueOf(10), btpAddress);
        contextMock.when(
                        () -> Context.call(eq(BigInteger.class), eq(irc2.getAddress()), eq("balanceOf"), eq(user.getAddress())))
                .thenReturn(BigInteger.ZERO);

        Map<String, BigInteger> balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(),
                TEST_TOKEN);
        assertEquals(BigInteger.TEN, balance.get("locked"));

        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_OK);
        message.setData(response.toBytes());

        BigInteger sn = (BigInteger) score.call("getSn");
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, sn, message.toBytes());
        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), TEST_TOKEN);
        assertEquals(BigInteger.ZERO, balance.get("locked"));
        assertEquals(BigInteger.ZERO, balance.get("refundable"));

        // another token transfer ,error response message, amount should go to refundable
        score.invoke(irc2, "tokenFallback", user.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(user, "transfer", TEST_TOKEN, BigInteger.valueOf(10), btpAddress);

        message = new BTSMessage();
        message.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
        response = new TransferResponse();
        response.setCode(TransferResponse.RC_ERR);
        message.setData(response.toBytes());

        sn = (BigInteger) score.call("getSn");
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, sn, message.toBytes());
        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), TEST_TOKEN);
        assertEquals(BigInteger.ZERO, balance.get("locked"));
        assertEquals(BigInteger.valueOf(9), balance.get("refundable"));

        // handleBTPError
        score.invoke(irc2, "tokenFallback", user.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(user, "transfer", TEST_TOKEN, BigInteger.valueOf(10), btpAddress);
        sn = (BigInteger) score.call("getSn");
        score.invoke(bmcMock, "handleBTPError", METACHAIN, SERVICE, sn, Long.valueOf(1), "");
        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), TEST_TOKEN);
        assertEquals(BigInteger.ZERO, balance.get("locked"));
        assertEquals(BigInteger.valueOf(18), balance.get("refundable"));

        // reclaim tokens
        contextMock.when(
                        () -> Context.call(irc2.getAddress(), "transfer", user.getAddress(), BigInteger.valueOf(18), null))
                .thenReturn(null);
        score.invoke(user, "reclaim", TEST_TOKEN, BigInteger.valueOf(18));
        balance = (Map<String, BigInteger>) score.call("balanceOf", user.getAddress(), TEST_TOKEN);
        assertEquals(BigInteger.ZERO, balance.get("refundable"));
        assertEquals(BigInteger.ZERO, balance.get("locked"));
    }

    @Test
    @Order(11)
    @Disabled
    public void transferBatch() {
        String[] coinNames = new String[]{"Token1", "Token2", "Token3", "Token4", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val, val, val, val, val};
        String destination = generateBTPAddress("harmony", ETH_ADDR);

        // none of them registered
        // user has not deposited any of them
        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "Not supported Token");

        BigInteger[] wrongValues = new BigInteger[]{val, val, val};

        call = () -> score.invoke(nonOwner, "transferBatch", coinNames, wrongValues, destination);
        expectErrorMessage(call, "Invalid arguments");

        String[] coinNamesBatch = new String[20];
        BigInteger[] coinValuesBatch = new BigInteger[20];
        Random rand = new Random();
        for (int i = 0; i < 20; i++) {
            coinNamesBatch[i] = String.valueOf(rand.nextInt());
            coinValuesBatch[i] = new BigInteger(String.valueOf(rand.nextInt()));
        }
        call = () -> score.invoke(nonOwner, "transferBatch", coinNamesBatch, coinValuesBatch, destination);
        expectErrorMessage(call,  "Cannot transfer over 15 tokens at once");

        // add native coin as well in batch
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));

        // register tokens
        Account token1 = Account.newScoreAccount(10);
        Account token2 = Account.newScoreAccount(11);
        Account token3 = Account.newScoreAccount(12);
        Account token4 = Account.newScoreAccount(13);

        // register irc2 token
        register(coinNames[0], token1.getAddress());
        register(coinNames[1], token2.getAddress());
        register(coinNames[2], token3.getAddress());
        register(coinNames[3], token4.getAddress());

        // register wrapped token "PARA"
        registerWrapped();

        assertEquals(BigInteger.ZERO, score.call("getSn"));

        score.invoke(token1, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(8), "0".getBytes());
        score.invoke(token2, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(token3, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(token4, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());

        // PARA wrapped transferFrom mock
        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);

        call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessageIn(call, "InSufficient Usable Balance");

        score.invoke(token1, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());

        // mock message to bmc
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("harmony"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        verify(scoreSpy).TransferStart(eq(nonOwner.getAddress()), eq(destination), eq(BigInteger.ONE), any());

        contextMock.when(() -> Context.call(eq(BigInteger.class), any(), eq("balanceOf"), eq(nonOwner.getAddress())))
                .thenReturn(BigInteger.valueOf(100));

        Map<String, BigInteger> bal = (Map<String, BigInteger>) score.call("balanceOf", nonOwner.getAddress(),
                coinNames[0]);
        assertEquals(bal.get("usable"), BigInteger.valueOf(8));
        assertEquals(bal.get("locked"), BigInteger.valueOf(10));

        // PARA wrapped token balance mock
        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(wrappedIRC2.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(50));

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(wrappedIRC2.getAddress()), eq("allowance"),
                eq(nonOwner.getAddress()), eq(score.getAddress()))).thenReturn(BigInteger.valueOf(40));

        bal = (Map<String, BigInteger>) score.call("balanceOf", nonOwner.getAddress(), coinNames[4]);
        // usable is min of balanceOf and allowed
        assertEquals(bal.get("usable"), BigInteger.valueOf(40));
        assertEquals(bal.get("locked"), BigInteger.valueOf(10));

        // message service metachain address 3er 1
        TransferTransaction txn = (TransferTransaction) score.call("getTransaction", BigInteger.ONE);
        assertEquals(txn.getFrom(), nonOwner.getAddress().toString());
        assertEquals(txn.getTo(), destination);
        assertEquals(txn.getAssets().length, 6);

        // this goes to bmc -> relayer -> solidity -> relayer -> bmc -> here again
    }

    @Test
    @Order(12)
    public void handleBTPMessage1() {
        // solidity -> relayer -> bmc -> bts -> transfer/mint on icon side
        Executable call = () -> score.invoke(owner, "handleBTPMessage", "from", SERVICE, BigInteger.ONE,
                "ehehehe".getBytes());
        expectErrorMessage(call, "Only BMC");

        // irc2, wrapped and native-coin
        Asset asset1 = new Asset(TEST_TOKEN, BigInteger.valueOf(50));
        Asset asset2 = new Asset(PARA, BigInteger.valueOf(30));
        Asset asset3 = new Asset(ICON, BigInteger.valueOf(20));

        // TransferRequest Message
        TransferRequest request = new TransferRequest();
        request.setFrom(bmcMock.getAddress().toString());
        request.setTo(nonOwner.getAddress().toString());
        request.setAssets(new Asset[]{asset1, asset2, asset3});

        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.REQUEST_COIN_TRANSFER);
        message.setData(request.toBytes());

        byte[] _msg = message.toBytes();

        // ICON preregistered, register remaining two
        Account token1 = Account.newScoreAccount(5);
        register(TEST_TOKEN, token1.getAddress());

        assertEquals(BigInteger.ZERO, score.call("getSn"));

        // wrapped token not registered yet
        call = () -> score.invoke(bmcMock, "handleBTPMessage", "from", SERVICE, BigInteger.valueOf(3), _msg);
        expectErrorMessage(call, "Invalid Token");

        // register wrapped token
        registerWrapped();

        assertEquals(BigInteger.ZERO, score.call("getSn"));

        // native-coin
        contextMock.when(() -> Context.transfer(any(Address.class), eq(BigInteger.valueOf(20))))
                .then(invocationOnMock -> null);
        // wrapped coin
        contextMock.when(() -> Context.call(any(), eq("mint"), eq(nonOwner.getAddress()), eq(BigInteger.valueOf(30))))
                .thenReturn(null);
        // irc2 transfer
        contextMock.when(() -> Context.call(eq(token1.getAddress()), eq("transfer"), eq(nonOwner.getAddress()),
                eq(BigInteger.valueOf(50)), any())).thenReturn(null);

        // mock bmc message
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("fromString"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.ONE, _msg);
        verify(scoreSpy).TransferReceived(eq("fromString"), eq(nonOwner.getAddress()), eq(BigInteger.ONE), any());
    }

    @Test
    @Order(13)
    @Disabled
    public void handleBTPMessage2() {

        // request plus response
        String[] coinNames = new String[]{"Token1", "Token2", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger val1 = BigInteger.valueOf(10000);
        BigInteger[] values = new BigInteger[]{val, val1, val};
        String destination = generateBTPAddress("0x1.bsc", ETH_ADDR);

        // add native coin as well in batch
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(50));

        Account token1 = Account.newScoreAccount(50);
        Account token2 = Account.newScoreAccount(51);

        register("Token1", token1.getAddress());
        register("Token2", token2.getAddress());
        registerWrapped();

        // transfer 50 of token 1
        score.invoke(token1, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(50), "0".getBytes());
        score.invoke(token2, "tokenFallback", nonOwner.getAddress(), val1, "0".getBytes());
        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);

        sendBTPMessageMock();
        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        verify(scoreSpy).TransferStart(any(), any(), any(), any());
        assertEquals(BigInteger.ONE, score.call("getSn"));

        // mock allowance and balanceOf for wrapped token

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(wrappedIRC2.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(wrappedIRC2.getAddress()), eq("allowance"),
                eq(nonOwner.getAddress()), eq(score.getAddress()))).thenReturn(BigInteger.valueOf(1000));

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(token1.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));
        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(token2.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));

        List<Map<String, BigInteger>> balances = (List<Map<String, BigInteger>>) score.call("balanceOfBatch",
                nonOwner.getAddress(), new String[]{"Token1", "Token2", PARA, ICON});

        // Token1
        assertEquals(balances.get(0).get("locked"), BigInteger.valueOf(10));
        assertEquals(balances.get(0).get("usable"), BigInteger.valueOf(40));
        // Token 2
        assertEquals(balances.get(1).get("locked"), val1);
        // PARA
        assertEquals(balances.get(2).get("locked"), BigInteger.valueOf(10));
        // ICON
        assertEquals(balances.get(3).get("locked"), BigInteger.valueOf(50));

        // BTS -> BMC -> RELAYER -> SOLIDITY -> RELAYER -> BMC -> BTS

        assertEquals(BigInteger.ONE, score.call(("getSn")));
        // successful case
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_OK);

        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
        message.setData(response.toBytes());

        byte[] _msg = message.toBytes();

        contextMock.when(() -> Context.call(eq(wrappedIRC2.getAddress()), eq("burn"), eq(BigInteger.valueOf(8))))
                .thenReturn(null);

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.ONE, _msg);
        verify(scoreSpy).TransferEnd(eq(nonOwner.getAddress()), eq(BigInteger.ONE), eq(TransferResponse.RC_OK), any());
        assertEquals(null, score.call("getTransaction", BigInteger.ONE));

        balances = (List<Map<String, BigInteger>>) score.call("balanceOfBatch", nonOwner.getAddress(),
                new String[]{"Token1", "Token2", PARA, ICON});

        // Token1
        assertEquals(balances.get(0).get("locked"), BigInteger.ZERO);
        assertEquals(balances.get(0).get("usable"), BigInteger.valueOf(40));
        // Token 2
        assertEquals(balances.get(1).get("locked"), BigInteger.ZERO);
        // PARA
        assertEquals(balances.get(2).get("locked"), BigInteger.ZERO);
        // ICON
        assertEquals(balances.get(3).get("locked"), BigInteger.ZERO);

        // fee for wrapped and native-coin set to zero
        Map<String, BigInteger> fees = (Map<String, BigInteger>) score.call("getAccumulatedFees");
        assertEquals(BigInteger.ONE, fees.get(coinNames[0]));
        assertEquals(BigInteger.valueOf(11), fees.get(coinNames[1]));
        assertEquals(BigInteger.TWO, fees.get(coinNames[2]));
        assertEquals(BigInteger.ZERO, fees.get(ICON));
    }

    @Test
    @Order(14)
    @DisplayName("unknown type")
    void handleBTPMessage3() {
        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.UNKNOWN_TYPE);
        message.setData("a for apple".getBytes());

        score.invoke(bmcMock, "handleBTPMessage", "from", SERVICE, BigInteger.ONE, message.toBytes());
        verify(scoreSpy).UnknownResponse(any(), eq(BigInteger.ONE));
    }

    @Test
    @Order(15)
    @DisplayName("not handled cases")
    void handleBTPMessage4() {
        BTSMessage message = new BTSMessage();
        message.setServiceType(100);
        message.setData("a for apple".getBytes());

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("fromString"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.ONE, message.toBytes());
    }

    @Test
    @Order(16)
    @DisplayName("change limit, check restrictions, check fees")
    public void changeLimit() {
        String TOKEN1 = "Token1";
        BigInteger TWO_HUNDRED_ICX = BigInteger.valueOf(200).multiply(ICX);
        BigInteger FIFTY_ICX = BigInteger.valueOf(50).multiply(ICX);
        Address user1 = sm.createAccount(20).getAddress();

        Account token1 = Account.newScoreAccount(50);
        Account token2 = Account.newScoreAccount(51);

        Executable call;

        score.invoke(owner, "register", TOKEN1, TOKEN1, 18, BigInteger.valueOf(10), ICX, token1.getAddress());

        assertEquals(UINT_CAP, score.call("getTokenLimit", TOKEN1));

        tokenLimitBTPMessage();
        score.invoke(owner, "setTokenLimit", new String[]{TOKEN1}, new BigInteger[]{BigInteger.valueOf(200)});

        assertEquals(BigInteger.valueOf(200), score.call("getTokenLimit", TOKEN1));
//
        call = () -> score.invoke(owner, "setTokenLimit", new String[]{TOKEN1},
                new BigInteger[]{TWO_HUNDRED_ICX.negate()});
        expectErrorMessage(call, "Invalid value");
        // can set for tokens not registered as well
        score.invoke(owner, "setTokenLimit", new String[]{"TokenBSH"}, new BigInteger[]{TWO_HUNDRED_ICX});
        assertEquals(TWO_HUNDRED_ICX, score.call("getTokenLimit", "TokenBSH"));

//        expectErrorMessage(call, "Not registered");

        score.invoke(owner, "setTokenLimit", new String[]{TOKEN1}, new BigInteger[]{TWO_HUNDRED_ICX});

        assertEquals(TWO_HUNDRED_ICX, score.call("getTokenLimit", TOKEN1));

//         send from Harmony to ICON
//         reciever adddress is blacklisted in icon

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), any(), any(), any(),
                any());

        Verification returnLinks = () -> Context.call(eq(String[].class), eq(bmcMock.getAddress()), eq("getLinks"),
                any());

        // these 3 BTP addresses are currently supported by ICON Bridge
        BTPAddress btpAddress1 = new BTPAddress("icon", "cx0E636c8cF214a9d702C5E4a6D8c020be217766D3");
        BTPAddress btpAddress2 = new BTPAddress(METACHAIN, "0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
        BTPAddress btpAddress3 = new BTPAddress("harmony", "0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
        String[] links = new String[]{btpAddress1.toString(), btpAddress2.toString(), btpAddress3.toString()};

        contextMock.when(sendMessage).thenReturn(null);
        contextMock.when(returnLinks).thenReturn(links);

        // blacklist user1 in icon
        blackListUser("icon", new String[]{user1.toString()}, BigInteger.valueOf(3));

        // handleRequest of coinTransfer coming from harmony
        // fee for this transfer will be handled in harmony side

        // value within limit
        Asset asset1 = new Asset(TOKEN1, FIFTY_ICX);

        // TransferRequest Message
        TransferRequest request = new TransferRequest();
        request.setFrom(ETH_ADDR);

        // user1 is blacklisted
        request.setTo(user1.toString());
        request.setAssets(new Asset[]{asset1});

        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.REQUEST_COIN_TRANSFER);
        message.setData(request.toBytes());

        byte[] _msg = message.toBytes();

        contextMock.when(() -> Context.call(eq(token1.getAddress()), eq("transfer"), eq(user1), eq(FIFTY_ICX), any()))
                .thenReturn(null);

        score.invoke(owner, "addRestriction");

        call = () -> score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.ONE, _msg);
        expectErrorMessage(call, "_to user is Blacklisted");

        removeBlackListedUser("icon", new String[]{user1.toString()}, BigInteger.valueOf(4));

        // check this
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(5), _msg);

        // another transfer, but over the limit
        // user not blacklisted

        Asset asset2 = new Asset(TOKEN1, FIFTY_ICX.add(TWO_HUNDRED_ICX));
        // TransferRequest Message
        TransferRequest request2 = new TransferRequest();
        request2.setFrom(ETH_ADDR);
        request2.setTo(user1.toString());
        request2.setAssets(new Asset[]{asset2});

        BTSMessage message2 = new BTSMessage();
        message2.setServiceType(BTSMessage.REQUEST_COIN_TRANSFER);
        message2.setData(request2.toBytes());

        byte[] _msg2 = message2.toBytes();
        contextMock.when(() -> Context.call(eq(token1.getAddress()), eq("transfer"), eq(user1),
                eq(FIFTY_ICX.add(TWO_HUNDRED_ICX)), any())).thenReturn(null);

        call = () -> score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(6), _msg2);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        // limit not set
        // any amount should be able to be transferred off chain
        String TOKEN2 = "Token 2";
        BigInteger MILLION = BigInteger.valueOf(10_000_000L).multiply(ICX);
        score.invoke(owner, "register", TOKEN2, TOKEN2, 18, BigInteger.valueOf(10), ICX, token2.getAddress());

        assertEquals(UINT_CAP, score.call("getTokenLimit", TOKEN2));

        Asset asset3 = new Asset(TOKEN2, MILLION);

        // TransferRequest Message
        TransferRequest request3 = new TransferRequest();
        request3.setFrom(ETH_ADDR);

        // user1 is blacklisted
        request3.setTo(user1.toString());
        request3.setAssets(new Asset[]{asset3});

        BTSMessage message3 = new BTSMessage();
        message3.setServiceType(BTSMessage.REQUEST_COIN_TRANSFER);
        message3.setData(request3.toBytes());

        byte[] _msg3 = message3.toBytes();

        contextMock.when(() -> Context.call(eq(token2.getAddress()), eq("transfer"), eq(user1), eq(MILLION), any()))
                .thenReturn(null);

        // change token limit
        score.invoke(owner, "setTokenLimit", new String[]{TOKEN2}, new BigInteger[]{TWO_HUNDRED_ICX});
        call = () -> score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(6), _msg3);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        // disable restrictions
        score.invoke(owner, "disableRestrictions");
        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, BigInteger.valueOf(7), _msg3);
    }

    @Test
    @Order(17)
    @DisplayName("add and remove owner")
    public void ownerTests() {
        String expectedErrorMessage = "caller is not owner";
        Account user1 = sm.createAccount(10);
        Account user2 = sm.createAccount(10);

        // Scenario 17: Non-Owner tries to add a new Owner
        Executable call = () -> score.invoke(nonOwner, "addOwner", owner.getAddress());
        expectErrorMessage(call, expectedErrorMessage);

        // owner tries to add themselves
        call = () -> score.invoke(owner, "addOwner", owner.getAddress());
        expectErrorMessage(call, "given address is score owner");

        // Scenario 18: Current Owner adds a new Owner
        score.invoke(owner, "addOwner", user1.getAddress());
        assertEquals(true, score.call("isOwner", user1.getAddress()));
        Address[] owners = (Address[]) score.call("getOwners");
        assertEquals(owner.getAddress(), owners[0]);
        assertEquals(user1.getAddress(), owners[1]);

        // Scenario 19: After adding a new Owner, owner registers a new coin
        Address rndmAddr = Account.newScoreAccount(13).getAddress();
        register("Random Token", rndmAddr);
        assertEquals(rndmAddr, score.call("coinId", "Random Token"));

        // Scenario 20: New Owner registers a new coin
        rndmAddr = Account.newScoreAccount(14).getAddress();
        register("wBTC", rndmAddr);
        assertEquals(rndmAddr, score.call("coinId", "wBTC"));

        // newly added owner tries to add owner
        score.invoke(user1, "addOwner", user2.getAddress());
        assertEquals(true, score.call("isOwner", user2.getAddress()));

        //Scenario 30: Current Owner removes another Owner
        score.invoke(user2, "removeOwner", user1.getAddress());
        assertEquals(false, score.call("isOwner", user1.getAddress()));

        // owner tries to add itself again
        call = () -> score.invoke(user2, "addOwner", user2.getAddress());
        expectErrorMessage(call, "already exists owner");

        // Scenario 31: The last Owner removes him/herself
        score.invoke(user2, "removeOwner", user2.getAddress());
        assertEquals(false, score.call("isOwner", user2.getAddress()));
    }

    @Test
    @Order(18)
    public void setFeeRatio() {
        registerWrapped();

        // Scenario 10: None-ownership role client updates a new fee ratio
        // Scenario 13: Non-ownership role client updates a new fixed_fee
        Executable call = () -> score.invoke(nonOwner, "setFeeRatio", PARA, BigInteger.valueOf(10), ICX);
        expectErrorMessage(call, "require owner access");

        call = () -> score.invoke(owner, "setFeeRatio", PARA, BigInteger.valueOf(10).negate(), ICX.negate());
        expectErrorMessageIn(call,
                "The feeNumerator should be less " + "than FEE_DENOMINATOR and feeNumerator should be greater than 1");

        // Scenario 11: Fee_numerator is set higher than fee_denominator
        call = () -> score.invoke(owner, "setFeeRatio", PARA, ICX, ICX.negate());
        expectErrorMessageIn(call,
                "The feeNumerator should be less " + "than FEE_DENOMINATOR and feeNumerator should be greater than 1");

        call = () -> score.invoke(owner, "setFeeRatio", PARA, BigInteger.valueOf(10), ICX.negate());
        expectErrorMessageIn(call, "Fixed fee cannot be less than zero");

        call = () -> score.invoke(owner, "setFeeRatio", "LAMB", BigInteger.valueOf(100), ICX);
        expectErrorMessage(call, "Not supported Coin");

        // Scenario 9: Contract’s owner updates a new fee ratio
        // Scenario 12: Contract’s owner updates fixed fee
        score.invoke(owner, "setFeeRatio", PARA, BigInteger.valueOf(100), ICX);
    }

    @Test
    @Order(19)
    public void coinDetails() {
        // Scenario 15: Query a valid supporting coin
        registerWrapped();
        assertEquals(wrappedIRC2.getAddress(), score.call("coinId", PARA));

        // Scenario 16: Query an invalid supporting coin
        assertEquals(null, score.call("coinId", "DUM"));

    }


    @Test
    @Order(20)
    @Disabled
    public void transferBatchNativecoinTokenLimit() {
        tokenLimitBTPMessage();
        score.invoke(owner, "setTokenLimit", new String[]{ICON}, new BigInteger[]{BigInteger.valueOf(1000)});
        assertEquals(BigInteger.valueOf(1000), score.call("getTokenLimit", ICON));

        String destination = generateBTPAddress("0x1.bsc", ETH_ADDR);
        Account token1 = Account.newScoreAccount(50);
        Account token2 = Account.newScoreAccount(51);

        register("Token1", token1.getAddress());
        register("Token2", token2.getAddress());

        score.invoke(token1, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(50), "0".getBytes());
        score.invoke(token2, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(50), "0".getBytes());
        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);

        // transferBatch 2 tokens, try to transfer ICX above 1000 should fail
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(1001));
        String[] coinNames = new String[]{"Token1", "Token2"};
        BigInteger[] values = new BigInteger[]{BigInteger.valueOf(50), BigInteger.valueOf(50)};

        sendBTPMessageMock();
        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        // 999 transfer should be successful
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(999));
        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);

    }

    @Test
    @Order(20)
    @Disabled
    public void handleFeeGathering() {
        String feeAggregator = generateBTPAddress("icon", "hx0000000000000000000000000000000000000000");

        // transferBatch
        String[] coinNames = new String[]{"Token1", "Token2", "Token3", "Token4", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val, val, val, val, val};
        String destination = generateBTPAddress(METACHAIN, ETH_ADDR);

        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));

        // register tokens
        Account token1 = Account.newScoreAccount(10);
        Account token2 = Account.newScoreAccount(11);
        Account token3 = Account.newScoreAccount(12);
        Account token4 = Account.newScoreAccount(13);

        // register irc2 token
        register(coinNames[0], token1.getAddress());
        register(coinNames[1], token2.getAddress());
        register(coinNames[2], token3.getAddress());
        register(coinNames[3], token4.getAddress());

        // register wrapped token "PARA"
        registerWrapped();

        assertEquals(BigInteger.ZERO, score.call("getSn"));

        score.invoke(token1, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(100), "0".getBytes());
        score.invoke(token2, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(100), "0".getBytes());
        score.invoke(token3, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(100), "0".getBytes());
        score.invoke(token4, "tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(100), "0".getBytes());

        // mock approved amount for PARA
        Verification transferFromMock = () -> Context.call(eq(Boolean.class),
                any(Address.class), eq("transferFrom"), eq(nonOwner.getAddress()), eq(score.getAddress()),
                any(), eq(null));
        contextMock.when(transferFromMock).thenReturn(true);

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq(METACHAIN),
                eq(SERVICE), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        verify(scoreSpy).TransferStart(eq(nonOwner.getAddress()), eq(destination), eq(BigInteger.ONE), any());

        // receive successful transfer response back
        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
        TransferResponse response = new TransferResponse();
        response.setCode(TransferResponse.RC_OK);
        message.setData(response.toBytes());

        BigInteger sn = (BigInteger) score.call("getSn");
        // burn PARA token
        Verification burnMock = () -> Context.call(
                any(Address.class), eq("burn"), eq(BigInteger.valueOf(8)));
        contextMock.when(burnMock).thenReturn(true);

        score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, sn, message.toBytes());

        // fees collected on BTS contract
        // [1,1,1,1,2]
        // send fees to fee gathering
        // mock transfer to fee aggregator for IRC2 Tokens

        Map<String, BigInteger> fees = (Map<String, BigInteger>) score.call("getAccumulatedFees");
        assertEquals(BigInteger.ONE, fees.get("Token1"));
        assertEquals(BigInteger.ONE, fees.get("Token2"));
        assertEquals(BigInteger.ONE, fees.get("Token3"));
        assertEquals(BigInteger.ONE, fees.get("Token4"));
        assertEquals(BigInteger.TWO, fees.get(PARA));

        Address ZERO_ADDR = Address.fromString("hx0000000000000000000000000000000000000000");
        Verification txferMock = () -> Context.call(
                any(Address.class), eq("transfer"),
                eq(ZERO_ADDR), eq(BigInteger.ONE), any());
        contextMock.when(txferMock).thenReturn(true);

        contextMock.when(() -> Context.call(
                eq(Boolean.class), any(Address.class), eq("approve"), eq(score.getAddress()), eq(BigInteger.TWO)
        )).thenReturn(true);

        transferFromMock = () -> Context.call(eq(Boolean.class),
                any(Address.class), eq("transferFrom"), eq(score.getAddress()), eq(ZERO_ADDR),
                eq(BigInteger.TWO), eq(null));
        contextMock.when(transferFromMock).thenReturn(true);

        score.invoke(bmcMock, "handleFeeGathering", feeAggregator, "bmc");
    }

    @Test
    @Order(21)
    void readonlyMethods() {
        // unregistered coin fee ratio
        Map<String, BigInteger> fee = (Map<String, BigInteger>) score.call("feeRatio", "ABCD");
        assertEquals(BigInteger.ZERO, fee.get("fixedFee"));
        assertEquals(BigInteger.ZERO, fee.get("feeNumerator"));
    }


    @Test
    @Order(22)
    void setTokenLimitResponse() {
        // set token limit for PARA and IRC2

        getLinksMock();
        tokenLimitBTPMessage();
        String[] coinNames = new String[] {ICON, PARA};
        BigInteger[] tokenLimits = new BigInteger[] {BigInteger.valueOf(100), BigInteger.valueOf(100)};

        score.invoke(owner, "setTokenLimit", coinNames, tokenLimits);

        assertEquals( false, score.call("tokenLimitStatus", METACHAIN, ICON));
        assertEquals( false, score.call("tokenLimitStatus", METACHAIN, PARA));
        assertEquals( false, score.call("tokenLimitStatus", "icon", ICON));
        assertEquals( false, score.call("tokenLimitStatus", "icon", PARA));
        assertEquals( false, score.call("tokenLimitStatus", "harmony", ICON));
        assertEquals( false, score.call("tokenLimitStatus", "harmony", PARA));

        BigInteger sn = (BigInteger) score.call("getSn");

        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.CHANGE_TOKEN_LIMIT);
        message.setData("tokenLimitSuccessfulResponse".getBytes());

        Executable call = () -> score.invoke(bmcMock, "handleBTPMessage", "fromString", SERVICE, sn, message.toBytes());
        expectErrorMessage(call, "Invalid change limit transaction");

        TokenLimitResponse response = new TokenLimitResponse();
        response.setCode(TokenLimitResponse.RC_OK);
        response.setMessage("tokenLimitSuccessfulResponse");
        message.setData(response.toBytes());
        score.invoke(bmcMock, "handleBTPMessage", METACHAIN, SERVICE, sn, message.toBytes());
        verify(scoreSpy).TokenLimitSet(eq(sn), eq("tokenLimitSuccessfulResponse".getBytes()));

        assertEquals( true, score.call("tokenLimitStatus", METACHAIN, ICON));
        assertEquals( true, score.call("tokenLimitStatus", METACHAIN, PARA));
        assertEquals( false, score.call("tokenLimitStatus", "icon", ICON));
        assertEquals( false, score.call("tokenLimitStatus", "icon", PARA));
        assertEquals( false, score.call("tokenLimitStatus", "harmony", ICON));
        assertEquals( false, score.call("tokenLimitStatus", "harmony", PARA));

        TokenLimitTransaction txn = ( TokenLimitTransaction ) score.call("getTokenLimitTxn", sn);
        String[] nets = txn.getNet();
        assertEquals("icon", nets[0]);
        assertEquals("harmony", nets[1]);

        score.invoke(bmcMock, "handleBTPMessage", "harmony", SERVICE, sn, message.toBytes());
        assertEquals( true, score.call("tokenLimitStatus", "harmony", ICON));
        assertEquals( true, score.call("tokenLimitStatus", "harmony", PARA));

        txn = ( TokenLimitTransaction ) score.call("getTokenLimitTxn", sn);
        nets = txn.getNet();
        assertEquals("icon", nets[0]);

        // failed response from "icon"
        score.invoke(bmcMock, "handleBTPError", "icon", SERVICE, sn, 1L, "message");
        assertEquals( false, score.call("tokenLimitStatus", "icon", ICON));
        assertEquals( false, score.call("tokenLimitStatus", "icon", PARA));

        // successful message from "icon"
        score.invoke(bmcMock, "handleBTPMessage", "icon", SERVICE, sn, message.toBytes());

        assertEquals(null, score.call("getTokenLimitTxn", sn));
        assertEquals( true, score.call("tokenLimitStatus", METACHAIN, ICON));
        assertEquals( true, score.call("tokenLimitStatus", METACHAIN, PARA));
        assertEquals( true, score.call("tokenLimitStatus", "icon", ICON));
        assertEquals( true, score.call("tokenLimitStatus", "icon", PARA));
        assertEquals( true, score.call("tokenLimitStatus", "harmony", ICON));
        assertEquals( true, score.call("tokenLimitStatus", "harmony", PARA));

        // suppose METACHAIN Sends response again
        score.invoke(bmcMock, "handleBTPMessage", "icon", SERVICE, sn, message.toBytes());

        score.invoke(bmcMock, "handleBTPError", "icon", SERVICE, sn, 1L, "message");
        // token limit txn is already set to null
        assertEquals( true, score.call("tokenLimitStatus", "icon", ICON));
        assertEquals( true, score.call("tokenLimitStatus", "icon", PARA));
    }

    // ICON BRIDGE MIGRATION TESTS

    @Test
    public void preventICXTransfer() {
        sendBTPMessageMock();
        String btpAddress = generateBTPAddress(METACHAIN, owner.getAddress().toString());
        // general condition
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));

        Executable call = () -> score.invoke(owner, "transferNativeCoin", btpAddress);
        expectErrorMessage(call, "Reverted(0): Cannot transfer ICX.");
    }

    @Test
    @Disabled
    public void preventIconTokenTransfer() {
        sendBTPMessageMock();
        String btpaddr = generateBTPAddress(METACHAIN, owner.getAddress().toString());

        register();

        // should not transfer native-coin
        Executable call = () -> score.invoke(nonOwner, "transfer", TEST_TOKEN, BigInteger.TEN, btpaddr);
        expectErrorMessage(call, "Cannnot transfer icon tokens anymore.");

    }

    @Test
    public void preventEthTokenTransfer() {
        sendBTPMessageMock();
        String btpaddr = generateBTPAddress(METACHAIN, owner.getAddress().toString());

        String tokenName = "btp-0x38.bsc-eth";

        Verification deployWrappedToken = () -> Context.deploy(any(), eq(tokenName),
                eq(tokenName),eq(18));
        contextMock.when(deployWrappedToken).thenReturn(wrappedIRC2.getAddress());

        score.invoke(owner, "register",tokenName, tokenName, 18, BigInteger.ZERO, BigInteger.TWO,
                    Address.fromString("cx0000000000000000000000000000000000000000"));


        // should not transfer native-coin
        Executable call = () -> score.invoke(nonOwner, "transfer", tokenName, BigInteger.TEN, btpaddr);
        expectErrorMessage(call, "NotETH");

    }

    @Test
    public void migrationRestrictionOnTransferBatch_includeICX() {

        String[] coinNames = new String[]{"Token1", "Token2", "Token3", "Token4", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val, val, val, val, val};
        String destination = generateBTPAddress("harmony", ETH_ADDR);

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("harmony"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);


        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));

        // register tokens
        Account token1 = Account.newScoreAccount(10);
        Account token2 = Account.newScoreAccount(11);
        Account token3 = Account.newScoreAccount(12);
        Account token4 = Account.newScoreAccount(13);

        // register irc2 token
        register(coinNames[0], token1.getAddress());
        register(coinNames[1], token2.getAddress());
        register(coinNames[2], token3.getAddress());
        register(coinNames[3], token4.getAddress());


        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "Reverted(0): Cannot transfer ICX.");

    }

    @Test
    public void migrationRestrictionOnTransferBatch_includeETHTokens() {
        String tokenName = "btp-0x38.bsc-eth";

        String[] coinNames = new String[]{tokenName,"Token1", "Token2", "Token3", "Token4", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val, val, val, val, val, val};
        String destination = generateBTPAddress("harmony", ETH_ADDR);

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("harmony"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);


        // register tokens
        Account token1 = Account.newScoreAccount(10);
        Account token2 = Account.newScoreAccount(11);
        Account token3 = Account.newScoreAccount(12);
        Account token4 = Account.newScoreAccount(13);

        // register irc2 token
        register(coinNames[1], token1.getAddress());
        register(coinNames[2], token2.getAddress());
        register(coinNames[3], token3.getAddress());
        register(coinNames[4], token4.getAddress());

        Verification deployWrappedToken = () -> Context.deploy(any(), eq(tokenName),
                eq(tokenName),eq(18));
        contextMock.when(deployWrappedToken).thenReturn(wrappedIRC2.getAddress());

        score.invoke(owner, "register",tokenName, tokenName, 18, BigInteger.ZERO, BigInteger.TWO,
                    Address.fromString("cx0000000000000000000000000000000000000000"));


        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "NotETH");
    }

    @Test
    public void migrationRestrictionOnTransferBatch_includeWrappedTokensOnly() {
        String tokenName = "META";

        String[] coinNames = new String[]{tokenName, PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val, val};
        String destination = generateBTPAddress("harmony", ETH_ADDR);

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("harmony"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        registerWrapped();

        Verification deployWrappedToken = () -> Context.deploy(any(), eq(tokenName),
                eq(tokenName),eq(18));
        contextMock.when(deployWrappedToken).thenReturn(wrappedIRC2.getAddress());

        score.invoke(owner, "register",tokenName, tokenName, 18, BigInteger.ZERO, BigInteger.TWO,
                    Address.fromString("cx0000000000000000000000000000000000000000"));

        // PARA and META registered

        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);


        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
    }

    @Test
    public void migrationRestrictionOnTransferBatch_includeIconTokens() {

        String[] coinNames = new String[]{"Token1", "Token2", "Token3", "Token4", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val, val, val, val, val};
        String destination = generateBTPAddress("harmony", ETH_ADDR);

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()), eq("sendMessage"), eq("harmony"),
                eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);


        // register tokens
        Account token1 = Account.newScoreAccount(10);
        Account token2 = Account.newScoreAccount(11);
        Account token3 = Account.newScoreAccount(12);
        Account token4 = Account.newScoreAccount(13);

        // register irc2 token
        register(coinNames[0], token1.getAddress());
        register(coinNames[1], token2.getAddress());
        register(coinNames[2], token3.getAddress());
        register(coinNames[3], token4.getAddress());


        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "Cannnot transfer icon tokens anymore.");
    }






}
