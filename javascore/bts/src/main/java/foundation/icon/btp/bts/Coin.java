package foundation.icon.btp.bts;

import score.Address;
import score.ObjectReader;
import score.ObjectWriter;

import java.math.BigInteger;

public class Coin {
    private String name;
    private String symbol;
    private int decimals;
    private BigInteger feeNumerator;
    private BigInteger fixedFee;
    private int coinType;
    private Address address;

    public Coin(Address address, String name, String symbol, int decimals, BigInteger feeNumerator, BigInteger fixedFee, int coinType) {
        this.address = address;
        this.name = name;
        this.symbol = symbol;
        this.decimals = decimals;
        this.feeNumerator = feeNumerator;
        this.fixedFee = fixedFee;
        this.coinType = coinType;
    }


    public static void writeObject(ObjectWriter w, Coin v) {
        w.beginList(6);
        w.write(v.getAddress());
        w.write(v.getName());
        w.write(v.getSymbol());
        w.write(v.getDecimals());
        w.write(v.getFeeNumerator());
        w.write(v.getFixedFee());
        w.write(v.getCoinType());
        w.end();
    }

    public static Coin readObject(ObjectReader r) {
        r.beginList();
        Coin result = new Coin(r.readAddress(), r.readString(), r.readString(), r.readInt(), r.readBigInteger(), r.readBigInteger(), r.readInt());
        r.end();
        return result;
    }

    public Address getAddress() {
        return address;
    }

    public void setAddress(Address address) {
        this.address = address;
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

    public int getDecimals() {
        return decimals;
    }

    public void setDecimals(int decimals) {
        this.decimals = decimals;
    }

    public BigInteger getFeeNumerator() {
        return feeNumerator;
    }

    public void setFeeNumerator(BigInteger feeNumerator) {
        this.feeNumerator = feeNumerator;
    }

    public BigInteger getFixedFee() {
        return fixedFee;
    }

    public void setFixedFee(BigInteger fixedFee) {
        this.fixedFee = fixedFee;
    }

    public int getCoinType() {
        return coinType;
    }

    public void setCoinType(int coinType) {
        this.coinType = coinType;
    }
}
