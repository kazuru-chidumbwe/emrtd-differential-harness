package org.emrtd.harness.jmrtd;

import com.google.gson.annotations.SerializedName;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;

public final class HarnessProfile {
    public String id;
    public String mechanism;
    public String condition;
    public String tier;
    public Mrz mrz;
    @SerializedName("card_access_hex")
    public String cardAccessHex;
    @SerializedName("dg14_hex_path")
    public String dg14HexPath;
    @SerializedName("dg15_hex_path")
    public String dg15HexPath;
    public Injection injection;
    @SerializedName("ca_injection")
    public CaInjection caInjection;
    @SerializedName("aa_injection")
    public AaInjection aaInjection;
    @SerializedName("ta_injection")
    public TaInjection taInjection;
    @SerializedName("peer_support")
    public Map<String, String> peerSupport;

    public static final class Mrz {
        @SerializedName("document_number")
        public String documentNumber;
        @SerializedName("date_of_birth")
        public String dateOfBirth;
        @SerializedName("date_of_expiry")
        public String dateOfExpiry;
    }

    public static final class Injection {
        @SerializedName("pace_sw")
        public String paceSw;
        @SerializedName("pace_fail_on")
        public String paceFailOn;
        @SerializedName("pace_channel")
        public String paceChannel;
    }

    public static final class CaInjection {
        @SerializedName("ca_sw")
        public String caSw;
    }

    public static final class AaInjection {
        @SerializedName("aa_sw")
        public String aaSw;
    }

    public static final class TaInjection {
        @SerializedName("ta_sw")
        public String taSw;
        @SerializedName("ta_fail_on")
        public String taFailOn;
    }

    public static HarnessProfile load(Path path) throws IOException {
        String raw = Files.readString(path);
        HarnessProfile profile = new com.google.gson.Gson().fromJson(raw, HarnessProfile.class);
        if (profile == null || profile.id == null || profile.mrz == null) {
            throw new IOException("profile missing required fields");
        }
        return profile;
    }
}
