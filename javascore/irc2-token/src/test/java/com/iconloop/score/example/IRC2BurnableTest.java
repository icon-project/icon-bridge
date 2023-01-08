/*
 * Copyright 2020 ICONLOOP Inc.
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

package com.iconloop.score.example;

import com.iconloop.score.test.Account;
import com.iconloop.score.test.Score;
import com.iconloop.score.test.ServiceManager;
import com.iconloop.score.test.TestBase;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import score.Address;

import java.math.BigInteger;

import static java.math.BigInteger.TEN;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;

public class IRC2BurnableTest extends TestBase {
    private static final String name = "MyIRC2Burnable";
    private static final String symbol = "MIB";
    private static final int decimals = 18;
    private static final BigInteger initialSupply = BigInteger.valueOf(1000);

    private static BigInteger totalSupply = initialSupply.multiply(TEN.pow(decimals));
    private static final ServiceManager sm = getServiceManager();
    private static final Account owner = sm.createAccount();
    private static Score tokenScore;
    private static IRC2BurnableToken tokenSpy;

    @BeforeAll
    public static void setup() throws Exception {
        tokenScore = sm.deploy(owner, IRC2BurnableToken.class,
                name, symbol, decimals, initialSupply);
        owner.addBalance(symbol, totalSupply);

        // setup spy object against the tokenScore object
        tokenSpy = (IRC2BurnableToken) spy(tokenScore.getInstance());
        tokenScore.setInstance(tokenSpy);
    }

    @Test
    void burn() {
        final Address zeroAddress = new Address(new byte[Address.LENGTH]);
        Account alice = sm.createAccount();
        alice.addBalance(symbol, transferToken(owner, alice, TEN));
        assertEquals(totalSupply, tokenScore.call("totalSupply"));

        // burn one token
        BigInteger amount = TEN.pow(decimals);
        tokenScore.invoke(alice, "burn", amount);
        alice.subtractBalance(symbol, amount);
        totalSupply = totalSupply.subtract(amount);
        assertEquals(alice.getBalance(symbol), tokenScore.call("balanceOf", alice.getAddress()));
        assertEquals(totalSupply, tokenScore.call("totalSupply"));
        verify(tokenSpy).Transfer(alice.getAddress(), zeroAddress, amount, "burn".getBytes());

        // burn all the tokens
        amount = (BigInteger) tokenScore.call("balanceOf", alice.getAddress());
        tokenScore.invoke(alice, "burn", amount);
        totalSupply = totalSupply.subtract(amount);
        assertEquals(BigInteger.ZERO, tokenScore.call("balanceOf", alice.getAddress()));
        assertEquals(totalSupply, tokenScore.call("totalSupply"));
        verify(tokenSpy).Transfer(alice.getAddress(), zeroAddress, amount, "burn".getBytes());
    }

    BigInteger transferToken(Account from, Account to, BigInteger tokenAmount) {
        BigInteger value = TEN.pow(decimals).multiply(tokenAmount);
        tokenScore.invoke(from, "transfer", to.getAddress(), value, "data".getBytes());
        from.subtractBalance(symbol, value);
        assertEquals(from.getBalance(symbol),
                tokenScore.call("balanceOf", from.getAddress()));
        assertEquals(value,
                tokenScore.call("balanceOf", to.getAddress()));
        return value;
    }
}
