package foundation.icon.btp.bts;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;

import com.iconloop.score.test.Account;
import foundation.icon.btp.lib.BTPAddress;
import java.math.BigInteger;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.Order;
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
        String[] addr = new String[] {"Hell "};
        Executable  call = () -> score.invoke(nonOwner, "addBlacklistAddress", "hello world ",addr);
        expectErrorMessage(call, "require owner access");

        blacklistMocks();
        call = () -> score.invoke(owner, "addBlacklistAddress", "networking", addr);
        expectErrorMessage(call, "Invalid link");

        String[] addr1 = new String[] {"all too well"};
        score.invoke(owner, "addBlacklistAddress", "network", addr1);

        byte[] _msg = blacklistSuccessfulResponse();

        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.ONE, _msg);


        String[] addr2 = new String[] {"   all too well       "};
        score.invoke(owner, "addBlacklistAddress", "network", addr2);

        List<String> actual = (List<String>) score.call("getBlackListedUsers", "network", 0, 10);
        List<String> expected = List.of("all too well");
        assertEquals(expected, actual);

        String[] addr3 = new String[] {"  you belong with me  "};
        score.invoke(owner, "addBlacklistAddress", "network", addr3);
        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.TWO, _msg);

        expected = List.of("all too well", "you belong with me");
        actual = (List<String>) score.call("getBlackListedUsers", "network", 0, 10);
        assertEquals(expected, actual);


        assertEquals(true, score.call("isUserBlackListed", "network", " you belong with me "));
        assertEquals(false, score.call("isUserBlackListed", "network"," yu belong with me "));

        assertEquals(2, score.call("blackListedUsersCount", "network"));

        assertEquals( BigInteger.valueOf(2), score.call("getSn"));

        score.invoke(owner, "addBlacklistAddress", "icon",new String[]{" invalid icon address "});
        assertEquals(0, score.call("blackListedUsersCount", "icon"));

        score.invoke(owner, "addBlacklistAddress", "icon",
                new String[]{" cx42bd7394a8272fdb8683a41b92921247c34c522a  "});
        assertEquals(1, score.call("blackListedUsersCount", "icon"));
        assertEquals(true, score.call("isUserBlackListed", "icon", " cx42bd7394a8272fdb8683a41b92921247c34c522a "));
    }

    @Order(3)
    @Test
    public void removeFromLocklist() {
        assertEquals(BigInteger.ZERO, score.call(("getSn")));

        blacklistMocks();

        String[] addr1 = new String[] {"all too well"};
        score.invoke(owner, "addBlacklistAddress", "network", addr1);
        String[] addr2 = new String[] {"  you belong with me  "};
        score.invoke(owner, "addBlacklistAddress", "network", addr2);
        byte[] _msg = blacklistSuccessfulResponse();

        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.ONE, _msg);
        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.TWO, _msg);

        Executable  call = () -> score.invoke(owner, "removeBlacklistAddress", "icx",new String[]{" hell "});
        expectErrorMessage(call, "Invalid link");

        // non-owner tries
        call = () -> score.invoke(nonOwner, "removeBlacklistAddress", "harmony",new String[]{"hell"});
        expectErrorMessage(call, "require owner access");

        // try to remove non blacklisted
        score.invoke(owner, "removeBlacklistAddress", "harmony",new String[]{"hell"});

        // remove legit user
        score.invoke(owner, "removeBlacklistAddress", "network", addr1);

        _msg = blacklistRemoveSuccessfulResponse();

        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.valueOf(3), _msg);
        List<String> expected = List.of("you belong with me");
        List<String> actual = (List<String>) score.call("getBlackListedUsers", "network", 0, 1);
        assertEquals(expected, actual);

        assertEquals(1, score.call("blackListedUsersCount", "network"));
        assertEquals(false, score.call("isUserBlackListed",  "network", " all too well "));


        assertEquals(BigInteger.valueOf(3), score.call(("getSn")));
    }


    @Test
    @Order(4)
    @DisplayName("register wrapped coin")
    public void registerWrappedToken() {
        tokenLimitBTPMessage();

        Verification deployWrappedToken = () -> Context.deploy(any(), eq(PARA),
                eq(PARA),eq(18));
        contextMock.when(deployWrappedToken).thenReturn(wrappedIRC2.getAddress());

        score.invoke(owner, "register",PARA, PARA, 18, BigInteger.ZERO, BigInteger.ZERO,
                Address.fromString("cx0000000000000000000000000000000000000000"));

        assertEquals(1, score.call("getRegisteredTokensCount"));

        List<String> registered = List.of(ICON, PARA);
        assertEquals(registered, score.call("coinNames"));
    }

    @Test
    @Order(5)
    public void registerIRC2Token() {
        // non owner tries to call (Scenario 2)
        Executable  call = () -> score.invoke(nonOwner, "register",
                TEST_TOKEN, "TTK", 18, ICX.divide(BigInteger.TEN), ICX, irc2.getAddress());
        expectErrorMessage(call, "require owner access");

        // owner registers new coin (Scenario 1)
        score.invoke(owner, "register",
                TEST_TOKEN, "TTK", 18, ICX.divide(BigInteger.TEN), ICX, irc2.getAddress());

        // owner registers coin that exists (Scenario 3)
        call = () -> score.invoke(owner, "register",
                TEST_TOKEN, "TTK", 18, ICX.divide(BigInteger.TEN), ICX, irc2.getAddress());
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
        Executable call = () -> score.invoke(IRC2Token,"tokenFallback", owner.getAddress(),
                BigInteger.valueOf(10), "0".getBytes());
        expectErrorMessage(call,"Token not registered");
    }

    @Test
    @Order(7)
    public void operationOnRegisteredTokens() {
        register();

        contextMock.when(() -> Context.call(eq(BigInteger.class),eq(irc2.getAddress()), eq("balanceOf"),
                eq(owner.getAddress()))).thenReturn(BigInteger.valueOf(100));
        score.invoke(irc2,"tokenFallback", owner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        Map<String, BigInteger> balance = (Map<String , BigInteger>) score.call("balanceOf", owner.getAddress(), TEST_TOKEN);
        assertEquals(balance.get("usable"), BigInteger.valueOf(10));
    }

    @Test
    @Order(8)
    public void reclaimDepositedTokens() {
        // for IRC2 tokens
        register();
        Executable call = () -> score.invoke(owner, "reclaim", "ABC", BigInteger.valueOf(0));
        expectErrorMessage(call, "_value must be positive");
        call = () -> score.invoke(owner, "reclaim", "ABC", BigInteger.valueOf(-10));
        expectErrorMessage(call, "_value must be positive");
        call = () -> score.invoke(owner, "reclaim", "ABC", BigInteger.valueOf(10));
        expectErrorMessage(call, "invalid value");

        deposit(BigInteger.valueOf(10));
        // locked = 10
        contextMock.when(() -> Context.call(eq(irc2.getAddress()), eq("transfer"),
                eq(owner.getAddress()), eq(BigInteger.valueOf(5)), any())).thenReturn(null);
        score.invoke(owner, "reclaim", TEST_TOKEN, BigInteger.valueOf(5));

        Map<String, BigInteger> balance = (Map<String , BigInteger>) score.call("balanceOf", owner.getAddress(), TEST_TOKEN);
        assertEquals(BigInteger.valueOf(5), balance.get("refundable"));
        assertEquals(BigInteger.valueOf(0), balance.get("locked"));
        assertEquals(BigInteger.valueOf(0), balance.get("usable"));

        // wrapped IRC2 tradable is not directly transferred

    }

    @Test
    @Order(9)
    public void transferNativeCoin() {
        sendBTPMessageMock();
        String btpAddress = generateBTPAddress("network", owner.getAddress().toString());
        Account user = sm.createAccount(1000);
        // general condition
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(100));
        score.invoke(user, "transferNativeCoin", btpAddress);

        // blacklist reciever
        blacklistMocks();
        score.invoke(owner, "addRestriction");
        String addr = generateBTPAddress("network", user.getAddress().toString());
        String[] finalAdddr = new String[]{user.getAddress().toString()};
        blackListUser("network", finalAdddr, BigInteger.TWO);

//        score.invoke(owner, "addBlacklistAddress", "network", user.getAddress().toString());

        Executable call = () ->
                score.invoke(owner, "transferNativeCoin", addr);
        expectErrorMessage(call, "_to user is Blacklisted");

        // blacklist sender
        String newAddr = generateBTPAddress("network", owner.getAddress().toString());
        call = () -> score.invoke(user, "transferNativeCoin", newAddr);
        expectErrorMessage(call, "_from user is Blacklisted");

        // remove blacklisted user
        removeBlackListedUser("network", finalAdddr, BigInteger.valueOf(3));

        score.invoke(owner, "setTokenLimit", new String[]{ICON},
                new BigInteger[]{BigInteger.valueOf(90)});

        call = () -> score.invoke(owner, "transferNativeCoin", addr);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        score.invoke(owner, "setTokenLimit", new String[]{ICON},
                new BigInteger[]{BigInteger.valueOf(100)});

        score.invoke(owner, "transferNativeCoin", addr);

        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(99));
        score.invoke(owner, "transferNativeCoin", addr);
    }

    @Test
    @Order(10)
    public void transfer() {
        String btpaddr = generateBTPAddress("harmony",ETH_ADDR);

        Executable call = () -> score.invoke(nonOwner, "transfer", "Token1",
                BigInteger.valueOf(-1), btpaddr);
        expectErrorMessage(call, "Invalid amount");

        call = () -> score.invoke(nonOwner, "transfer", "Token1",
                BigInteger.ZERO, btpaddr);
        expectErrorMessage(call, "Invalid amount");

        call = () -> score.invoke(nonOwner, "transfer", "Token1",
                BigInteger.valueOf(10), btpaddr);
        expectErrorMessage(call, "Not supported Token");

        register();

        // mock message to bmc
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), eq("harmony"), eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(irc2,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(nonOwner, "transfer", TEST_TOKEN, BigInteger.valueOf(10), btpaddr);

        contextMock.when(() -> Context.call(
                eq(BigInteger.class),
                eq(irc2.getAddress()),
                eq("balanceOf"),
                eq(nonOwner.getAddress())
        )).thenReturn(BigInteger.valueOf(100));
        Map<String, BigInteger> balance = (Map<String , BigInteger>) score.call("balanceOf", nonOwner.getAddress(), TEST_TOKEN);
        // currently, still in locked state
        // will be updated once response comes from relayer via bmc
        assertEquals(BigInteger.valueOf(10), balance.get("locked"));
    }

    @Test
    @Order(11)
    public void transferBatch() {
        String[] coinNames = new String[]{"Token1","Token2","Token3","Token4", PARA};
        BigInteger val = BigInteger.valueOf(10);
        BigInteger[] values = new BigInteger[]{val,val, val, val, val};
        String destination = generateBTPAddress("harmony", ETH_ADDR);

        // none of them registered
        // user has not deposited any of them
        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "Not supported Token");

        BigInteger[] wrongValues = new BigInteger[]{val,val, val};

        call = () -> score.invoke(nonOwner, "transferBatch", coinNames, wrongValues, destination);
        expectErrorMessage(call, "Invalid arguments");

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

        score.invoke(token1,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(8), "0".getBytes());
        score.invoke(token2,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(token3,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());
        score.invoke(token4,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());

        // PARA wrapped transferFrom mock
        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);

        call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessageIn(call, "InSufficient Usable Balance");

        score.invoke(token1,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(10), "0".getBytes());

        // mock message to bmc
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), eq("harmony"), eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        verify(scoreSpy).TransferStart(eq(nonOwner.getAddress()), eq(destination), eq(BigInteger.ONE), any());

        contextMock.when(() -> Context.call(
                eq(BigInteger.class),
                any(),
                eq("balanceOf"),
                eq(nonOwner.getAddress())
        )).thenReturn(BigInteger.valueOf(100));

        Map<String, BigInteger> bal = (Map<String, BigInteger>) score.call("balanceOf", nonOwner.getAddress(), coinNames[0]);
        assertEquals(bal.get("usable"), BigInteger.valueOf(8));
        assertEquals(bal.get("locked"), BigInteger.valueOf(10));

        // PARA wrapped token balance mock
        contextMock.when(() -> Context.call(eq(BigInteger.class),eq(wrappedIRC2.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(50));

        contextMock.when(() -> Context.call(eq(BigInteger.class),eq(wrappedIRC2.getAddress()), eq("allowance"),
                eq(nonOwner.getAddress()), eq(score.getAddress()))).thenReturn(BigInteger.valueOf(40));

        bal = (Map<String, BigInteger>) score.call("balanceOf", nonOwner.getAddress(), coinNames[4]);
        // usable is min of balanceOf and allowed
        assertEquals(bal.get("usable"), BigInteger.valueOf(40));
        assertEquals(bal.get("locked"), BigInteger.valueOf(10));

        // message service number 1
        TransferTransaction txn = (TransferTransaction) score.call("getTransaction", BigInteger.ONE);
        assertEquals(txn.getFrom(), nonOwner.getAddress().toString());
        assertEquals(txn.getTo(), destination);
        assertEquals(txn.getAssets().length,6);

        // this goes to bmc -> relayer -> solidity -> relayer -> bmc -> here again
    }

    @Test
    @Order(12)
    public void handleBTPMessage1() {
        // solidity -> relayer -> bmc -> bts -> transfer/mint on icon side
        Executable call = () -> score.invoke(owner, "handleBTPMessage",
                "from",SERVICE,BigInteger.ONE, "ehehehe".getBytes());
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
        call = () -> score.invoke(bmcMock, "handleBTPMessage",
                "from",SERVICE,BigInteger.valueOf(3), _msg);
        expectErrorMessage(call, "Invalid Token");

        // register wrapped token
        registerWrapped();

        assertEquals(BigInteger.ZERO, score.call("getSn"));

        // native-coin
        contextMock.when(()->Context.transfer(any(Address.class),
                eq(BigInteger.valueOf(20)))).then(invocationOnMock -> null);
        // wrapped coin
        contextMock.when(() ->Context.call(any(), eq("mint"),
                eq(nonOwner.getAddress()), eq(BigInteger.valueOf(30)))).thenReturn(null);
        // irc2 transfer
        contextMock.when(() ->Context.call(eq(token1.getAddress()), eq("transfer"),
                eq(nonOwner.getAddress()), eq(BigInteger.valueOf(50)), any())).thenReturn(null);

        // mock bmc message
        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), eq("fromString"), eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.ONE, _msg);
        verify(scoreSpy).TransferReceived(eq("fromString"), eq(nonOwner.getAddress()), eq(BigInteger.ONE), any());
    }

    @Test
    @Order(13)
    public void handleBTPMessage2() {

        // request plus response
        String[] coinNames = new String[]{"Token1","Token2",PARA};
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
        score.invoke(token1,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(50), "0".getBytes());
        score.invoke(token2,"tokenFallback", nonOwner.getAddress(), val1, "0".getBytes());
        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);

        sendBTPMessageMock();
        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        verify(scoreSpy).TransferStart(any(), any(), any(), any());
        assertEquals(BigInteger.ONE, score.call("getSn"));

        // mock allowance and balanceOf for wrapped token

        contextMock.when(() -> Context.call(eq(BigInteger.class),eq(wrappedIRC2.getAddress()), eq("balanceOf"),
                eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));

        contextMock.when(() -> Context.call(eq(BigInteger.class),eq(wrappedIRC2.getAddress()), eq("allowance"),
                eq(nonOwner.getAddress()), eq(score.getAddress()))).thenReturn(BigInteger.valueOf(1000));

        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(token1.getAddress()),
                eq("balanceOf"), eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));
        contextMock.when(() -> Context.call(eq(BigInteger.class), eq(token2.getAddress()),
                eq("balanceOf"),eq(nonOwner.getAddress()))).thenReturn(BigInteger.valueOf(100));

        List<Map<String, BigInteger>> balances = (List<Map<String, BigInteger>>) score.call(
                "balanceOfBatch", nonOwner.getAddress(), new String[]{"Token1","Token2", PARA, ICON});

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

        contextMock.when(() ->Context.call(eq(wrappedIRC2.getAddress()), eq("burn"),
                eq(BigInteger.valueOf(10)))).thenReturn(null);

        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.ONE, _msg);
        verify(scoreSpy).TransferEnd(eq(nonOwner.getAddress()), eq(BigInteger.ONE), eq(TransferResponse.RC_OK), any());
        assertEquals(null, score.call("getTransaction", BigInteger.ONE));

        balances = (List<Map<String, BigInteger>>) score.call(
                "balanceOfBatch", nonOwner.getAddress(), new String[]{"Token1","Token2", PARA, ICON});

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
        assertEquals(BigInteger.ONE,fees.get(coinNames[0]));
        assertEquals(BigInteger.valueOf(11),fees.get(coinNames[1]));
        assertEquals(BigInteger.ZERO,fees.get(coinNames[2]));
        assertEquals(BigInteger.ZERO,fees.get(ICON));
    }

    @Test
    @Order(14)
    @DisplayName("unknown type")
    void handleBTPMessage3() {
        BTSMessage message = new BTSMessage();
        message.setServiceType(BTSMessage.UNKNOWN_TYPE);
        message.setData("a for apple".getBytes());

        score.invoke(bmcMock, "handleBTPMessage",
                "from",SERVICE,BigInteger.ONE, message.toBytes());
        verify(scoreSpy).UnknownResponse(any(), eq(BigInteger.ONE));
    }

    @Test
    @Order(15)
    @DisplayName("not handled cases")
    void handleBTPMessage4() {
        BTSMessage message = new BTSMessage();
        message.setServiceType(100);
        message.setData("a for apple".getBytes());

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), eq("fromString"), eq("bts"), eq(BigInteger.ONE), any());
        contextMock.when(sendMessage).thenReturn(null);

        score.invoke(bmcMock, "handleBTPMessage",
                "fromString",SERVICE,BigInteger.ONE, message.toBytes());
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

        Executable call;

        score.invoke(owner, "register",
                TOKEN1, TOKEN1, 18, BigInteger.valueOf(10),
                ICX, token1.getAddress());

        assertEquals(UINT_CAP, score.call("getTokenLimit", TOKEN1));

        tokenLimitBTPMessage();
        score.invoke(owner, "setTokenLimit", new String[]{TOKEN1}, new BigInteger[]{BigInteger.valueOf(200)});

        assertEquals(BigInteger.valueOf(200), score.call("getTokenLimit", TOKEN1));
//
        call = () -> score.invoke(owner, "setTokenLimit", new String[]{TOKEN1},
                new BigInteger[]{TWO_HUNDRED_ICX.negate()});
        expectErrorMessage(call, "Invalid value");
        // can set for tokens not registered as well
        score.invoke(owner, "setTokenLimit", new String[]{"TokenBSH"},
                new BigInteger[]{TWO_HUNDRED_ICX});
        assertEquals(TWO_HUNDRED_ICX, score.call("getTokenLimit","TokenBSH"));

//        expectErrorMessage(call, "Not registered");

        score.invoke(owner, "setTokenLimit",  new String[]{TOKEN1}, new BigInteger[]{TWO_HUNDRED_ICX});

        assertEquals(TWO_HUNDRED_ICX, score.call("getTokenLimit", TOKEN1));

//         send from Harmony to ICON
//         reciever adddress is blacklisted in icon

        Verification sendMessage = () -> Context.call(eq(bmcMock.getAddress()),
                eq("sendMessage"), any(), any(), any(), any());

        Verification returnLinks = () -> Context.call(eq(String[].class),
                eq(bmcMock.getAddress()), eq("getLinks"), any());

        // these 3 BTP addresses are currently supported by ICON Bridge
        BTPAddress btpAddress1 = new BTPAddress("icon","cx0E636c8cF214a9d702C5E4a6D8c020be217766D3");
        BTPAddress btpAddress2 = new BTPAddress("network","0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
        BTPAddress btpAddress3 = new BTPAddress("harmony","0x0E636c8cF214a9d702C5E4a6D8c020be217766D3");
        String[] links = new String[] {btpAddress1.toString(), btpAddress2.toString(), btpAddress3.toString()};

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

        contextMock.when(() -> Context.call(eq(token1.getAddress()), eq("transfer"),
                eq(user1), eq(FIFTY_ICX), any())).thenReturn(null);

        score.invoke(owner, "addRestriction");

        call = () -> score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.ONE, _msg);
        expectErrorMessage(call, "_to user is Blacklisted");

        removeBlackListedUser("icon", new String[]{user1.toString()}, BigInteger.valueOf(4) );

        // check this
        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.valueOf(5), _msg);

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
        contextMock.when(() -> Context.call(eq(token1.getAddress()), eq("transfer"),
                eq(user1), eq(FIFTY_ICX.add(TWO_HUNDRED_ICX)), any())).thenReturn(null);

        call = () -> score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.valueOf(6), _msg2);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        // limit not set
        // any amount should be able to be transferred off chain
        String TOKEN2 = "Token 2";
        BigInteger MILLION = BigInteger.valueOf(10_000_000L).multiply(ICX);
        score.invoke(owner, "register",
                TOKEN2, TOKEN2, 18, BigInteger.valueOf(10),
                ICX, token1.getAddress());

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

        contextMock.when(() -> Context.call(eq(token1.getAddress()), eq("transfer"),
                eq(user1), eq(MILLION), any())).thenReturn(null);


        // change token limit
        score.invoke(owner, "setTokenLimit",  new String[]{TOKEN2}, new BigInteger[]{TWO_HUNDRED_ICX});
        call = () -> score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.valueOf(6), _msg3);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        // disable restrictions
        score.invoke(owner, "disableRestrictions");
        score.invoke(bmcMock, "handleBTPMessage",
                "fromString", SERVICE, BigInteger.valueOf(7), _msg3);
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
        expectErrorMessage(call,"already exists owner" );

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
        Executable call = () -> score.invoke(nonOwner, "setFeeRatio",
                PARA, BigInteger.valueOf(10), ICX);
        expectErrorMessage(call, "require owner access" );

        call = () -> score.invoke(owner, "setFeeRatio",
                PARA, BigInteger.valueOf(10).negate(), ICX.negate());
        expectErrorMessageIn(call, "The feeNumerator should be less "
                + "than FEE_DENOMINATOR and feeNumerator should be greater than 1");

        // Scenario 11: Fee_numerator is set higher than fee_denominator
        call = () -> score.invoke(owner, "setFeeRatio",
                PARA, ICX, ICX.negate());
        expectErrorMessageIn(call, "The feeNumerator should be less "
                + "than FEE_DENOMINATOR and feeNumerator should be greater than 1");

        call = () -> score.invoke(owner, "setFeeRatio",
                PARA, BigInteger.valueOf(10), ICX.negate());
        expectErrorMessageIn(call, "Fixed fee cannot be less than zero");

        call = () -> score.invoke(owner, "setFeeRatio",
                "LAMB", BigInteger.valueOf(100), ICX);
        expectErrorMessage(call, "Not supported Coin");

        // Scenario 9: Contract’s owner updates a new fee ratio
        // Scenario 12: Contract’s owner updates fixed fee
        score.invoke(owner, "setFeeRatio",
                PARA, BigInteger.valueOf(100), ICX);
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
    public void transferBatchNativecoinTokenLimit() {
        tokenLimitBTPMessage();
        score.invoke(owner, "setTokenLimit",  new String[]{ICON}, new BigInteger[]{BigInteger.valueOf(1000)});
        assertEquals(BigInteger.valueOf(1000), score.call("getTokenLimit", ICON));

        String destination = generateBTPAddress("0x1.bsc", ETH_ADDR);
        Account token1 = Account.newScoreAccount(50);
        Account token2 = Account.newScoreAccount(51);

        register("Token1", token1.getAddress());
        register("Token2", token2.getAddress());

        score.invoke(token1,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(50), "0".getBytes());
        score.invoke(token2,"tokenFallback", nonOwner.getAddress(), BigInteger.valueOf(50), "0".getBytes());
        Verification transferFromMock = () -> Context.call(eq(Boolean.class), any(), eq("transferFrom"),
                eq(nonOwner.getAddress()), eq(score.getAddress()), eq(BigInteger.valueOf(10)), any());
        contextMock.when(transferFromMock).thenReturn(true);

        // transferBatch 2 tokens, try to transfer ICX above 1000 should fail
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(1001));
        String[] coinNames = new String[]{"Token1", "Token2"};
        BigInteger[] values = new BigInteger[]{ BigInteger.valueOf(50),  BigInteger.valueOf(50)};

        sendBTPMessageMock();
        Executable call = () -> score.invoke(nonOwner, "transferBatch", coinNames, values, destination);
        expectErrorMessage(call, "Transfer amount exceeds the transaction limit");

        // 999 transfer should be successful
        contextMock.when(sendICX()).thenReturn(BigInteger.valueOf(999));
        score.invoke(nonOwner, "transferBatch", coinNames, values, destination);

    }
}