package org.emrtd.harness.jmrtd;

import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.jmrtd.BACKey;
import org.jmrtd.Util;

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.SecretKey;
import javax.crypto.spec.IvParameterSpec;
import java.security.GeneralSecurityException;
import java.security.Security;
import java.util.Arrays;

/** Dynamic BAC chip side aligned with JMRTD PassportApduService.sendMutualAuth (0.8.6). */
final class BacChipSimulator {
    private static final byte[] RND_ICC = hex("4608F91988702212");
    private static final byte[] K_ICC = hex("0B4F80323EB3191CB04970CB4052790B");
    private static final IvParameterSpec ZERO_IV = new IvParameterSpec(new byte[8]);

    static {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    private final BACKey bacKey;
    private ChipSecureMessaging session;

    BacChipSimulator(BACKey bacKey) {
        this.bacKey = bacKey;
    }

    byte[] getChallenge() {
        return Arrays.copyOf(RND_ICC, RND_ICC.length);
    }

    ChipSecureMessaging session() {
        return session;
    }

    boolean sessionEstablished() {
        return session != null;
    }

    byte[] mutualAuthResponse(byte[] cmd) throws GeneralSecurityException {
        if (cmd.length != 40) {
            throw new GeneralSecurityException("expected 40-byte EXTERNAL AUTHENTICATE payload");
        }

        byte[] keySeed = Util.computeKeySeed(
                bacKey.getDocumentNumber(),
                bacKey.getDateOfBirth(),
                bacKey.getDateOfExpiry(),
                "SHA-1",
                true);
        SecretKey kEnc = Util.deriveKey(keySeed, Util.ENC_MODE);
        SecretKey kMac = Util.deriveKey(keySeed, Util.MAC_MODE);

        byte[] eIfd = Arrays.copyOfRange(cmd, 0, 32);
        byte[] mIfd = Arrays.copyOfRange(cmd, 32, 40);

        Mac mac = Mac.getInstance("ISO9797Alg3Mac", BouncyCastleProvider.PROVIDER_NAME);
        mac.init(kMac);
        byte[] expMac = Arrays.copyOf(mac.doFinal(Util.pad(eIfd, 8)), 8);
        if (!Arrays.equals(mIfd, expMac)) {
            throw new GeneralSecurityException("MAC mismatch on EXTERNAL AUTHENTICATE");
        }

        Cipher cipher = Cipher.getInstance("DESede/CBC/NoPadding");
        cipher.init(Cipher.DECRYPT_MODE, kEnc, ZERO_IV);
        byte[] plain = cipher.doFinal(eIfd);
        byte[] rndIfd = Arrays.copyOfRange(plain, 0, 8);
        byte[] kIfd = Arrays.copyOfRange(plain, 16, 32);

        byte[] s = new byte[32];
        System.arraycopy(RND_ICC, 0, s, 0, 8);
        System.arraycopy(rndIfd, 0, s, 8, 8);
        System.arraycopy(K_ICC, 0, s, 16, 16);

        cipher.init(Cipher.ENCRYPT_MODE, kEnc, ZERO_IV);
        byte[] eIcc = cipher.doFinal(s);
        mac.init(kMac);
        byte[] mIcc = Arrays.copyOf(mac.doFinal(Util.pad(eIcc, 8)), 8);

        session = ChipSecureMessaging.establish(kIfd, K_ICC, RND_ICC, rndIfd);
        return concat(eIcc, mIcc);
    }

    private static byte[] concat(byte[] a, byte[] b) {
        byte[] out = new byte[a.length + b.length];
        System.arraycopy(a, 0, out, 0, a.length);
        System.arraycopy(b, 0, out, a.length, b.length);
        return out;
    }

    private static byte[] hex(String s) {
        int len = s.length();
        byte[] out = new byte[len / 2];
        for (int i = 0; i < len; i += 2) {
            out[i / 2] = (byte) Integer.parseInt(s.substring(i, i + 2), 16);
        }
        return out;
    }
}
