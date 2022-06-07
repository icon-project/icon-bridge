/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package foundation.icon.btp.bsh.test;

import com.iconloop.testsvc.Account;
import com.iconloop.testsvc.Score;
import com.iconloop.testsvc.ServiceManager;
import com.iconloop.testsvc.TestBase;
import foundation.icon.btp.bsh.BMCMock;
import foundation.icon.btp.bsh.ServiceHandler;
import foundation.icon.btp.bsh.types.Asset;
import foundation.icon.btp.bsh.types.BTPAddress;
import foundation.icon.btp.bsh.types.Balance;
import foundation.icon.btp.bsh.types.TransferAsset;
import foundation.icon.btp.irc2.IRC2Basic;
import foundation.icon.btp.restrictions.Restrictions;
import org.bouncycastle.util.encoders.Hex;
import org.junit.jupiter.api.*;
import score.ByteArrayObjectWriter;
import score.Context;
import scorex.util.ArrayList;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import static java.math.BigInteger.TEN;
import static org.junit.jupiter.api.Assertions.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class ServiceHandlerTest extends TestBase {
    final static String RLPn = "RLPn";
    final static String _svc = "TokenBSH";
    final static String _net = "icon";
    final static String tokenName = "ETH";
    final static String symbol = "ETH";
    final static int decimals = 18;
    final static BigInteger fees = BigInteger.valueOf(100);
    final static BigInteger transferAmount = new BigInteger("10000000000000000000");
    private static final BigInteger initialSupply = BigInteger.valueOf(2000);
    private static final BigInteger totalSupply = initialSupply.multiply(TEN.pow(decimals));
    private static final int REQUEST_TOKEN_TRANSFER = 0;
    private static final int REQUEST_TOKEN_REGISTER = 1;
    private static final int RESPONSE_HANDLE_SERVICE = 2;
    private static final int RESPONSE_UNKNOWN_ = 3;
    private static final ServiceManager sm = getServiceManager();
    private static Account[] owners;
    private static Score bsh;
    private static Score token;
    private static Score bmc;
    private static Score irc2Basic;
    private ServiceHandler bshSyp;
    private static Score restrictons;

    @BeforeAll
    public static void setup() throws Exception {
        // setup accounts and deploy
        owners = new Account[3];
        for (int i = 0; i < owners.length; i++) {
            owners[i] = sm.createAccount(100);
        }
        String initialOwners = Arrays.stream(owners)
                .map(a -> a.getAddress().toString())
                .collect(Collectors.joining(","));
        bmc = sm.deploy(owners[0], BMCMock.class);
        bsh = sm.deploy(owners[0], ServiceHandler.class, bmc.getAddress().toString());
        token = sm.deploy(owners[0], IRC2Basic.class, tokenName, symbol, decimals, initialSupply);
        restrictons = sm.deploy(owners[0], Restrictions.class);

        BigInteger balance = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        assertEquals(totalSupply, balance);
        // setup spy
       /* bshSyp = (ServiceHandler) spy(bshScore.getInstance());
        bshScore.setInstance(bshSyp);*/
    }
/*
    @Order(1)
    @Test
    public void handleBTPMessageFromHexBytesTest() {
        String _from = "btp://0x97.bsc/0x7D66b33f2b2d2Cd565e5024E651B6c6bE491c493";
        //The message with Bigint decoding issue
        String _msg="f87c00b879f877b3307864303833393232613432373765376264626164633939363335316531653338353637336164373936000000000000000000aa687834363963396162383031663861303766663561343561636164383535643539333933343630343938d7d683455448888963dd8c2c5e000088016345785d8a0000";
        bmc.invoke(owners[0], "addService", _svc, bsh.getAddress());
        bsh.invoke(owners[0], "register", tokenName, symbol, BigInteger.valueOf(decimals), fees, token.getAddress());
        token.invoke(owners[0],"transfer",bsh.getAddress(),new BigInteger("100000000000000000000"),"transfer to Receiver".getBytes());
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ONE, Hex.decode(_msg));
    }
*/
/*    @Order(1)
    @Test
    public void handleBTPMessageFromHexBytesTest() {
        String _from = "btp://0x97.bsc/0x7D66b33f2b2d2Cd565e5024E651B6c6bE491c493";
        String _msg="f87b00b878f876b3307865626362643461393334613638353130653231626132356232613832373133383234386136336535000000000000000000aa687832376131356131633361396263303933343038323939643732623965633233383036356337346137d6d583455448880dbd2fc137a30000872386f26fc10000";
        restrictons.invoke(owners[0], "addBlacklistedUser", "0xebcbd4a934a68510e21ba25b2a827138248a63e5");
        boolean isrest = (boolean) restrictons.call("isUserBlackListed","0xebcbd4a934a68510e21ba25b2a827138248a63e5");
        bmc.invoke(owners[0], "addService", _svc, bsh.getAddress());
        bsh.invoke(owners[0], "register", tokenName, symbol, BigInteger.valueOf(decimals), fees, token.getAddress());
        bsh.invoke(owners[0], "addRestrictor", restrictons.getAddress());
        token.invoke(owners[0],"transfer",bsh.getAddress(),new BigInteger("100000000000000000000"),"transfer to Receiver".getBytes());
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ONE, Hex.decode(_msg));
    }*/

    /**
     * Scenario 1: Receiving address is an invalid address - fail
     */
    @Order(1)
    @Test
    public void scenario1() {
        String _from = "0x12345678";
        String _to = "0x1234567890123456789";
        bmc.invoke(owners[0], "addService", _svc, bsh.getAddress());
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ZERO, handleBTPRequestBtpMsg(_from, _to))
        );
        assertTrue(thrown.getMessage().contains("Invalid Address format"));
    }

    /**
     * Scenario 2:  User creates a transfer, but a token_name has not yet registered - fail
     */
    @Order(2)
    @Test
    public void scenario2() {
        String _to = owners[0].getAddress().toString();
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "transfer", tokenName, transferAmount, _to)
        );
        assertTrue(thrown.getMessage().contains("Token not registered"));

    }

    /**
     * Scenario 3:  Register Token without permission - fail
     */
    @Order(3)
    @Test
    public void scenario3() {
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[1], "register", tokenName, symbol, BigInteger.valueOf(decimals), fees, token.getAddress()));
        assertTrue(thrown.getMessage().contains("No Permission"));
    }

    /**
     * Scenario 4:  Register Token with permission - Success
     */
    @Order(4)
    @Test
    public void scenario4() {
        String _to = "btp://bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        bsh.invoke(owners[0], "register", tokenName, symbol, BigInteger.valueOf(decimals), fees, token.getAddress());
    }


    /**
     * Scenario 5:  Register Token - Token already exists - Failure
     */
    @Order(5)
    @Test
    public void scenario5() {
        String _to = "btp://bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "register", tokenName, symbol, BigInteger.valueOf(decimals), fees, token.getAddress()));
        assertTrue(thrown.getMessage().contains("Token with same name exists already"));
    }


    /**
     * Scenario 6:  User does not have enough balance - fail
     */
    @Order(6)
    @Test
    public void scenario6() {
        String _to = "btp://bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "transfer", tokenName, transferAmount, _to)
        );
        assertTrue(thrown.getMessage().contains("Overdrawn"));

    }


    /**
     * Scenario #:  Invalid amount specified(transfer amount = 0) - fail
     */
    @Order(7)
    @Test
    public void scenario7() {
        String _to = "btp://bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "transfer", tokenName, BigInteger.ZERO, _to)
        );
        assertTrue(thrown.getMessage().contains("Invalid amount specified"));

    }

    /**
     * Secnario#: Transfer IRC2 tokens from Token contract to BSH via fallback - success
     */
    @Order(8)
    @Test
    public void scenario8() {
        token.invoke(owners[0], "transfer", bsh.getAddress(), transferAmount, new byte[0]);
        //TODO: assert the balance after
    }


    /**
     * Scenario #:  User transfers to an invalid BTP address - fail
     */
    @Order(9)
    @Test
    public void scenario9a() {
        String _to = "btp://0x1.bsc:0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        //bsh.invoke(owners[0],"transfer", tokenName, transferAmount,_to);
        assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "transfer", tokenName, transferAmount, _to)
        );
    }

    /**
     * Scenario #:  User transfers to an invalid BTP address - fail
     */

    @Order(9)
    @Test
    public void scenario9b() {
        String _to = "btp://0x1.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        bsh.invoke(owners[0],"addRestrictor",restrictons.getAddress());
        restrictons.invoke(owners[0], "addBlacklistedUser", BTPAddress.fromString(_to).getContract());
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "transfer", tokenName, transferAmount, _to)
        );
        assertTrue(thrown.getMessage().contains("_to user is Blacklisted"));
        restrictons.invoke(owners[0], "removeBlacklistedUser", BTPAddress.fromString(_to).getContract());
    }


    /**
     * Scenario #:   All requirements are qualified and BSH initiates Transfer start - Success
     */
    @Order(10)
    @Test
    public void scenario10() {
        String _to = "btp://0x1.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        Balance balanceBefore = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        bsh.invoke(owners[0], "transfer", tokenName, transferAmount, _to);
        Balance balanceAfter = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        assertEquals(balanceBefore.getLocked().add(transferAmount), balanceAfter.getLocked());
    }


    /**
     * Scenario #:  All requirements are qualified and BSH receives a failed message - Success
     */
    @Order(11)
    @Test
    public void scenario11() {
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        Balance balanceBefore = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ONE, handleBTPResponseBtpMsg(1, "Transfer Failed"));
        Balance balanceAfter = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        assertEquals(balanceAfter.getRefundable(), balanceBefore.getRefundable().add(transferAmount));
    }

    /**
     * Scenario #: withdraw refund amount - invalid - failure
     */
    @Order(12)
    @Test
    public void scenario11a() {

        Balance balanceBefore = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "withdraw", tokenName, transferAmount.add(BigInteger.ONE))
        );
        assertTrue(thrown.getMessage().contains("imbalance"));

        thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "withdraw", tokenName, BigInteger.ZERO)
        );
        assertTrue(thrown.getMessage().contains("_value must be positive"));

        thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "withdraw", "BNB", transferAmount)
        );
        assertTrue(thrown.getMessage().contains("Token not registered"));
    }

    /**
     * Scenario #: withdraw refund amount - valid - success
     */
    @Order(12)
    @Test
    public void scenario11b() {
        BigInteger balanceBefore = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        bsh.invoke(owners[0], "withdraw", tokenName, transferAmount);
        BigInteger balanceAfter = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        assertEquals(balanceAfter, balanceBefore.add(transferAmount));
    }

    /**
     * Scenario #:  AAll requirements are qualified and BSH receives a successful message - Success
     */
    @Order(12)
    @Test
    public void scenario12() {
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        String _to = "btp://0x1.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        token.invoke(owners[0], "transfer", bsh.getAddress(), transferAmount, new byte[0]);
        Balance balanceBefore = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        bsh.invoke(owners[0], "transfer", tokenName, transferAmount, _to);
        Balance balanceAfter = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        assertEquals(balanceBefore.getLocked().add(transferAmount), balanceAfter.getLocked());
        balanceBefore = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.TWO, handleBTPResponseBtpMsg(0, "Transfer Success"));
        balanceAfter = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        assertEquals(balanceAfter.getLocked().add(transferAmount), balanceBefore.getLocked());
    }

    /**
     * Scenario #:  handleBTPMessage request mint balance for the user, blacklisted user: to- failed
     */
    @Test
    @Order(13)
    public void scenario13a() {
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        restrictons.invoke(owners[0], "addBlacklistedUser", owners[1].getAddress().toString());
        BigInteger balanceBefore = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ZERO, handleBTPRequestBtpMsg(_from, owners[1].getAddress().toString()))
        );
        assertTrue(thrown.getMessage().contains("_to user is Blacklisted"));
    }

    /**
     * Scenario #:  handleBTPMessage request mint balance for the user, blacklisted user: from- failed
     */
    @Test
    @Order(13)
    public void scenario13b() {
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        restrictons.invoke(owners[0], "addBlacklistedUser", BTPAddress.fromString(_from).getContract());
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ZERO, handleBTPRequestBtpMsg( BTPAddress.fromString(_from).getContract(), owners[0].getAddress().toString()))
        );
        assertTrue(thrown.getMessage().contains("_from user is Blacklisted"));
        restrictons.invoke(owners[0], "removeBlacklistedUser", BTPAddress.fromString(_from).getContract());
    }

    /**
     * Scenario #:  handleBTPMessage request mint balance for the user, Exceeds Transaction amount - failed
     */
    @Test
    @Order(13)
    public void scenario13c() {
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        restrictons.invoke(owners[0], "registerTokenLimit", tokenName, tokenName, token.getAddress(), transferAmount.subtract(BigInteger.valueOf(1)));
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ZERO, handleBTPRequestBtpMsg( BTPAddress.fromString(_from).getContract(), owners[0].getAddress().toString()))
        );
        assertTrue(thrown.getMessage().contains("Transfer amount exceeds the transaction limit"));
        //set restriction back to transfer amount
        //restrictons.invoke(owners[0], "registerTokenLimit", tokenName, tokenName, token.getAddress(), transferAmount);
        bsh.invoke(owners[0], "disableRestrictions");
    }

    /**
     * Scenario #:  All requirements are qualified handleBTPMessage mints balance for the user- Success
     */
    @Test
    @Order(13)
    public void scenario13d() {
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        //Balance balanceBefore = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        BigInteger balanceBefore = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.ZERO, handleBTPRequestBtpMsg(_from, owners[0].getAddress().toString()));
        //Balance balanceAfter = (Balance) bsh.call("getBalance", owners[0].getAddress(), tokenName);
        BigInteger balanceAfter = (BigInteger) token.call("balanceOf", owners[0].getAddress());
        assertEquals(balanceBefore.add(transferAmount), balanceAfter);
    }

    /**
     * scenario #: Add remove owner - min owner - should fail
     */
    @Test
    @Order(14)
    public void scenario14() {
        assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[0], "removeOwner", owners[0].getAddress())
        );
    }

    /**
     * scenario #: Add add owner - without permission - should fail
     */
    @Test
    @Order(15)
    public void scenario15() {
        AssertionError thrown = assertThrows(AssertionError.class, () ->
                bsh.invoke(owners[1], "addOwner", owners[1].getAddress())
        );
        assertTrue(thrown.getMessage().contains("No Permission"));
    }

    /**
     * scenario #: Add add owner - with permission - success
     */
    @Test
    @Order(16)
    public void scenario16() {
        bsh.invoke(owners[0], "addOwner", owners[1].getAddress());
    }

    /**
     * Scenario #:  Register Token with permission new owner - Success
     */
    @Order(17)
    @Test
    public void scenario17() {
        String _to = "btp://bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        bsh.invoke(owners[1], "register", "BNB", "BNB", BigInteger.valueOf(decimals), fees, token.getAddress());
    }

    /**
     * Scenario #:  Handle Accumulated Fees - Failure
     */
    @Order(18)
    @Test
    public void scenario18() {
        String _fa = "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b";
        List<Map<String, BigInteger>> _assets = (List<Map<String, BigInteger>>) bsh.call("getAccumulatedFees");
        BigInteger _fees = transferAmount.multiply(fees).divide(BigInteger.valueOf(10000));
        assertEquals(_assets.get(0).get(tokenName), _fees);
        bmc.invoke(owners[1], "handleFeeGathering", _fa, _svc);
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.valueOf(3), handleBTPResponseBtpMsg(1, "Trasnfer Failed"));
        _assets = (List<Map<String, BigInteger>>) bsh.call("getAccumulatedFees");
        // should still have the fees in accumulator after failure handleresponse
        assertEquals(_assets.get(0).get(tokenName), _fees);
    }

    /**
     * Scenario #:  Handle Accumulated Fees - Success
     */
    @Order(19)
    @Test
    public void scenario19() {
        String _fa = "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b";
        List<Map<String, BigInteger>> _assets = (List<Map<String, BigInteger>>) bsh.call("getAccumulatedFees");
        BigInteger _fees = transferAmount.multiply(fees).divide(BigInteger.valueOf(10000));
        assertEquals(_assets.get(0).get(tokenName), _fees);
        bmc.invoke(owners[1], "handleFeeGathering", _fa, _svc);
        String _from = "btp://0x97.bsc/0xa36a32c114ee13090e35cb086459a690f5c1f8e8";
        bmc.invoke(owners[0], "handleBTPMessage", _from, _svc, BigInteger.valueOf(4), handleBTPResponseBtpMsg(0, "Trasnfer Success"));
        _assets = (List<Map<String, BigInteger>>) bsh.call("getAccumulatedFees");
        // Should not have any fees left in FeeAccumulator db
        assertEquals(_assets.size(), 0);
    }


    public byte[] handleBTPRequestBtpMsg(String from, String to) {
        ByteArrayObjectWriter writer = Context.newByteArrayObjectWriter(RLPn);
        writer.beginList(2);
        writer.write(REQUEST_TOKEN_TRANSFER);//ActionType
        List<Asset> assets = new ArrayList<Asset>();
        assets.add(new Asset(tokenName, transferAmount, BigInteger.ZERO));
        ByteArrayObjectWriter writerTa = Context.newByteArrayObjectWriter(RLPn);
        TransferAsset _ta = new TransferAsset(from, to, assets);
        TransferAsset.writeObject(writerTa, _ta);
        writer.write(writerTa.toByteArray());
        writer.end();
        return writer.toByteArray();
    }

    public byte[] handleBTPResponseBtpMsg(int code, String msg) {
        ByteArrayObjectWriter writer = Context.newByteArrayObjectWriter(RLPn);
        writer.beginList(2);
        writer.write(RESPONSE_HANDLE_SERVICE);//ActionType
        ByteArrayObjectWriter writerRespMsg = Context.newByteArrayObjectWriter(RLPn);
        writerRespMsg.beginList(2);
        writerRespMsg.write(code);//Code
        writerRespMsg.write(msg);//Msg
        writerRespMsg.end();
        writer.write(writerRespMsg.toByteArray());
        writer.end();
        return writer.toByteArray();
    }

}
