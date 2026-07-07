package org.emrtd.harness.jmrtd;

import com.google.gson.annotations.SerializedName;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

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
    public Injection injection;
    @SerializedName("ca_injection")
    public CaInjection caInjection;

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
    }

    public static final class CaInjection {
        @SerializedName("ca_sw")
        public String caSw;
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
