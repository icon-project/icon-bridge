package foundation.icon.btp.restrictions;

import foundation.icon.btp.lib.OwnerManager;
import foundation.icon.btp.lib.OwnerManagerImpl;
import score.*;
import score.annotation.External;

import java.math.BigInteger;

public class Restrictions implements RestrictionsManager {
    private final DictDB<String, Address> tokenAddrDb = Context.newDictDB("token_addr", Address.class);
    private final ArrayDB<String> tokenNameDb = Context.newArrayDB("token_name", String.class);
    private final DictDB<Address, TokenLimit> tokenDb = Context.newDictDB("tokens", TokenLimit.class);
    private final OwnerManager ownerManager = new OwnerManagerImpl("owners");
    private final VarDB<Integer> noOfBlacklistedUsers = Context.newVarDB("noOfBlacklistedUsers", Integer.class);
    DictDB<String, Boolean> blacklistDb = Context.newDictDB("blacklistedUsers", Boolean.class);

    public Restrictions() {
        noOfBlacklistedUsers.set(0);
    }

    /**
     * @param _name    name of the token
     * @param _symbol  symbol of the token
     * @param _address Address of the token contract
     * @param _limit   Max transaction limit
     */
    @External
    public void registerTokenLimit(String _name, String _symbol, Address _address, BigInteger _limit) {
        requireOwnerAccess();
        tokenAddrDb.set(_name, _address);
        tokenNameDb.add(_name);
        tokenDb.set(_address, new TokenLimit(_name, _symbol, _limit));
    }

    /**
     * @param _address Address of the user to blacklist
     */
    @External
    public void addBlacklistedUser(String _address) {
        requireOwnerAccess();
        blacklistDb.set(_address.trim().toLowerCase(), true);
        noOfBlacklistedUsers.set(noOfBlacklistedUsers.get() + 1);
    }

    /**
     * @param _address address of an user to remove from blacklist
     */
    @External
    public void removeBlacklistedUser(String _address) {
        requireOwnerAccess();
        blacklistDb.set(_address, false);
        noOfBlacklistedUsers.set(noOfBlacklistedUsers.get() - 1);
    }

    /**
     * @param _address address of an user to check if blacklististed
     */
    @External(readonly = true)
    public boolean isUserBlackListed(String _address) {
        if (blacklistDb.get(_address.trim().toLowerCase()) != null) {
            return blacklistDb.get(_address.trim().toLowerCase());
        }
        return false;
    }


    /**
     * @param _token token/coin name
     * @param _from  Address String  of account to be transferred value from
     * @param _to    Address String of the account to which the value will be transferred to
     * @param _value Amount of tokens/coins to be transferred
     */
    @External
    public void validateRestriction(String _token, String _from, String _to, BigInteger _value) {
        if (isUserBlackListed(_from)) {
            throw RestrictionsException.restricted("_from user is Blacklisted");
        }
        if (isUserBlackListed(_to)) {
            throw RestrictionsException.restricted("_to user is Blacklisted");
        }
        Address token = tokenAddrDb.get(_token);
        if (token != null) {
            TokenLimit _tokenLimit = tokenDb.get(token);
            if (_tokenLimit.getLimit().compareTo(_value) < 0) {
                throw RestrictionsException.reverted("Transfer amount exceeds the transaction limit");
            }
        }
    }

    /* Delegate OwnerManager */
    private void requireOwnerAccess() {
        if (!ownerManager.isOwner(Context.getCaller())) {
            throw RestrictionsException.unauthorized("require owner access");
        }
    }

    @External
    public void addOwner(Address _addr) {
        try {
            ownerManager.addOwner(_addr);
        } catch (IllegalStateException e) {
            throw RestrictionsException.unauthorized(e.getMessage());
        } catch (IllegalArgumentException e) {
            throw RestrictionsException.unknown(e.getMessage());
        }
    }

    @External
    public void removeOwner(Address _addr) {
        try {
            ownerManager.removeOwner(_addr);
        } catch (IllegalStateException e) {
            throw RestrictionsException.unauthorized(e.getMessage());
        } catch (IllegalArgumentException e) {
            throw RestrictionsException.unknown(e.getMessage());
        }
    }

    @External(readonly = true)
    public Address[] getOwners() {
        return ownerManager.getOwners();
    }

    @External(readonly = true)
    public boolean isOwner(Address _addr) {
        return ownerManager.isOwner(_addr);
    }
}