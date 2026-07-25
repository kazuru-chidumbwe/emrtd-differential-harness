package org.emrtd.harness.jmrtd;

import java.util.LinkedHashMap;
import java.util.Locale;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/** Cross-library failure taxonomy for AA/TA/EAC result JSON (matches Go internal/normfail). */
public final class NormalizedFailure {
    private NormalizedFailure() {}

    public static final String CLASS_CHIP_SW_REJECT = "chip_sw_reject";
    public static final String CLASS_PEER_UNSUPPORTED = "peer_unsupported";
    public static final String CLASS_PROTOCOL_EXCEPTION = "protocol_exception";

    private static final Pattern HEX_SW = Pattern.compile("(?i)(?:SW\\s*=\\s*0x|0x)([0-9a-f]{4})");
    private static final Pattern STATUS = Pattern.compile("(?i)status[:\\s]*([0-9a-f]{4})");

    public static String extractSw(String msg) {
        if (msg == null || msg.isEmpty()) {
            return "";
        }
        Matcher m = HEX_SW.matcher(msg);
        if (m.find()) {
            return m.group(1).toUpperCase(Locale.ROOT);
        }
        m = STATUS.matcher(msg);
        if (m.find()) {
            return m.group(1).toUpperCase(Locale.ROOT);
        }
        return "";
    }

    public static Map<String, Object> chipSw(String mechanism, String step, String sw, boolean surfaced) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("mechanism", mechanism);
        m.put("step", step);
        m.put("iso7816_sw", sw == null ? "" : sw.toUpperCase(Locale.ROOT));
        m.put("failure_class", CLASS_CHIP_SW_REJECT);
        m.put("surfaced", surfaced);
        return m;
    }

    public static Map<String, Object> fromErr(String mechanism, String step, String errMsg, boolean surfaced) {
        String sw = extractSw(errMsg);
        if (!sw.isEmpty()) {
            return chipSw(mechanism, step, sw, surfaced);
        }
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("mechanism", mechanism);
        m.put("step", step);
        m.put("failure_class", CLASS_PROTOCOL_EXCEPTION);
        m.put("surfaced", surfaced);
        return m;
    }

    public static Map<String, Object> peerUnsupported(String mechanism) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("mechanism", mechanism);
        m.put("step", "n/a");
        m.put("failure_class", CLASS_PEER_UNSUPPORTED);
        m.put("surfaced", true);
        return m;
    }
}
