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

import com.iconloop.score.test.TestBase;
import org.bouncycastle.util.encoders.Hex;
import org.junit.jupiter.api.Test;

class BTPMessageCenterTest extends TestBase  {

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
