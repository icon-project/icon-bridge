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
import org.junit.jupiter.api.function.Executable;
import org.junit.jupiter.api.Test;
import org.mockito.MockedStatic.Verification;
import score.Address;
import score.Context;

class BTPMessageCenterTest extends AbstractBTPMessageCenterTest {


    @Test
    void name() {
        assertEquals("BTP Message Center", score.call("name"));
    }

    @Test
    void getBTPAddress() {
        assertEquals(getBTPAddress(score.getAddress()), score.call("getBtpAddress"));
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
        contextMock.when(mockBlockHeight()).thenReturn(getCurrentBlockHeight()+43200L);
        contextMock.when(mockICXBalance(score.getAddress())).thenReturn(BigInteger.valueOf(10000));
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
        assertThrows(AssertionError.class, () -> score.invoke(owner, "sendMessage", to, svc, sn, msg ));

        addLink(to);
        addRoute(DESTINATION_NETWORK, to);
        addService(svc, bts.getAddress());

        assertThrows(AssertionError.class, () -> score.invoke(owner, "sendMessage", to, svc, sn, msg));
        assertThrows(AssertionError.class, () -> score.invoke(bts, "sendMessage", to, svc, BigInteger.ONE.negate(), msg));

        score.invoke(bts, "sendMessage", DESTINATION_NETWORK, svc, sn, msg);

        BTPAddress source = new BTPAddress(NETWORK, score.getAddress().toString());
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
    public void bmcDecodeMessage() {
//        String str = "0xf8a9b8406274703a2f2f307836333536346334302e686d6e792f307861363937313261333831336430353035626244353541654433666438343731426332663732324444b8396274703a2f2f3078312e69636f6e2f6378393937383439643339323064333338656438313830303833336662623237306337383565373433649a576f6e6465726c616e64546f6b656e53616c655365727669636589008963dd8c2c5e000086c50283c20100";
//        String str = "0xf8a9b8406274703a2f2f307836333536346334302e686d6e792f307861363937313261333831336430353035626244353541654433666438343731426332663732324444b8396274703a2f2f3078312e69636f6e2f6378393937383439643339323064333338656438313830303833336662623237306337383565373433649a576f6e6465726c616e64546f6b656e53616c655365727669636589ff769c2273d3a2000086c50283c20100";
        String str = "0xf8a8b8406274703a2f2f307836333536346334302e686d6e792f307861363937313261333831336430353035626244353541654433666438343731426332663732324444b8396274703a2f2f3078312e69636f6e2f6378393937383439643339323064333338656438313830303833336662623237306337383565373433649a576f6e6465726c616e64546f6b656e53616c6553657276696365886f05b59d3b20000086c50283c20100";
        byte[] msg = Hex.decode(str.substring(2));

        BTPMessage btpMsg = BTPMessage.fromBytes(msg);

        System.out.println(btpMsg);
        System.out.println(btpMsg.getSn());
    }
}
