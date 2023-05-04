package foundation.icon.btp.token;

import com.iconloop.score.token.irc2.IRC2Basic;

import java.math.BigInteger;

import score.Context;

public class TestToken extends IRC2Basic {
    public TestToken (String _name, String _symbol, int _decimals, BigInteger _amount) {
        super(_name, _symbol, _decimals);
        _mint(Context.getCaller(), _amount);
    }
}