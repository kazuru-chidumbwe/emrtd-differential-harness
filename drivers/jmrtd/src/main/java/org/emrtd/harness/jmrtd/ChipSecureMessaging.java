package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CommandAPDU;
import net.sf.scuba.smartcards.ResponseAPDU;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.jmrtd.Util;

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.SecretKey;
import javax.crypto.spec.IvParameterSpec;
import java.nio.ByteBuffer;
import java.security.GeneralSecurityException;
import java.security.Security;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * Chip-side BAC secure messaging (TDES), paired with JMRTD DESedeSecureMessagingWrapper.
 * SSC increments once per inbound SM command and once per outbound SM response.
 */
final class ChipSecureMessaging {
    private static final IvParameterSpec ZERO_IV = new IvParameterSpec(new byte[8]);

    static {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    private final SecretKey ksEnc;
    private final SecretKey ksMac;
    private long ssc;

    ChipSecureMessaging(SecretKey ksEnc, SecretKey ksMac, long ssc) {
        this.ksEnc = ksEnc;
        this.ksMac = ksMac;
        this.ssc = ssc;
    }

    static ChipSecureMessaging establish(byte[] kIfd, byte[] kIcc, byte[] rndIcc, byte[] rndIfd)
            throws GeneralSecurityException {
        byte[] keySeed = new byte[16];
        for (int i = 0; i < 16; i++) {
            keySeed[i] = (byte) (kIfd[i] ^ kIcc[i]);
        }
        SecretKey ksEnc = Util.deriveKey(keySeed, Util.ENC_MODE);
        SecretKey ksMac = Util.deriveKey(keySeed, Util.MAC_MODE);
        // ICAO / JMRTD: SSC = RND.ICC[4:8] || RND.IFD[4:8] as big-endian long.
        long ssc = ByteBuffer.wrap(concatPair(
                Arrays.copyOfRange(rndIcc, 4, 8),
                Arrays.copyOfRange(rndIfd, 4, 8))).getLong();
        return new ChipSecureMessaging(ksEnc, ksMac, ssc);
    }

    private static byte[] concatPair(byte[] a, byte[] b) {
        byte[] out = new byte[a.length + b.length];
        System.arraycopy(a, 0, out, 0, a.length);
        System.arraycopy(b, 0, out, a.length, b.length);
        return out;
    }

    CommandAPDU unwrapCommand(CommandAPDU smCapdu) throws GeneralSecurityException {
        ssc++;
        byte[] data = smCapdu.getData();
        List<Tlv> tlvs = parseTlvs(data);
        byte[] actMac = null;
        List<Tlv> macTlvs = new ArrayList<>();
        for (Tlv t : tlvs) {
            if (t.tag == 0x8E) {
                actMac = t.value;
            } else {
                macTlvs.add(t);
            }
        }
        if (actMac == null) {
            throw new GeneralSecurityException("missing DO'8E'");
        }
        byte[] header = new byte[] {(byte) 0x0C, (byte) smCapdu.getINS(), (byte) smCapdu.getP1(), (byte) smCapdu.getP2()};
        byte[] macData = concat(sscBytes(), Util.pad(header, 8), encodeTlvs(macTlvs));
        byte[] expMac = mac8(Util.pad(macData, 8));
        if (!Arrays.equals(expMac, actMac)) {
            throw new GeneralSecurityException("SM command MAC mismatch");
        }
        byte[] plainData = null;
        int ne = smCapdu.getNe();
        for (Tlv t : macTlvs) {
            if (t.tag == 0x87 || t.tag == 0x85) {
                if (t.value.length < 1 || t.value[0] != 0x01) {
                    throw new GeneralSecurityException("DO'87/85' version");
                }
                Cipher cipher = Cipher.getInstance("DESede/CBC/NoPadding");
                cipher.init(Cipher.DECRYPT_MODE, ksEnc, ZERO_IV);
                plainData = Util.unpad(cipher.doFinal(Arrays.copyOfRange(t.value, 1, t.value.length)));
            } else if (t.tag == 0x97) {
                if (t.value.length == 1) {
                    ne = (t.value[0] & 0xFF) == 0 ? 256 : (t.value[0] & 0xFF);
                } else if (t.value.length == 2) {
                    int v = ((t.value[0] & 0xFF) << 8) | (t.value[1] & 0xFF);
                    ne = v == 0 ? 65536 : v;
                }
            }
        }
        return new CommandAPDU(0x00, smCapdu.getINS(), smCapdu.getP1(), smCapdu.getP2(),
                plainData == null ? new byte[0] : plainData, ne);
    }

    ResponseAPDU wrapResponse(int sw, byte[] data) throws GeneralSecurityException {
        ssc++;
        List<Tlv> nodes = new ArrayList<>();
        if (data != null && data.length > 0) {
            Cipher cipher = Cipher.getInstance("DESede/CBC/NoPadding");
            cipher.init(Cipher.ENCRYPT_MODE, ksEnc, ZERO_IV);
            byte[] enc = cipher.doFinal(Util.pad(data, 8));
            byte[] val = new byte[1 + enc.length];
            val[0] = 0x01;
            System.arraycopy(enc, 0, val, 1, enc.length);
            nodes.add(new Tlv(0x87, val));
        }
        nodes.add(new Tlv(0x99, new byte[] {(byte) ((sw >> 8) & 0xFF), (byte) (sw & 0xFF)}));
        byte[] macPayload = concat(sscBytes(),
                findEncode(nodes, 0x85), findEncode(nodes, 0x87), findEncode(nodes, 0x99));
        byte[] mac = mac8(Util.pad(macPayload, 8));
        nodes.add(new Tlv(0x8E, mac));
        byte[] body = encodeTlvs(nodes);
        byte[] out = new byte[body.length + 2];
        System.arraycopy(body, 0, out, 0, body.length);
        out[out.length - 2] = (byte) ((sw >> 8) & 0xFF);
        out[out.length - 1] = (byte) (sw & 0xFF);
        return new ResponseAPDU(out);
    }

    private byte[] sscBytes() {
        ByteBuffer bb = ByteBuffer.allocate(8);
        bb.putLong(ssc);
        return bb.array();
    }

    private byte[] mac8(byte[] data) throws GeneralSecurityException {
        Mac mac = Mac.getInstance("ISO9797Alg3Mac", BouncyCastleProvider.PROVIDER_NAME);
        mac.init(ksMac);
        return Arrays.copyOf(mac.doFinal(data), 8);
    }

    private static byte[] findEncode(List<Tlv> nodes, int tag) {
        for (Tlv t : nodes) {
            if (t.tag == tag) {
                return t.encode();
            }
        }
        return new byte[0];
    }

    private static byte[] concat(byte[]... parts) {
        int n = 0;
        for (byte[] p : parts) {
            n += p.length;
        }
        byte[] out = new byte[n];
        int o = 0;
        for (byte[] p : parts) {
            System.arraycopy(p, 0, out, o, p.length);
            o += p.length;
        }
        return out;
    }

    private static List<Tlv> parseTlvs(byte[] data) {
        List<Tlv> out = new ArrayList<>();
        int i = 0;
        while (i + 2 <= data.length) {
            int tag = data[i++] & 0xFF;
            int len = data[i++] & 0xFF;
            if (len == 0x81) {
                len = data[i++] & 0xFF;
            } else if (len == 0x82) {
                len = ((data[i++] & 0xFF) << 8) | (data[i++] & 0xFF);
            }
            if (i + len > data.length) {
                break;
            }
            byte[] val = Arrays.copyOfRange(data, i, i + len);
            i += len;
            out.add(new Tlv(tag, val));
        }
        return out;
    }

    private static byte[] encodeTlvs(List<Tlv> tlvs) {
        int n = 0;
        for (Tlv t : tlvs) {
            n += t.encode().length;
        }
        byte[] out = new byte[n];
        int o = 0;
        for (Tlv t : tlvs) {
            byte[] e = t.encode();
            System.arraycopy(e, 0, out, o, e.length);
            o += e.length;
        }
        return out;
    }

    private static final class Tlv {
        final int tag;
        final byte[] value;

        Tlv(int tag, byte[] value) {
            this.tag = tag;
            this.value = value;
        }

        byte[] encode() {
            byte[] lenBytes;
            if (value.length < 0x80) {
                lenBytes = new byte[] {(byte) value.length};
            } else if (value.length <= 0xFF) {
                lenBytes = new byte[] {(byte) 0x81, (byte) value.length};
            } else {
                lenBytes = new byte[] {(byte) 0x82, (byte) ((value.length >> 8) & 0xFF), (byte) (value.length & 0xFF)};
            }
            byte[] out = new byte[1 + lenBytes.length + value.length];
            out[0] = (byte) tag;
            System.arraycopy(lenBytes, 0, out, 1, lenBytes.length);
            System.arraycopy(value, 0, out, 1 + lenBytes.length, value.length);
            return out;
        }
    }
}
