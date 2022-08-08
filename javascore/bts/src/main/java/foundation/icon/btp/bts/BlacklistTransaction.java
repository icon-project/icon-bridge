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

package foundation.icon.btp.bts;

import java.util.List;
import score.ByteArrayObjectWriter;
import score.Context;
import score.ObjectReader;
import score.ObjectWriter;
import scorex.util.ArrayList;

public class BlacklistTransaction {
    private String[] address;
    private String net;

    public String[] getAddress() {
        return address;
    }

    public String getNet() {
        return net;
    }

    public void setAddress(String[] address) {
        this.address = address;
    }

    public void setNet(String net) {
        this.net = net;
    }

    public BlacklistTransaction() {}

    public BlacklistTransaction(String[] address, String net) {
        this.address = address;
        this.net = net;
    }

    public static BlacklistTransaction readObject(ObjectReader reader) {
        BlacklistTransaction obj = new BlacklistTransaction();
        reader.beginList();
        if (reader.beginNullableList()) {
            String[] addreses = null;
            List<String> addressList = new ArrayList<>();
            while (reader.hasNext()) {
                addressList.add(reader.readNullable(String.class));
            }
            addreses = new String[addressList.size()];
            for (int i = 0; i < addreses.length; i++) {
                addreses[i] = addressList.get(i);
            }
            obj.setAddress(addreses);
            reader.end();
        }
        obj.setNet(reader.readNullable(String.class));
        reader.end();
        return obj;
    }

    public static void writeObject(ObjectWriter writer, BlacklistTransaction obj) {
        obj.writeObject(writer);
    }

    public void writeObject(ObjectWriter writer) {
        writer.beginList(2);

        String[] addresses = this.getAddress();
        if (addresses != null) {
            writer.beginNullableList(addresses.length);
            for(String s : addresses) {
                writer.writeNullable(s);
            }
            writer.end();
        } else {
            writer.writeNull();
        }
        writer.writeNullable(this.getNet());
        writer.end();
    }

    public byte[] toBytes() {
        ByteArrayObjectWriter writer = Context.newByteArrayObjectWriter("RLPn");
        BlacklistTransaction.writeObject(writer, this);
        return writer.toByteArray();
    }
}