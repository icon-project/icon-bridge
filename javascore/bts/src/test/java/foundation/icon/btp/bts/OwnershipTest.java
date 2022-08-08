// /*
//  * Copyright 2021 ICON Foundation
//  *
//  * Licensed under the Apache License, Version 2.0 (the "License");
//  * you may not use this file except in compliance with the License.
//  * You may obtain a copy of the License at
//  *
//  *     http://www.apache.org/licenses/LICENSE-2.0
//  *
//  * Unless required by applicable law or agreed to in writing, software
//  * distributed under the License is distributed on an "AS IS" BASIS,
//  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  * See the License for the specific language governing permissions and
//  * limitations under the License.
//  */

// package foundation.icon.btp.bts;

// import foundation.icon.jsonrpc.Address;
// import foundation.icon.score.test.ScoreIntegrationTest;
// import org.junit.jupiter.api.TestInfo;

// import java.math.BigInteger;

// import static org.junit.jupiter.api.Assertions.assertFalse;
// import static org.junit.jupiter.api.Assertions.assertTrue;

// public class OwnershipTest implements BTSIntegrationTest {
//     static Address address = ScoreIntegrationTest.Faker.address(Address.Type.EOA);
//     static String string = "";
//     static BigInteger bigInteger = BigInteger.ONE;

//     static boolean isExistsOwner(Address address) {
//         return btsOwnerManager.isOwner(address);
//     }

//     static void addOwner(Address address) {
//         btsOwnerManager.addOwner(address);
//         assertTrue(isExistsOwner(address));
//     }

//     static void removeOwner(Address address) {
//         btsOwnerManager.removeOwner(address);
//         assertFalse(isExistsOwner(address));
//     }

//     static void clearOwner(Address address) {
//         if (isExistsOwner(address)) {
//             System.out.println("clear owner address:"+address);
//             removeOwner(address);
//         }
//     }

//     @Override
//     public void clearIfExists(TestInfo testInfo) {
//         String testMethod = testInfo.getTestMethod().orElseThrow().getName();
//         if (!testMethod.endsWith("RevertUnauthorized")) {
//             clearOwner(address);
//         }
//     }

//     // @Test
//     // void addOwnerShouldSuccess() {
//     //     addOwner(address);
//     // }

//     // static void assertAlreadyExists(Executable executable) {
//     //     AssertBTSException.assertUnknown(executable);
//     // }

//     // @Test
//     // void addOwnerShouldRevertAlreadyExists() {
//     //     addOwner(address);

//     //     assertAlreadyExists(() -> addOwner(address));
//     // }

//     // @Test
//     // void removeOwnerShouldSuccess() {
//     //     addOwner(address);

//     //     removeOwner(address);
//     // }

//     // static void assertNotExists(Executable executable) {
//     //     AssertBTSException.assertUnknown(executable);
//     // }

//     // @Test
//     // void removeOwnerShouldRevertNotExists() {
//     //     assertNotExists(() -> removeOwner(address));
//     // }

//     // static void assertUnauthorized(Executable executable) {
//     //     AssertBTSException.assertUnauthorized(executable);
//     // }

//     // @Test
//     // void addOwnerShouldRevertUnauthorized() {
//     //     assertUnauthorized(() -> btsOwnerManagerWithTester.addOwner(address));
//     // }

//     // @Test
//     // void removeOwnerShouldRevertUnauthorized() {
//     //     assertUnauthorized(() -> btsOwnerManagerWithTester.removeOwner(address));
//     // }

//     // @Test
//     // void registerShouldRevertUnauthorized() {
//     //     assertUnauthorized(() -> btsWithTester.register(string));
//     // }

//     // @Test
//     // void setFeeRateShouldRevertUnauthorized() {
//     //     assertUnauthorized(() -> btsWithTester.setFeeRatio(bigInteger));
//     // }

// }
