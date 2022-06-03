package foundation.icon.btp.restrictions;

import score.ObjectReader;
import score.ObjectWriter;

import java.math.BigInteger;

public class TokenLimit {
    private String name;
    private String symbol;
    private BigInteger limit;

    public TokenLimit(String name, String symbol, BigInteger limit) {
        this.name = name;
        this.symbol = symbol;
        this.limit = limit;
    }


    public static void writeObject(ObjectWriter w, TokenLimit v) {
        w.beginList(4);
        w.write(v.getName());
        w.write(v.getSymbol());
        w.write(v.getLimit());
        w.end();
    }

    public static TokenLimit readObject(ObjectReader r) {
        r.beginList();
        TokenLimit result = new TokenLimit(r.readString(), r.readString(), r.readBigInteger());
        r.end();
        return result;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getSymbol() {
        return symbol;
    }

    public void setSymbol(String symbol) {
        this.symbol = symbol;
    }

    public BigInteger getLimit() {
        return limit;
    }

    public void setLimit(BigInteger limit) {
        this.limit = limit;
    }
}
