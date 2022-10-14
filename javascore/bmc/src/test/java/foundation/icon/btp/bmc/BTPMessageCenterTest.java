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

package foundation.icon.btp.bmc;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;

import com.iconloop.score.test.Account;
import foundation.icon.btp.lib.BMCStatus;
import foundation.icon.btp.lib.BTPAddress;
import java.math.BigInteger;
import java.util.Map;
import org.bouncycastle.util.encoders.Hex;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.MethodOrderer.OrderAnnotation;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.TestMethodOrder;
import org.junit.jupiter.api.function.Executable;
import org.junit.jupiter.api.Test;
import org.mockito.MockedStatic.Verification;
import score.Address;
import score.Context;

@TestMethodOrder(OrderAnnotation.class)
class BTPMessageCenterTest extends AbstractBTPMessageCenterTest {


    @Test
    void name() {
        assertEquals("BTP Message Center", score.call("name"));
    }

    @Test
    void getBTPAddress() {
        assertEquals(getBTPAddress(BMC_SCORE), score.call("getBtpAddress"));
    }

    @Test
    void addRemoveService() {
        String ADD_SERVICE = "addService";
        Executable call = () -> score.invoke(nonOwner, ADD_SERVICE, "BTS", score.getAddress());
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);
        call = () -> score.invoke(owner, ADD_SERVICE, "!@!@!", score.getAddress());
        expectErrorMessage(call, "invalid service name");
        assertThrows(AssertionError.class, () -> score.invoke(owner, ADD_SERVICE, "bmc", score.getAddress()));

        score.invoke(owner, ADD_SERVICE, BTS, BTSScore.getAddress());
        Map<String, Address> map = Map.of(BTS, BTSScore.getAddress());
        Map<String, Address> response = (Map<String, Address>) score.call("getServices");
        assertEquals(map, response);

        // add a service candidate
        score.invoke(nonOwner, "addServiceCandidate", "BKT", score.getAddress());
        ServiceCandidate[] candidates = (ServiceCandidate[]) score.call("getServiceCandidates");
        assertEquals(1, candidates.length);
        assertEquals("BKT", candidates[0].getSvc());
        assertEquals(score.getAddress(), candidates[0].getAddress());

        // service candidate should be removed if it is registered as service
        score.invoke(owner, ADD_SERVICE, "BKT", score.getAddress());
        candidates = (ServiceCandidate[]) score.call("getServiceCandidates");
        assertEquals(0, candidates.length);

        String REMOVE_SERVICE = "removeService";
        call = () -> score.invoke(nonOwner, REMOVE_SERVICE, BTS);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        assertThrows(AssertionError.class, () -> score.invoke(owner, REMOVE_SERVICE, "BNS"));

        score.invoke(owner, REMOVE_SERVICE, BTS);
        score.invoke(owner, REMOVE_SERVICE, "BKT");
        map = Map.of();
        response = (Map<String, Address>) score.call("getServices");
        assertEquals(map, response);
    }

    @Test
    void addRemoveLinks() {
        String ADD_LINK = "addLink";
        String link = getBTPAddress("0xa1442c90120A891c3de9793caC70968Cab113234");
        Executable call = () -> score.invoke(nonOwner, ADD_LINK, link);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        score.invoke(owner, ADD_LINK, link);
        // generates Message event <verify BTP Message format>
        verify(scoreSpy).Message(eq(link), eq(BigInteger.ONE), any());
        String[] expected = new String[]{link};
        String[] actual = (String[]) score.call("getLinks");
        assertEquals(expected[0], actual[0]);

        // exception on adding to array again
        assertThrows(AssertionError.class, () -> score.invoke(owner, ADD_LINK, link));

        // check status
        BMCStatus status = (BMCStatus) score.call("getStatus", link);
        assertEquals(BigInteger.ONE, status.getTx_seq());
        assertEquals(BigInteger.ZERO, status.getRx_seq());

        call = () -> score.invoke(nonOwner, "removeLink", link);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        score.invoke(owner, "removeLink", link);
        actual = (String[]) score.call("getLinks");
        assertEquals(0, actual.length);

        assertThrows(AssertionError.class, () -> score.invoke(owner, "removeLink", link));

        // add the same link again
        score.invoke(owner, ADD_LINK, link);
        status = (BMCStatus) score.call("getStatus", link);
        assertEquals(BigInteger.ONE, status.getTx_seq());
        assertEquals(BigInteger.ZERO, status.getRx_seq());

        // set link Rx height
        long height = 1000L;
        call = () -> score.invoke(nonOwner, "setLinkRxHeight", link, height);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        score.invoke(owner, "setLinkRxHeight", link, height);
        // assert Rx Height

        int blockInterval = 100;
        int maxAgg = 1000;
        call = () -> score.invoke(nonOwner, "setLinkRotateTerm", link, blockInterval, maxAgg);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, "setLinkRotateTerm", link, -1, maxAgg);
        expectErrorMessage(call, "invalid param");

        call = () -> score.invoke(owner, "setLinkRotateTerm", link, blockInterval, 0);
        expectErrorMessage(call, "invalid param");

        score.invoke(owner, "setLinkRotateTerm", link, blockInterval, maxAgg);
        // assert link rotate term

        int linkDelayLimit = 1000;
        call = () -> score.invoke(nonOwner, "setLinkDelayLimit", link, linkDelayLimit);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, "setLinkDelayLimit", link, -1);
        expectErrorMessage(call, "invalid param");

        score.invoke(owner, "setLinkDelayLimit", link, linkDelayLimit);
        // assert link delay limit

        int linkSackTerm = 1000;
        call = () -> score.invoke(nonOwner, "setLinkSackTerm", link, linkSackTerm);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, "setLinkSackTerm", link, -1);
        expectErrorMessage(call, "invalid param");

        score.invoke(owner, "setLinkSackTerm", link, linkSackTerm);
        // assert link delay limit

    }

    @Test
    void addRemoveRoutes() {
        String ADD_ROUTE = "addRoute";
        String NETWORK = "0x38.bsc";
        String link = getBTPAddress("0xa1442c90120A891c3de9793caC70968Cab113234");
        // non owner tries to add route
        assertThrows(AssertionError.class, () -> score.invoke(nonOwner, ADD_ROUTE, NETWORK, link));

        // add route for unregistered link
        assertThrows(AssertionError.class, () -> score.invoke(owner, ADD_ROUTE, NETWORK, link));

        addLink(link);
        score.invoke(owner, ADD_ROUTE, NETWORK, link);

        Map<String, String> expected = Map.of(NETWORK, link);
        Map<String, String> actual = (Map<String, String>) score.call("getRoutes");
        assertEquals(expected, actual);

        // route already exists
        assertThrows(AssertionError.class, () -> score.invoke(owner, ADD_ROUTE, NETWORK, link));

        // remove route
        Executable call = () -> score.invoke(nonOwner, "removeRoute", NETWORK);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);
        score.invoke(owner, "removeRoute", NETWORK);
        actual = (Map<String, String>) score.call("getRoutes");
        assertEquals(Map.of(), actual);
    }

    @Test
    void serviceCandidates() {
        String ADD_SERVICE_CANDIDATE = "addServiceCandidate";
        Address addr1 = sm.createAccount().getAddress();
        Address addr2 = sm.createAccount().getAddress();

        // anyone can call this, no owner check required
        score.invoke(owner, ADD_SERVICE_CANDIDATE, "bmc", addr1);
        score.invoke(nonOwner, ADD_SERVICE_CANDIDATE, "bmc", addr2);

        ServiceCandidate[] candidates = (ServiceCandidate[]) score.call("getServiceCandidates");

        ServiceCandidate candidate1 = new ServiceCandidate();
        candidate1.setAddress(addr1);
        candidate1.setOwner(owner.getAddress());
        candidate1.setSvc("bmc");

        assertEquals(candidate1, candidates[0]);

        ServiceCandidate candidate2 = new ServiceCandidate();
        candidate2.setAddress(addr2);
        candidate2.setOwner(nonOwner.getAddress());
        candidate2.setSvc("bmc");
        assertEquals(candidate2, candidates[1]);

        assertThrows(AssertionError.class, () -> score.invoke(owner, ADD_SERVICE_CANDIDATE, "bmc", addr2));

        Executable call = () -> score.invoke(nonOwner, "removeServiceCandidate", "bmc", addr1);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        score.invoke(owner, "removeServiceCandidate", "bmc", addr1);
        candidates = (ServiceCandidate[]) score.call("getServiceCandidates");
        assertEquals(candidate2, candidates[0]);
        assertEquals(1, candidates.length);

        assertThrows(AssertionError.class, () -> score.invoke(owner, "removeServiceCandidate", "bts", addr1));
    }

    @Test
    void addRemoveRelay() {
        String ADD_RELAY = "addRelay";
        Address addr = Address.fromString("hxa1442c90120A891c3de9793caC70968Cab113234");
        String link = getBTPAddress(addr);

        Executable call = () -> score.invoke(nonOwner, ADD_RELAY, link, addr);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, ADD_RELAY, link, addr);
        expectErrorMessage(call, "NotExistsLink");

        addLink(link);
        score.invoke(owner, ADD_RELAY, link, addr);

        Address addr2 = sm.createAccount().getAddress();
        score.invoke(owner, ADD_RELAY, link, addr2);

        Address[] relays = (Address[]) score.call("getRelays", link);
        assertEquals(relays.length, 2);
        assertEquals(addr, relays[0]);
        assertEquals(addr2, relays[1]);

        String REMOVE_RELAY = "removeRelay";
        call = () -> score.invoke(nonOwner, REMOVE_RELAY, link, score.getAddress());
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        assertThrows(AssertionError.class, () -> score.invoke(owner, REMOVE_RELAY, link, score.getAddress()));

        score.invoke(owner, REMOVE_RELAY, link, addr);
        relays = (Address[]) score.call("getRelays", link);
        assertEquals(relays.length, 1);
        assertEquals(addr2, relays[0]);
    }

    @Test
    void feeDetails() {
        // fee gathering term
        Executable call = () -> score.invoke(nonOwner, "setFeeGatheringTerm", 1L);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, "setFeeGatheringTerm", -1L);
        expectErrorMessage(call, "invalid param");

        score.invoke(owner, "setFeeGatheringTerm", 42300L);
        assertEquals(42300L, score.call("getFeeGatheringTerm"));

        // fee aggregator
        Address aggregator = sm.createAccount().getAddress();
        call = () -> score.invoke(nonOwner, "setFeeAggregator", aggregator);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        score.invoke(owner, "setFeeAggregator", aggregator);
        assertEquals(aggregator, score.call("getFeeAggregator"));
    }

    @Test
    void relayerDetails() {
        String SET_RELAYER_MIN_BOND = "setRelayerMinBond";
        BigInteger bond = BigInteger.valueOf(100);

        Executable call = () -> score.invoke(nonOwner, SET_RELAYER_MIN_BOND, bond);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, SET_RELAYER_MIN_BOND, BigInteger.ONE.negate());
        expectErrorMessage(call, "minBond must be positive");

        score.invoke(owner, SET_RELAYER_MIN_BOND, bond);
        assertEquals(bond, score.call("getRelayerMinBond"));

        String SET_RELAYER_TERM = "setRelayerTerm";
        long term = 100L;

        call = () -> score.invoke(nonOwner, SET_RELAYER_TERM, term);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, SET_RELAYER_TERM, 0L);
        expectErrorMessage(call, "term must be positive");

        score.invoke(owner, SET_RELAYER_TERM, term);
        assertEquals(term, score.call("getRelayerTerm"));

        String SET_RELAYER_REWARD_RANK = "setRelayerRewardRank";
        int rank = 10;

        call = () -> score.invoke(nonOwner, SET_RELAYER_REWARD_RANK, rank);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        call = () -> score.invoke(owner, SET_RELAYER_REWARD_RANK, 0);
        expectErrorMessage(call, "rewardRank must be positive");

        score.invoke(owner, SET_RELAYER_REWARD_RANK, rank);
        assertEquals(rank, score.call("getRelayerRewardRank"));

        String SET_NEXT_REWARD_DISTRIBUTION = "setNextRewardDistribution";
        long height = 1000000000L;

        call = () -> score.invoke(nonOwner, SET_NEXT_REWARD_DISTRIBUTION, height);
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        score.invoke(owner, SET_NEXT_REWARD_DISTRIBUTION, height);

        RelayersProperties properties = (RelayersProperties) score.call("getRelayersProperties");
        assertEquals(height, properties.getNextRewardDistribution());
        assertEquals(rank, properties.getRelayerRewardRank());
        assertEquals(term, properties.getRelayerTerm());
        assertEquals(bond, properties.getRelayerMinBond());
    }

    @Test
    @Order(1)
    public void addRelayer() {
        Account alice = sm.createAccount(10);
        contextMock.when(mockICXSent()).thenReturn(BigInteger.valueOf(0));
        Executable call = () -> score.invoke(alice, "registerRelayer", "Hey, I want to be a relayer");
        expectErrorMessage(call, "require bond at least 1 icx");

        contextMock.when(mockICXSent()).thenReturn(BigInteger.valueOf(50));
        score.invoke(alice, "registerRelayer", "Hey, I want to be a relayer");

        call = () -> score.invoke(alice, "registerRelayer", "Hey, I want to be a relayer");
        expectErrorMessage(call, "already registered relayer");

        score.invoke(owner, "setRelayerMinBond", BigInteger.valueOf(100));

        Account bob = sm.createAccount(1000);
        contextMock.when(mockICXSent()).thenReturn(BigInteger.valueOf(90));
        call = () -> score.invoke(bob, "registerRelayer", "Hey, I want to be a relayer");
        expectErrorMessage(call, "require bond at least 100 icx");

        contextMock.when(mockICXSent()).thenReturn(BigInteger.valueOf(200));
        score.invoke(bob, "registerRelayer", "Hey, I want to be a relayer too");

        call = () -> score.invoke(alice, "claimRelayerReward");
        expectErrorMessage(call, "reward is not remained");

        // mocks for distributeRelayerReward
        contextMock.when(mockBlockHeight()).thenReturn(getCurrentBlockHeight() + 43200L);
        contextMock.when(mockICXBalance(BMC_SCORE)).thenReturn(BigInteger.valueOf(10000));
        score.invoke(owner, "setNextRewardDistribution", 43200L);
        score.invoke(owner, "distributeRelayerReward");

        /*
         * Total ICX balance: 10000
         * Bond             : 50 + 200 = 250
         * Amount for reward: 10000 - 250 = 9750
         * Alice reward     : 50 / 250 * 9750 = 1950
         * Bob reward       : 200 / 250 * 9750 = 7800
         */

        Map<String, Relayer> relayers = (Map<String, Relayer>) score.call("getRelayers");

        assertEquals(alice.getAddress(), relayers.get(alice.getAddress().toString()).getAddr());
        assertEquals(BigInteger.valueOf(50), relayers.get(alice.getAddress().toString()).getBond());
        assertEquals(BigInteger.valueOf(1950), relayers.get(alice.getAddress().toString()).getReward());

        assertEquals(bob.getAddress(), relayers.get(bob.getAddress().toString()).getAddr());
        assertEquals(BigInteger.valueOf(200), relayers.get(bob.getAddress().toString()).getBond());
        assertEquals(BigInteger.valueOf(7800), relayers.get(bob.getAddress().toString()).getReward());

        Verification sendICXToUser = () -> Context.transfer(alice.getAddress(), BigInteger.valueOf(1950));
        contextMock.when(sendICXToUser).then(invocation -> null);
        score.invoke(alice, "claimRelayerReward");

        sendICXToUser = () -> Context.transfer(alice.getAddress(), BigInteger.valueOf(50));
        contextMock.when(sendICXToUser).then(invocationOnMock -> null);
        score.invoke(alice, "unregisterRelayer");
        relayers = (Map<String, Relayer>) score.call("getRelayers");
        assertEquals(null, relayers.get(alice.getAddress().toString()));

        call = () -> score.invoke(alice, "unregisterRelayer");
        expectErrorMessage(call, "not found registered relayer");

        call = () -> score.invoke(alice, "claimRelayerReward");
        expectErrorMessage(call, "not found registered relayer");

        call = () -> score.invoke(alice, "removeRelayer", bob.getAddress(), alice.getAddress());
        expectErrorMessage(call, REQUIRE_OWNER_ACCESS);

        // 200 is the bond
        sendICXToUser = () -> Context.transfer(bob.getAddress(), BigInteger.valueOf(200));
        contextMock.when(sendICXToUser).then(invocationOnMock -> null);

        // 7800 is the reward
        sendICXToUser = () -> Context.transfer(bob.getAddress(), BigInteger.valueOf(7800));
        contextMock.when(sendICXToUser).then(invocationOnMock -> null);

        score.invoke(owner, "removeRelayer", bob.getAddress(), bob.getAddress());

        relayers = (Map<String, Relayer>) score.call("getRelayers");
        assertEquals(null, relayers.get(bob.getAddress().toString()));
    }

    @Test
    void messageEventLogContents() {

        Account bts = sm.createAccount();
        String DESTINATION_NETWORK = "0x38.bmc";
        String DESTINATION_BMC = "0x034AaDE86BF402F023Aa17E5725fABC4ab9E9798";
        String to = getDestinationBTPAddress(DESTINATION_NETWORK, DESTINATION_BMC);
        String svc = "bts";
        BigInteger sn = BigInteger.ONE;
        byte[] msg = "Message Received From BTS".getBytes();

        // before registering BTS to BMC
        assertThrows(AssertionError.class, () -> score.invoke(owner, "sendMessage", to, svc, sn, msg));

        addLink(to);
        addRoute(DESTINATION_NETWORK, to);
        addService(svc, bts.getAddress());

        assertThrows(AssertionError.class, () -> score.invoke(owner, "sendMessage", to, svc, sn, msg));
        assertThrows(AssertionError.class,
                () -> score.invoke(bts, "sendMessage", to, svc, BigInteger.ONE.negate(), msg));

        score.invoke(bts, "sendMessage", DESTINATION_NETWORK, svc, sn, msg);
        BTPAddress source = new BTPAddress(NETWORK, BMC_SCORE.toString());
        BTPAddress destination = new BTPAddress(DESTINATION_NETWORK, DESTINATION_BMC);

        BTPMessage message = new BTPMessage();
        message.setSrc(source);
        message.setDst(destination);
        message.setSvc(svc);
        message.setSn(sn);
        message.setPayload(msg);

        verify(scoreSpy).Message(to, BigInteger.TWO, message.toBytes());
    }

    @Test
    @DisplayName("For external service: BSH")
    public void handleRelayMessage() {

        String prevLink = getDestinationBTPAddress("0x38.bsc", "0xa1442c90120A891c3de9793caC70968Cab113234");
        Account relay = registerRelayer();

        addLink(prevLink);

        Executable call = () -> score.invoke(owner, "handleRelayMessage", prevLink, "message");
        expectErrorMessage(call, "relay not registered: " + owner.getAddress());

        score.invoke(owner, "addRelay", prevLink, relay.getAddress());

        call = () -> score.invoke(relay, "handleRelayMessage", prevLink, "@@@");
        expectErrorMessageIn(call, "failed to decode base64 relay message");

        /*
         * nextBMC  : btp://0x1.icon/cx23a91ee3dd290486a9113a6a42429825d813de53
         * sn       : 90
         */
        String str = "-QFN-QFKuQFH-QFEVbkBO_kBOPkBNbg5YnRwOi8vMHgxLmljb24vY3gyM2E5MWVlM2RkMjkwNDg2YTkxMTNhNmE0MjQyOTgyNWQ4MTNkZTUzWrj3-PW4OWJ0cDovLzB4MzguYnNjLzB4MDM0QWFERTg2QkY0MDJGMDIzQWExN0U1NzI1ZkFCQzRhYjlFOTc5OLg5YnRwOi8vMHgxLmljb24vY3gyM2E5MWVlM2RkMjkwNDg2YTkxMTNhNmE0MjQyOTgyNWQ4MTNkZTUzg2J0cy-4ePh2ALhz-HGqMHg3QTQzNDFBZjQ5OTU4ODQ1NDZCY2Y3ZTA5ZUI5OGJlRDNlRDI2RDI4qmh4OTM3NTE3YWMwNDJkMGExNGYwOWQ0Njc3ZDMwMmJiMjExMTg0YWM1ZtrZkWJ0cC0weDM4LmJzYy1CVENChjjX6kxoAIQBSv5z";
        call = () -> score.invoke(relay, "handleRelayMessage", prevLink, str);
        expectErrorMessageIn(call, "Invalid Next BMC");

        /*
         * nextBMC  : btp://0x1.icon/cx0000000000000000000000000000000000000004
         *
         */
        String validBMCBase64 = "-QFN-QFKuQFH-QFEVbkBO_kBOPkBNbg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0Abj3-PW4OWJ0cDovLzB4MzguYnNjLzB4MDM0QWFERTg2QkY0MDJGMDIzQWExN0U1NzI1ZkFCQzRhYjlFOTc5OLg5YnRwOi8vMHgxLmljb24vY3gyM2E5MWVlM2RkMjkwNDg2YTkxMTNhNmE0MjQyOTgyNWQ4MTNkZTUzg2J0cy-4ePh2ALhz-HGqMHg3QTQzNDFBZjQ5OTU4ODQ1NDZCY2Y3ZTA5ZUI5OGJlRDNlRDI2RDI4qmh4OTM3NTE3YWMwNDJkMGExNGYwOWQ0Njc3ZDMwMmJiMjExMTg0YWM1ZtrZkWJ0cC0weDM4LmJzYy1CVENChjjX6kxoAIQBSv5z";
        score.invoke(relay, "handleRelayMessage", prevLink, validBMCBase64);

        // valid next bmc, but sequence ev.getSeq() = 5, when rxSeq is currently 1.
        String invalidSeq = "-QFN-QFKuQFH-QFEVbkBO_kBOPkBNbg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0Bbj3-PW4OWJ0cDovLzB4MzguYnNjLzB4MDM0QWFERTg2QkY0MDJGMDIzQWExN0U1NzI1ZkFCQzRhYjlFOTc5OLg5YnRwOi8vMHgxLmljb24vY3gyM2E5MWVlM2RkMjkwNDg2YTkxMTNhNmE0MjQyOTgyNWQ4MTNkZTUzg2J0cy-4ePh2ALhz-HGqMHg3QTQzNDFBZjQ5OTU4ODQ1NDZCY2Y3ZTA5ZUI5OGJlRDNlRDI2RDI4qmh4OTM3NTE3YWMwNDJkMGExNGYwOWQ0Njc3ZDMwMmJiMjExMTg0YWM1ZtrZkWJ0cC0weDM4LmJzYy1CVENChjjX6kxoAIQBSv5z";
        assertThrows(AssertionError.class, () -> score.invoke(relay, "handleRelayMessage", prevLink, invalidSeq));

        // valid next bmc, valid sequence, but message is not for icon
        String invalidBTPMessage = "-QFc-QFZuQFW-QFTVbkBSvkBR_kBRLg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0ArkBBfj1uDlidHA6Ly8weDM4LmJzYy8weDAzNEFhREU4NkJGNDAyRjAyM0FhMTdFNTcyNWZBQkM0YWI5RTk3OTi4OWJ0cDovLzB4MS5pY29uL2N4MjNhOTFlZTNkZDI5MDQ4NmE5MTEzYTZhNDI0Mjk4MjVkODEzZGU1M4NidHMvuHj4dgC4c_hxqjB4N0E0MzQxQWY0OTk1ODg0NTQ2QmNmN2UwOWVCOThiZUQzZUQyNkQyOKpoeDkzNzUxN2FjMDQyqqqqqqqqqqqqqqqqqqpkMGExNGYwOWQ0Njc3ZDMwMmJiMjExMTg0YWM1ZtrZkWJ0cC0weDM4LmJzYy1CVENChjjX6kxoAIQBSv5z";
        score.invoke(relay, "handleRelayMessage", prevLink, invalidBTPMessage);
        verify(scoreSpy).Message(eq(prevLink), eq(BigInteger.valueOf(3)), any());

        // message for ICON chain itself
        String validBTPMessage = "-ND4zrjM-MpVuML4wPi-uDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDQDuID4frg5YnRwOi8vMHgzOC5ic2MvMHhhMTQ0MmM5MDEyMEE4OTFjM2RlOTc5M2NhQzcwOTY4Q2FiMTEzMjM0uDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYnRzAYIweIQBSv5z";

        // BTS Service not added, so sendError (BTPException e), fail to get service
        score.invoke(relay, "handleRelayMessage", prevLink, validBTPMessage);
        verify(scoreSpy).Message(eq(prevLink), eq(BigInteger.valueOf(4)), any());

        // add bts service
        addService("bts", Address.fromString("cxcef70e92b89f2d8191a0582de966280358713c32"));

        // handleBTPMessage not mocked, so sendError (Exception e), IllegalStateException
        validBTPMessage = "-ND4zrjM-MpVuML4wPi-uDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDQEuID4frg5YnRwOi8vMHgzOC5ic2MvMHhhMTQ0MmM5MDEyMEE4OTFjM2RlOTc5M2NhQzcwOTY4Q2FiMTEzMjM0uDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYnRzAYIweIQBSv5z";
        score.invoke(relay, "handleRelayMessage", prevLink, validBTPMessage);
        verify(scoreSpy).Message(eq(prevLink), eq(BigInteger.valueOf(5)), any());

        // mock handleBTPMessage
        validBTPMessage = "-ND4zrjM-MpVuML4wPi-uDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDQFuID4frg5YnRwOi8vMHgzOC5ic2MvMHhhMTQ0MmM5MDEyMEE4OTFjM2RlOTc5M2NhQzcwOTY4Q2FiMTEzMjM0uDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYnRzAYIweIQBSv5z";
        Verification mock_handleBTPMessage = () -> Context.call(any(), eq("handleBTPMessage"), any(Object.class));
        contextMock.when(mock_handleBTPMessage).thenReturn(null);

        score.invoke(relay, "handleRelayMessage", prevLink, validBTPMessage);
    }

    @Test
    @Order(2)
    public void initService() {
        String prevLink = getDestinationBTPAddress("0x228.arctic", "0x1111111111111111111111111111111111111111");
        Account relay = registerRelayer();

        addLink(prevLink);
        score.invoke(owner, "addRelay", prevLink, relay.getAddress());

        /**
         * Internal init
         * sn = 1
         * 0x38.bsc is reachable through 0x228.arctic
         * initMessage contains 0x38.bsc as payload
         */
        String initMessage = "-QEh-QEeuQEb-QEYVbkBD_kBDPkBCbg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0AbjL-Mm4PWJ0cDovLzB4MjI4LmFyY3RpYy8weDExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTG4OWJ0cDovLzB4MS5pY29uL2N4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwNINibWMBuEj4RoRJbml0uD_4Pfg7uDlidHA6Ly8weDM4LmJzYy8weGExNDQyYzkwMTIwQTg5MWMzZGU5NzkzY2FDNzA5NjhDYWIxMTMyMzWEAUr-cw==";
        score.invoke(relay, "handleRelayMessage", prevLink, initMessage);
        score.invoke(owner, "addService", BTS, BTSScore.getAddress());

        score.invoke(BTSScore, "sendMessage", "0x38.bsc", BTS, BigInteger.valueOf(2), new byte[0]);
        verify(scoreSpy).Message(eq(prevLink), eq(BigInteger.TWO), any());

        String invalidInitMessage = "-QEf-QEcuQEZ-QEWVbkBDfkBCvkBB7g5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0ArjJ-Me4O2J0cDovLzB4MjI4LmF2YXgvMHgxMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExuDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYm1jAbhI-EaESW5pdLg_-D34O7g5YnRwOi8vMHgzOC5ic2MvMHhhMTQ0MmM5MDEyMEE4OTFjM2RlOTc5M2NhQzcwOTY4Q2FiMTEzMjM1hAFK_nM=";
        Executable call = () -> score.invoke(relay, "handleRelayMessage", prevLink, invalidInitMessage);
        expectErrorMessageIn(call, "internal message not allowed from ");

        // not supported internal type, only rewards distributed
        String invalidType = "-QEg-QEduQEa-QEXVbkBDvkBC_kBCLg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0ArjK-Mi4O2J0cDovLzB4MjI4LmF2YXgvMHgxMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExuDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYm1jAbhJ-EeFVGlyZWS4P_g9-Du4OWJ0cDovLzB4MzguYnNjLzB4YTE0NDJjOTAxMjBBODkxYzNkZTk3OTNjYUM3MDk2OENhYjExMzIzNYQBSv5z";
        score.invoke(relay, "handleRelayMessage", prevLink, invalidType);
    }

    @Test
    @Order(2)
    public void linkUnlink() {
        String prevLink = getDestinationBTPAddress("0x228.avax", "0x1111111111111111111111111111111111111111");
        Account relay = registerRelayer();

        addLink(prevLink);
        score.invoke(owner, "addRelay", prevLink, relay.getAddress());

        String linkMessage = "-QEd-QEauQEX-QEUVbkBC_kBCPkBBbg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0AbjH-MW4O2J0cDovLzB4MjI4LmF2YXgvMHgxMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExuDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYm1jAbhG-ESETGlua7g9-Du4OWJ0cDovLzB4MzguYnNjLzB4YTE0NDJjOTAxMjBBODkxYzNkZTk3OTNjYUM3MDk2OENhYjExMzIzNYQBSv5z";
        score.invoke(relay, "handleRelayMessage", prevLink, linkMessage);
        score.invoke(owner, "addService", BTS, BTSScore.getAddress());

        // link : reachable
        score.invoke(BTSScore, "sendMessage", "0x38.bsc", BTS, BigInteger.valueOf(2), new byte[0]);
        verify(scoreSpy).Message(eq(prevLink), eq(BigInteger.TWO), any());

        // unlink : reachable
        String unlinkMessage = "-QEf-QEcuQEZ-QEWVbkBDfkBCvkBB7g5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0ArjJ-Me4O2J0cDovLzB4MjI4LmF2YXgvMHgxMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExuDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDSDYm1jAbhI-EaGVW5saW5ruD34O7g5YnRwOi8vMHgzOC5ic2MvMHhhMTQ0MmM5MDEyMEE4OTFjM2RlOTc5M2NhQzcwOTY4Q2FiMTEzMjM1hAFK_nM=";
        score.invoke(relay, "handleRelayMessage", prevLink, unlinkMessage);

        // unreachable exception
        assertThrows(AssertionError.class,
                () -> score.invoke(BTSScore, "sendMessage", "0x38.bsc", BTS, BigInteger.TEN, new byte[0]));
    }

    @Test
    @Order(2)
    public void feeGatheringMessage() {
        String prevLink = getDestinationBTPAddress("0x228.arctic", "0x1111111111111111111111111111111111111111");
        Account relay = registerRelayer();

        addLink(prevLink);
        score.invoke(owner, "addRelay", prevLink, relay.getAddress());
        score.invoke(owner, "addService", BTS, BTSScore.getAddress());
        String feeMessage;
        contextMock.when(() -> Context.call(eq(BTSScore.getAddress()), eq("handleFeeGathering"), any())).thenReturn(null);

        feeMessage = "-QEw-QEtuQEq-QEnVbkBHvkBG_kBGLg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0Abja-Ni4PWJ0cDovLzB4MjI4LmFyY3RpYy8weDExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTG4OWJ0cDovLzB4MS5pY29uL2N4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwNINibWMBuFf4VYxGZWVHYXRoZXJpbme4RvhEuD1idHA6Ly8weDIyOC5hcmN0aWMvMHhhMTQ0MmM5MDEyMEE4OTFjM2RlOTc5M2NhQzcwOTY4Q2FiMTEzMjM1xINidHOEAUr-cw==";
        score.invoke(relay, "handleRelayMessage", prevLink, feeMessage);

        // not allowed GatherFeeMessage from random chain
        String invalidMessage = "-QEs-QEpuQEm-QEjVbkBGvkBF_kBFLg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0ArjW-NS4PWJ0cDovLzB4MjI4LmFyY3RpYy8weDExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTG4OWJ0cDovLzB4MS5pY29uL2N4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwNINibWMBuFP4UYxGZWVHYXRoZXJpbme4QvhAuDlidHA6Ly8weDM4LmJzYy8weGExNDQyYzkwMTIwQTg5MWMzZGU5NzkzY2FDNzA5NjhDYWIxMTMyMzXEg2J0c4QBSv5z";
        score.invoke(relay, "handleRelayMessage", prevLink, invalidMessage);
    }

    @Test
    @Order(2)
    public void handleErrorMessage() {
        String prevLink = getDestinationBTPAddress("0x228.arctic", "0x1111111111111111111111111111111111111111");
        Account relay = registerRelayer();

        addLink(prevLink);
        score.invoke(owner, "addRelay", prevLink, relay.getAddress());
        score.invoke(owner, "addService", "candice", BTSScore.getAddress());

        String err = "-PD47rjs-OpVuOL44PjeuDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDQBuKD4nrg9YnRwOi8vMHgyMjguYXJjdGljLzB4MTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMbg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0h2NhbmRpY2WB9pnYAJZBZGQgdG8gYmxhY2tsaXN0IGVycm9yhAFK_nM=";
        score.invoke(relay, "handleRelayMessage", prevLink, err);
        verify(scoreSpy).ErrorOnBTPError(eq("candice"), eq(BigInteger.TEN), eq(-1L), eq(null) ,eq(-1L), any());

        contextMock.when(() ->
                Context.call(eq(BTSScore.getAddress()), eq("handleBTPError"), any()
                )).thenReturn(null);

        String err2 = "-PD47rjs-OpVuOL44PjeuDlidHA6Ly8weDEuaWNvbi9jeDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDQCuKD4nrg9YnRwOi8vMHgyMjguYXJjdGljLzB4MTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMbg5YnRwOi8vMHgxLmljb24vY3gwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0h2NhbmRpY2WB9pnYAJZBZGQgdG8gYmxhY2tsaXN0IGVycm9yhAFK_nM=";
        score.invoke(relay, "handleRelayMessage", prevLink, err2);
    }

    @Test
    public void ownerTests() {
        String expectedErrorMessage = "caller is not owner";
        Account user1 = sm.createAccount(10);
        Account user2 = sm.createAccount(10);

        // Non-Owner tries to add a new Owner
        Executable call = () -> score.invoke(nonOwner, "addOwner", owner.getAddress());
        expectErrorMessage(call, expectedErrorMessage);

        // owner tries to add themselves
        call = () -> score.invoke(owner, "addOwner", owner.getAddress());
        expectErrorMessage(call, "given address is score owner");

        // Current Owner adds a new Owner
        score.invoke(owner, "addOwner", user1.getAddress());
        assertEquals(true, score.call("isOwner", user1.getAddress()));
        Address[] owners = (Address[]) score.call("getOwners");
        assertEquals(owner.getAddress(), owners[0]);
        assertEquals(user1.getAddress(), owners[1]);

        // newly added owner tries to add owner
        score.invoke(user1, "addOwner", user2.getAddress());
        assertEquals(true, score.call("isOwner", user2.getAddress()));

        //Current Owner removes another Owner
        score.invoke(user2, "removeOwner", user1.getAddress());
        assertEquals(false, score.call("isOwner", user1.getAddress()));

        // owner tries to add itself again
        call = () -> score.invoke(user2, "addOwner", user2.getAddress());
        expectErrorMessage(call, "already exists owner");

        // The last Owner removes him/herself
        score.invoke(user2, "removeOwner", user2.getAddress());
        assertEquals(false, score.call("isOwner", user2.getAddress()));
    }


    @Test
    public void bmcDecodeMessage() {
//        String str = "0xf8a9b8406274703a2f2f307836333536346334302e686d6e792f307861363937313261333831336430353035626244353541654433666438343731426332663732324444b8396274703a2f2f3078312e69636f6e2f6378393937383439643339323064333338656438313830303833336662623237306337383565373433649a576f6e6465726c616e64546f6b656e53616c655365727669636589008963dd8c2c5e000086c50283c20100";
//        String str = "0xf8a9b8406274703a2f2f307836333536346334302e686d6e792f307861363937313261333831336430353035626244353541654433666438343731426332663732324444b8396274703a2f2f3078312e69636f6e2f6378393937383439643339323064333338656438313830303833336662623237306337383565373433649a576f6e6465726c616e64546f6b656e53616c655365727669636589ff769c2273d3a2000086c50283c20100";
        String str = "0xf8a8b8406274703a2f2f307836333536346334302e686d6e792f307861363937313261333831336430353035626244353541654433666438343731426332663732324444b8396274703a2f2f3078312e69636f6e2f6378393937383439643339323064333338656438313830303833336662623237306337383565373433649a576f6e6465726c616e64546f6b656e53616c6553657276696365886f05b59d3b20000086c50283c20100";
        byte[] msg = Hex.decode(str.substring(2));

        BTPMessage btpMsg = BTPMessage.fromBytes(msg);

        System.out.println(btpMsg);
        System.out.println(btpMsg.getSn());
    }

    @Test
    public void generateBTPMessage() {

        FeeGatheringMessage feem = new FeeGatheringMessage();
        BTPAddress addr = new BTPAddress("0x228.arctic", "0xa1442c90120A891c3de9793caC70968Cab113235");
        feem.setFa(addr);
        feem.setSvcs(new String[]{"bts"});

        BMCMessage bmcm = new BMCMessage();
        bmcm.setType("FeeGathering");
        bmcm.setPayload(feem.toBytes());

//        ErrorMessage em = new ErrorMessage();
//        em.setCode(0L);
//        em.setMsg("Add to blacklist error");

        BTPAddress d = new BTPAddress("0x1.icon", "cx0000000000000000000000000000000000000004");
        BTPAddress s = new BTPAddress("0x228.arctic", "0x1111111111111111111111111111111111111111");

        BTPMessage m = new BTPMessage();
        m.setSrc(s);
        m.setDst(d);
        m.setSn(BigInteger.valueOf(10).negate());
        m.setSvc("bts");
        m.setPayload(feem.toBytes());
        byte[] byt = m.toBytes();
        String aa = Hex.toHexString(byt);
        System.out.println(aa);
    }
}
