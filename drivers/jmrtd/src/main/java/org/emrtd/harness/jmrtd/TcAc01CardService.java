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
 * Synthetic PACE-fail / BAC-ok card service.
 *
 * <p>paceChannel {@code sw} (default) returns configured status word.
 * {@code timeout} / {@code no_response} / {@code transport_abort} throw
 * CardServiceException (APDU-boundary adversarial-channel abstractions; not RF).
 * BAC remains capable after PACE failure.
 */
public final class TcAc01CardService extends CardService {
    private static final long serialVersionUID = 1L;

    public static final class TraceEntry {
        public final String label;
        public final String capdu;
        public final String rapdu;
        public final boolean success;

        TraceEntry(String label, String capdu, String rapdu, boolean success) {
            this.label = label;
            this.capdu = capdu;
            this.rapdu = rapdu;
            this.success = success;
        }
    }

    private final BacChipSimulator bac;
    private final int paceSw;
    private final String paceFailOn;
    private final String paceChannel;
    private boolean open;
    private boolean paceFailed;
    private final List<TraceEntry> trace = new ArrayList<>();

    public TcAc01CardService(BACKey bacKey, String paceStatusWordHex) {
        this(bacKey, paceStatusWordHex, "mse_set_at", "sw");
    }

    public TcAc01CardService(BACKey bacKey, String paceStatusWordHex, String paceFailOn, String paceChannel) {
        this.bac = new BacChipSimulator(bacKey);
        this.paceSw = Integer.parseInt(paceStatusWordHex, 16);
        this.paceFailOn = "general_authenticate".equals(paceFailOn) ? "general_authenticate" : "mse_set_at";
        this.paceChannel = (paceChannel == null || paceChannel.isEmpty()) ? "sw" : paceChannel;
    }

    public List<TraceEntry> trace() {
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

    /** Deliver PACE failure via SW return body, or throw for channel modes. */
    private byte[] failPace(String label, CommandAPDU capdu) throws CardServiceException {
        paceFailed = true;
        if ("timeout".equals(paceChannel)) {
            try {
                Thread.sleep(50L);
            } catch (InterruptedException ie) {
                Thread.currentThread().interrupt();
            }
            trace.add(new TraceEntry(label, toHex(capdu.getBytes()), "CHANNEL_TIMEOUT", false));
            throw new CardServiceException("synthetic PACE timeout (adversarial-channel arm)");
        }
        if ("no_response".equals(paceChannel) || "transport_abort".equals(paceChannel)) {
            String tag = "no_response".equals(paceChannel) ? "CHANNEL_NO_RESPONSE" : "CHANNEL_TRANSPORT_ABORT";
            trace.add(new TraceEntry(label, toHex(capdu.getBytes()), tag, false));
            throw new CardServiceException("synthetic PACE " + paceChannel + " (adversarial-channel arm)");
        }
        return swBytes(paceSw);
    }

    @Override
    public ResponseAPDU transmit(CommandAPDU capdu) throws CardServiceException {
        int ins = capdu.getINS() & 0xFF;
        String label;
        byte[] body;
        try {
            if (ins == 0xA4) {
                label = "SELECT";
                body = swBytes(ISO7816.SW_NO_ERROR);
            } else if (ins == 0x22) {
                label = "MSE:Set AT";
                if ("mse_set_at".equals(paceFailOn) && !paceFailed) {
                    body = failPace(label, capdu);
                } else if ("general_authenticate".equals(paceFailOn)) {
                    body = swBytes(ISO7816.SW_NO_ERROR);
                } else {
                    body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
                }
            } else if (ins == 0x86) {
                label = "General Authenticate";
                if ("general_authenticate".equals(paceFailOn) && !paceFailed) {
                    body = failPace(label, capdu);
                } else if (paceFailed) {
                    // Defensive: PACE already failed; reject further GA.
                    if (!"sw".equals(paceChannel)) {
                        throw new CardServiceException("synthetic PACE already aborted");
                    }
                    body = swBytes(paceSw);
                } else {
                    body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
                }
            } else if (ins == 0x84) {
                label = "Get Challenge";
                body = concat(bac.getChallenge(), swBytes(ISO7816.SW_NO_ERROR));
            } else if (ins == 0x82) {
                label = "External Authenticate";
                byte[] cmd = capdu.getData();
                body = concat(bac.mutualAuthResponse(cmd), swBytes(ISO7816.SW_NO_ERROR));
            } else {
                label = "Unknown";
                body = swBytes(ISO7816.SW_INS_NOT_SUPPORTED);
            }
        } catch (CardServiceException e) {
            throw e;
        } catch (Exception e) {
            throw new CardServiceException(e.toString());
        }

        ResponseAPDU rapdu = new ResponseAPDU(body);
        boolean success = rapdu.getSW() == 0x9000;
        trace.add(new TraceEntry(label, toHex(capdu.getBytes()), toHex(body), success));
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
        return e != null && e.getMessage() != null && e.getMessage().contains("adversarial-channel");
    }
}
