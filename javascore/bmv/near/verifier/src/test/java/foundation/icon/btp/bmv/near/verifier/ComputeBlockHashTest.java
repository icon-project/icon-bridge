/*
 * Copyright 2022 ICONLOOP Inc.
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

package foundation.icon.btp.bmv.near.verifier;

import org.junit.jupiter.api.Test;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import io.ipfs.multibase.Base58;
import foundation.icon.btp.bmv.near.verifier.types.BlockHeader;
import static org.junit.jupiter.api.Assertions.assertArrayEquals;

class ComputeBlockHashTest {
    private static final String validBlock = "f918a9a04c036d0dd5576df86b2cfe72e517bc528053cf067c7063d9ade807d8400333e2f8d48404b3a217a00b8dfd8e3ce3659aeaaf0716a3eaf9337905975540f75f48b65134b7c2c487c3a04c036d0dd5576df86b2cfe72e517bc528053cf067c7063d9ade807d8400333e2a0e38e65c8729d83e98945f353f59d7dc6b8e77837ba2d27b1b93bee6d6a837063a073f1824f31655b978d61c37aa337f47a02261631eef432cc380e4578340fc0648816ca2a06023bb194a0a4bf6111c30a84ac347a1f70721e7bfe0bd7937a14ed9a0a0c5795d695a6cd0da05920736baf29610598821d451e5d45d368d75b979e0b6e51f72645d7aec87217f9176ca06000e23e31e545c2058507e22f943c8999c413e5d2d1ce573ffc09f8fc777c49a0cb9e2f0d1777c0129c29cdbb936204bd37e1a8c1e7993e3860f2c5cbb4bf15a0a01223349a40d2ee10bd1bebb5889ef8018c8bc13359ed94b387810af96c6e4268a00000000000000000000000000000000000000000000000000000000000000000a0d0843018e762612b48d2418a6f0809928a276fbe878d17168365754cbe679e27c0c40101010189313030303030303030a232313732363839323931353034333131343530323831363336393936393137353433c0a0eb4fced6ebd34543d636cd942c5c79ab227897bb06866cf41eda1992afe88300a04c036d0dd5576df86b2cfe72e517bc528053cf067c7063d9ade807d8400333e28402157df58404b3a216a0f6497264d55add073555a79ca81debe2370e3a094c642b7dcc209bb22f57fb80f91622b841002ddba5abde4a958fd3baae592a888817a05efa328ef0e0645bc73b8c969c640cf77b8ac0eac5e102e553cf1fc3e30307ba1a6141dbf776e2c93fcb6f96376904b841002ddba5abde4a958fd3baae592a888817a05efa328ef0e0645bc73b8c969c640cf77b8ac0eac5e102e553cf1fc3e30307ba1a6141dbf776e2c93fcb6f96376904b841002ddba5abde4a958fd3baae592a888817a05efa328ef0e0645bc73b8c969c640cf77b8ac0eac5e102e553cf1fc3e30307ba1a6141dbf776e2c93fcb6f96376904b841002ddba5abde4a958fd3baae592a888817a05efa328ef0e0645bc73b8c969c640cf77b8ac0eac5e102e553cf1fc3e30307ba1a6141dbf776e2c93fcb6f96376904f800b84100b130ec92309068a81eab0c5e97d956fe9d4aab8ae11ddf6fc3b1c9127abca9edbce898c6112b4827cc9b80293069366b55cf6996e76382e9f7aa7e0da94fdb0ab84100e7a74fa1c1c9ad373ab3461774bc0553a6c531ffda14493e900fd55e7d9776ab44efd00afd16f9cc6b9c23e965de9574c3b5ed589c229d26963f5c230f11260bb8410066a56d32954eed3447a590d03a594b0a0438a875495492fa44aa3fda9688351c7b920f13ee088aa40893240b3041e0d9c1fd144eff5b863f0e31772f8c6c7107b841002b7ca1f730b7b4ad7674e1937161c07a55e0ff3d8d0a3789a776a4293d0d70b1751a52295b474b6c211870ec6d95fd4756e4f5e8148c6ab6247b7aa95d59c802f800b84100493db1ef9ffd5ffa44c9bb1f417a9a3eda22f654c2221d44870b5c44bffc95dcecfb43fde4d31b7855343bf341eda171e75c2f5ebe6d70e06f8c5562fc540504b84100d1d0be28686ceb2744352518f423c843b7cb6a61637039e49bfa6d605855397427d8bf02b6a7bbc120d349637722fe5148a1bb0539e058fb49390b2040b65d0ab84100b00d447999926789de463c81ba6947d88123f23f0aabc4a5f1a290bd0b2c05a9a8b7f4e28750c7a427c9dd8a65fbeb3e77d95de782167c554dda594c61b1d90bb841003b1d79bd357371f05a3d968900b7b0aeae68c79c05e215a1fcf3841cc6166107e4ad74ac3f8cf08405f71a43683fefb00d06ee6e1abdbd8ef9339a9645e97d0db8410000f80f74d9711dc1ebf29d8fcaf61a5c4e86a606bc2dba1ae5fc7f3a694b48516c86adfbcf00cc27d902f1a87611437437a5675219620fe164380159d1660f08b84100c18457aeb6a66286aeb2d8fc50abe94cde95a250c5f8e8ededaa72c556b7c1b426303049e61e00131977fac8034083a51f38f760149d688de9c4b6d3e9bd560ab84100cb5ead952fb596da3d44125ecf1288b7f71ca88c45f30124fd0336605246892ee6bdade1b7d62f74caf3baddde5f18dcb4ebe300cba48521ccdf96b9c803c106b841007077f32e37da8f9881cdd8ff3cfbde310943d5c3d008348e315df2674ce76ec12f6e5e404cdb593453f1f64b6b865e6a276981819a72a5c891280d3b648c9d0ab841004968caccfb5816d185165518a961438ad13a0f1d21cc2fe66697d155eda9bd7880fe6ce124001f30978bc274df44e4dfefc171fffb500719bb8cba7a8a287b0cb841001b4b5013da02f749357e9fb6ce05a2a5c6694da6ac6b413d8bed25afc4a0680b3061b5647fb5ba2ef4e4048c02233d81079091cc4e53686802681e891637130af800b841009c534d3d2e174eb7f9fe351f323af10302867f223f5a8a4d0e51f618414512327d6772ae1f2bb444170dc7079ddb839eaae9f7236d51eed5ad31e97fa5731701f800b84100ad62475584eb6f0316a011ec4ef30f60e264553f47cb295caa16f4de614e5aa5796688e3161654709e27c486ff014163a50b98d23e4904b98b05128ec550c90ab841006cc4cb383969590f7ec2a60dc244345f8bd2fb672cb436f5db90fa1f262424ab7eaad89ff1b299a2f90d0c7924d89605f65bc36ad1c792d49e630c05ead5e30ff800b84100a0970084bf96ca6eb60ec8263c968d5b6de879541ea0e4d1d2b6f74f2c1cc371e8ffb5027f4195f80a50752b645ce8d0582294fdc98519df58fb7378375d750cf800b84100400eeaf80cbbb28605bce3e4f24d341790da3b5a05a16cc63e0fb5357496e8f93df2df8a00d85863b8a3630d54b3662e5a0bd53925ba0814a54868d63e6dc902f800b84100762ff6e376bee7a071fd0358d309666ace5dedeb3452157c0b3fcba7a5926d9080cf0da81ae9057e2a41b8bd84fc78e568a207cff6bb73764bb460613ce93100b8410072097bbe233da4782c744a724a38f9db25759821ca95a61f8cac46538a499709ea98657ce2dcdb692dd46ace1b54a311871c778244a435b2158bc6a58b10ae08b8410049c513400841b44599c92947a0b342f2b111d29772f858dfdb30804bd1ad6f030588634e3fb2e7f2fd15c194c837506c428c9150a1fbb0265f6c2de290d7dc0bb84100346d0377e8849fc1199c0d0184c083dd854fbff17ab95304193ebd8f748fe950e80ec52d1cdbeba369e4605b5ccfdca1c012615af9bc32bc57e4fa042d617e0ef800b84100d79b472605b1ea93ba002c9e82fd3d9c39fd6df49c9872c0fa6c8094bb5b7fecbd4cc4b5b00910bde9ee9caeb1a616fcece77cafc7a33bff94d129622752120ab84100bca1135e3606b0fe3dd1f3bd81496bc86253daaf2b31510a4bfc7fb9b3aa8ad4ea970dc4bfe360472dfcbe8d234568750e1745c1da84dde84a0e3c3d3812760fb84100f0f2f92734ebc71d6ea3ae499f8578c2a3641ad384ac13de90eed36976d3249207aaeb99e743ebf8d2a8262c4d44f1ef26e3d38273ced905b564304f7c864705b8410002c17f9727fcd09acf04370a5e097a35593effecf1e556373bb938104a9a158a0893af80b9d2af9f0eac2312a8af57cd5f10ab83611f944144f6076918db5d09b84100e4e5cca03b102e081c87fa7f6c26028ef43f6016899e92d00e0f91b919e80793ae3e9747faf303bfb37d6489d87678c5f765ddffda84744d9ec993a99fda1006b841005b33e86eb301639f9366c808d6d5b85d2f6eb8fd97bc75884597fe201bc2ef30be6b220cec3f0883b32e3c73e43223f4feee01ddc4150b807b0d63be5aa8470cb84100c4fc02aa63a1b27315c00ae4e36fae41a772129e933f516b135669f30c44fafdb9bba394156e42e01993190275019964ceb81035feb66f29649eb17656135b03b8410020cea9574338209e54a6b71fb7fba5c49d948bdf26beed2b6f1f80f87b4a88bf85078aea4f30e35f7a6cb459ce5926e1b40938c2877535df2721718a07f53b0bb841002f204f2691a11624943cb52cae7cc75f5039163119a356025285d6505efb0ac8bf719a5b46e6f721206a5976d640e8e0eb893900983469a64c96493f47a9ee03b84100eb50a9c597299eab83874b0273fbcde58a969f5287c83fcc406a86d1323ebf799b9fb57b575f09fddba3567b519c87dd8e9a4111034b9b79f9d527d466f7df0ab84100d2140bf65ca24764cfe0e1a64e20c0fe390481c2cf18f2edeaa9f30c71137fffea9c4ad16d0354ef47eb73e56e749f6db4221d4ae80b94a8f9365e761fac7e03b841003ae1d053f25b93663e17cfd42e58aa687d4f12495e001f7f4f14cc72cf862ca38a3977e919eb6d9c70093fdba0175469e6683a3e0baa072766fcd57f635efd05b84100e68ff9d2acbae892e0f58bdbf86f3370b38d5ea2421108b60894204bdfa2b3035b8758b69ef9f0d7283afab0a57d7adf8d458dff83d34e50165ec5cdea3b5e07b84100e2015639ed395708b77ec2dc947a8bec834ae1109f878adf1089d7d1a3cda21bdb09fccb07375e1588b4025795a33205cdf8edcb2bceb724b50f4c8811975504b8410018e60e19b19e954074b17946721014bf7b9f9972725ed88f869edf75c89cfda78ea893cfa9333826a0ed74e9ea5b4f9ecd582d9d8aea0ec48bd95ffd2dcfc00df800b8410040bb05ce3b38e98b43e2d29fdee17f208f1180e8179b6a185e8cce2f65956fdd197b2c1c6d8826fa9cd04bf83e9ed232c5d2ef369e306690657f2168ae91b70cf800b84100790c95730f63d1a088a9aa946dec78d122d8831b57e973b5387a3e72c6d427403499854f01dc4cbc85fff9c2d85bf2b7e5b443fbd765ba69e613a8bc38bbb40cb8410008430cc97483ba62651122203434b22679d1aafcd295de224105befa3ce5aa587561930f4212d57e182fa4a28b46749697713931d701b509f08299f7e155a700b8410006cdacf98f06148a4b8a7fb8c765c6ecb099b9149e7585ca9be838e7aeb5be76b8e89fde1d7f97cbc5e7336c4f1d61927de07122f684ae61ce44356214859e01b841002fa3a810119c80b4cb1fde6347a4fca6c2620a9ce0ac4600127f10a179fa1f175f5ed8344169e8309bd5a630002b0434b56b6187ac82e7539265aaf9d3645f03b84100d0609460ac29c8e05000328cb48e8b03178e64b5a6c51acf560865fc110e4c6eadbf58c836c2c77dda310e0e54d083b4d9cc95683889b0a91c88706ddc8ae50cf800b84100002ffe3b708d4cfa3267edc0df1358f57b2de26dec353b00acdd251726da7c1c4104e08046c4bb867453c84e89bf9a6938beacac89aaf582bd7da76d18f25800b841003a5b29dbb239e776f7d600925faf1a7ef14cd62f936098c2392f176e6c53503f00078351fbd55c128937e8eddb8b312a3160417734bd4e7b3c245a123a4f6f06f800b84100a1bedc023e7d1f2ed2cdcef2a2ec93bdc7e8b9584d9a80a458a2588958a7e7d1ee17fd9cc0af4dfe6679b894e51d64dfd79770e1ce50c447c975dc4559025005f800b841008241c13c7234106c25c6ddf69d9436ff2775d1e60faaa188e509809963ab50d2ed4f4c3a6477f578a38db29f69f340288d3d9ba297a0702f3b96a1e6c1e0f902b84100a80e48801982a919b58452cecc7e20bfa93306e57e18cefd337ecafb92c75394c8c7a201f70ee162b0f78ebeddefb8a841a582df884660eedaf13b92b7b2df0ab84100a2421be7d968a24c6bd65a29f516c1cac0f1c4b2facab654ff34455fe7eaf94cdc2df74cab367bc29aa57fb78813885ddfe527b46650f2dce3d70c81f4498f00b841002541844bbb28ea069036d4d13ecef855714eda73eddc4710cf3b0b008039a5f372c198c5375f2bf2613f2dd5ff762ef66a42c5668e64f4ecf0ed9e6275ac1404b84100575f139878b839ea1841500328fc6cbc11dd7d4edfa3849d34b6118e498a478b5f33d5c3ebcd552d70a0b7152ef84001628196e96f1a694a5d4ce2cadc72330db841006c4ef2d0bef13fbeacb4f1a2675daf86b78240969abe3b3314a758332ee97c34d1d69256ec2cad8f05aea621f140ae96728970096193a9f599284c88fd8c7902b84100cbb45706c35b61e582d28cc0a05efc50855ef5a9611aa22820eb47bd09ed1f961fb767edc0c06fc5bfc0a478af44d5f78ef96eb982b18bfa8a8a6ede6d8f2000f800f800b8410001e5b66791f39c7be18c0ba7076cb6f170ea95dabe8f08178f1aeee7913732b8aaf7113f28a71633d60631e99deedeb289a0afe7ec8351f2cc2b039a62e0ae09f800b84100ff47910a53c33803193fdc5fc703a5c5de8368dc87d75ab653cb980d43fc2e16991578410bdeb0f4bbddd4ed66999ec3a2d58c3cceced603fe5a33fbdc00e80db8410015743f09b2be06b1854eb59427e746bb6b6788d378eadc033adb7f95710f8aa7c3d87fb7f4707005a043b17842619d1d48d6e92ce572da4eb2e50c18ed0ff404b841005e8491baa7e9d1969a2182102b2f713d1bd4ca7b74b2b542418b3c2cd63f80753d1cfd397abffd166f74ff364f4a4492e0043151bd7b9cf110d60f93087ba302b8410045225b7725035062647f46a45edbe75d620b07753b5896f57904227891c850124a7e2a5afd6ae0397276daca886f63ba0c1c5c3c066a1d90520b61b0aaf2c900b84100d8553a62b2aa2c9fc05bf50740f3359c4fc49c9ad0d41ba844c4de3c5ce58a06f44ca69c1646a6bfc3ce23ee0262b7bef711961034d6252dddf10ffe5b79a707b84100adc9e094bc4df4645e0821bcea40b3f94f499fb2a2d746bb1e397f1ef40974675e0b1208d4452d309fbe85ff1ec08ccb73cdd7d0305117b5b950afca42e03d0ab8410096a3749a7ac92ce60fd45d1a54c340e0ccc1d918e1054b9a9587d6a7006e1b9ff847c539053e4c003a36d62a692810d2ff71bc76965305b2f43096091a573b0bf800b84100bd70eb40c5df5e4629275d32e0ed393fb6b4a1b547cd0f547733f4e4fe4d7bf21c7b109603bc36a7c79d26fcb7230804ee1b0befb266684a45b8042f94985d0cf800b84100a7c48de4138f0965a41961f1bdb1477989b00b71db8bb30b712d45dfc46729d98773ea961a9f94c9ccfada0f47dd85cfa7cc48559a0fd370e3e8122377c58b09b84100849ef695e9230a9bf77b09b0375f9eaef64c751da38333567b7db5d6f3042f0f0f666ecb0a962c6508bf314f0d0b300113eaae2ad663ebafa3ab14bb9202a301b8410083e84bd138443f839d80a74f55bee3450fcf769ca02c9aa4e1f8040d6156c814451751e6e468364d844b5ade354ac98beeb12cb1c45ab2486b1f7f568b2d3d0cb8410080de01f41b748ab7eb1451ae58e1d44efd6e171fbbbc1a1569e015fc88ce87f890df6c3fc066b6991237616a78c0299667f3c3f1ad0cfdc570ac0e4e4ffdb004b841008a9574e11563be7b578647e9516dc9076f830a4438c7de5598ef3d5140e918af00eb9522c72fda189c0a0b4ba1725b6268cbb4e2f12f877c58052fdb608ff505b8410000afe917cd205757b010dbaeeeb2481a406c2154fb7432fa0f7ab39863d3476df2223e7aa969665ed443bbb74a37eeb33a297129283fad41cb156d73f89e0f04b841001e69aa64882438c75af2cf314724aa626aa46bf720ff0bdeaef7c6b3e494e910ce6962900d2c037752b196a299b602b17f796e7db1dfdc4e4f029e3441fc0000b8410021b2a428806c193d6ab624a2ee590d775a2d1989614b3e3dd6761e0927581d0e5f2118e16b0a3e8fc00b4665087d693b6a65cf1eafbc1683b0f78eecccfa9b0ab8410005e3b5bca9748cd5a148be81528fa7374749a3f0fecf8b0f00c40cc16e9e20627871ad2cc676df87d1f88d63608aa7ce9ea5d845130a139355da375374545006b84100ab2f3fd4b1411f4a24721df0fb766ec95d809ca07a3ff5e2dfc26864e14b0583774459b211cd29efcbb29f30db05a16d4b4e49d7d5ad177b2a4b5918643ae804b84100b8ebec22809c407225dbdab9fdfa0c397001375478e9acfa7f0ee347a15179b8b2c9fdfe320e929a191a64e1ac928464d71684cdc23cce232aa35288d2f7f306b841003bb506cbe049cc858a8de96228b2af38016c5dafeac9c90fa9ef43cf49767540bee7469a57934a02b1a11a42e93b12067296c4996c2197fd0acb1aaf8938b60cb8410095ffb6088671dcc7e57bd676fc82a1a56a69661fb2b4b0d5c0d13531d92eb2bf7cbf0e50e8a4b7e84800685b3ad484b1e7377243cd1222f4328b042b75105e03b84100a10b7ed86cf99c76d2f22c4271feb75e980949192534c0a1441c2ac26a7b80d71ead075f75dbaa5fc119d3603472b20f305d8f5b97e783a909ec700197b7c70bb84100637452149528d78d6b161f59744d3e0057fd92e5030e591067f8ef26629001410e84dcb373a2d12477e864e86536edc2acc1dbac0e5e4514fa645e2e616f6903b84100268973cfe5b5dcadb025dcd03a36e6068c8a058e678a257740c45cbe05b48b812952fc3b07d491b65120f40069a49c4bbab9c458a60d6cad77a5afa6ec2a750ab84100a3f8d390ef24c8a975decb50440082086d8adfbaaf8027443b154aac791301a0cfb1e89ad12a2356c3566eea8f81c151f1a4348b62e09d83ce8df6715a835105f80032b84100cd0c7c95ef5bdffc8b4e29b472aa713803c86bd12c92ee0a6b3bf0a4af2a7f14323a0c5b5e18436935e48fc28a3cc5342713c556037d4d157e5cef7583836803";
    
    @Test
    void validBlockHeaderWithPreviousHashProvided() throws DecoderException {
        byte[] bytes = Hex.decodeHex(validBlock.toCharArray());
        try {
            BlockHeader blockHeader = BlockHeader.fromBytes(bytes);

            assertArrayEquals(Base58.decode("6gvUukWemPD9tmNoCP4UXoKaxA6uADmv2u9WH6jV7RRD"),
                    blockHeader.hash(Base58.decode("67j1sSg3QrhmYGkMKhnvr64Sb5ww2PSTSM5Z99VwLawj")));

        } catch (Exception e) {
            e.printStackTrace();
            System.out.println(e);
        }

    }

    @Test
    void validBlockHeaderWithoutPreviousHashProvided() throws DecoderException {
        byte[] bytes = Hex.decodeHex(validBlock.toCharArray());
        try {
            BlockHeader blockHeader = BlockHeader.fromBytes(bytes);

            assertArrayEquals(Base58.decode("6gvUukWemPD9tmNoCP4UXoKaxA6uADmv2u9WH6jV7RRD"),
                    blockHeader.hash());

                    System.out.println(blockHeader.toString());
        } catch (Exception e) {
            e.printStackTrace();
            System.out.println(e);
        }

    }
}
