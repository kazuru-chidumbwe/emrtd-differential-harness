package org.emrtd.harness.jmrtd;

import net.sf.scuba.smartcards.CardService;
import net.sf.scuba.smartcards.CardServiceException;
import net.sf.scuba.smartcards.CommandAPDU;
import net.sf.scuba.smartcards.ISO7816;
import net.sf.scuba.smartcards.ResponseAPDU;
import org.jmrtd.BACKey;

import java.util.ArrayList;
import java.util.List;

/**
 * Synthetic chip: BAC succeeds and establishes SM; CA MSE:Set AT fails with injected SW
 * (SM-wrapped). Post-CA continue-check requires SM-validated READ BINARY.
 */
public final class TcCa01CardService extends CardService {
    private static final long serialVersionUID = 1L;

    private final BacChipSimulator bac;
    private final int caSw;
    private boolean open;
    private boolean caFailed;
    private final List<TcAc01CardService.TraceEntry> trace = new ArrayList<>();

    public TcCa01CardService(BACKey bacKey, String caStatusWordHex) {
        this.bac = new BacChipSimulator(bacKey);
        this.caSw = Integer.parseInt(caStatusWordHex, 16);
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
        int cla = capdu.getCLA() & 0xFF;
        int ins = capdu.getINS() & 0xFF;
        int p1 = capdu.getP1() & 0xFF;
        int p2 = capdu.getP2() & 0xFF;
        String label;
        byte[] body;
        boolean ok;
        try {
            boolean smCla = (cla & 0x0C) == 0x0C;

            // ISO 7816-4: 6987 = Expected SM data objects missing.
            if (bac.sessionEstablished() && ins == 0xB0 && !smCla) {
                label = "READ BINARY (unprotected rejected)";
                body = swBytes(0x6987);
                ok = false;
            } else if (bac.sessionEstablished() && smCla) {
                ChipSecureMessaging sm = bac.session();
                CommandAPDU plain = sm.unwrapCommand(capdu);
                int pIns = plain.getINS() & 0xFF;
                int pP1 = plain.getP1() & 0xFF;
                int pP2 = plain.getP2() & 0xFF;
                if (pIns == 0x22 && pP1 == 0x41 && pP2 == 0xA4) {
                    label = "MSE:Set AT (CA)";
                    if (!caFailed) {
                        caFailed = true;
                        body = sm.wrapResponse(caSw, new byte[0]).getBytes();
                    } else {
                        body = sm.wrapResponse(ISO7816.SW_INS_NOT_SUPPORTED, new byte[0]).getBytes();
                    }
                    ok = false;
                } else if (caFailed && pIns == 0xB0) {
                    label = "READ BINARY (post-CA SM continue)";
                    body = sm.wrapResponse(ISO7816.SW_NO_ERROR,
                            new byte[] {0x61, 0x03, 0x5F, 0x2E, 0x00}).getBytes();
                    ok = true;
                } else {
                    label = "SM unknown INS";
                    body = sm.wrapResponse(ISO7816.SW_INS_NOT_SUPPORTED, new byte[0]).getBytes();
                    ok = false;
                }
            } else if (ins == 0xA4) {
                label = "SELECT";
                body = swBytes(ISO7816.SW_NO_ERROR);
                ok = true;
            } else if (ins == 0x22 && p1 == 0x41 && p2 == 0xA4) {
                label = "MSE:Set AT (CA)";
                if (!caFailed) {
                    caFailed = true;
                    body = swBytes(caSw);
                } else {
                    body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
                }
                ok = false;
            } else if (ins == 0x22 || ins == 0x86) {
                label = ins == 0x22 ? "MSE:Set AT" : "General Authenticate";
                body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
                ok = false;
            } else if (ins == 0x84) {
                label = "Get Challenge";
                body = concat(bac.getChallenge(), swBytes(ISO7816.SW_NO_ERROR));
                ok = true;
            } else if (ins == 0x82) {
                label = "External Authenticate";
                body = concat(bac.mutualAuthResponse(capdu.getData()), swBytes(ISO7816.SW_NO_ERROR));
                ok = true;
            } else {
                label = "Unknown";
                body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
                ok = false;
            }
        } catch (Exception e) {
            throw new CardServiceException(e.toString());
        }
        ResponseAPDU rapdu = new ResponseAPDU(body);
        trace.add(new TcAc01CardService.TraceEntry(label, toHex(capdu.getBytes()), toHex(body), ok));
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
