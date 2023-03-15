package foundation.icon.btp.token;

import java.lang.reflect.Method;
import java.math.BigInteger;
import java.security.SecureRandom;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ConversionTest {

    public static BigInteger generateRandomBigInteger(BigInteger min, BigInteger max) {
        SecureRandom rnd = new SecureRandom();
        int numBits = max.bitLength();
        BigInteger result;
        
        do {
            result = new BigInteger(numBits, rnd);
        } while (result.compareTo(min) < 0 || result.compareTo(max) > 0);
        
        return result;
    }

    @Test
    public void testBigIntToByteArray() {
        BigInteger min = new BigInteger("1");
        BigInteger max = Conversion.maxUint64;
        
        for (int i = 0; i < 10; i++) {
            BigInteger randomBigInt = generateRandomBigInteger(min, max);

            byte[] bytes = Conversion.bigIntToByteArray(randomBigInt);
            BigInteger bigInt = Conversion.byteArraytoBigInt(bytes);
    
            assertEquals(8, bytes.length);
            assertEquals(randomBigInt, bigInt);

            System.out.println(randomBigInt);
            System.out.println(bigInt);
        }
    }
}