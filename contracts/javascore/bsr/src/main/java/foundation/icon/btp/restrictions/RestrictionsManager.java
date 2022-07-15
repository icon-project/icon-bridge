/*
 * Copyright 2022 ICON Foundation
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

package foundation.icon.btp.restrictions;

import score.Address;
import score.annotation.External;

import java.math.BigInteger;

public interface RestrictionsManager {
    /**
     * @param _name    name of the token
     * @param _symbol  symbol of the token
     * @param _address Address of the token contract
     * @param _limit   Max transaction limit
     */
    @External
    void registerTokenLimit(String _name, String _symbol, Address _address, BigInteger _limit);

    /**
     * @param _address Address of the user to blacklist
     */
    @External
    public void addBlacklistedUser(String _address);

    /**
     * @param _address address of an user to remove from blacklist
     */
    @External
    public void removeBlacklistedUser(String _address);

    /**
     * @param _address address of an user to check if blacklististed
     * @return boolean true if the user is blacklisted, else false
     */
    @External
    public boolean isUserBlackListed(String _address);

    /**
     * @param _token token/coin name
     * @param _from  Address String  of account to be transferred value from
     * @param _to    Address String of the account to which the value will be transferred to
     * @param _value Amount of tokens/coins to be transferred
     */
    @External
    public void validateRestriction(String _token, String _from, String _to, BigInteger _value);
}
