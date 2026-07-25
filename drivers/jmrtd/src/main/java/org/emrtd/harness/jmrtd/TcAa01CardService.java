package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardService;
import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.smartcards.CommandAPDU;
import net.sf.scuba.smartcards.ISO7816;
import net.sf.scuba.smartcards.ResponseAPDU;
import org.jmrtd.BACKey;

import java.util.ArrayList;
import java.util.List;

/** Synthetic chip: BAC succeeds; first INTERNAL AUTHENTICATE returns injected status word. */
public final class TcAa01CardService extends CardService {
    private static final long serialVersionUID = 1L;

    private final BacChipSimulator bac;
    private final int aaSw;
    private boolean open;
    private boolean aaFailed;
    private final List<TcAc01CardService.TraceEntry> trace = new ArrayList<>();

    public TcAa01CardService(BACKey bacKey, String aaStatusWordHex) {
        this.bac = new BacChipSimulator(bacKey);
        this.aaSw = Integer.parseInt(aaStatusWordHex, 16);
    }

    public List<TcAc01CardService.TraceEntry> trace() {
        return trace;
    }

    @Override
    public void open() {
        open = true;
    }

    @Override
    public boolean isOpen() {
        return open;
    }

    @Override
    public byte[] getATR() {
        return new byte[] {0x3B, (byte) 0x88, (byte) 0x80, 0x01};
    }

    @Override
    public void close() {
        open = false;
    }

    @Override
    public ResponseAPDU transmit(CommandAPDU capdu) throws CardServiceException {
        int ins = capdu.getINS() & 0xFF;
        int p1 = capdu.getP1() & 0xFF;
        int p2 = capdu.getP2() & 0xFF;
        String label;
        byte[] body;
        try {
            if (ins == 0xA4) {
                label = "SELECT";
                body = swBytes(ISO7816.SW_NO_ERROR);
            } else if (ins == 0x88) {
                label = "Internal Authenticate (AA)";
                if (!aaFailed) {
                    aaFailed = true;
                    body = swBytes(aaSw);
                } else {
                    body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
                }
            } else if (ins == 0x22 || ins == 0x86) {
                label = ins == 0x22 ? "MSE:Set AT" : "General Authenticate";
                body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
            } else if (ins == 0x84) {
                label = "Get Challenge";
                body = concat(bac.getChallenge(), swBytes(ISO7816.SW_NO_ERROR));
            } else if (ins == 0x82) {
                label = "External Authenticate";
                body = concat(bac.mutualAuthResponse(capdu.getData()), swBytes(ISO7816.SW_NO_ERROR));
            } else {
                label = "Unknown";
                body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
            }
        } catch (Exception e) {
            throw new CardServiceException(e.toString());
        }
        ResponseAPDU rapdu = new ResponseAPDU(body);
        trace.add(new TcAc01CardService.TraceEntry(label, toHex(capdu.getBytes()), toHex(body), rapdu.getSW() == 0x9000));
        return rapdu;
    }

    private static byte[] swBytes(int sw) {
        return new byte[] {(byte) ((sw >> 8) & 0xFF), (byte) (sw & 0xFF)};
    }

    private static byte[] concat(byte[] a, byte[] b) {
        byte[] out = new byte[a.length + b.length];
        System.arraycopy(a, 0, out, 0, a.length);
        System.arraycopy(b, 0, out, a.length, b.length);
        return out;
    }

    private static String toHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder(bytes.length * 2);
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    @Override
    public boolean isConnectionLost(Exception e) {
        return false;
    }
}
