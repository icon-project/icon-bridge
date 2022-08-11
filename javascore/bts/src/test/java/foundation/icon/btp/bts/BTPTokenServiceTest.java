///*
// * Copyright 2021 ICON Foundation
// *
// * Licensed under the Apache License, Version 2.0 (the "License");
// * you may not use this file except in compliance with the License.
// * You may obtain a copy of the License at
// *
// *     http://www.apache.org/licenses/LICENSE-2.0
// *
// * Unless required by applicable law or agreed to in writing, software
// * distributed under the License is distributed on an "AS IS" BASIS,
// * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// * See the License for the specific language governing permissions and
// * limitations under the License.
// */
//
//package foundation.icon.btp.bts;
//
//import foundation.icon.btp.lib.BTPAddress;
//import foundation.icon.btp.test.BTPIntegrationTest;
//import foundation.icon.btp.test.MockBMCIntegrationTest;
//import foundation.icon.btp.test.SendMessageEventLog;
//import foundation.icon.jsonrpc.Address;
//import foundation.icon.jsonrpc.model.TransactionResult;
//import foundation.icon.score.client.RevertedException;
//import foundation.icon.score.test.ScoreIntegrationTest;
//import org.junit.jupiter.api.BeforeAll;
//import org.junit.jupiter.api.Test;
//import org.junit.jupiter.api.function.Executable;
//import score.UserRevertedException;
//
//import java.math.BigInteger;
//import java.util.ArrayList;
//import java.util.Arrays;
//import java.util.List;
//import java.util.function.Consumer;
//import java.util.stream.Collectors;
//
//import static org.junit.jupiter.api.Assertions.assertEquals;
//import static org.junit.jupiter.api.Assertions.assertTrue;
//
//class BTPTokenServiceTest implements BTSIntegrationTest {
//
//    static Address btsAddress = btsClient._address();
//    static Address owner = Address.of(btsClient._wallet());
//    static BTPAddress link = BTPIntegrationTest.Faker.btpLink();
//    static String linkNet = BTPIntegrationTest.Faker.btpNetwork();
//    static Address testerAddress = Address.of(tester);
//    static BTPAddress to = new BTPAddress(BTPAddress.PROTOCOL_BTP, linkNet, testerAddress.toString());
//    static String nativeCoinName = bts.coinNames()[0];
//    static BigInteger nativeValue = BigInteger.valueOf(10);
//    static String coinName = "coin";
//    static BigInteger coinValue = BigInteger.valueOf(10);
//    static BigInteger coinId;
//    static long code = 1;
//    static String msg = "err";
//
//    static String net = MockBMCIntegrationTest.mockBMC.getNet();
//    static BTPAddress fa = new BTPAddress(BTPAddress.PROTOCOL_BTP, net, testerAddress.toString());
//    static BTPAddress linkFa = new BTPAddress(BTPAddress.PROTOCOL_BTP, linkNet, testerAddress.toString());
//    static BigInteger feeRatio = BigInteger.TEN;
//
//    // static boolean isExistsCoin(String name) {
//    //     return ScoreIntegrationTest.indexOf(bts.coinNames(), name) >= 0;
//    // }
//
//    // static void register(String name) {
//    //     bts.register(name);
//    //     assertTrue(isExistsCoin(name));
//    // }
//
//    // static BigInteger transferBatch(TransferTransaction transaction) {
//    //     BigInteger nativeCoinValue = nativeCoinValue(transaction.getAssets());
//    //     List<Asset> assets = coinAssets(transaction.getAssets());
//    //     BigInteger[] coinIds = coinIds(assets);
//    //     String[] coinNames = coinNames(assets);
//    //     BigInteger[] coinValues = coinValues(assets);
//    //     Address from = new Address(transaction.getFrom());
//    //     BigInteger[] snContainer = new BigInteger[1];
//    //     Consumer<TransactionResult> checker = transferStartEventLogChecker(transaction)
//    //             .andThen(sendMessageEventLogChecker(transaction, snContainer))
//    //             .andThen(IRC31SupplierTest.transferFromBatchChecker(
//    //                     btsAddress, from, btsAddress, coinIds, coinValues));
//    //     Executable executable = () -> ((BTSScoreClient) bts).transferBatch(
//    //             checker,
//    //             nativeCoinValue,
//    //             coinNames, coinValues, to.toString());
//    //     ScoreIntegrationTest.balanceCheck(btsAddress, nativeCoinValue, () ->
//    //             IRC31SupplierTest.balanceBatchCheck(btsAddress, coinIds, coinValues, () ->
//    //                     balanceBatchCheck(from, transaction.getAssets(), executable,
//    //                             BalanceCheckType.lock)));
//    //     return snContainer[0];
//    // }
//
//
//    // static void handleTransferResponse(TransferTransaction transaction, BigInteger sn) {
//    //     TransferResponse response = new TransferResponse();
//    //     response.setCode(TransferResponse.RC_OK);
//    //     response.setMessage(TransferResponse.OK_MSG);
//    //     BTSMessage btsMessage = new BTSMessage();
//    //     btsMessage.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
//    //     btsMessage.setData(response.toBytes());
//
//    //     List<Asset> assets = coinAssets(transferRequest(transaction).getAssets());
//    //     BigInteger[] coinIds = coinIds(assets);
//    //     BigInteger[] coinValues = coinValues(assets);
//    //     Address from = new Address(transaction.getFrom());
//    //     Consumer<TransactionResult> checker = IRC31SupplierTest.burnBatchChecker(
//    //             btsAddress, btsAddress, coinIds, coinValues)
//    //             .andThen(transferEndEventLogChecker(from, sn, response));
//    //     balanceBatchCheck(from, transaction.getAssets(), () ->
//    //                     ((MockBMCScoreClient) MockBMCIntegrationTest.mockBMC).intercallHandleBTPMessage(
//    //                             checker,
//    //                             btsAddress, linkNet, BTPTokenService.SERVICE, sn, btsMessage.toBytes()),
//    //             BalanceCheckType.unlock);
//    // }
///*
//    static void handleBTPError(TransferTransaction transaction, BigInteger sn, long code, String msg) {
//        TransferResponse response = new TransferResponse();
//        response.setCode(TransferResponse.RC_ERR);
//        response.setMessage(("BTPError [code:" + code + ",msg:" + msg));
//        BTSMessage btsMessage = new BTSMessage();
//        btsMessage.setServiceType(BTSMessage.REPONSE_HANDLE_SERVICE);
//        btsMessage.setData(response.toBytes());
//
//        Address from = new Address(transaction.getFrom());
//        balanceBatchCheck(
//                from,
//                transaction.getAssets(),
//                () -> ((MockBMCScoreClient) MockBMCIntegrationTest.mockBMC).intercallHandleBTPError(
//                        transferEndEventLogChecker(from, sn, response),
//                        btsAddress, link.toString(), BTPTokenService.SERVICE, sn, code, msg),
//                BalanceCheckType.refund);
//    }*/
//
//    static Consumer<TransactionResult> transferStartEventLogChecker(TransferTransaction transaction) {
//        return BTSIntegrationTest.eventLogChecker(TransferStartEventLog::eventLogs, (el) -> {
//            System.out.println(el);
//            assertEquals(transaction.getFrom(), el.getFrom().toString());
//            assertEquals(transaction.getTo(), el.getTo());
//            AssertBTS.assertEqualsAssetTransferDetails(transaction.getAssets(), el.getAssets());
//        });
//    }
//
//    static Consumer<TransactionResult> transferEndEventLogChecker(Address from, BigInteger sn, TransferResponse response) {
//        return BTSIntegrationTest.eventLogChecker(TransferEndEventLog::eventLogs, (el) -> {
//            assertEquals(from, el.getFrom());
//            assertEquals(sn, el.getSn());
//            assertEquals(response.getCode(), el.getCode());
//            assertEquals(response.getMessage(), new String(el.getMsg()));
//        });
//    }
//
//    static Consumer<TransactionResult> sendMessageEventLogChecker(TransferTransaction transaction) {
//        return sendMessageEventLogChecker(transaction, null);
//    }
//
//    static Consumer<TransactionResult> sendMessageEventLogChecker(TransferTransaction transaction, BigInteger[] snContainer) {
//        return MockBMCIntegrationTest.eventLogChecker(SendMessageEventLog::eventLogs, (el) -> {
//            assertEquals(BTPAddress.valueOf(transaction.getTo()).net(), el.getTo());
//            assertEquals(BTPTokenService.SERVICE, el.getSvc());
//            if (snContainer != null) {
//                snContainer[0] = el.getSn();
//            }
////            assertEquals(sn, el.getSn());
//            BTSMessage btsMessage = BTSMessage.fromBytes(el.getMsg());
//            assertEquals(BTSMessage.REQUEST_COIN_TRANSFER, btsMessage.getServiceType());
//            AssertBTS.assertEqualsTransferRequest(transaction, TransferRequest.fromBytes(btsMessage.getData()));
//        });
//    }
//
//    static void lockedBalanceCheck(Address address, Asset asset, Executable executable) {
//        Balance balance = bts.balanceOf(address, asset.getCoinName());
//        try {
//            executable.execute();
//        } catch (UserRevertedException | RevertedException e) {
//            throw e;
//        } catch (Throwable e) {
//            throw new RuntimeException(e);
//        }
//        assertEquals(balance.getLocked().add(asset.getAmount()),
//                bts.balanceOf(address, asset.getCoinName()).getLocked());
//    }
//
//    enum BalanceCheckType {
//        lock, unlock, refund
//    }
//
//    static void balanceBatchCheck(Address address, AssetTransferDetail[] assetDetails, Executable executable, BalanceCheckType type) {
//        List<AssetTransferDetail> list = Arrays.asList(assetDetails);
//        String[] coinNames = list.stream()
//                .map(AssetTransferDetail::getCoinName).toArray(String[]::new);
//        BigInteger[] values = list.stream()
//                .map((a) -> a.getAmount().add(a.getFee())).toArray(BigInteger[]::new);
//
//        Balance[] balances = bts.balanceOfBatch(address, coinNames);
//        assertEquals(coinNames.length, balances.length);
//        try {
//            executable.execute();
//        } catch (UserRevertedException | RevertedException e) {
//            throw e;
//        } catch (Throwable e) {
//            throw new RuntimeException(e);
//        }
//        Balance[] actual = bts.balanceOfBatch(address, coinNames);
//        for (int i = 0; i < coinNames.length; i++) {
//            BigInteger locked = balances[i].getLocked();
//            if (BalanceCheckType.lock.equals(type)) {
//                locked = locked.add(values[i]);
//            } else if (BalanceCheckType.unlock.equals(type)) {
//                locked = locked.subtract(values[i]);
//            } else if (BalanceCheckType.refund.equals(type)) {
//                locked = locked.subtract(values[i]);
//                assertEquals(balances[i].getRefundable().add(values[i]), actual[i].getRefundable());
//            }
//            assertEquals(locked, actual[i].getLocked());
//        }
//    }
//
//    static TransferRequest transferRequest(TransferTransaction transaction) {
//        TransferRequest request = new TransferRequest();
//        request.setFrom(transaction.getFrom());
//        request.setTo(BTPAddress.valueOf(transaction.getTo()).account());
//
//        AssetTransferDetail[] assetDetails = transaction.getAssets();
//        Asset[] assets = new Asset[assetDetails.length];
//        for (int i = 0; i < assetDetails.length; i++) {
//            assets[i] = new Asset(assetDetails[i]);
//        }
//        request.setAssets(assets);
//        return request;
//    }
//
//    static TransferTransaction transferTransaction(Address from, BTPAddress to, BigInteger feeRatio, Asset... assets) {
//        TransferTransaction transaction = new TransferTransaction();
//        transaction.setFrom(from.toString());
//        transaction.setTo(to.toString());
//        AssetTransferDetail[] assetDetails = new AssetTransferDetail[assets.length];
//        for (int i = 0; i < assets.length; i++) {
//            Asset asset = assets[i];
//            AssetTransferDetail assetDetail = new AssetTransferDetail();
//            assetDetail.setCoinName(asset.getCoinName());
//            BigInteger fee = asset.getAmount().multiply(feeRatio).divide(BTPTokenService.FEE_DENOMINATOR);
//            if (feeRatio.compareTo(BigInteger.ZERO) > 0 && fee.compareTo(BigInteger.ZERO) == 0) {
//                fee = BigInteger.ONE;
//            }
//            assetDetail.setFee(fee);
//            assetDetail.setAmount(asset.getAmount().subtract(assetDetail.getFee()));
//            assetDetails[i] = assetDetail;
//        }
//        transaction.setAssets(assetDetails);
//        return transaction;
//    }
//
//    static BigInteger nativeCoinValue(AssetTransferDetail[] assetDetails) {
//        return Arrays.stream(assetDetails)
//                .filter((a) -> a.getCoinName().equals(nativeCoinName))
//                .map((a) -> a.getAmount().add(a.getFee()))
//                .findAny().orElse(BigInteger.ZERO);
//    }
//
//    static BigInteger nativeCoinValue(Asset[] assets) {
//        return Arrays.stream(assets)
//                .filter((a) -> a.getCoinName().equals(nativeCoinName))
//                .map(Asset::getAmount)
//                .findAny().orElse(BigInteger.ZERO);
//    }
//
//    static List<Asset> coinAssets(AssetTransferDetail[] assetDetails) {
//        return Arrays.stream(assetDetails)
//                .filter((a) -> !a.getCoinName().equals(nativeCoinName))
//                .map((a) -> {
//                    Asset asset = new Asset(a);
//                    asset.setAmount(a.getAmount().add(a.getFee()));
//                    return asset;
//                }).collect(Collectors.toList());
//    }
//
//    static List<Asset> coinAssets(Asset[] assets) {
//        return Arrays.stream(assets)
//                .filter((a) -> !a.getCoinName().equals(nativeCoinName))
//                .collect(Collectors.toList());
//    }
//
//    // static BigInteger[] coinIds(List<Asset> assets) {
//    //     return assets.stream()
//    //             .map((a) -> bts.coinId(a.getCoinName())).toArray(BigInteger[]::new);
//    // }
//
//    static String[] coinNames(List<Asset> assets) {
//        return assets.stream()
//                .map(Asset::getCoinName).toArray(String[]::new);
//    }
//
//    static BigInteger[] coinValues(List<Asset> assets) {
//        return assets.stream()
//                .map(Asset::getAmount).toArray(BigInteger[]::new);
//    }
//
//    static Asset[] feeAssets() {
//        String[] coinNames = bts.coinNames();
//        Balance[] balances = bts.balanceOfBatch(btsAddress, coinNames);
//        List<Asset> feeAssets = new ArrayList<>();
//        for (int i = 0; i < coinNames.length; i++) {
//            BigInteger coinFee = balances[i].getRefundable();
//            if (coinFee.compareTo(BigInteger.ZERO) > 0) {
//                String coinName = coinNames[i];
//                feeAssets.add(new Asset(coinName, coinFee));
//            }
//        }
//        return feeAssets.toArray(Asset[]::new);
//    }
//
//    // @BeforeAll
//    // static void beforeAll() {
//    //     if (!isExistsCoin(coinName)) {
//    //         register(coinName);
//    //         coinId = bts.coinId(coinName);
//    //     }
//    //     if (!bts.feeRatio().equals(feeRatio)) {
//    //         bts.setFeeRatio(feeRatio);
//    //     }
//    // }
//
//    // @Test
//    // void registerShouldSuccess() {
//    //     register(ScoreIntegrationTest.Faker.faker.name().name());
//    // }
//
//    // @Test
//    // void transferNativeCoinShouldMakeEventLogAndLockBalance() {
//    //     Asset asset = new Asset(nativeCoinName, nativeValue);
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio, asset);
//    //     Consumer<TransactionResult> checker = transferStartEventLogChecker(transaction)
//    //             .andThen(sendMessageEventLogChecker(transaction));
//    //     ScoreIntegrationTest.balanceCheck(btsAddress, asset.getAmount(), () ->
//    //             lockedBalanceCheck(new Address(transaction.getFrom()), asset, () ->
//    //                     ((BTSScoreClient) bts).transferNativeCoin(
//    //                             checker,
//    //                             asset.getAmount(),
//    //                             transaction.getTo())));
//    // }
//
//    // @Test
//    // void transferShouldMakeEventLogAndLockBalance() {
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     Asset asset = new Asset(coinName, coinValue);
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio, asset);
//    //     Address from = new Address(transaction.getFrom());
//    //     Consumer<TransactionResult> checker = transferStartEventLogChecker(transaction)
//    //             .andThen(sendMessageEventLogChecker(transaction))
//    //             .andThen(IRC31SupplierTest.transferFromChecker(
//    //                     btsAddress, from, btsAddress, coinId, asset.getAmount()));
//    //     IRC31SupplierTest.balanceCheck(btsAddress, coinId, coinValue, () ->
//    //             lockedBalanceCheck(from, asset, () ->
//    //                     ((BTSScoreClient) bts).transfer(
//    //                             checker,
//    //                             asset.getCoinName(), asset.getAmount(), transaction.getTo())
//    //             )
//    //     );
//    // }
//
//    // @Test
//    // void transferBatchShouldShouldMakeEventLogAndLockBalance() {
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//
//    //     transferBatch(transaction);
//    // }
//
//    // @Test
//    // void handleTransferRequestShouldIRC31MintBatchAndResponse() {
//    //     //mint to tester
//    //     IRC31SupplierTest.mint(testerAddress, coinId, coinValue);
//
//    //     //transfer owner to tester
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//    //     TransferRequest request = transferRequest(transaction);
//    //     List<Asset> assets = coinAssets(request.getAssets());
//
//    //     BigInteger[] coinIds = coinIds(assets);
//    //     BigInteger[] coinValues = coinValues(assets);
//
//    //     BTSMessage btsMessage = new BTSMessage();
//    //     btsMessage.setServiceType(BTSMessage.REQUEST_COIN_TRANSFER);
//    //     btsMessage.setData(request.toBytes());
//
//    //     Address to = new Address(BTPAddress.valueOf(transaction.getTo()).account());
//
//    //     BigInteger sn = BigInteger.ONE;
//    //     Consumer<TransactionResult> checker = IRC31SupplierTest.mintBatchChecker(
//    //             btsAddress, to, coinIds, coinValues).andThen(
//    //             MockBMCIntegrationTest.eventLogChecker(SendMessageEventLog::eventLogs, (el) -> {
//    //                 assertEquals(linkNet, el.getTo());
//    //                 assertEquals(BTPTokenService.SERVICE, el.getSvc());
//    //                 assertEquals(sn, el.getSn());
//    //                 BTSMessage btsMsg = BTSMessage.fromBytes(el.getMsg());
//    //                 assertEquals(BTSMessage.REPONSE_HANDLE_SERVICE, btsMsg.getServiceType());
//    //                 TransferResponse response = TransferResponse.fromBytes(btsMsg.getData());
//    //                 assertEquals(TransferResponse.RC_OK, response.getCode());
//    //                 assertEquals(TransferResponse.OK_MSG, response.getMessage());
//    //             }));
//    //     Executable executable = () -> ((MockBMCScoreClient) MockBMCIntegrationTest.mockBMC)
//    //             .intercallHandleBTPMessage(
//    //                     checker,
//    //                     btsAddress, linkNet, BTPTokenService.SERVICE, sn, btsMessage.toBytes());
//    //     ScoreIntegrationTest.balanceCheck(to, nativeCoinValue(request.getAssets()), () ->
//    //             IRC31SupplierTest.balanceBatchCheck(to, coinIds, coinValues, executable));
//    // }
//
//    // @Test
//    // void handleTransferResponseShouldIRC31BurnBatchAndMakeEventLog() {
//    //     //mint
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     //transferBatch
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//    //     BigInteger sn = transferBatch(transaction);
//
//    //     //
//    //     handleTransferResponse(transaction, sn);
//    // }
//
//    // @Test
//    // void handleUnknownResponseShouldMakeEventLog() {
//    //     TransferResponse response = new TransferResponse();
//    //     response.setCode(TransferResponse.RC_ERR);
//    //     response.setMessage(TransferResponse.ERR_MSG_UNKNOWN_TYPE);
//    //     BTSMessage btsMessage = new BTSMessage();
//    //     btsMessage.setServiceType(BTSMessage.UNKNOWN_TYPE);
//    //     btsMessage.setData(response.toBytes());
//
//    //     BigInteger sn = BigInteger.ONE;
//    //     ((MockBMCScoreClient) MockBMCIntegrationTest.mockBMC).intercallHandleBTPMessage(
//    //             BTSIntegrationTest.eventLogChecker(UnknownResponseEventLog::eventLogs, (el) -> {
//    //                 assertEquals(linkNet, el.getFrom());
//    //                 assertEquals(sn, el.getSn());
//    //             }),
//    //             btsAddress, linkNet, BTPTokenService.SERVICE, sn, btsMessage.toBytes());
//    // }
//
//    // @Test
//    // void handleBTPMessageShouldRevert() {
//    //     AssertBTSException.assertUnknown(() ->
//    //             btsBSH.handleBTPMessage(linkNet, BTPTokenService.SERVICE, BigInteger.ONE, new byte[]{}));
//    // }
//
//    // @Test
//    // void handleBTPErrorShouldMakeEventLogAndAddRefundableBalance() {
//    //     //mint
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     //transferBatch
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//    //     BigInteger sn = transferBatch(transaction);
//
//    //     //handleBTPError
//    //     handleBTPError(transaction, sn, code, msg);
//    // }
//
//    // @Test
//    // void reclaim() {
//    //     //mint
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     //transferBatch
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//    //     BigInteger sn = transferBatch(transaction);
//
//    //     //handleBTPError for make refundable balance
//    //     handleBTPError(transaction, sn, code, msg);
//
//    //     //reclaim
//    //     ScoreIntegrationTest.balanceCheck(owner, nativeValue, () ->
//    //             bts.reclaim(nativeCoinName, nativeValue));
//    //     IRC31SupplierTest.balanceCheck(owner, coinId, coinValue, () ->
//    //             bts.reclaim(coinName, coinValue));
//    // }
//
//    // @Test
//    // void handleBTPErrorShouldRevert() {
//    //     AssertBTSException.assertUnknown(() ->
//    //             btsBSH.handleBTPError(linkNet,
//    //                     BTPTokenService.SERVICE, BigInteger.ONE, 0, ""));
//    // }
//
//    // @Test
//    // void handleFeeGatheringShouldIRC31Transfer() {//how to clear feeBalances as zero?
//    //     //mint
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     //transferBatch
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//    //     BigInteger sn = transferBatch(transaction);
//
//    //     //handleTransferResponse
//    //     handleTransferResponse(transaction, sn);
//
//    //     //
//    //     Asset[] feeAssets = feeAssets();
//    //     System.out.println(Arrays.toString(feeAssets));
//    //     BigInteger nativeFee = nativeCoinValue(feeAssets);
//    //     List<Asset> coinAssets = coinAssets(feeAssets);
//    //     BigInteger[] coinIds = coinIds(coinAssets);
//    //     BigInteger[] coinValues = coinValues(coinAssets);
//
//    //     Address faAddr = new Address(fa.account());
//    //     Executable executable = () -> ((MockBMCScoreClient) MockBMCIntegrationTest.mockBMC)
//    //             .intercallHandleFeeGathering(
//    //                     IRC31SupplierTest.transferFromBatchChecker(
//    //                             btsAddress, btsAddress, faAddr, coinIds, coinValues),
//    //                     btsAddress, fa.toString(), BTPTokenService.SERVICE);
//    //     ScoreIntegrationTest.balanceCheck(faAddr, nativeFee, () ->
//    //             IRC31SupplierTest.balanceBatchCheck(faAddr, coinIds, coinValues, executable));
//    // }
//
//    // @Test
//    // void handleFeeGatheringShouldTransferStart() {
//    //     //mint
//    //     IRC31SupplierTest.mint(owner, coinId, coinValue);
//
//    //     //transferBatch
//    //     TransferTransaction transaction = transferTransaction(owner, to, feeRatio,
//    //             new Asset(nativeCoinName, nativeValue),
//    //             new Asset(coinName, coinValue));
//    //     BigInteger sn = transferBatch(transaction);
//
//    //     //handleTransferResponse
//    //     handleTransferResponse(transaction, sn);
//
//    //     //
//    //     Asset[] feeAssets = feeAssets();
//    //     System.out.println(Arrays.toString(feeAssets));
//
//    //     TransferTransaction feeTransaction = transferTransaction(
//    //             btsAddress, linkFa, BigInteger.ZERO, feeAssets);
//    //     ((MockBMCScoreClient) MockBMCIntegrationTest.mockBMC)
//    //             .intercallHandleFeeGathering(
//    //                     transferStartEventLogChecker(feeTransaction)
//    //                             .andThen(sendMessageEventLogChecker(feeTransaction)),
//    //                     btsAddress, linkFa.toString(), BTPTokenService.SERVICE);
//    // }
//
//    // @Test
//    // void handleFeeGatheringShouldRevert() {
//    //     AssertBTSException.assertUnknown(() ->
//    //             btsBSH.handleFeeGathering(link.toString(), BTPTokenService.SERVICE));
//    // }
//}