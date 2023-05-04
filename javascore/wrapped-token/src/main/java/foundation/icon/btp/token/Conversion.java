package foundation.icon.btp.token;

import java.math.BigInteger;
import java.util.Arrays;

public class Conversion {
    static BigInteger maxUint64 = new BigInteger("18446744073709551615");

    public static byte[] bigIntToByteArray(BigInteger num) {
        byte[] bytes = num.toByteArray();

        if (bytes.length <= 8) {
            byte[] paddedArray = new byte[8];
            System.arraycopy(bytes, 0, paddedArray, 8 - bytes.length, bytes.length);
            bytes = paddedArray;

        } else {
            bytes = Arrays.copyOfRange(bytes, bytes.length - 8, bytes.length);
        }

        return bytes;
    }

    public static BigInteger byteArraytoBigInt(byte[] bytes) {
        int trailingZeroes = 0;
        for (int i = 0; i < 8; i++) {
            if (bytes[7 - i] == 0) {
                trailingZeroes++;
            } else {
                break;
            }
        }
        
        byte[] ba = new byte[9 - trailingZeroes];
        for (int i = 0; i < ba.length - 1; i++) {
            ba[i + 1] = bytes[i];
        }

        return new BigInteger(1, ba);
    }
}