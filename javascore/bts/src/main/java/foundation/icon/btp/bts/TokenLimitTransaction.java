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

import java.math.BigInteger;
import java.util.List;
import score.ByteArrayObjectWriter;
import score.Context;
import score.ObjectReader;
import score.ObjectWriter;
import scorex.util.ArrayList;

public class TokenLimitTransaction {
    private String[] coinName;
    private BigInteger[] tokenLimit;
    private String net;

    public TokenLimitTransaction(){}

    public TokenLimitTransaction(String[] coinName, BigInteger[] tokenLimit, String net) {
        this.coinName = coinName;
        this.tokenLimit = tokenLimit;
        this.net = net;
    }
    public TokenLimitTransaction(String[] coinName, BigInteger[] tokenLimit) {
        this.coinName = coinName;
        this.tokenLimit = tokenLimit;
        this.net = "ALL";
    }

    public String[] getCoinName() {
        return coinName;
    }

    public BigInteger[] getTokenLimit() {
        return tokenLimit;
    }

    public String getNet() {
        return net;
    }

    public void setCoinName(String[] coinName) {
        this.coinName = coinName;
    }

    public void setTokenLimit(BigInteger[] tokenLimit) {
        this.tokenLimit = tokenLimit;
    }

    public void setNet(String net) {
        this.net = net;
    }

    public static TokenLimitTransaction readObject(ObjectReader reader) {
        TokenLimitTransaction obj = new TokenLimitTransaction();
        reader.beginList();
        if (reader.beginNullableList()) {
            String[] coinNames = null;
            List<String> coinNamesList = new ArrayList<>();
            while (reader.hasNext()) {
                coinNamesList.add(reader.readNullable(String.class));
            }
            coinNames = new String[coinNamesList.size()];
            for (int i = 0; i < coinNames.length; i++) {
                coinNames[i] = coinNamesList.get(i);
            }
            obj.setCoinName(coinNames);
            reader.end();
        }

        if (reader.beginNullableList()) {
            BigInteger[] tokenLimits = null;
            List<BigInteger> tokenLimitList = new ArrayList<>();
            while (reader.hasNext()) {
                tokenLimitList.add(reader.readNullable(BigInteger.class));
            }
            tokenLimits = new BigInteger[tokenLimitList.size()];
            for (int i = 0; i < tokenLimits.length; i++) {
                tokenLimits[i] = tokenLimitList.get(i);
            }
            obj.setTokenLimit(tokenLimits);
            reader.end();
        }

        obj.setNet(reader.readNullable(String.class));
        reader.end();
        return obj;
    }

    public static void writeObject(ObjectWriter writer, TokenLimitTransaction obj) {
        obj.writeObject(writer);
    }

    public void writeObject(ObjectWriter writer) {
        writer.beginList(2);
        String[] coinNames = this.getCoinName();
        if (coinNames != null) {
            writer.beginNullableList(coinNames.length);
            for(String s : coinNames) {
                writer.writeNullable(s);
            }
            writer.end();
        } else {
            writer.writeNull();
        }

        BigInteger[] tokenLimits = this.getTokenLimit();
        if (tokenLimits != null) {
            writer.beginNullableList(tokenLimits.length);
            for(BigInteger s : tokenLimits) {
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
        TokenLimitTransaction.writeObject(writer, this);
        return writer.toByteArray();
    }
}