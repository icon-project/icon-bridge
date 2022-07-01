package foundation.icon.btp.nativecoin.irc2;

import foundation.icon.btp.irc2.IRC2;
import score.Address;
import score.Context;
import score.annotation.Optional;

import java.math.BigInteger;

public class IRC2ScoreInterface implements IRC2 {

    protected final Address address;

    public IRC2ScoreInterface(Address address) {
        this.address = address;
    }

    public Address _address() {
        return this.address;
    }

    public String name() {
        return Context.call(String.class, this.address, "name");
    }

    public String symbol() {
        return Context.call(String.class, this.address, "symbol");
    }

    public int decimals() {
        return Context.call(Integer.class, this.address, "symbol");
    }

    public BigInteger totalSupply() {
        return Context.call(BigInteger.class, this.address, "totalSupply");
    }

    public BigInteger balanceOf(Address _owner) {
        return Context.call(BigInteger.class, this.address, "balanceOf", _owner);
    }

    public void transfer(Address _to, BigInteger _value, @Optional byte[] _data) {
        Context.call(this.address, "transfer", _to, _value, _data);
    }

    /**
     * @deprecated Do not use this method, this is generated only for preventing compile error. not supported EventLog method
     * @throws RuntimeException
     */
    @Deprecated
    public void Transfer(Address _from, Address _to, BigInteger _value, byte[] _data) {
        throw new RuntimeException("not supported EventLog method");
    }

}
