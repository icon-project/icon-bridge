package foundation.icon.btp.restrictions;

import foundation.icon.btp.lib.BTPException;

public class RestrictionsException extends BTPException.BSH {

    public RestrictionsException(Code c) {
        super(c, c.name());
    }

    public RestrictionsException(Code c, String message) {
        super(c, message);
    }

    public static RestrictionsException unknown(String message) {
        return new RestrictionsException(Code.Unknown, message);
    }

    public static RestrictionsException unauthorized() {
        return new RestrictionsException(Code.Unauthorized);
    }

    public static RestrictionsException unauthorized(String message) {
        return new RestrictionsException(Code.Unauthorized, message);
    }

    public static RestrictionsException restricted() {
        return new RestrictionsException(Code.Restricted);
    }

    public static RestrictionsException restricted(String message) {
        return new RestrictionsException(Code.Restricted, message);
    }

    public static RestrictionsException failed() {
        return new RestrictionsException(Code.Failed);
    }

    public static RestrictionsException failed(String message) {
        return new RestrictionsException(Code.Failed, message);
    }

    public static RestrictionsException reverted() {
        return new RestrictionsException(Code.Reverted);
    }

    public static RestrictionsException reverted(String message) {
        return new RestrictionsException(Code.Reverted, message);
    }

    //BTPException.BSH => 40 ~ 54
    // BlacklistException => 50 ~ 54
    public enum Code implements Coded {
        Unknown(10),
        Unauthorized(11),
        Restricted(12),
        Reverted(13),
        Failed(14);


        final int code;

        Code(int code) {
            this.code = code;
        }

        @Override
        public int code() {
            return code;
        }

    }
}