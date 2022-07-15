package foundation.icon.btp.restrictions;

import score.Address;
import score.Context;

import java.math.BigInteger;

public final class RestrictionsScoreInterface implements RestrictionsManager {

    private final Address scoreAddress;

    public RestrictionsScoreInterface(Address address) {
        this.scoreAddress = address;
    }

    @Override
    public void registerTokenLimit(String _name, String _symbol, Address _address, BigInteger _limit) {
        Context.call(this.scoreAddress, "registerTokenLimit", _name, _symbol, _address, _limit);
    }

    @Override
    public void addBlacklistedUser(String _address) {
        Context.call(this.scoreAddress, "addBlacklistedUser", _address);
    }

    @Override
    public void removeBlacklistedUser(String _address) {
        Context.call(this.scoreAddress, "removeBlacklistedUser", _address);
    }

    @Override
    public boolean isUserBlackListed(String _address) {
        return Context.call(Boolean.class, this.scoreAddress, "removeBlacklistedUser", _address);
    }

    @Override
    public void validateRestriction(String _token, String _from, String _to, BigInteger _value) {
        Context.call(this.scoreAddress, "validateRestriction", _token, _from, _to, _value);
    }
}
