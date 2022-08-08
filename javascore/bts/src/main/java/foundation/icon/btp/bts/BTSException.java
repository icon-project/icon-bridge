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

import foundation.icon.btp.lib.BTPException;

public class BTSException extends BTPException.BSH {

    public BTSException(Code c) {
        super(c, c.name());
    }

    public BTSException(Code c, String message) {
        super(c, message);
    }

    public static BTSException unknown(String message) {
        return new BTSException(Code.Unknown, message);
    }

    public static BTSException unauthorized() {
        return new BTSException(Code.Unauthorized);
    }
    public static BTSException unauthorized(String message) {
        return new BTSException(Code.Unauthorized, message);
    }

    public static BTSException irc31Failure() {
        return new BTSException(Code.IRC31Failure);
    }
    public static BTSException irc31Failure(String message) {
        return new BTSException(Code.IRC31Failure, message);
    }

    public static BTSException irc31Reverted() {
        return new BTSException(Code.IRC31Reverted);
    }
    public static BTSException irc31Reverted(String message) {
        return new BTSException(Code.IRC31Reverted, message);
    }

    public static BTSException restricted() {
        return new BTSException(Code.Restricted);
    }
    public static BTSException restricted(String message) {
        return new BTSException(Code.Restricted, message);
    }


    //BTPException.BSH => 40 ~ 54
    public enum Code implements BTPException.Coded{
        Unknown(0),
        Unauthorized(1),
        IRC31Failure(2),
        IRC31Reverted(3),
        Restricted(4);

        final int code;
        Code(int code){ this.code = code; }

        @Override
        public int code() { return code; }

    }
}